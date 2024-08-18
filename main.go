package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reft-go/parser"
	"runtime/debug"
	"sync"

	// Adjust the import path based on your module name and structure
	"github.com/antlr4-go/antlr/v4" // Ensure this import path is correct based on your setup
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "reft [directory]",
	Short: "Process .nf files in a directory",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	debug.SetGCPercent(-1)
	dir := args[0]

	var wg sync.WaitGroup

	// Walk the directory and process each .nf or .groovy file
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .nf or .groovy extension
		// TODO: handle groovy extension
		if !info.IsDir() && (filepath.Ext(path) == ".nf") {
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

func processFile(filePath string) {
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		fmt.Printf("Failed to open file %s: %s\n", filePath, err)
		return
	}

	// Create a new instance of the lexer
	l := parser.NewGroovyLexer(input)
	l.RemoveErrorListeners()
	errorListener := parser.NewCustomErrorListener(filePath)
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
		p := parser.NewGroovyParser(stream)
		tree := p.CompilationUnit()
		fmt.Println("Parsed Successfully")
		builder := parser.NewASTBuilder(filePath)
		ast := builder.Visit(tree).(*parser.ModuleNode)
		_ = ast
		//builder.VisitCompilationUnit(unit.(*parser.CompilationUnitContext))
		//antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
	}
}
