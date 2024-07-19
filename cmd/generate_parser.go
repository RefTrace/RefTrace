package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	// Change to the parser directory
	err := os.Chdir("../parser")
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	// Run ANTLR4
	cmd := exec.Command("sh", "-c", `java -Xmx500M -cp "./antlr-4.13.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Go -no-listener -visitor -package parser *.g4`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running ANTLR4:", err)
		return
	}

	// Post-processing
	filename := "groovyparser_base_visitor.go"

	// Read the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Process the lines
	structRegex := regexp.MustCompile(`type BaseGroovyParserVisitor struct {`)
	visitRegex := regexp.MustCompile(`func \(v \*BaseGroovyParserVisitor\) Visit(\w+)\(ctx \*(\w+)Context\) interface{} {`)
	for i, line := range lines {
		if match := structRegex.FindStringSubmatch(line); match != nil {
			lines[i+1] = "\t*antlr.BaseParseTreeVisitor"
			lines = append(lines[:i+2], append([]string{"\tVisitChildren func(node antlr.RuleNode) interface{}"}, lines[i+2:]...)...)
		}
		if match := visitRegex.FindStringSubmatch(line); match != nil {
			lines[i+1] = "\treturn v.VisitChildren(ctx)"
		}
	}

	// Write the modified content back to the file
	output, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	fmt.Println("File processed successfully.")
}
