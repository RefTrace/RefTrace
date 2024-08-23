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

func buildCST(filePath string) (ParseResult, error) {
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		return ParseResult{Mode: "Failed"}, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	stream.Fill()
	parser := NewGroovyParser(stream)
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

	// Check for panics or errors in SLL mode
	if parseErr != nil || len(strategy.GetErrors()) > 0 {
		// If SLL failed, try LL mode
		stream.Seek(0)
		parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLL)
		strategy.ClearErrors()
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
			return ParseResult{Mode: "Failed"}, fmt.Errorf("parsing failed due to panic in both SLL and LL modes: %w", parseErr)
		}

		if len(strategy.GetErrors()) > 0 {
			return ParseResult{Mode: "Failed"}, fmt.Errorf("parsing failed in both SLL and LL modes: %v", strategy.GetErrors())
		}

		return ParseResult{Tree: result, Mode: "LL"}, nil
	}

	return ParseResult{Tree: result, Mode: "SLL"}, nil
}
