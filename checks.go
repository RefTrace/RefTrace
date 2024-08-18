package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"reft-go/nf"
	"reft-go/parser"
	"strconv"

	"github.com/antlr4-go/antlr/v4"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var checkCmd = &cobra.Command{
	Use:   "check [nf_file] [checks_file]",
	Short: "Check a .nf file using rules from a checks.nf file",
	Args:  cobra.ExactArgs(2),
	Run:   runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
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

func runCheck(cmd *cobra.Command, args []string) {
	nfFile := args[0]
	checksFile := args[1]

	// Read the checks.nf file
	checksContent, err := os.ReadFile(checksFile)
	if err != nil {
		log.Fatalf("Error reading checks.nf file: %v", err)
	}

	// Create a new Starlark thread
	thread := &starlark.Thread{Name: "check_thread"}

	fo := &syntax.FileOptions{}

	// Parse the Starlark code without executing it
	f, err := fo.Parse(checksFile, checksContent, 0)
	if err != nil {
		log.Fatalf("Error parsing checks program: %v", err)
	}

	// Compile the parsed code
	prog, err := starlark.FileProgram(f, func(name string) bool {
		return false
	})
	if err != nil {
		log.Fatalf("Error compiling checks program: %v", err)
	}

	// Create predefined variables for the Starlark environment
	predefined := starlark.StringDict{
		// Add any other predefined variables or functions here
	}

	// Execute the compiled program
	globals, err := prog.Init(thread, predefined)
	if err != nil {
		log.Fatalf("Error initializing checks program: %v", err)
	}

	// Check if main function is defined
	mainFunc, ok := globals["main"]
	if !ok {
		log.Fatal("main function not defined in checks program")
	}

	// Create the params list
	params, includes := parse(nfFile)
	starlarkParams := make([]starlark.Value, len(params))
	for i, param := range params {
		starlarkParams[i] = NewStarlarkParamInfo(param)
	}
	starlarkParamsList := starlark.NewList(starlarkParams)

	starlarkIncludes := make([]starlark.Value, len(includes))
	for i, include := range includes {
		starlarkIncludes[i] = NewStarlarkIncludeInfo(include)
	}
	starlarkIncludesList := starlark.NewList(starlarkIncludes)

	// Call the main function
	_, err = starlark.Call(thread, mainFunc, starlark.Tuple{starlarkParamsList, starlarkIncludesList}, nil)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			fmt.Printf("Execution failed: %s\n", evalErr.Msg)
		} else {
			log.Fatalf("Error calling main function: %v\n", err)
		}
		return
	}
}
