package main

import (
	"fmt"
	"os"
	"reft-go/parser" // Adjust the import path based on your module name and structure

	"github.com/antlr4-go/antlr/v4" // Ensure this import path is correct based on your setup
)

func main() {
	// Ensure the input file path is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_input_file>")
		os.Exit(1)
	}

	input, err := antlr.NewFileStream(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err)
		os.Exit(1)
	}

	// Create a new instance of the lexer
	lexer := parser.NewGroovyLexer(input)
	tokens := lexer.GetAllTokens()
	println(len(tokens))
}
