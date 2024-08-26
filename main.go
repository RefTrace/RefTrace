package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reft-go/parser"
	"sync"
	"sync/atomic"
	"time"

	// Adjust the import path based on your module name and structure
	// Ensure this import path is correct based on your setup
	"github.com/spf13/cobra"
)

var (
	version = "0.2.0"
)

var rootCmd = &cobra.Command{
	Use:  "reft",
	Long: `RefTrace - Static analysis for bioinformatics pipelines`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf("reft version %s\n", version)
			return
		}
		cmd.Help()
	},
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ProcessDirectory(dir string) (int64, int64, error) {
	var totalFiles, totalLines int64
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".nf" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				fileLines, err := processFile(path)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("error processing file %s", path))
					mu.Unlock()
					return
				}
				atomic.AddInt64(&totalFiles, 1)
				atomic.AddInt64(&totalLines, int64(fileLines))
			}(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return totalFiles, totalLines, err
	}

	if len(errors) > 0 {
		return totalFiles, totalLines, fmt.Errorf("encountered %d errors during processing: %v", len(errors), errors)
	}

	return totalFiles, totalLines, nil
}

func run(cmd *cobra.Command, args []string) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("Total execution time: %v\n", elapsed)
	}()

	dir := args[0]

	totalFiles, totalLines, err := ProcessDirectory(dir)

	fmt.Printf("Total files parsed: %d\n", totalFiles)
	fmt.Printf("Total lines processed: %d\n", totalLines)

	if err != nil {
		fmt.Printf("Errors encountered during processing:\n%v\n", err)
	}
}

func processFile(filePath string) (int, error) {
	_, err := parser.BuildAST(filePath)
	if err != nil {
		return 0, err
	}
	lineCount := countLines(filePath)
	return lineCount, nil
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
