package parser

import (
	"fmt"
	"strings"

	ast "reft-go/parser/ast"

	"github.com/antlr4-go/antlr/v4"
)

const (
	PACKAGE_INFO           = "package-info"
	PACKAGE_INFO_FILE_NAME = PACKAGE_INFO + ".groovy"
)

type NumberFormatError struct {
	Context   *antlr.ParserRuleContext
	Exception error
}

// Define the SyntaxException struct
type SyntaxException struct {
	Message           string
	StartLine         int
	StartCharPosition int
	StopLine          int
	StopCharPosition  int
}

// Implement the String method for SyntaxException
func (e *SyntaxException) String() string {
	return fmt.Sprintf("SyntaxException: %s (line %d, char %d to line %d, char %d)",
		e.Message, e.StartLine, e.StartCharPosition, e.StopLine, e.StopCharPosition)
}

// Function to create a CompilationFailedException
func createParsingFailedException(msg string, ctx *antlr.ParserRuleContext) *SyntaxException {
	start := (*ctx).GetStart()
	stop := (*ctx).GetStop()

	syntaxException := &SyntaxException{
		Message:           msg,
		StartLine:         start.GetLine(),
		StartCharPosition: start.GetTokenSource().GetCharPositionInLine() + 1,
		StopLine:          stop.GetLine(),
		StopCharPosition:  stop.GetTokenSource().GetCharPositionInLine() + 1 + len(stop.GetText()),
	}

	return syntaxException
}

type ASTBuilder struct {
	BaseGroovyParserVisitor
	moduleNode        *ModuleNode
	classNodeList     []*ast.ClassNode
	numberFormatError *NumberFormatError
	sourceUnitName    string
}

func (builder *ASTBuilder) VisitCompilationUnit(ctx *CompilationUnitContext) *ModuleNode {
	builder.Visit(ctx.PackageDeclaration())

	for _, node := range builder.VisitScriptStatements(ctx.ScriptStatements().(*ScriptStatementsContext)) {
		switch n := node.(type) {
		case *DeclarationListStatement:
			for _, stmt := range n.GetDeclarationStatements() {
				builder.moduleNode.AddStatement(&stmt.Statement)
			}
		case *Statement:
			builder.moduleNode.AddStatement(n)
		case *MethodNode:
			builder.moduleNode.AddMethod(n)
		}
	}

	for _, node := range builder.classNodeList {
		builder.moduleNode.AddClass(node)
	}

	if builder.isPackageInfoDeclaration() {
		packageInfo := ClassHelper.Make(builder.moduleNode.GetPackageName() + PACKAGE_INFO)
		if !builder.moduleNode.GetClasses().Contains(packageInfo) {
			builder.moduleNode.AddClass(packageInfo)
		}
	} else if builder.isBlankScript() {
		builder.moduleNode.AddStatement(ReturnStatement.RETURN_NULL_OR_VOID)
	}

	builder.configureScriptClassNode()

	if builder.numberFormatError != nil {
		panic(createParsingFailedException(builder.numberFormatError.Exception.Error(), builder.numberFormatError.Context))
	}

	return builder.moduleNode
}

func (builder *ASTBuilder) VisitScriptStatements(ctx *ScriptStatementsContext) []ast.ASTNode {
	if ctx == nil {
		return []ast.ASTNode{}
	}

	var nodes []ast.ASTNode
	for _, stmt := range ctx.AllScriptStatement() {
		nodes = append(nodes, builder.Visit(stmt).(ast.ASTNode))
	}

	return nodes
}

func (builder *ASTBuilder) isPackageInfoDeclaration() bool {
	name := builder.sourceUnitName
	return name != "" && strings.HasSuffix(name, PACKAGE_INFO_FILE_NAME)
}

func (builder *ASTBuilder) isBlankScript() bool {
	return len(builder.moduleNode.GetStatementBlock()) == 0 &&
		len(builder.moduleNode.GetMethods()) == 0 &&
		len(builder.moduleNode.GetClasses()) == 0
}
