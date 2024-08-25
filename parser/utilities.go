package parser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
)

var invalidJavaIdentifiers = map[string]bool{
	"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true,
	"case": true, "catch": true, "char": true, "class": true, "const": true,
	"continue": true, "default": true, "do": true, "double": true, "else": true,
	"enum": true, "extends": true, "final": true, "finally": true, "float": true,
	"for": true, "goto": true, "if": true, "implements": true, "import": true,
	"instanceof": true, "int": true, "interface": true, "long": true, "native": true,
	"new": true, "package": true, "private": true, "protected": true, "public": true,
	"short": true, "static": true, "strictfp": true, "super": true, "switch": true,
	"synchronized": true, "this": true, "throw": true, "throws": true,
	"transient": true, "try": true, "void": true, "volatile": true, "while": true,
	"true": true, "false": true, "null": true,
}

// RepeatString returns a string made up of repetitions of the specified string.
func RepeatString(pattern string, repeats int) string {
	return strings.Repeat(pattern, repeats)
}

// EOL returns the end-of-line marker.
func EOL() string {
	return "\n"
}

// IsJavaIdentifier tells if the given string is a valid Java identifier.
func IsJavaIdentifier(name string) bool {
	if len(name) == 0 || invalidJavaIdentifiers[name] {
		return false
	}
	runes := []rune(name)
	if !unicode.IsLetter(runes[0]) && runes[0] != '_' {
		return false
	}
	for _, r := range runes[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

type ParseResult struct {
	Tree antlr.ParseTree
	Mode string // "SLL", "LL", or "Failed"
}

func BuildCST(filePath string) (ParseResult, error) {
	// Try SLL mode first
	result, err := TryBuildCST(filePath, antlr.PredictionModeSLL)
	if err == nil {
		return ParseResult{Tree: result, Mode: "SLL"}, nil
	}

	// If SLL failed, try LL mode
	result, err = TryBuildCST(filePath, antlr.PredictionModeLL)
	if err == nil {
		return ParseResult{Tree: result, Mode: "LL"}, nil
	}

	// If both modes failed, return the error
	return ParseResult{Mode: "Failed"}, fmt.Errorf("parsing failed in both SLL and LL modes: %w", err)
}

func TryBuildCST(filePath string, mode int) (antlr.ParseTree, error) {
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	lexerErrorListener := NewCustomErrorListener(filePath)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrorListener)

	stream := antlr.NewCommonTokenStream(lexer, 0)
	stream.Fill()

	parser := NewGroovyParser(stream)
	parser.RemoveErrorListeners()

	strategy := NewCustomErrorStrategy()
	parser.SetErrorHandler(strategy)

	parser.GetInterpreter().SetPredictionMode(mode)

	var result antlr.ParseTree
	var parseErr error

	var modeStr string
	if mode == antlr.PredictionModeSLL {
		modeStr = "SLL"
	} else {
		modeStr = "LL"
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				parseErr = fmt.Errorf("%s parsing panicked: %v", modeStr, r)
			}
		}()
		result = parser.CompilationUnit()
	}()

	if parseErr != nil {
		return nil, parseErr
	}

	var allErrors []string
	for _, err := range strategy.GetErrors() {
		allErrors = append(allErrors, err.Error())
	}
	if len(allErrors) > 0 {
		return nil, fmt.Errorf("parsing failed in %s mode: %v", modeStr, allErrors)
	}

	return result, nil
}

// BuildCSTTest builds a CST without error handling, for testing purposes.
// It panics on any errors, which is useful for stack trace testing.
func BuildCSTTest(filePath string) ParseResult {
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to open file %s: %v", filePath, err))
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	stream.Fill()

	parser := NewGroovyParser(stream)
	strategy := NewCustomErrorStrategy()
	parser.SetErrorHandler(strategy)

	// Try SLL mode
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	result := parser.CompilationUnit()

	if len(strategy.GetErrors()) > 0 {
		// If SLL failed, try LL mode
		stream.Seek(0)
		parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLL)
		strategy.ClearErrors()
		result = parser.CompilationUnit()

		if len(strategy.GetErrors()) > 0 {
			panic("parsing failed in both SLL and LL modes")
		}

		return ParseResult{Tree: result, Mode: "LL"}
	}

	return ParseResult{Tree: result, Mode: "SLL"}
}

// BuildAST builds the Abstract Syntax Tree (AST) for the given file.
func BuildAST(filePath string) (*ModuleNode, error) {
	// First, build the Concrete Syntax Tree (CST)
	parseResult, err := BuildCST(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to build CST: %w", err)
	}

	// Now, build the AST
	var ast *ModuleNode
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic while building AST: %v", r)
			}
		}()
		builder := NewASTBuilder(filePath)
		ast = builder.Visit(parseResult.Tree).(*ModuleNode)
	}()

	if err != nil {
		return nil, err
	}

	return ast, nil
}

// BuildASTTest builds the Abstract Syntax Tree (AST) for the given file without error handling.
// It panics on any errors, which is useful for stack trace testing.
func BuildASTTest(filePath string) *ModuleNode {
	// First, build the Concrete Syntax Tree (CST)
	parseResult := BuildCSTTest(filePath)

	// Now, build the AST
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(parseResult.Tree).(*ModuleNode)

	return ast
}

func processFile(filePath string) error {
	_, err := BuildAST(filePath)
	if err != nil {
		return err
	}
	return nil
}

func processDirectory(dir string) (int64, error) {
	var totalFiles int64
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".nf" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				err := processFile(path)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("error processing file %s", path))
					mu.Unlock()
					return
				}
				atomic.AddInt64(&totalFiles, 1)
			}(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return totalFiles, err
	}

	if len(errors) > 0 {
		return totalFiles, fmt.Errorf("encountered %d errors during processing: %v", len(errors), errors)
	}

	return totalFiles, nil
}
