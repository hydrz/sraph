//go:build ignore

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
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
	schemaPath  string
	headerPaths StringListFlag
	yamlPaths   StringListFlag
	extPrefix   bool
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
	ExtPrefix  string
	HeaderName string
	*Yml
}

// Utility functions
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

func ConstantCase(v string) string {
	return strings.ToUpper(v)
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

func Singularize(s string) string {
	switch s {
	case "entries":
		return "entry"
	default:
		return strings.TrimSuffix(s, "s")
	}
}

func TrimTypePrefix(s string) string {
	switch {
	case strings.HasPrefix(s, "enum."):
		return strings.TrimPrefix(s, "enum.")
	case strings.HasPrefix(s, "bitflag."):
		return strings.TrimPrefix(s, "bitflag.")
	case strings.HasPrefix(s, "struct."):
		return strings.TrimPrefix(s, "struct.")
	case strings.HasPrefix(s, "callback."):
		return strings.TrimPrefix(s, "callback.")
	case strings.HasPrefix(s, "object."):
		return strings.TrimPrefix(s, "object.")
	default:
		return ""
	}
}

func UpperConcatCase(s string) string {
	return strings.ToUpper(PascalCase(s))
}

// Struct sorting function
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

// Validation functions
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

// Sorting and transformation functions
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

	// Add add_ref and release methods for objects
	for i, o := range yml.Objects {
		if !o.Extended && !o.IsStruct {
			yml.Objects[i].Methods = append(yml.Objects[i].Methods,
				Function{
					Base: Base{
						Name:      "add_ref",
						Namespace: o.Namespace,
						Doc:       "TODO",
					},
				},
				Function{
					Base: Base{
						Name:      "release",
						Namespace: o.Namespace,
						Doc:       "TODO",
					},
				})
		}
	}
}

// Generator methods
func (g *Generator) Gen(dst io.Writer) error {
	t := template.
		New("").
		Funcs(template.FuncMap{
			"SComment":  func(v string, indent int) string { return Comment(v, CommentTypeSingleLine, indent, true) },
			"MComment":  func(v string, indent int) string { return Comment(v, CommentTypeMultiLine, indent, true) },
			"SCommentN": func(v string, indent int) string { return Comment(v, CommentTypeSingleLine, indent, false) },
			"MCommentN": func(v string, indent int) string { return Comment(v, CommentTypeMultiLine, indent, false) },
			"MCommentMainPage": func(v string, indent int) string {
				if v == "" || strings.TrimSpace(v) == "TODO" {
					return ""
				}
				return Comment("\\mainpage\n\n"+strings.TrimSpace(v), CommentTypeMultiLine, indent, true)
			},
			"MCommentEnumValue": func(v string, indent int, e Enum, entryIndex int) string {
				var s string
				v = strings.TrimSpace(v)
				if v != "" && v != "TODO" {
					s += v
				}
				value, _ := g.EnumValue32(e, entryIndex)
				if value == 0 {
					s = "`0`. " + s
				}
				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"MCommentBitflagType": func(v string, indent int) string {
				var s string
				v = strings.TrimSpace(v)
				if v != "" && v != "TODO" {
					s += v
				}
				s += "\n\nFor reserved non-standard bitflag values, see @ref BitflagRegistry."
				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"MCommentBitflagValue": func(v string, indent int, b Bitflag, entryIndex int) string {
				value, _ := g.BitflagValue(b, entryIndex, true)
				s := value + "\n"
				v = strings.TrimSpace(v)
				if v != "" && v != "TODO" {
					s += v
				}
				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"MCommentFunction": func(fn *Function, indent int) string {
				var s string
				{
					var funcDoc = strings.TrimSpace(fn.Doc)
					if funcDoc != "" && funcDoc != "TODO" {
						s += funcDoc
					}
				}
				for _, arg := range fn.Args {
					argDoc := strings.TrimSpace(arg.Doc)
					var sArg string
					if argDoc != "" && argDoc != "TODO" {
						sArg = argDoc
					}

					if arg.PassedWithOwnership != nil {
						if *arg.PassedWithOwnership {
							sArg += "\nThis parameter is @ref ReturnedWithOwnership."
						} else {
							panic("invalid")
						}
					}

					sArg = strings.TrimSpace(sArg)
					if sArg != "" {
						s += "\n\n@param " + CamelCase(arg.Name) + "\n" + sArg
					}
				}
				if fn.Returns != nil {
					returnsDoc := strings.TrimSpace(fn.Returns.Doc)
					var sRet string
					if returnsDoc != "" && returnsDoc != "TODO" {
						sRet = returnsDoc
					}

					if fn.Returns.PassedWithOwnership != nil {
						if *fn.Returns.PassedWithOwnership {
							sRet += "\nThis value is @ref ReturnedWithOwnership."
						} else {
							panic("invalid")
						}
					}

					sRet = strings.TrimSpace(sRet)
					if sRet != "" {
						s += "\n\n@returns\n" + sRet
					}
				}
				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"MCommentCallback": func(cb *Callback, indent int) string {
				var s string
				{
					var funcDoc = strings.TrimSpace(cb.Doc)
					if funcDoc != "" && funcDoc != "TODO" {
						s += funcDoc
					}
					s += "\n\nSee also @ref CallbackError."
				}
				for _, arg := range cb.Args {
					var argDoc = strings.TrimSpace(arg.Doc)
					var sArg string
					if argDoc != "" && argDoc != "TODO" {
						sArg += argDoc
					}

					if arg.PassedWithOwnership != nil {
						if *arg.PassedWithOwnership {
							sArg += "\nThis parameter is @ref PassedWithOwnership."
						} else {
							sArg += "\nThis parameter is @ref PassedWithoutOwnership."
						}
					}

					sArg = strings.TrimSpace(sArg)
					if sArg != "" {
						s += "\n\n@param " + CamelCase(arg.Name) + "\n" + sArg
					}
				}
				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"MCommentMember": func(member *ParameterType, indent int) string {
				var s string

				var srcDoc = strings.TrimSpace(member.Doc)
				if srcDoc != "" && srcDoc != "TODO" {
					s += srcDoc
				}

				switch member.Type {
				case "nullable_string":
					s += "\n\nThis is a \\ref NullableInputString."
				case "string_with_default_empty":
					s += "\n\nThis is a \\ref NonNullInputString."
				case "out_string":
					s += "\n\nThis is an \\ref OutputString."
				}

				s += "\n\nThe `INIT` macro sets this to " + g.DefaultValue(*member, true /* isDocString */) + "."

				if member.PassedWithOwnership != nil {
					panic("invalid")
				}

				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"MCommentStruct": func(st *Struct, indent int) string {
				var s string

				var srcDoc = strings.TrimSpace(st.Doc)
				if srcDoc != "" && srcDoc != "TODO" {
					s += srcDoc
				}

				if st.Type == "extensible_callback_arg" {
					s += "\n\nThis is an @ref ImplementationAllocatedStructChain root.\nArbitrary chains must be handled gracefully by the application!"
				}

				s += "\n\nDefault values can be set using @ref WGPU_" + g.ConstantCaseName(st.Base) + "_INIT as initializer."

				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"MCommentProcPointer": func(name string, indent int) string {
				var s string
				s += "Proc pointer type for @ref wgpu" + name + ":\n"
				s += "> @copydoc wgpu" + name
				return Comment(strings.TrimSpace(s), CommentTypeMultiLine, indent, true)
			},
			"ConstantCase":     ConstantCase,
			"PascalCase":       PascalCase,
			"CamelCase":        CamelCase,
			"ConstantCaseName": g.ConstantCaseName,
			"PascalCaseName":   g.PascalCaseName,
			"CEnumName":        g.CEnumName,
			"CMethodName":      g.CMethodName,
			"CType":            g.CType,
			"CValue":           g.CValue,
			"EnumValue32":      g.EnumValue32,
			"BitflagValue": func(b Bitflag, entryIndex int) (string, error) {
				return g.BitflagValue(b, entryIndex, false)
			},
			"IsArray": func(typ string) bool {
				return arrayTypeRegexp.Match([]byte(typ))
			},
			"ArrayType": func(typ string, pointer PointerType) string {
				matches := arrayTypeRegexp.FindStringSubmatch(typ)
				if len(matches) == 2 {
					return g.CType(matches[1], pointer)
				}
				return ""
			},
			"Singularize":             Singularize,
			"IsLast":                  func(i int, s any) bool { return i == reflect.ValueOf(s).Len()-1 },
			"FunctionReturns":         g.FunctionReturns,
			"FunctionArgs":            g.FunctionArgs,
			"CallbackArgs":            g.CallbackArgs,
			"StructMember":            g.StructMember,
			"StructMemberArrayCount":  g.StructMemberArrayCount,
			"StructMemberArrayData":   g.StructMemberArrayData,
			"StructMemberInitializer": g.StructMemberInitializer,
		})
	t, err := t.Parse(tmpl)
	if err != nil {
		return fmt.Errorf("GenCHeader: failed to parse template: %w", err)
	}
	if err := t.Execute(dst, g); err != nil {
		return fmt.Errorf("GenCHeader: failed to execute template: %w", err)
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

func (g *Generator) ConstantCaseName(b Base) string {
	if b.Extended {
		return ConstantCase(b.Name)
	}

	prefix := g.PrefixForNamespace(b.Namespace)
	switch prefix {
	case "":
		return ConstantCase(b.Name)
	default:
		return ConstantCase(prefix + "_" + b.Name)
	}
}

func (g *Generator) PascalCaseName(b Base) string {
	if b.Extended {
		return PascalCase(b.Name)
	}

	prefix := g.PrefixForNamespace(b.Namespace)
	switch prefix {
	case "":
		return PascalCase(b.Name)
	default:
		return PascalCase(prefix + "_" + b.Name)
	}
}

func (g *Generator) CEnumName(typ Base, entry Base) string {
	if !typ.Extended {
		return g.CType(typ, "") + "_" + PascalCase(entry.Name)
	} else {
		return g.CType(typ, "") + "_" + g.PascalCaseName(entry)
	}
}

func (g *Generator) CMethodName(o Object, m Function) string {
	if !o.Extended {
		return g.PascalCaseName(o.Base) + PascalCase(m.Name)
	} else {
		return PascalCase(o.Name) + g.PascalCaseName(m.Base)
	}
}

func (g *Generator) CValue(s string) (string, error) {
	switch s {
	case "usize_max":
		return "SIZE_MAX", nil
	case "uint32_max":
		return "UINT32_MAX", nil
	case "uint64_max":
		return "UINT64_MAX", nil
	case "nan":
		return "NAN", nil
	default:
		var num string
		var base int
		if strings.HasPrefix(s, "0x") {
			base = 16
			num = strings.TrimPrefix(s, "0x")
		} else {
			base = 10
			num = s
		}
		v, err := strconv.ParseUint(num, base, 64)
		if err != nil {
			return "", fmt.Errorf("CValue: failed to parse \"%s\": %w", s, err)
		}
		var suffix string
		if v <= math.MaxUint32 {
			suffix = "UL"
		} else {
			suffix = "ULL"
		}
		return "0x" + strconv.FormatUint(v, 16) + suffix, nil
	}
}

func (g *Generator) CType(typ any, pointerType PointerType) string {
	appendModifiers := func(s string, pointerType PointerType) string {
		var sb strings.Builder
		sb.WriteString(s)
		switch pointerType {
		case PointerTypeImmutable:
			sb.WriteString(" const *")
		case PointerTypeMutable:
			sb.WriteString(" *")
		}
		return sb.String()
	}

	var ctype string
	switch t := typ.(type) {
	case string:
		{
			switch t {
			case "bool":
				ctype = "WGPUBool"
			case "nullable_string", "string_with_default_empty", "out_string":
				ctype = "WGPUStringView"
			case "uint16":
				ctype = "uint16_t"
			case "uint32":
				ctype = "uint32_t"
			case "uint64":
				ctype = "uint64_t"
			case "usize":
				ctype = "size_t"
			case "int16":
				ctype = "int16_t"
			case "int32":
				ctype = "int32_t"
			case "float32", "nullable_float32":
				ctype = "float"
			case "float64", "float64_supertype":
				ctype = "double"
			case "c_void":
				ctype = "void"
			default:
				// Handle type names prefixed with the type category.
				return g.CType(g.FindBaseType(t), pointerType)
			}
		}
	case Base:
		{
			ctype = "WGPU" + g.PascalCaseName(t)
		}
	default:
		panic("Unknown input for type")
	}

	return appendModifiers(ctype, pointerType)
}

func (g *Generator) FunctionReturns(f Function) string {
	if f.Callback != nil {
		return "WGPUFuture"
	}
	if f.Returns != nil {
		sb := &strings.Builder{}
		if f.Returns.Optional {
			sb.WriteString("WGPU_NULLABLE ")
		}
		sb.WriteString(g.CType(f.Returns.Type, f.Returns.Pointer))
		return sb.String()
	}
	return "void"
}

func (g *Generator) FunctionArgs(f Function, o *Object) string {
	sb := &strings.Builder{}
	if o != nil {
		if len(f.Args) > 0 {
			fmt.Fprintf(sb, "%s %s, ", g.CType(o.Base, ""), CamelCase(o.Name))
		} else {
			fmt.Fprintf(sb, "%s %s", g.CType(o.Base, ""), CamelCase(o.Name))
		}
	}
	for i, arg := range f.Args {
		if arg.Optional {
			sb.WriteString("WGPU_NULLABLE ")
		}
		matches := arrayTypeRegexp.FindStringSubmatch(arg.Type)
		if len(matches) == 2 {
			fmt.Fprintf(sb, "size_t %sCount, ", CamelCase(Singularize(arg.Name)))
			fmt.Fprintf(sb, "%s %s", g.CType(matches[1], arg.Pointer), CamelCase(arg.Name))
		} else {
			fmt.Fprintf(sb, "%s %s", g.CType(arg.Type, arg.Pointer), CamelCase(arg.Name))
		}
		if i != len(f.Args)-1 {
			sb.WriteString(", ")
		}
	}
	if f.Callback != nil {
		fmt.Fprintf(sb, ", %sCallbackInfo callbackInfo", g.CType(*f.Callback, ""))
	}
	return sb.String()
}

func (g *Generator) CallbackArgs(f Callback) string {
	sb := &strings.Builder{}
	for _, arg := range f.Args {
		if arg.Optional {
			sb.WriteString("WGPU_NULLABLE ")
		}
		var structPrefix string
		if strings.HasPrefix(arg.Type, "struct.") {
			structPrefix = "struct "
		}
		matches := arrayTypeRegexp.FindStringSubmatch(arg.Type)
		if len(matches) == 2 {
			fmt.Fprintf(sb, "size_t %sCount, ", CamelCase(Singularize(arg.Name)))
			fmt.Fprintf(sb, "%s%s %s, ", structPrefix, g.CType(matches[1], arg.Pointer), CamelCase(arg.Name))
		} else {
			fmt.Fprintf(sb, "%s%s %s, ", structPrefix, g.CType(arg.Type, arg.Pointer), CamelCase(arg.Name))
		}
	}
	sb.WriteString("WGPU_NULLABLE void* userdata1, WGPU_NULLABLE void* userdata2")
	return sb.String()
}

func (g *Generator) EnumValue16(e Enum, entryIndex int) (uint16, error) {
	entry := e.Entries[entryIndex]
	if entry.Value == "" {
		return uint16(entryIndex), nil
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
		value, err := strconv.ParseUint(num, base, 16)
		if err != nil {
			return 0, err
		}
		return uint16(value), nil
	}
}

func (g *Generator) EnumValue32(e Enum, entryIndex int) (uint32, error) {
	value16, err := g.EnumValue16(e, entryIndex)
	if err != nil {
		return 0, err
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
	var entryComment string
	if len(entry.ValueCombination) > 0 {
		if entry.Value != "" {
			return "", fmt.Errorf("BitflagValue: found conflicting 'value' and 'value_combination' in '%s'", b.Name)
		}
		entryComment += "`"
		for valueIndex, v := range entry.ValueCombination {
			// find the value by searching in b, bitwise-OR it into the result
			for searchIndex, search := range b.Entries {
				if search.Name == v {
					searchValue, err := bitflagEntryValue(search, searchIndex)
					if err != nil {
						return "", nil
					}
					value |= searchValue
					break
				}
			}
			// construct comment
			idx := slices.IndexFunc(b.Entries, func(e BitflagEntry) bool { return e.Name == v })
			if idx != -1 {
				entryComment += g.PascalCaseName(b.Entries[idx].Base)
			} else {
				entryComment += PascalCase(v)
			}
			if valueIndex != len(entry.ValueCombination)-1 {
				entryComment += " | "
			}
		}
		entryComment += "`."
	} else {
		var err error
		value, err = bitflagEntryValue(entry, entryIndex)
		if err != nil {
			return "", nil
		}
		if value == 0 {
			entryComment = "`0`."
		}
	}
	if isDocString {
		return entryComment, nil
	} else {
		return fmt.Sprintf("0x%.16X", value), nil
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

func (g *Generator) StructMember(s Struct, memberIndex int) (string, error) {
	member := s.Members[memberIndex]

	matches := arrayTypeRegexp.FindStringSubmatch(member.Type)
	if len(matches) == 2 {
		panic("StructMember used on array type")
	}

	sb := &strings.Builder{}
	if member.Optional {
		sb.WriteString("WGPU_NULLABLE ")
	}
	if strings.HasPrefix(member.Type, "callback.") {
		fmt.Fprintf(sb, "%sCallbackInfo %s;", g.CType(member.Type, ""), CamelCase(member.Name))
	} else {
		fmt.Fprintf(sb, "%s %s;", g.CType(member.Type, member.Pointer), CamelCase(member.Name))
	}
	return sb.String(), nil
}

func (g *Generator) StructMemberArrayCount(s Struct, memberIndex int) (string, error) {
	member := s.Members[memberIndex]

	matches := arrayTypeRegexp.FindStringSubmatch(member.Type)
	if len(matches) != 2 {
		panic("StructMemberArrayCount used on non-array")
	}

	return fmt.Sprintf("size_t %sCount;", CamelCase(Singularize(member.Name))), nil
}

func (g *Generator) StructMemberArrayData(s Struct, memberIndex int) (string, error) {
	member := s.Members[memberIndex]

	matches := arrayTypeRegexp.FindStringSubmatch(member.Type)
	if len(matches) != 2 {
		panic("StructMemberArrayCount used on non-array")
	}

	sb := &strings.Builder{}
	if member.Optional {
		sb.WriteString("WGPU_NULLABLE ")
	}
	fmt.Fprintf(sb, "%s %s;", g.CType(matches[1], member.Pointer), CamelCase(member.Name))
	return sb.String(), nil
}

func (g *Generator) StructMemberInitializer(s Struct, memberIndex int) (string, error) {
	member := s.Members[memberIndex]
	sb := &strings.Builder{}
	matches := arrayTypeRegexp.FindStringSubmatch(member.Type)
	if len(matches) == 2 {
		fmt.Fprintf(sb, "/*.%sCount=*/0 _wgpu_COMMA \\\n", CamelCase(Singularize(member.Name)))
		fmt.Fprintf(sb, "    /*.%s=*/NULL _wgpu_COMMA \\", CamelCase(member.Name))
	} else {
		fmt.Fprintf(sb, "/*.%s=*/%s _wgpu_COMMA \\", CamelCase(member.Name), g.DefaultValue(member, false /* isDocString */))
	}
	return sb.String(), nil
}

func (g *Generator) DefaultValue(member ParameterType, isDocString bool) string {
	ref := func(s string) string {
		if isDocString {
			return "@ref " + s
		} else {
			return s
		}
	}
	literal := func(s string) string {
		if isDocString {
			return "`" + s + "`"
		} else {
			return s
		}
	}

	switch {
	case member.Pointer != "":
		if member.Default != nil {
			panic("pointer type should not have a default")
		}
		return literal("NULL")

	// Cases that may have member.Default
	case strings.HasPrefix(member.Type, "enum."):
		if member.Default == nil {
			if member.Type == "enum.optional_bool" {
				// This Undefined is a special one that is not the zero-value, so that
				// a stdbool.h bool cast correctly to WGPUOptionalBool; this means we
				// must explicitly initialize it
				return ref("WGPUOptionalBool_Undefined")
			} else if isDocString {
				return "(@ref " + g.CType(member.Type, "") + ")0"
			} else {
				return "_wgpu_ENUM_ZERO_INIT(" + g.CType(member.Type, "") + ")"
			}
		} else {
			return ref(g.CType(member.Type, "") + "_" + PascalCase(*member.Default))
		}
	case strings.HasPrefix(member.Type, "bitflag."):
		if member.Default == nil {
			return ref(g.CType(member.Type, "") + "_None")
		} else {
			return ref(g.CType(member.Type, "") + "_" + PascalCase(*member.Default))
		}
	case member.Type == "uint16", member.Type == "uint32", member.Type == "uint64", member.Type == "usize", member.Type == "int32":
		if member.Default == nil {
			return literal("0")
		} else if strings.HasPrefix(*member.Default, "constant.") {
			return ref("WGPU_" + g.ConstantCaseName(g.FindBaseType(*member.Default)))
		} else {
			return literal(*member.Default)
		}
	case member.Type == "float32" || member.Type == "nullable_float32":
		if member.Default == nil {
			return literal("0.f")
		} else if strings.HasPrefix(*member.Default, "constant.") {
			return ref("WGPU_" + g.ConstantCaseName(g.FindBaseType(*member.Default)))
		} else if strings.Contains(*member.Default, ".") {
			return literal(*member.Default + "f")
		} else {
			return literal(*member.Default + ".f")
		}
	case member.Type == "float64" || member.Type == "float64_supertype":
		if member.Default == nil {
			return literal("0.")
		} else if strings.HasPrefix(*member.Default, "constant.") {
			return ref("WGPU_" + g.ConstantCaseName(g.FindBaseType(*member.Default)))
		} else {
			return literal(*member.Default)
		}
	case member.Type == "bool":
		if member.Default == nil {
			return literal("WGPU_FALSE")
		} else if strings.HasPrefix(*member.Default, "constant.") {
			return ref("WGPU_" + g.ConstantCaseName(g.FindBaseType(*member.Default)))
		} else if *member.Default == "true" {
			return literal("WGPU_TRUE")
		} else if *member.Default == "false" {
			return literal("WGPU_FALSE")
		} else {
			return *member.Default
		}
	case strings.HasPrefix(member.Type, "struct."):
		if member.Optional {
			return literal("NULL")
		} else if member.Default == nil {
			return ref("WGPU_" + g.ConstantCaseName(g.FindBaseType(member.Type)) + "_INIT")
		} else if *member.Default == "zero" {
			if isDocString {
				return "zero (which sets the entry to `BindingNotUsed`)"
			} else {
				return literal("_wgpu_STRUCT_ZERO_INIT")
			}
		} else {
			panic("unknown default for struct type")
		}
	case member.Default != nil:
		panic(fmt.Errorf("type %s should not have a default", member.Type))

	// Cases that should not have member.Default
	case strings.HasPrefix(member.Type, "callback."):
		return ref("WGPU_" + g.ConstantCaseName(g.FindBaseType(member.Type)) + "_CALLBACK_INFO_INIT")
	case strings.HasPrefix(member.Type, "object."):
		return literal("NULL")
	case strings.HasPrefix(member.Type, "array<"):
		return literal("NULL")
	case member.Type == "out_string", member.Type == "string_with_default_empty", member.Type == "nullable_string":
		return ref("WGPU_STRING_VIEW_INIT")
	case member.Type == "c_void":
		return literal("NULL")
	default:
		panic("invalid prefix: " + member.Type + " in member " + member.Name)
	}
}

// Main function
func main() {
	flag.StringVar(&schemaPath, "schema", "", "path of the json schema")
	flag.Var(&yamlPaths, "yaml", "path of the yaml spec")
	flag.Var(&headerPaths, "header", "output path of the header")
	flag.BoolVar(&extPrefix, "extprefix", true, "append prefix to extension identifiers")
	flag.Parse()
	if schemaPath == "" || len(headerPaths) == 0 || len(yamlPaths) == 0 || len(headerPaths) != len(yamlPaths) {
		flag.Usage()
		os.Exit(1)
	}

	// Order matters for validation steps, so enforce it.
	if len(yamlPaths) > 1 && filepath.Base(yamlPaths[0]) != "webgpu.yml" {
		panic(`"webgpu.yml" must be the first sequence in the order`)
	}

	// Validate the yaml files (jsonschema, duplications)
	if err := ValidateYamls(schemaPath, yamlPaths); err != nil {
		panic(err)
	}

	// Generate the header files
	for i, yamlPath := range yamlPaths {
		headerPath := headerPaths[i]
		headerFileName := filepath.Base(headerPath)
		headerFileNameSplit := strings.Split(headerFileName, ".")
		if len(headerFileNameSplit) != 2 {
			panic("got invalid header file name: " + headerFileName)
		}

		src, err := os.ReadFile(yamlPath)
		if err != nil {
			panic(err)
		}

		dst, err := os.Create(headerPath)
		if err != nil {
			panic(err)
		}

		var yml Yml
		if err := yaml.Unmarshal(src, &yml); err != nil {
			panic(err)
		}

		SortAndTransform(&yml)

		prefix := ""
		if yml.Name != "webgpu" && extPrefix {
			prefix = yml.Name
		}
		g := &Generator{
			Yml:        &yml,
			HeaderName: headerFileNameSplit[0],
			ExtPrefix:  prefix,
		}
		if err := g.Gen(dst); err != nil {
			panic(err)
		}
	}
}
