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

func TestGroovyLexer2(t *testing.T) {
	s := `
    def subject = "[$workflow.manifest.name] Successful: $workflow.runName"
    if (!workflow.success) {
        subject = "[$workflow.manifest.name] FAILED: $workflow.runName"
    }
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
