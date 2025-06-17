package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	_ "embed"

	"github.com/goccy/go-yaml"
	"golang.org/x/tools/imports"
)

//go:embed wgpu.go.tpl
var tmpl string

// Constants
const (
	defaultURL    = "https://raw.githubusercontent.com/webgpu-native/webgpu-headers/refs/heads/main/webgpu.yml"
	enumShift     = 16
	bitflagBase   = 2
	expectedParts = 2
)

// Regular expressions
var arrayTypeRegexp = regexp.MustCompile(`array<([a-zA-Z0-9._]+)>`)

// Command line flags
var (
	goPaths   StringListFlag
	yamlPaths StringListFlag
	extPrefix bool
)

// StringListFlag implements flag.Value interface for string slices
type StringListFlag []string

func (f *StringListFlag) String() string     { return fmt.Sprintf("%#v", f) }
func (f *StringListFlag) Set(v string) error { *f = append(*f, v); return nil }

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	parseFlags()

	if err := validateFlags(); err != nil {
		return err
	}

	yamlDatas, err := loadYAMLFiles()
	if err != nil {
		return fmt.Errorf("failed to load YAML files: %w", err)
	}

	return generateGoFiles(yamlDatas)
}

func parseFlags() {
	flag.Var(&yamlPaths, "yaml", "path or URL of the yaml spec")
	flag.Var(&goPaths, "go", "output path of the go")
	flag.BoolVar(&extPrefix, "extprefix", true, "append prefix to extension identifiers")
	flag.Parse()

	// Use default URL if not provided
	if len(yamlPaths) == 0 {
		yamlPaths = append(yamlPaths, defaultURL)
	}
}

func validateFlags() error {
	if len(goPaths) == 0 || len(goPaths) != len(yamlPaths) {
		flag.Usage()
		return fmt.Errorf("number of go paths must match yaml paths")
	}
	return nil
}

func loadYAMLFiles() ([][]byte, error) {
	yamlDatas := make([][]byte, len(yamlPaths))

	for i, yamlPath := range yamlPaths {
		data, err := loadFile(yamlPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s: %w", yamlPath, err)
		}
		yamlDatas[i] = data
	}

	return yamlDatas, nil
}

func loadFile(path string) ([]byte, error) {
	if isURL(path) {
		return fetchFile(path)
	}
	return os.ReadFile(path)
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func generateGoFiles(yamlDatas [][]byte) error {
	for i, yamlData := range yamlDatas {
		if err := generateSingleGoFile(yamlData, goPaths[i]); err != nil {
			return fmt.Errorf("failed to generate %s: %w", goPaths[i], err)
		}
	}
	return nil
}

func generateSingleGoFile(yamlData []byte, goPath string) error {
	goName, err := extractGoName(goPath)
	if err != nil {
		return err
	}

	dst, err := os.Create(goPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	var yml Yml
	if err := yaml.Unmarshal(yamlData, &yml); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	SortAndTransform(&yml)

	generator := NewGenerator(&yml, goName, extPrefix)
	return generator.Generate(dst)
}

func extractGoName(goPath string) (string, error) {
	goFileName := filepath.Base(goPath)
	parts := strings.Split(goFileName, ".")
	if len(parts) != expectedParts {
		return "", fmt.Errorf("invalid go file name: %s", goFileName)
	}
	return parts[0], nil
}

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

// Comment types
type CommentType uint8

const (
	CommentTypeSingleLine CommentType = iota
	CommentTypeMultiLine
)

// Generator creates a new code generator
func NewGenerator(yml *Yml, goName string, useExtPrefix bool) *Generator {
	prefix := ""
	if yml.Name != "webgpu" && useExtPrefix {
		prefix = yml.Name
	}

	return &Generator{
		Yml:       yml,
		GoName:    goName,
		ExtPrefix: prefix,
	}
}

// Generator structure
type Generator struct {
	ExtPrefix string
	GoName    string
	*Yml
}

// Generate produces Go code from the YAML specification
func (g *Generator) Generate(dst io.Writer) error {
	tmpl, err := g.parseTemplate()
	if err != nil {
		return err
	}

	code, err := g.executeTemplate(tmpl)
	if err != nil {
		return err
	}

	formatted, err := formatCode(code)
	if err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}

	_, err = dst.Write(formatted)
	return err
}

func (g *Generator) parseTemplate() (*template.Template, error) {
	funcMap := g.createTemplateFuncMap()
	t := template.New("").Funcs(funcMap)

	parsed, err := t.Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return parsed, nil
}

func (g *Generator) createTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"SComment":            func(v string, indent int) string { return Comment(v, CommentTypeSingleLine, indent, true) },
		"SCommentN":           func(v string, indent int) string { return Comment(v, CommentTypeSingleLine, indent, false) },
		"IsArray":             func(typ string) bool { return arrayTypeRegexp.MatchString(typ) },
		"IsCallback":          func(typ string) bool { return strings.HasPrefix(typ, "callback.") },
		"HasSuffix":           strings.HasSuffix,
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
		"BitflagValue":        func(b Bitflag, entryIndex int) (string, error) { return g.BitflagValue(b, entryIndex, false) },
	}
}

func (g *Generator) executeTemplate(tmpl *template.Template) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}

// formatCode applies goimports and gofmt to the generated code
func formatCode(code []byte) ([]byte, error) {
	// Apply goimports first
	processed, err := imports.Process("", code, nil)
	if err != nil {
		processed = code // fallback to original code
	}

	// Apply gofmt
	formatted, err := format.Source(processed)
	if err != nil {
		return processed, nil // return processed version as fallback
	}

	return formatted, nil
}

// Type finder methods with improved error handling
func (g *Generator) FindBaseType(typ string) Base {
	category, name, found := strings.Cut(typ, ".")
	if !found {
		return Base{} // return zero value instead of panic
	}

	finder := map[string]func(string) (Base, bool){
		"constant": g.findConstant,
		"typedef":  g.findTypedef,
		"enum":     g.findEnum,
		"bitflag":  g.findBitflag,
		"struct":   g.findStruct,
		"callback": g.findCallback,
		"object":   g.findObject,
	}

	if findFunc, exists := finder[category]; exists {
		if base, found := findFunc(name); found {
			return base
		}
	}

	// Return default base with namespace prefix
	return Base{Name: name, Namespace: g.PrefixForNamespace("")}
}

func (g *Generator) findConstant(name string) (Base, bool) {
	idx := slices.IndexFunc(g.Constants, func(c Constant) bool { return c.Name == name })
	if idx == -1 {
		return Base{}, false
	}
	return g.Constants[idx].Base, true
}

func (g *Generator) findTypedef(name string) (Base, bool) {
	idx := slices.IndexFunc(g.Typedefs, func(t Typedef) bool { return t.Name == name })
	if idx == -1 {
		return Base{}, false
	}
	return g.Typedefs[idx].Base, true
}

func (g *Generator) findEnum(name string) (Base, bool) {
	idx := slices.IndexFunc(g.Enums, func(e Enum) bool { return e.Name == name })
	if idx == -1 {
		return Base{}, false
	}
	return g.Enums[idx].Base, true
}

func (g *Generator) findBitflag(name string) (Base, bool) {
	idx := slices.IndexFunc(g.Bitflags, func(b Bitflag) bool { return b.Name == name })
	if idx == -1 {
		return Base{}, false
	}
	return g.Bitflags[idx].Base, true
}

func (g *Generator) findStruct(name string) (Base, bool) {
	idx := slices.IndexFunc(g.Structs, func(s Struct) bool { return s.Name == name })
	if idx == -1 {
		return Base{}, false
	}
	return g.Structs[idx].Base, true
}

func (g *Generator) findCallback(name string) (Base, bool) {
	idx := slices.IndexFunc(g.Callbacks, func(c Callback) bool { return c.Name == name })
	if idx == -1 {
		return Base{}, false
	}
	return g.Callbacks[idx].Base, true
}

func (g *Generator) findObject(name string) (Base, bool) {
	idx := slices.IndexFunc(g.Objects, func(o Object) bool { return o.Name == name })
	if idx == -1 {
		return Base{}, false
	}
	return g.Objects[idx].Base, true
}

func (g *Generator) FindCallback(typ string) Callback {
	category, name, found := strings.Cut(typ, ".")
	if !found || category != "callback" {
		return Callback{} // return zero value instead of panic
	}

	idx := slices.IndexFunc(g.Callbacks, func(c Callback) bool { return c.Name == name })
	if idx == -1 {
		return Callback{} // return zero value instead of panic
	}
	return g.Callbacks[idx]
}

func (g *Generator) PrefixForNamespace(namespace string) string {
	namespaceMap := map[string]string{
		"":       g.ExtPrefix,
		"webgpu": "",
	}

	if prefix, exists := namespaceMap[namespace]; exists {
		return prefix
	}
	return namespace
}

// Value calculation methods with improved error handling
func (g *Generator) EnumValue32(e Enum, entryIndex int) (uint32, error) {
	if entryIndex >= len(e.Entries) {
		return 0, fmt.Errorf("entry index %d out of range", entryIndex)
	}

	entry := e.Entries[entryIndex]
	value16, err := g.parseValue16(entry.Value, entryIndex)
	if err != nil {
		return 0, err
	}

	return uint32(g.EnumPrefix)<<enumShift | uint32(value16), nil
}

func (g *Generator) parseValue16(value string, defaultIndex int) (uint16, error) {
	if value == "" {
		return uint16(defaultIndex), nil
	}

	base, numStr := parseValueBase(value)
	v, err := strconv.ParseUint(numStr, base, 16)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value %s: %w", value, err)
	}

	return uint16(v), nil
}

func parseValueBase(value string) (base int, numStr string) {
	if strings.HasPrefix(value, "0x") {
		return 16, strings.TrimPrefix(value, "0x")
	}
	return 10, value
}

// Generator methods
func (g *Generator) GoConstantName(b Base) string {
	return PascalCase(b.Name)
}

func (g *Generator) GoTypeName(b Base) string {
	return PascalCase(b.Name)
}

func (g *Generator) GoEnumName(typ Base, entry Base) string {
	return PascalCase(typ.Name) + PascalCase(entry.Name)
}

var goKeywordMap = map[string]string{
	"type":      "errorType",
	"func":      "function",
	"var":       "variable",
	"const":     "constant",
	"struct":    "structure",
	"interface": "iface",
	"package":   "pkg",
	"import":    "imp",
	"return":    "ret",
	"if":        "condition",
	"else":      "alternative",
	"for":       "loop",
	"range":     "rng",
	"switch":    "switchValue",
	"case":      "caseValue",
	"default":   "defaultValue",
	"go":        "routine",
	"defer":     "deferred",
	"select":    "selector",
	"chan":      "channel",
	"map":       "mapping",
}

var goValueMap = map[string]string{
	"usize_max":  "^uintptr(0)",
	"uint32_max": "^uint32(0)",
	"uint64_max": "^uint64(0)",
	"nan":        "math.NaN()",
}

var goTypeMap = map[string]string{
	"bool":                      "bool",
	"uint16":                    "uint16",
	"uint32":                    "uint32",
	"uint64":                    "uint64",
	"usize":                     "uintptr",
	"int16":                     "int16",
	"int32":                     "int32",
	"float32":                   "float32",
	"nullable_float32":          "float32",
	"float64":                   "float64",
	"float64_supertype":         "float64",
	"nullable_string":           "string",
	"string_with_default_empty": "string",
	"out_string":                "string",
	"c_void":                    "unsafe.Pointer",
}

// Go-specific generator methods - only the ones used in templates
func (g *Generator) GoValue(s string) string {
	if replacement, exists := goValueMap[s]; exists {
		return replacement
	}
	return s
}

func (g *Generator) GoType(typ string) string {
	if goType, exists := goTypeMap[typ]; exists {
		return goType
	}
	// Handle type names prefixed with the type category
	return g.GoTypeName(g.FindBaseType(typ))
}

func (g *Generator) GoFunctionName(b Base) string {
	if strings.HasPrefix(b.Name, "get_") {
		// For getter functions, we use the name without the "get_" prefix
		return PascalCase(strings.TrimPrefix(b.Name, "get_"))
	}

	return PascalCase(b.Name)
}

func (g *Generator) GoParameterName(name string) string {
	camelName := CamelCase(name)
	if replacement, isKeyword := goKeywordMap[camelName]; isKeyword {
		return replacement
	}
	return camelName
}

func (g *Generator) GoFunctionArgs(f Function) string {
	sb := &strings.Builder{}
	for _, arg := range f.Args {
		if strings.HasPrefix(f.Name, "get_") && arg.Pointer == PointerTypeMutable {
			// If the function is a getter, we don't want to include mutable pointers in the arguments.
			// This is because getters are expected to return values, not modify them.
			continue
		}

		if sb.Len() > 0 {
			sb.WriteString(", ")
		}

		pointer := ""
		if arg.Pointer == PointerTypeMutable {
			pointer = "*"
		}

		matches := arrayTypeRegexp.FindStringSubmatch(arg.Type)
		if len(matches) == 2 {
			fmt.Fprintf(sb, "%s %s[]%s", g.GoParameterName(arg.Name), pointer, g.GoType(matches[1]))
		} else {
			fmt.Fprintf(sb, "%s %s%s", g.GoParameterName(arg.Name), pointer, g.GoType(arg.Type))
		}
	}

	return sb.String()
}

func (g *Generator) GoFunctionReturns(f Function) string {
	sb := &strings.Builder{}

	for _, arg := range f.Args {
		if strings.HasPrefix(f.Name, "get_") && arg.Pointer == PointerTypeMutable {
			if sb.Len() > 0 {
				sb.WriteString(", ")
			}
			pointer := ""
			if arg.Pointer == PointerTypeMutable {
				pointer = "*"
			}
			matches := arrayTypeRegexp.FindStringSubmatch(arg.Type)
			if len(matches) == 2 {
				fmt.Fprintf(sb, "%s[]%s", pointer, g.GoType(matches[1]))
			} else {
				fmt.Fprintf(sb, "%s%s", pointer, g.GoType(arg.Type))
			}
		}
	}

	if f.Callback != nil {
		cb := g.FindCallback(*f.Callback)
		for _, arg := range cb.Args {
			if !strings.HasPrefix(arg.Type, "object.") {
				continue
			}

			if sb.Len() > 0 {
				sb.WriteString(", ")
			}
			fmt.Fprintf(sb, "%s, error", g.GoTypeName(g.FindBaseType(arg.Type)))
		}
	}

	// In go, we don't use "enum.status" as a return type, we use "error" instead.
	if f.Returns != nil {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}

		if f.Returns.Type == "enum.status" {
			sb.WriteString("error")
		} else {
			sb.WriteString(g.GoType(f.Returns.Type))
		}
	}

	if sb.Len() == 0 {
		return ""
	}

	return "(" + sb.String() + ")"
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

func PascalCase(s string) string {
	return camelCase(s, true)
}

func CamelCase(s string) string {
	return camelCase(s, false)
}

func camelCase(s string, initCase bool) string {
	var out strings.Builder
	out.Grow(len(s))
	nextUpper := initCase
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

func BitflagEntryValue(entry BitflagEntry, entryIndex int) (uint64, error) {
	if entry.Value != "" {
		base, numStr := parseValueBase(entry.Value)
		return strconv.ParseUint(numStr, base, 64)
	}

	if entryIndex <= 0 {
		return 0, nil
	}

	return uint64(math.Pow(bitflagBase, float64(entryIndex-1))), nil
}

func (g *Generator) BitflagValue(b Bitflag, entryIndex int, isDocString bool) (string, error) {
	if entryIndex >= len(b.Entries) {
		return "", fmt.Errorf("entry index %d out of range for bitflag %s", entryIndex, b.Name)
	}

	entry := b.Entries[entryIndex]

	var value uint64
	if len(entry.ValueCombination) > 0 {
		if entry.Value != "" {
			return "", fmt.Errorf("BitflagValue: found conflicting 'value' and 'value_combination' in '%s'", b.Name)
		}
		for _, v := range entry.ValueCombination {
			// Find the value by searching in b, bitwise-OR it into the result
			for searchIndex, search := range b.Entries {
				if search.Name == v {
					searchValue, err := BitflagEntryValue(search, searchIndex)
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
		value, err = BitflagEntryValue(entry, entryIndex)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("0x%.16X", value), nil
}
