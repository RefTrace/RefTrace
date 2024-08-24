package parser

import (
	"fmt"
	"strings"
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
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		return ParseResult{Mode: "Failed"}, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	lexerErrorListener := NewCustomErrorListener(filePath)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrorListener)

	stream := antlr.NewCommonTokenStream(lexer, 0)
	stream.Fill()

	parser := NewGroovyParser(stream)
	parserErrorListener := NewCustomErrorListener(filePath)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(parserErrorListener)

	strategy := NewCustomErrorStrategy()
	parser.SetErrorHandler(strategy)

	var result antlr.ParseTree
	var parseErr error

	// First, try SLL mode
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)

	func() {
		defer func() {
			if r := recover(); r != nil {
				parseErr = fmt.Errorf("SLL parsing panicked: %v", r)
			}
		}()
		result = parser.CompilationUnit()
	}()

	// Check for panics, or parser errors in SLL mode
	if parseErr != nil || len(strategy.GetErrors()) > 0 {
		// If SLL failed or there were lexer errors, try LL mode
		stream.Seek(0)
		parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLL)
		strategy.ClearErrors()
		lexerErrorListener = NewCustomErrorListener(filePath)
		parserErrorListener = NewCustomErrorListener(filePath)
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(lexerErrorListener)
		parser.RemoveErrorListeners()
		parser.AddErrorListener(parserErrorListener)
		parseErr = nil // Reset the panic error

		func() {
			defer func() {
				if r := recover(); r != nil {
					parseErr = fmt.Errorf("LL parsing panicked: %v", r)
				}
			}()
			result = parser.CompilationUnit()
		}()

		if parseErr != nil {
			return ParseResult{Mode: "Failed"}, fmt.Errorf("parsing failed in both SLL and LL modes: %w", parseErr)
		}

		/*
			var allErrors []string
			if len(strategy.GetErrors()) > 0 {
				for _, err := range strategy.GetErrors() {
					allErrors = append(allErrors, err.Error())
				}
			}
			if lexerErrorListener.HasErrors() {
				allErrors = append(allErrors, lexerErrorListener.GetErrors()...)
			}
			if parserErrorListener.HasErrors() {
				allErrors = append(allErrors, parserErrorListener.GetErrors()...)
			}

			if len(allErrors) > 0 {
				allErrors := append(lexerErrorListener.GetErrors(), parserErrorListener.GetErrors()...)
				return ParseResult{Mode: "Failed"}, fmt.Errorf("parsing failed in both SLL and LL modes: %v", allErrors)
			}
		*/

		return ParseResult{Tree: result, Mode: "LL"}, nil
	}

	return ParseResult{Tree: result, Mode: "SLL"}, nil
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
