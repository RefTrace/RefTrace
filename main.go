package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reft-go/parser" // Adjust the import path based on your module name and structure

	"github.com/antlr4-go/antlr/v4" // Ensure this import path is correct based on your setup
)

func main() {
	// Ensure the input directory path is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	// Walk the directory and process each .nf or .groovy file
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .nf or .groovy extension
		if !info.IsDir() && (filepath.Ext(path) == ".nf" || filepath.Ext(path) == ".groovy") {
			processFile(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", dir, err)
		os.Exit(1)
	}
}

func processFile(filePath string) {
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		fmt.Printf("Failed to open file %s: %s\n", filePath, err)
		return
	}

	// Create a new instance of the lexer
	lexer := parser.NewGroovyLexer(input)
	lexer.RemoveErrorListeners()
	errorListener := parser.NewCustomErrorListener(filePath)
	lexer.AddErrorListener(errorListener)
	_ = lexer.GetAllTokens()

	// Check for lexing errors
	if errorListener.HasError() {
		fmt.Printf("File: %s has errors.\n", filePath)
	}
}
