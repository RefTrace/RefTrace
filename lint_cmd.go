package main

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"reft-go/nf"
	"reft-go/parser"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var lintCmd = &cobra.Command{
	Use:   "lint [rules file] [directory]",
	Short: "Lint a directory using rules from a rules.py file",
	Args:  cobra.ExactArgs(2),
	Run:   runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)
}

type StarlarkParamInfo struct {
	paramInfo nf.ParamInfo
}

func (p *StarlarkParamInfo) String() string {
	return fmt.Sprintf("Param(name=%s, line=%d)", p.paramInfo.Name, p.paramInfo.LineNumber)
}

func (p *StarlarkParamInfo) Type() string {
	return "Param"
}

func (p *StarlarkParamInfo) Freeze() {}

func (p *StarlarkParamInfo) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (p *StarlarkParamInfo) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(p.paramInfo.Name))
	h.Write([]byte(strconv.Itoa(p.paramInfo.LineNumber)))
	return h.Sum32(), nil
}

func (p *StarlarkParamInfo) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(p.paramInfo.Name), nil
	case "line":
		return starlark.MakeInt(p.paramInfo.LineNumber), nil
	default:
		return nil, nil
	}
}

func (p *StarlarkParamInfo) AttrNames() []string {
	return []string{"name", "line"}
}

func NewStarlarkParamInfo(p nf.ParamInfo) *StarlarkParamInfo {
	return &StarlarkParamInfo{paramInfo: p}
}

type StarlarkIncludeInfo struct {
	includeInfo nf.IncludeInfo
}

func (i *StarlarkIncludeInfo) String() string {
	return fmt.Sprintf("Include(name=%s, from_=%s, line=%d)", i.includeInfo.Name, i.includeInfo.From, i.includeInfo.LineNumber)
}

func (i *StarlarkIncludeInfo) Type() string {
	return "Include"
}

func (i *StarlarkIncludeInfo) Freeze() {}

func (i *StarlarkIncludeInfo) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (i *StarlarkIncludeInfo) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(i.includeInfo.Name))
	h.Write([]byte(i.includeInfo.From))
	h.Write([]byte(strconv.Itoa(i.includeInfo.LineNumber)))
	return h.Sum32(), nil
}

func (i *StarlarkIncludeInfo) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(i.includeInfo.Name), nil
	case "from_":
		return starlark.String(i.includeInfo.From), nil
	case "line":
		return starlark.MakeInt(i.includeInfo.LineNumber), nil
	default:
		return nil, nil
	}
}

func (i *StarlarkIncludeInfo) AttrNames() []string {
	return []string{"name", "from_", "line"}
}

func NewStarlarkIncludeInfo(i nf.IncludeInfo) *StarlarkIncludeInfo {
	return &StarlarkIncludeInfo{includeInfo: i}
}

func parse(filePath string) ([]nf.ParamInfo, []nf.IncludeInfo) {
	//debug.SetGCPercent(-1)
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		log.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := parser.NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	stream.Fill()
	groovyParser := parser.NewGroovyParser(stream)
	//parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	// Parse the file
	tree := groovyParser.CompilationUnit()
	builder := parser.NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*parser.ModuleNode)
	paramVisitor := nf.NewParamVisitor()
	paramVisitor.VisitBlockStatement(ast.StatementBlock)
	params := paramVisitor.GetSortedParams()

	includeVisitor := nf.NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.GetSortedIncludes()

	return params, includes
}

// failnowFunc is the implementation of the failnow function for Starlark
func failnowFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sep := " "
	if err := starlark.UnpackArgs("failnow", nil, kwargs, "sep?", &sep); err != nil {
		return nil, err
	}
	buf := new(strings.Builder)
	buf.WriteString("failnow: ")
	for i, v := range args {
		if i > 0 {
			buf.WriteString(sep)
		}
		if s, ok := starlark.AsString(v); ok {
			buf.WriteString(s)
		} else {
			buf.WriteString(v.String())
		}
	}

	return nil, errors.New(buf.String())
}

func runLint(cmd *cobra.Command, args []string) {
	rulesFile := args[0]
	dir := args[1]

	// Read the rules.py file
	rulesContent, err := os.ReadFile(rulesFile)
	if err != nil {
		log.Fatalf("Error reading rules.py file: %v", err)
	}

	// Remove the 'fail' function from the Universe
	delete(starlark.Universe, "fail")

	// Create a new Starlark thread
	thread := &starlark.Thread{Name: "lint_thread"}

	fo := &syntax.FileOptions{}

	// Parse the Starlark code without executing it
	f, err := fo.Parse(rulesFile, rulesContent, 0)
	if err != nil {
		log.Fatalf("Error parsing rules program: %v", err)
	}

	// Compile the parsed code
	prog, err := starlark.FileProgram(f, func(name string) bool {
		return false
	})
	if err != nil {
		log.Fatalf("Error compiling rules program: %v", err)
	}

	// Create predefined variables for the Starlark environment
	predefined := starlark.StringDict{
		"failnow": starlark.NewBuiltin("failnow", failnowFunc),
		// Add any other predefined variables or functions here
	}

	// Execute the compiled program
	globals, err := prog.Init(thread, predefined)
	if err != nil {
		log.Fatalf("Error initializing rules program: %v", err)
	}

	// Collect rules (functions starting with "rule_")
	rules := make(map[string]starlark.Callable)
	for name, value := range globals {
		if strings.HasPrefix(name, "rule_") {
			if callable, ok := value.(starlark.Callable); ok {
				rules[name] = callable
			}
		}
	}

	// Create the params list
	/*
		params, includes := parse(rulesFile)
		starlarkParams := make([]starlark.Value, len(params))
		for i, param := range params {
			starlarkParams[i] = NewStarlarkParamInfo(param)
		}
		_ = starlark.NewList(starlarkParams)
	*/

	/*
		starlarkIncludes := make([]starlark.Value, len(includes))
		for i, include := range includes {
			starlarkIncludes[i] = NewStarlarkIncludeInfo(include)
		}
		_ = starlark.NewList(starlarkIncludes)
	*/

	// Parse the directory and get the modules
	modules, err := nf.ProcessDirectory(dir)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Execute each rule
	for ruleName, ruleFunc := range rules {
		for _, module := range modules {
			starlarkModule := nf.ConvertToStarlarkModule(module)
			_, err := starlark.Call(thread, ruleFunc, starlark.Tuple{starlarkModule}, nil)
			if err != nil {
				if evalErr, ok := err.(*starlark.EvalError); ok {
					fmt.Printf("Rule %s execution failed: %s\n", ruleName, evalErr.Msg)
				} else {
					log.Fatalf("Error calling rule %s: %v\n", ruleName, err)
				}
			}
		}
	}
}
