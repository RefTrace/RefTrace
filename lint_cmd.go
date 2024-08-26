package main

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"reft-go/nf"
	"reft-go/parser"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/antlr4-go/antlr/v4"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var (
	rulesFile string
	dir       string
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a directory using rules from a rules.py file",
	Run:   runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)

	lintCmd.Flags().StringVarP(&rulesFile, "rules", "r", "rules.py", "Path to the rules file")
	lintCmd.Flags().StringVarP(&dir, "directory", "d", ".", "Directory to lint")
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
func fatalFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sep := " "
	if err := starlark.UnpackArgs("failnow", nil, kwargs, "sep?", &sep); err != nil {
		return nil, err
	}
	buf := new(strings.Builder)
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

type RuleModuleOutput struct {
	Errors  []string
	Outputs []string
}

// rule -> module -> output
type GroupedOutput map[string]map[string]RuleModuleOutput

func runLint(cmd *cobra.Command, args []string) {
	// Check if rules file exists
	if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
		log.Fatalf("Rules file not found: %s", rulesFile)
	}

	// Convert relative paths to absolute paths
	rulesFile, _ = filepath.Abs(rulesFile)
	dir, _ = filepath.Abs(dir)

	// Read the rules.py file
	rulesContent, err := os.ReadFile(rulesFile)
	if err != nil {
		log.Fatalf("Error reading rules.py file: %v", err)
	}

	// Remove the 'fail' function from the Universe
	delete(starlark.Universe, "fail")

	var outputMutex sync.Mutex
	groupedOutput := make(GroupedOutput)

	// Create a new Starlark thread with a custom print function
	thread := &starlark.Thread{
		Name: "lint_thread",
		Print: func(thread *starlark.Thread, msg string) {
			outputMutex.Lock()
			defer outputMutex.Unlock()

			ruleName := thread.Local("current_rule").(string)
			moduleName := thread.Local("current_module").(string)

			entry := groupedOutput[ruleName][moduleName]
			entry.Outputs = append(entry.Outputs, msg)
			groupedOutput[ruleName][moduleName] = entry
		},
	}

	fo := &syntax.FileOptions{}

	// Parse the Starlark code without executing it
	f, err := fo.Parse(rulesFile, rulesContent, 0)
	if err != nil {
		log.Fatalf("Error parsing rules program: %v", err)
	}

	// Compile the parsed code
	prog, err := starlark.FileProgram(f, func(name string) bool {
		if name == "fatal" || name == "error" {
			return true
		}
		return false
	})
	if err != nil {
		log.Fatalf("Error compiling rules program: %v", err)
	}

	// Create predefined variables for the Starlark environment
	predefined := starlark.StringDict{
		"fatal": starlark.NewBuiltin("fatal", fatalFunc),
		"error": starlark.NewBuiltin("error", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			sep := " "
			if err := starlark.UnpackArgs("error", nil, kwargs, "sep?", &sep); err != nil {
				return nil, err
			}
			buf := new(strings.Builder)
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

			outputMutex.Lock()
			defer outputMutex.Unlock()

			ruleName := thread.Local("current_rule").(string)
			moduleName := thread.Local("current_module").(string)

			entry := groupedOutput[ruleName][moduleName]
			entry.Errors = append(entry.Errors, buf.String())
			groupedOutput[ruleName][moduleName] = entry

			return starlark.None, nil
		}),
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
		groupedOutput[ruleName] = make(map[string]RuleModuleOutput)
		for _, module := range modules {
			groupedOutput[ruleName][module.Path] = RuleModuleOutput{}
			starlarkModule := nf.ConvertToStarlarkModule(module)

			// Set the current rule and module context
			thread.SetLocal("current_rule", ruleName)
			thread.SetLocal("current_module", module.Path)

			_, err := starlark.Call(thread, ruleFunc, starlark.Tuple{starlarkModule}, nil)
			if err != nil {
				if evalErr, ok := err.(*starlark.EvalError); ok {
					entry := groupedOutput[ruleName][module.Path]
					entry.Errors = append(entry.Errors, evalErr.Msg)
					groupedOutput[ruleName][module.Path] = entry
					//fmt.Printf("Rule %s execution failed: %s\n", ruleName, evalErr.Msg)
				} else {
					log.Fatalf("Error calling rule %s: %v\n", ruleName, err)
				}
			}
		}
	}

	hasErrors := printGroupedOutput(groupedOutput)
	if hasErrors {
		os.Exit(1)
	}
}

func printGroupedOutput(groupedOutput GroupedOutput) bool {
	hasErrors := false
	ruleNames := make([]string, 0, len(groupedOutput))
	for ruleName := range groupedOutput {
		ruleNames = append(ruleNames, ruleName)
	}
	sort.Strings(ruleNames)

	rulePrinter := color.New(color.FgCyan, color.Bold)
	modulePrinter := color.New(color.FgYellow)
	errorPrinter := color.New(color.FgRed)
	outputPrinter := color.New(color.FgGreen)

	fmt.Println() // Start with a blank line

	for _, ruleName := range ruleNames {
		rulePrinter.Printf("Rule: %s\n", ruleName)

		moduleNames := make([]string, 0, len(groupedOutput[ruleName]))
		for moduleName := range groupedOutput[ruleName] {
			moduleNames = append(moduleNames, moduleName)
		}
		sort.Strings(moduleNames)

		for _, moduleName := range moduleNames {
			entry := groupedOutput[ruleName][moduleName]
			if len(entry.Errors) > 0 || len(entry.Outputs) > 0 {
				modulePrinter.Printf("  Module: %s\n", moduleName)
				for _, e := range entry.Errors {
					hasErrors = true
					errorPrinter.Printf("    Error: %s\n", e)
				}
				for _, o := range entry.Outputs {
					outputPrinter.Printf("    Output: %s\n", o)
				}
			}
		}
		fmt.Println() // Add a blank line between rules
	}
	return hasErrors
}
