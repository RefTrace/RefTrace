package main

import (
	"fmt"
	"github.com/antlr4-go/antlr" // Ensure this import path is correct based on your setup
	"os"
	"reft-go/parser" // Adjust the import path based on your module name and structure
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

	// Use the lexer to create a token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Assuming you have a parser that can handle the token stream
	// For example, if you have a GroovyParser, you would do something like this:
	p := parser.NewGroovyParser(stream) // Adjust according to your actual parser
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.StartRule() // Replace 'StartRule' with the actual entry rule of your parser

	// Walk the parse tree
	listener := NewTreeShapeListener() // Define or import TreeShapeListener if you have one
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
}
