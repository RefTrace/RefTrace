package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
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
	processBaseVisitor()

	// Apply diffs to groovy_parser.go
	applyFirstDiffToGroovyParser()
	applySecondDiffToGroovyParser()

	fmt.Println("Files processed successfully.")
}

func processBaseVisitor() {
	// Read the file
	file, err := os.Open("groovyparser_base_visitor.go")
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
	output, err := os.Create("groovyparser_base_visitor.go")
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

func applyFirstDiffToGroovyParser() {
	filename := "groovy_parser.go"

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Define the old line, new lines, and the replacement
	oldLine := "return !isFollowingArgumentsOrClosure(localctx.(*CommandExpressionContext).Get_expression())"
	newLines := `if cmdExprCtx, ok := localctx.(*CommandExpressionContext); ok {
	return !isFollowingArgumentsOrClosure(cmdExprCtx.Get_expression())
}
return !isFollowingArgumentsOrClosure(localctx)`

	// Find the oldLine and insert the newLines
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == oldLine {
			// Get the indentation of the current line
			indent := strings.Repeat("\t", strings.Count(line, "\t"))

			// Insert the new lines with proper indentation
			indentedNewLines := indent + strings.ReplaceAll(newLines, "\n", "\n"+indent)
			lines[i] = indentedNewLines
			break
		}
	}

	// Join the lines back into a single string
	newContent := strings.Join(lines, "\n")

	// Write the modified content back to the file
	err = os.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("First diff applied to groovy_parser.go successfully.")
}

func applySecondDiffToGroovyParser() {
	filename := "groovy_parser.go"

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	// Second replacement
	oldLines := []string{
		"var t *CommandExpressionContext = nil",
		"if localctx != nil {",
		"t = localctx.(*CommandExpressionContext)",
		"}",
		"return p.CommandExpression_Sempred(t, predIndex)",
	}
	newLine := "return p.CommandExpression_Sempred(localctx, predIndex)"

	for i := 0; i < len(lines); i++ {
		if i+len(oldLines) <= len(lines) {
			match := true
			for j, oldLine := range oldLines {
				if strings.TrimSpace(lines[i+j]) != strings.TrimSpace(oldLine) {
					match = false
					break
				}
			}
			if match {
				indent := strings.Repeat("\t", strings.Count(lines[i], "\t"))
				lines[i] = indent + newLine
				lines = append(lines[:i+1], lines[i+len(oldLines):]...)
			}
		}
	}

	// Join the lines back into a single string
	newContent := strings.Join(lines, "\n")

	// Write the modified content back to the file
	err = os.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Second diff applied to groovy_parser.go successfully.")
}
