package parser

import (
	"path/filepath"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

func TestGroovyLexer(t *testing.T) {
	input := antlr.NewInputStream(`def foo = "bar"`)
	lexer := NewGroovyLexer(input)

	for {
		token := lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		if token.GetTokenType() == antlr.TokenInvalidType {
			t.Fatalf("Token recognition error: %s", token.GetText())
		}
		t.Logf("Token: %s, Type: %d", token.GetText(), token.GetTokenType())
	}
}

func TestGString(t *testing.T) {
	input := antlr.NewInputStream(`"Hello, ${name}!"`)
	lexer := NewGroovyLexer(input)

	for {
		token := lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		if token.GetTokenType() == antlr.TokenInvalidType {
			t.Fatalf("Token recognition error: %s", token.GetText())
		}
		t.Logf("Token: %s, Type: %d", token.GetText(), token.GetTokenType())
	}
}

func TestGroovyLexer2(t *testing.T) {
	s := `
    def subject = "$workflow.runName"
	`
	input := antlr.NewInputStream(s)
	lexer := NewGroovyLexer(input)

	for {
		token := lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		if token.GetTokenType() == antlr.TokenInvalidType {
			t.Fatalf("Token recognition error: %s", token.GetText())
		}
		t.Logf("Token: %s, Type: %d", token.GetText(), token.GetTokenType())
	}
}

func TestGStringFile(t *testing.T) {
	filePath := filepath.Join("testdata", "gstring.groovy")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	for {
		token := lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		if token.GetTokenType() == antlr.TokenInvalidType {
			t.Fatalf("Token recognition error: %s", token.GetText())
		}
		t.Logf("Token: %s, Type: %d", token.GetText(), token.GetTokenType())
	}
}

func TestBWAMem2File(t *testing.T) {
	filePath := filepath.Join("testdata", "bwamem2_mem.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	for {
		token := lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		if token.GetTokenType() == antlr.TokenInvalidType {
			t.Fatalf("Token recognition error: %s", token.GetText())
		}
		t.Logf("Token: %s, Type: %d", token.GetText(), token.GetTokenType())
	}
}

func TestGroovyLexerFromFile(t *testing.T) {
	filePath := filepath.Join("testdata", "utils_nfcore_pipeline.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	for {
		token := lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		if token.GetTokenType() == antlr.TokenInvalidType {
			t.Fatalf("Token recognition error: %s", token.GetText())
		}
		t.Logf("Token: %s, Type: %d", token.GetText(), token.GetTokenType())
	}
}

func TestIncludeLexer(t *testing.T) {
	filePath := filepath.Join("testdata", "include.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	stream.Fill()
	tokens := stream.GetAllTokens()

	expectedTokenTypes := []int{
		GroovyLexerIdentifier,
		GroovyLexerLBRACE,
		GroovyLexerCapitalizedIdentifier,
		GroovyLexerRBRACE,
		GroovyLexerIdentifier,
		GroovyLexerStringLiteral,
		antlr.TokenEOF,
	}

	expectedTokenTexts := []string{
		"include",
		"{",
		"SAREK",
		"}",
		"from",
		"'./workflows/sarek'",
		"<EOF>",
	}

	if len(tokens) != len(expectedTokenTypes) {
		t.Fatalf("Expected %d tokens, but got %d", len(expectedTokenTypes), len(tokens))
	}

	for i, token := range tokens {
		if token.GetTokenType() != expectedTokenTypes[i] {
			t.Errorf("Token %d: expected type %d, but got %d", i, expectedTokenTypes[i], token.GetTokenType())
		}
		if token.GetText() != expectedTokenTexts[i] {
			t.Errorf("Token %d: expected text '%s', but got '%s'", i, expectedTokenTexts[i], token.GetText())
		}
	}
}
