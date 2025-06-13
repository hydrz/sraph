//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	_ "embed"

	"github.com/goccy/go-yaml"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

//go:embed webgpu.tmpl
var tmpl string

const (
	defaultSchemaURL = "https://raw.githubusercontent.com/webgpu-native/webgpu-headers/refs/heads/main/schema.json"
	defaultYamlURL   = "https://raw.githubusercontent.com/webgpu-native/webgpu-headers/refs/heads/main/webgpu.yml"
)

// Yml structure and related types
type Yml struct {
	Copyright  string `yaml:"copyright"`
	Name       string `yaml:"name"`
	Doc        string `yaml:"doc"`
	EnumPrefix uint16 `yaml:"enum_prefix"`

	Constants []Constant `yaml:"constants"`
	Typedefs  []Typedef  `yaml:"typedefs"`
	Enums     []Enum     `yaml:"enums"`
	Bitflags  []Bitflag  `yaml:"bitflags"`
	Structs   []Struct   `yaml:"structs"`
	Callbacks []Callback `yaml:"callbacks"`
	Functions []Function `yaml:"functions"`
	Objects   []Object   `yaml:"objects"`
}

type Base struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Doc       string `yaml:"doc"`
	Extended  bool   `yaml:"extended"`
}

type Constant struct {
	Base  `yaml:",inline"`
	Value string `yaml:"value"`
}

type Typedef struct {
	Base `yaml:",inline"`
	Type string `yaml:"type"`
}

type Enum struct {
	Base    `yaml:",inline"`
	Entries []*EnumEntry `yaml:"entries"`
}
type EnumEntry struct {
	Base  `yaml:",inline"`
	Value string `yaml:"value"`
}

type Bitflag struct {
	Base    `yaml:",inline"`
	Entries []BitflagEntry `yaml:"entries"`
}
type BitflagEntry struct {
	Base             `yaml:",inline"`
	Value            string   `yaml:"value"`
	ValueCombination []string `yaml:"value_combination"`
}

type PointerType string

const (
	PointerTypeMutable   PointerType = "mutable"
	PointerTypeImmutable PointerType = "immutable"
)

type ParameterType struct {
	Name                string      `yaml:"name"`
	Namespace           string      `yaml:"namespace"`
	Doc                 string      `yaml:"doc"`
	Type                string      `yaml:"type"`
	PassedWithOwnership *bool       `yaml:"passed_with_ownership"`
	Pointer             PointerType `yaml:"pointer"`
	Optional            bool        `yaml:"optional"`
	Default             *string     `yaml:"default"`
}

type Callback struct {
	Base  `yaml:",inline"`
	Style string          `yaml:"style"`
	Args  []ParameterType `yaml:"args"`
}

type Function struct {
	Base     `yaml:",inline"`
	Returns  *ParameterType  `yaml:"returns"`
	Callback *string         `yaml:"callback"`
	Args     []ParameterType `yaml:"args"`
}

type Struct struct {
	Base        `yaml:",inline"`
	Type        string          `yaml:"type"`
	FreeMembers bool            `yaml:"free_members"`
	Members     []ParameterType `yaml:"members"`
	Extends     []string        `yaml:"extends"`
}

type Object struct {
	Base    `yaml:",inline"`
	Methods []Function `yaml:"methods"`

	IsStruct bool `yaml:"-"`
}

// Global variables
var arrayTypeRegexp = regexp.MustCompile(`array<([a-zA-Z0-9._]+)>`)

var (
	schemaPath string
	goPaths    StringListFlag
	yamlPaths  StringListFlag
	extPrefix  bool
)

// StringListFlag implementation
type StringListFlag []string

var StringListFlagI flag.Value = &StringListFlag{}

func (f *StringListFlag) String() string     { return fmt.Sprintf("%#v", f) }
func (f *StringListFlag) Set(v string) error { *f = append(*f, v); return nil }

// Comment types
type CommentType uint8

const (
	CommentTypeSingleLine CommentType = iota
	CommentTypeMultiLine
)

// Generator structure
type Generator struct {
	ExtPrefix string
	GoName    string
	*Yml
}

// Utility functions - only keep those used in templates
func Comment(in string, mode CommentType, indent int, newline bool) string {
	if in == "" || strings.TrimSpace(in) == "TODO" {
		return ""
	}

	const space = ' '
	var out strings.Builder
	if newline {
		out.WriteString("\n")
	}
	if mode == CommentTypeMultiLine {
		for i := 0; i < indent; i++ {
			out.WriteRune(space)
		}
		out.WriteString("/**\n")
	}
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(in)))
	for sc.Scan() {
		line := sc.Text()
		for i := 0; i < indent; i++ {
			out.WriteRune(space)
		}
		switch mode {
		case CommentTypeSingleLine:
			out.WriteString("//")
		case CommentTypeMultiLine:
			out.WriteString(" *")
		default:
			panic("unreachable")
		}
		if line != "" {
			out.WriteString(" ")
			out.WriteString(line)
		}
		out.WriteString("\n")
	}
	if mode == CommentTypeMultiLine {
		for i := 0; i < indent; i++ {
			out.WriteRune(space)
		}
		out.WriteString(" */")
	}
	return out.String()
}

func PascalCase(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	nextUpper := true
	for _, c := range s {
		if nextUpper {
			out.WriteRune(unicode.ToUpper(c))
			nextUpper = false
		} else {
			if c == '_' {
				nextUpper = true
			} else {
				out.WriteRune(c)
			}
		}
	}
	return out.String()
}

func CamelCase(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	nextUpper := false
	for _, c := range s {
		if nextUpper {
			out.WriteRune(unicode.ToUpper(c))
			nextUpper = false
		} else {
			if c == '_' {
				nextUpper = true
			} else {
				out.WriteRune(c)
			}
		}
	}
	return out.String()
}

func UpperConcatCase(s string) string {
	return strings.ToUpper(PascalCase(s))
}

// Struct sorting function - needed for SortAndTransform
func SortStructs(structs []Struct) {
	type node struct {
		visited    bool
		depth      int
		Extensions []string
		Struct
	}
	nodeMap := make(map[string]node, len(structs))
	for _, s := range structs {
		node := nodeMap[s.Name]
		node.Struct = s
		nodeMap[s.Name] = node
		if s.Type == "extension" {
			for _, extend := range s.Extends {
				parent := nodeMap[extend]
				parent.Extensions = append(parent.Extensions, s.Name)
				nodeMap[extend] = parent
			}
		}
	}

	slices.SortStableFunc(structs, func(a, b Struct) int {
		return strings.Compare(PascalCase(a.Name), PascalCase(b.Name))
	})

	var computeDepth func(string) int
	computeDepth = func(name string) int {
		node, ok := nodeMap[name]
		if !ok {
			panic("found invalid non-existing type: " + name)
		}

		if node.visited {
			return node.depth
		}

		maxDependentDepth := 0
		for _, member := range node.Members {
			if strings.HasPrefix(member.Type, "struct.") {
				dependentDepth := computeDepth(strings.TrimPrefix(member.Type, "struct."))
				if dependentDepth+1 > maxDependentDepth {
					maxDependentDepth = dependentDepth + 1
				}
			} else {
				matches := arrayTypeRegexp.FindStringSubmatch(member.Type)
				if len(matches) == 2 {
					typ := matches[1]
					if strings.HasPrefix(typ, "struct.") {
						dependentDepth := computeDepth(strings.TrimPrefix(typ, "struct."))
						if dependentDepth+1 > maxDependentDepth {
							maxDependentDepth = dependentDepth + 1
						}
					}
				}
			}
		}
		for _, extension := range node.Extensions {
			dependentDepth := computeDepth(extension)
			if dependentDepth+1 > maxDependentDepth {
				maxDependentDepth = dependentDepth + 1
			}
		}

		node.depth = maxDependentDepth
		node.visited = true
		nodeMap[name] = node
		return maxDependentDepth
	}

	for _, s := range structs {
		computeDepth(s.Name)
	}
	slices.SortStableFunc(structs, func(a, b Struct) int {
		return nodeMap[a.Name].depth - nodeMap[b.Name].depth
	})
}

// Validation functions - needed for main
func ValidateYamls(schemaPath string, yamlPaths []string) error {
	// Validation through json schema
	for _, yamlPath := range yamlPaths {
		yamlFile, err := os.ReadFile(yamlPath)
		if err != nil {
			return fmt.Errorf("ValidateYaml: %w", err)
		}
		var yml map[string]any
		if err := yaml.Unmarshal(yamlFile, &yml); err != nil {
			return fmt.Errorf("ValidateYaml: %w", err)
		}

		schema := jsonschema.MustCompile(schemaPath)
		if err := schema.Validate(yml); err != nil {
			return fmt.Errorf("ValidateYaml: %w", err)
		}
	}

	// Validation of possible duplication of entries across multiple yaml files
	if err := mergeAndValidateDuplicates(yamlPaths); err != nil {
		panic(err)
	}

	// TODO: add dependency check validations
	return nil
}

func mergeAndValidateDuplicates(yamlPaths []string) (errs error) {
	constants := make(map[string]Constant)
	enums := make(map[string]Enum)
	bitflags := make(map[string]Bitflag)
	structs := make(map[string]Struct)
	callbacks := make(map[string]Callback)
	functions := make(map[string]Function)
	objects := make(map[string]Object)

	for _, yamlPath := range yamlPaths {
		src, err := os.ReadFile(yamlPath)
		if err != nil {
			panic(err)
		}

		var data Yml
		if err := yaml.Unmarshal(src, &data); err != nil {
			panic(err)
		}

		for _, c := range data.Constants {
			if _, ok := constants[c.Name]; ok {
				errs = errors.Join(errs, fmt.Errorf("merge: constants.%s in %s was already found previously while parsing, duplicates are not allowed", c.Name, yamlPath))
			}
			constants[c.Name] = c
		}
		for _, e := range data.Enums {
			if prevEnum, ok := enums[e.Name]; ok {
				if !e.Extended {
					errs = errors.Join(errs, fmt.Errorf("merge: enums.%s in %s is being extended but isn't marked as one", e.Name, yamlPath))
				}
				for _, entry := range e.Entries {
					if entry != nil {
						if slices.ContainsFunc(prevEnum.Entries, func(e *EnumEntry) bool { return e != nil && e.Name == entry.Name }) {
							errs = errors.Join(errs, fmt.Errorf("merge: enums.%s.%s in %s was already found previously while parsing, duplicates are not allowed", e.Name, entry.Name, yamlPath))
						}
						prevEnum.Entries = append(prevEnum.Entries, entry)
					}
				}
				enums[e.Name] = prevEnum
			} else {
				enums[e.Name] = e
			}
		}
		for _, bf := range data.Bitflags {
			if prevBf, ok := bitflags[bf.Name]; ok {
				if !bf.Extended {
					errs = errors.Join(errs, fmt.Errorf("merge: bitflags.%s in %s is being extended but isn't marked as one", bf.Name, yamlPath))
				}
				for _, entry := range bf.Entries {
					if slices.ContainsFunc(prevBf.Entries, func(e BitflagEntry) bool { return e.Name == entry.Name }) {
						errs = errors.Join(errs, fmt.Errorf("merge: bitflags.%s.%s in %s was already found previously while parsing, duplicates are not allowed", bf.Name, entry.Name, yamlPath))
					}
					if entry.Value == "" && len(entry.ValueCombination) == 0 {
						errs = errors.Join(errs, fmt.Errorf("merge: bitflags.%s.%s in %s was extended but doesn't have a value or value_combination, extended bitflag entries must have an explicit value", bf.Name, entry.Name, yamlPath))
					}
					prevBf.Entries = append(prevBf.Entries, entry)
				}
				bitflags[bf.Name] = prevBf
			} else {
				bitflags[bf.Name] = bf
			}
		}
		for _, c := range data.Callbacks {
			if _, ok := callbacks[c.Name]; ok {
				errs = errors.Join(errs, fmt.Errorf("merge: callbacks.%s in %s was already found previously while parsing, duplicates are not allowed", c.Name, yamlPath))
			}
			callbacks[c.Name] = c
		}
		for _, s := range data.Structs {
			if _, ok := structs[s.Name]; ok {
				errs = errors.Join(errs, fmt.Errorf("merge: structs.%s in %s was already found previously while parsing, duplicates are not allowed", s.Name, yamlPath))
			}
			structs[s.Name] = s
		}
		for _, f := range data.Functions {
			if _, ok := functions[f.Name]; ok {
				errs = errors.Join(errs, fmt.Errorf("merge: functions.%s in %s was already found previously while parsing, duplicates are not allowed", f.Name, yamlPath))
			}
			functions[f.Name] = f
		}
		for _, o := range data.Objects {
			if prevObj, ok := objects[o.Name]; ok {
				if !o.Extended {
					errs = errors.Join(errs, fmt.Errorf("merge: objects.%s in %s is being extended but isn't marked as one", o.Name, yamlPath))
				}
				for _, method := range o.Methods {
					if slices.ContainsFunc(prevObj.Methods, func(f Function) bool { return f.Name == method.Name }) {
						errs = errors.Join(errs, fmt.Errorf("merge: objects.%s.%s in %s was already found previously while parsing, duplicates are not allowed", o.Name, method.Name, yamlPath))
					}
					prevObj.Methods = append(prevObj.Methods, method)
				}
				objects[o.Name] = prevObj
			} else {
				objects[o.Name] = o
			}
		}
	}
	return
}

// Sorting and transformation functions - needed for main
func SortAndTransform(yml *Yml) {
	// Sort structs
	SortStructs(yml.Structs)

	// Sort constants
	slices.SortStableFunc(yml.Constants, func(a, b Constant) int {
		return strings.Compare(UpperConcatCase(a.Name), UpperConcatCase(b.Name))
	})

	// Sort enums
	slices.SortStableFunc(yml.Enums, func(a, b Enum) int {
		// We want to generate extended enum declarations before the normal ones.
		if a.Extended && !b.Extended {
			return -1
		} else if !a.Extended && b.Extended {
			return 1
		}
		return strings.Compare(UpperConcatCase(a.Name), UpperConcatCase(b.Name))
	})

	// Sort bitflags
	slices.SortStableFunc(yml.Bitflags, func(a, b Bitflag) int {
		// We want to generate extended bitflag declarations before the normal ones.
		if a.Extended && !b.Extended {
			return -1
		} else if !a.Extended && b.Extended {
			return 1
		}
		return strings.Compare(UpperConcatCase(a.Name), UpperConcatCase(b.Name))
	})

	// Add free_member function for relevant structs
	for _, s := range yml.Structs {
		if s.FreeMembers {
			yml.Objects = append(yml.Objects, Object{
				Base:     Base{Name: s.Name, Namespace: s.Namespace},
				IsStruct: true,
				Methods: []Function{{
					Base: Base{
						Name:      "free_members",
						Namespace: s.Namespace,
						Doc:       "Frees members which were allocated by the API."},
				}},
			})
		}
	}

	// Sort objects
	slices.SortStableFunc(yml.Objects, func(a, b Object) int {
		return strings.Compare(UpperConcatCase(a.Name), UpperConcatCase(b.Name))
	})

	// Sort methods
	for _, obj := range yml.Objects {
		slices.SortStableFunc(obj.Methods, func(a, b Function) int {
			return strings.Compare(UpperConcatCase(a.Name), UpperConcatCase(b.Name))
		})
	}
}

// Generator methods
func (g *Generator) Gen(dst io.Writer) error {
	t := template.
		New("").
		Funcs(template.FuncMap{
			"SComment":  func(v string, indent int) string { return Comment(v, CommentTypeSingleLine, indent, true) },
			"SCommentN": func(v string, indent int) string { return Comment(v, CommentTypeSingleLine, indent, false) },
			"IsArray": func(typ string) bool {
				return arrayTypeRegexp.Match([]byte(typ))
			},
			"IsLast": func(i int, s any) bool { return i == reflect.ValueOf(s).Len()-1 },
			// Go template functions - only the ones actually used
			"GoConstantName":      g.GoConstantName,
			"GoTypeName":          g.GoTypeName,
			"GoEnumName":          g.GoEnumName,
			"GoValue":             g.GoValue,
			"GoType":              g.GoType,
			"GoFunctionName":      g.GoFunctionName,
			"GoParameterName":     g.GoParameterName,
			"GoFunctionArgs":      g.GoFunctionArgs,
			"GoFunctionReturns":   g.GoFunctionReturns,
			"GoStructMember":      g.GoStructMember,
			"GoStructMemberArray": g.GoStructMemberArray,
			"EnumValue32":         g.EnumValue32,
			"BitflagValue": func(b Bitflag, entryIndex int) (string, error) {
				return g.BitflagValue(b, entryIndex, false)
			},
		})
	t, err := t.Parse(tmpl)
	if err != nil {
		return fmt.Errorf("GenCHeader: failed to parse template: %w", err)
	}

	// Render template to buffer
	var buf bytes.Buffer
	if err := t.Execute(&buf, g); err != nil {
		return fmt.Errorf("GenCHeader: failed to execute template: %w", err)
	}

	// Format Go code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write unformatted code for easier debugging
		_, _ = dst.Write(buf.Bytes())
		return fmt.Errorf("GenCHeader: gofmt failed: %w", err)
	}

	_, err = dst.Write(formatted)
	if err != nil {
		return fmt.Errorf("GenCHeader: failed to write formatted code: %w", err)
	}
	return nil
}

func (g *Generator) FindBaseType(typ string) Base {
	// Handle type names prefixed with the type category.
	category, name, found := strings.Cut(typ, ".")
	if !found {
		panic("Cannot find base type for invalid type identifier: " + typ)
	}

	switch category {
	case "constant":
		idx := slices.IndexFunc(g.Constants, func(c Constant) bool { return c.Name == name })
		if idx == -1 {
			return Base{Name: name, Namespace: g.PrefixForNamespace("")}
		}
		return g.Constants[idx].Base
	case "typedef":
		idx := slices.IndexFunc(g.Typedefs, func(t Typedef) bool { return t.Name == name })
		if idx == -1 {
			return Base{Name: name, Namespace: g.PrefixForNamespace("")}
		}
		return g.Typedefs[idx].Base
	case "enum":
		idx := slices.IndexFunc(g.Enums, func(e Enum) bool { return e.Name == name })
		if idx == -1 {
			return Base{Name: name, Namespace: g.PrefixForNamespace("")}
		}
		return g.Enums[idx].Base
	case "bitflag":
		idx := slices.IndexFunc(g.Bitflags, func(b Bitflag) bool { return b.Name == name })
		if idx == -1 {
			return Base{Name: name, Namespace: g.PrefixForNamespace("")}
		}
		return g.Bitflags[idx].Base
	case "struct":
		idx := slices.IndexFunc(g.Structs, func(s Struct) bool { return s.Name == name })
		if idx == -1 {
			return Base{Name: name, Namespace: g.PrefixForNamespace("")}
		}
		return g.Structs[idx].Base
	case "callback":
		idx := slices.IndexFunc(g.Callbacks, func(c Callback) bool { return c.Name == name })
		if idx == -1 {
			return Base{Name: name, Namespace: g.PrefixForNamespace("")}
		}
		return g.Callbacks[idx].Base
	case "object":
		idx := slices.IndexFunc(g.Objects, func(o Object) bool { return o.Name == name })
		if idx == -1 {
			return Base{Name: name, Namespace: g.PrefixForNamespace("")}
		}
		return g.Objects[idx].Base
	default:
		panic("Unable to find unknown category type: " + category + " for identifier: " + typ)
	}
}

func (g *Generator) PrefixForNamespace(namespace string) string {
	switch namespace {
	case "":
		return g.ExtPrefix
	case "webgpu":
		return ""
	default:
		return namespace
	}
}

func (g *Generator) EnumValue32(e Enum, entryIndex int) (uint32, error) {
	entry := e.Entries[entryIndex]
	var value16 uint16
	if entry.Value == "" {
		value16 = uint16(entryIndex)
	} else {
		var num string
		var base int
		if strings.HasPrefix(entry.Value, "0x") {
			base = 16
			num = strings.TrimPrefix(entry.Value, "0x")
		} else {
			base = 10
			num = entry.Value
		}
		v, err := strconv.ParseUint(num, base, 16)
		if err != nil {
			return 0, err
		}
		value16 = uint16(v)
	}
	return uint32(g.EnumPrefix)<<16 | uint32(value16), nil
}

func bitflagEntryValue(entry BitflagEntry, entryIndex int) (uint64, error) {
	if entry.Value == "" {
		value := uint64(math.Pow(2, float64(entryIndex-1)))
		return value, nil
	} else {
		var num string
		var base int
		if strings.HasPrefix(entry.Value, "0x") {
			base = 16
			num = strings.TrimPrefix(entry.Value, "0x")
		} else {
			base = 10
			num = entry.Value
		}
		return strconv.ParseUint(num, base, 64)
	}
}

func (g *Generator) BitflagValue(b Bitflag, entryIndex int, isDocString bool) (string, error) {
	entry := b.Entries[entryIndex]

	var value uint64
	if len(entry.ValueCombination) > 0 {
		if entry.Value != "" {
			return "", fmt.Errorf("BitflagValue: found conflicting 'value' and 'value_combination' in '%s'", b.Name)
		}
		for _, v := range entry.ValueCombination {
			// find the value by searching in b, bitwise-OR it into the result
			for searchIndex, search := range b.Entries {
				if search.Name == v {
					searchValue, err := bitflagEntryValue(search, searchIndex)
					if err != nil {
						return "", err
					}
					value |= searchValue
					break
				}
			}
		}
	} else {
		var err error
		value, err = bitflagEntryValue(entry, entryIndex)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("0x%.16X", value), nil
}

// Go-specific generator methods - only the ones used in templates
func (g *Generator) GoConstantName(b Base) string {
	return PascalCase(b.Name)
}

func (g *Generator) GoTypeName(b Base) string {
	return PascalCase(b.Name)
}

func (g *Generator) GoEnumName(typ Base, entry Base) string {
	return PascalCase(typ.Name) + PascalCase(entry.Name)
}

func (g *Generator) GoValue(s string) string {
	switch s {
	case "usize_max":
		return "^uintptr(0)"
	case "uint32_max":
		return "^uint32(0)"
	case "uint64_max":
		return "^uint64(0)"
	case "nan":
		// Use a variable instead of function call in constant declaration
		return "math.NaN()"
	default:
		return s
	}
}

func (g *Generator) GoType(typ string) string {
	switch typ {
	case "bool":
		return "bool"
	case "uint16":
		return "uint16"
	case "uint32":
		return "uint32"
	case "uint64":
		return "uint64"
	case "usize":
		return "uintptr"
	case "int16":
		return "int16"
	case "int32":
		return "int32"
	case "float32", "nullable_float32":
		return "float32"
	case "float64", "float64_supertype":
		return "float64"
	case "nullable_string", "string_with_default_empty", "out_string":
		return "string"
	case "c_void":
		return "unsafe.Pointer"
	default:
		// Handle type names prefixed with the type category
		return g.GoTypeName(g.FindBaseType(typ))
	}
}

func (g *Generator) GoFunctionName(b Base) string {
	return PascalCase(b.Name)
}

func (g *Generator) GoParameterName(name string) string {
	// Convert to camelCase first
	camelName := CamelCase(name)

	// Handle Go keywords by appending underscore or using alternative names
	switch camelName {
	case "type":
		return "errorType"
	case "func":
		return "function"
	case "var":
		return "variable"
	case "const":
		return "constant"
	case "struct":
		return "structure"
	case "interface":
		return "iface"
	case "package":
		return "pkg"
	case "import":
		return "imp"
	case "return":
		return "ret"
	case "if":
		return "condition"
	case "else":
		return "alternative"
	case "for":
		return "loop"
	case "range":
		return "rng"
	case "switch":
		return "switchValue"
	case "case":
		return "caseValue"
	case "default":
		return "defaultValue"
	case "go":
		return "routine"
	case "defer":
		return "deferred"
	case "select":
		return "selector"
	case "chan":
		return "channel"
	case "map":
		return "mapping"
	default:
		return camelName
	}
}

func (g *Generator) GoFunctionArgs(f Function, o *Object) string {
	sb := &strings.Builder{}
	for _, arg := range f.Args {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		matches := arrayTypeRegexp.FindStringSubmatch(arg.Type)
		if len(matches) == 2 {
			fmt.Fprintf(sb, "%s []%s", g.GoParameterName(arg.Name), g.GoType(matches[1]))
		} else {
			fmt.Fprintf(sb, "%s %s", g.GoParameterName(arg.Name), g.GoType(arg.Type))
		}
	}

	if f.Callback != nil {
		if o != nil && o.IsStruct {
			// If the function is a method of an object, use the object type as the first argument
			if sb.Len() > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s *%s", g.GoParameterName("this"), g.GoType(o.Name)))
		} else {
			// If the function is not a method, use the callback type as the first argument
			if sb.Len() > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s %sCallbackInfo", g.GoParameterName("callback"), g.GoTypeName(g.FindBaseType(*f.Callback))))
		}
	}

	return sb.String()
}

func (g *Generator) GoFunctionReturns(f Function) string {
	if f.Callback != nil {
		return "Future"
	}
	if f.Returns != nil {
		return fmt.Sprintf("(%s, error)", g.GoType(f.Returns.Type))
	}
	return "error"
}

func (g *Generator) GoStructMember(s Struct, memberIndex int) string {
	member := s.Members[memberIndex]

	matches := arrayTypeRegexp.FindStringSubmatch(member.Type)
	if len(matches) == 2 {
		panic("GoStructMember used on array type")
	}

	// Handle callback types
	if strings.HasPrefix(member.Type, "callback.") {
		return fmt.Sprintf("%s %sCallbackInfo", PascalCase(member.Name), g.GoTypeName(g.FindBaseType(member.Type)))
	}

	return fmt.Sprintf("%s %s", PascalCase(member.Name), g.GoType(member.Type))
}

func (g *Generator) GoStructMemberArray(s Struct, memberIndex int) string {
	member := s.Members[memberIndex]

	matches := arrayTypeRegexp.FindStringSubmatch(member.Type)
	if len(matches) != 2 {
		panic("GoStructMemberArray used on non-array")
	}

	return fmt.Sprintf("%s []%s", PascalCase(member.Name), g.GoType(matches[1]))
}

func fetchFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: status %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// Main function
func main() {
	flag.StringVar(&schemaPath, "schema", "", "path or URL of the json schema")
	flag.Var(&yamlPaths, "yaml", "path or URL of the yaml spec")
	flag.Var(&goPaths, "go", "output path of the go")
	flag.BoolVar(&extPrefix, "extprefix", true, "append prefix to extension identifiers")
	flag.Parse()

	// Use default URLs if not provided
	if schemaPath == "" {
		schemaPath = defaultSchemaURL
	}
	if len(yamlPaths) == 0 {
		yamlPaths = append(yamlPaths, defaultYamlURL)
	}
	if len(goPaths) == 0 || len(goPaths) != len(yamlPaths) {
		flag.Usage()
		os.Exit(1)
	}

	// Order matters for validation steps, so enforce it.
	if len(yamlPaths) > 1 && filepath.Base(yamlPaths[0]) != "webgpu.yml" {
		panic(`"webgpu.yml" must be the first sequence in the order`)
	}

	// Download schema if it's a URL
	var schemaData []byte
	if strings.HasPrefix(schemaPath, "http://") || strings.HasPrefix(schemaPath, "https://") {
		var err error
		schemaData, err = fetchFile(schemaPath)
		if err != nil {
			panic(err)
		}
		// Save to a temporary file for jsonschema.MustCompile
		tmpFile, err := os.CreateTemp("", "schema-*.json")
		if err != nil {
			panic(err)
		}
		defer os.Remove(tmpFile.Name())
		if _, err := tmpFile.Write(schemaData); err != nil {
			panic(err)
		}
		if err := tmpFile.Close(); err != nil {
			panic(err)
		}
		schemaPath = tmpFile.Name()
	}

	// Download yamls if any are URLs
	yamlDatas := make([][]byte, len(yamlPaths))
	for i, yamlPath := range yamlPaths {
		if strings.HasPrefix(yamlPath, "http://") || strings.HasPrefix(yamlPath, "https://") {
			data, err := fetchFile(yamlPath)
			if err != nil {
				panic(err)
			}
			yamlDatas[i] = data
		} else {
			data, err := os.ReadFile(yamlPath)
			if err != nil {
				panic(err)
			}
			yamlDatas[i] = data
		}
	}

	// Validate the yaml files (jsonschema, duplications)
	// Save yamls to temp files for validation and further processing
	tmpYamlPaths := make([]string, len(yamlDatas))
	for i, data := range yamlDatas {
		tmpFile, err := os.CreateTemp("", "yaml-*.yml")
		if err != nil {
			panic(err)
		}
		defer os.Remove(tmpFile.Name())
		if _, err := tmpFile.Write(data); err != nil {
			panic(err)
		}
		if err := tmpFile.Close(); err != nil {
			panic(err)
		}
		tmpYamlPaths[i] = tmpFile.Name()
	}
	if err := ValidateYamls(schemaPath, tmpYamlPaths); err != nil {
		panic(err)
	}

	// Generate the go files
	for i := range yamlDatas {
		goPath := goPaths[i]
		goFileName := filepath.Base(goPath)
		goFileNameSplit := strings.Split(goFileName, ".")
		if len(goFileNameSplit) != 2 {
			panic("got invalid go file name: " + goFileName)
		}

		dst, err := os.Create(goPath)
		if err != nil {
			panic(err)
		}

		var yml Yml
		if err := yaml.Unmarshal(yamlDatas[i], &yml); err != nil {
			panic(err)
		}

		SortAndTransform(&yml)

		prefix := ""
		if yml.Name != "webgpu" && extPrefix {
			prefix = yml.Name
		}
		g := &Generator{
			Yml:       &yml,
			GoName:    goFileNameSplit[0],
			ExtPrefix: prefix,
		}
		if err := g.Gen(dst); err != nil {
			panic(err)
		}
	}
}
