package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reft-go/lexer"
	"sync"

	// Adjust the import path based on your module name and structure
	"github.com/antlr4-go/antlr/v4" // Ensure this import path is correct based on your setup
)

func main() {
	// Ensure the input directory path is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	var wg sync.WaitGroup

	// Walk the directory and process each .nf or .groovy file
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .nf or .groovy extension
		if !info.IsDir() && (filepath.Ext(path) == ".nf" || filepath.Ext(path) == ".groovy") {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				processFile(path)
			}(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", dir, err)
		os.Exit(1)
	}

	wg.Wait()
}

type TreeShapeListener struct {
	*lexer.BaseGroovyParserListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (tsl *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func processFile(filePath string) {
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		fmt.Printf("Failed to open file %s: %s\n", filePath, err)
		return
	}

	// Create a new instance of the lexer
	l := lexer.NewGroovyLexer(input)
	l.RemoveErrorListeners()
	errorListener := lexer.NewCustomErrorListener(filePath)
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
		fmt.Printf("File: %s has no errors.\n", filePath)
		//tokenStream := lexer.NewPreloadedTokenStream(tokens, l)
		p := lexer.NewGroovyParser(stream)
		p.CompilationUnit()
		fmt.Println("Parsed Successfully")
		//antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
	}
}
