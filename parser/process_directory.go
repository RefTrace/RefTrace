package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/antlr4-go/antlr/v4"
)

func ProcessDirectory(dir string) (int64, int64, error) {
	var totalFiles, totalLines int64
	var wg sync.WaitGroup

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (filepath.Ext(path) == ".nf") {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				fileLines := processFile(path)
				atomic.AddInt64(&totalFiles, 1)
				atomic.AddInt64(&totalLines, int64(fileLines))
			}(path)
		}
		return nil
	})

	wg.Wait()

	return totalFiles, totalLines, err
}

func processFile(filePath string) int {
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		fmt.Printf("Failed to open file %s: %s\n", filePath, err)
		return 0
	}

	lineCount := countLines(filePath)

	// Create a new instance of the lexer
	l := NewGroovyLexer(input)
	l.RemoveErrorListeners()
	errorListener := NewCustomErrorListener(filePath)
	l.AddErrorListener(errorListener)
	//tokens := l.GetAllTokens()
	stream := antlr.NewCommonTokenStream(l, 0)
	stream.Fill()

	// Print the token type for each token
	/*
		for _, token := range tokens {
			fmt.Printf("Token: %s, Type: %d\n", token.GetText(), token.GetTokenType())
		}
	*/

	// Check for lexing errors
	if !errorListener.HasError() {
		//fmt.Printf("File: %s has no errors.\n", filePath)
		//tokenStream := lexer.NewPreloadedTokenStream(tokens, l)
		p := NewGroovyParser(stream)
		tree := p.CompilationUnit()
		//fmt.Println("Parsed Successfully")
		builder := NewASTBuilder(filePath)
		ast := builder.Visit(tree).(*ModuleNode)
		_ = ast
		//builder.VisitCompilationUnit(unit.(*parser.CompilationUnitContext))
		//antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
	}
	return lineCount
}

func countLines(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file for line counting %s: %s\n", filePath, err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error counting lines in %s: %s\n", filePath, err)
	}

	return lineCount
}
