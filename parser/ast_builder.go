package parser

import (
	"fmt"
	"strings"

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

// SyntaxException Define the SyntaxException struct
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
	classNodeList     []*ClassNode
	numberFormatError *NumberFormatError
	sourceUnitName    string
}

func (builder *ASTBuilder) VisitCompilationUnit(ctx *CompilationUnitContext) *ModuleNode {
	builder.Visit(ctx.PackageDeclaration())

	for _, node := range builder.VisitScriptStatements(ctx.ScriptStatements().(*ScriptStatementsContext)) {
		switch n := node.(type) {
		case *DeclarationListStatement:
			for _, stmt := range n.GetDeclarationStatements() {
				builder.moduleNode.AddStatement(stmt.Statement)
			}
		case Statement:
			builder.moduleNode.AddStatement(n)
		case *MethodNode:
			builder.moduleNode.AddMethod(n)
		}
	}

	for _, node := range builder.classNodeList {
		builder.moduleNode.AddClass(node)
	}

	if builder.isPackageInfoDeclaration() {
		packageInfo := MakeFromString(builder.moduleNode.GetPackageName() + PACKAGE_INFO)
		if !builder.moduleNode.Contains(packageInfo) {
			builder.moduleNode.AddClass(packageInfo)
		}
	} else if builder.isBlankScript() {
		builder.moduleNode.AddStatement(RETURN_NULL_OR_VOID)
	}

	builder.configureScriptClassNode()

	if builder.numberFormatError != nil {
		panic(createParsingFailedException(builder.numberFormatError.Exception.Error(), builder.numberFormatError.Context))
	}

	return builder.moduleNode
}

func (builder *ASTBuilder) configureScriptClassNode() {
	scriptClassNode := builder.moduleNode.GetScriptClassDummy()
	if scriptClassNode != nil {
		statements := builder.moduleNode.GetStatementBlock().GetStatements()
		if len(statements) > 0 {
			firstStatement := statements[0]
			scriptClassNode.SetSourcePosition(firstStatement)
			lastStatement := statements[len(statements)-1]
			scriptClassNode.SetLastLineNumber(lastStatement.GetLastLineNumber())
			scriptClassNode.SetLastColumnNumber(lastStatement.GetLastColumnNumber())
		}
	}
}

func (builder *ASTBuilder) VisitScriptStatements(ctx *ScriptStatementsContext) []ASTNode {
	if ctx == nil {
		return []ASTNode{}
	}

	var nodes []ASTNode
	for _, stmt := range ctx.AllScriptStatement() {
		nodes = append(nodes, builder.Visit(stmt).(ASTNode))
	}

	return nodes
}

func (builder *ASTBuilder) isPackageInfoDeclaration() bool {
	name := builder.sourceUnitName
	return name != "" && strings.HasSuffix(name, PACKAGE_INFO_FILE_NAME)
}

func (builder *ASTBuilder) isBlankScript() bool {
	return len(builder.moduleNode.GetStatementBlock().GetStatements()) == 0 &&
		len(builder.moduleNode.GetMethods()) == 0 &&
		len(builder.moduleNode.GetClasses()) == 0
}

// DeclarationListStatement represents a list of declaration statements
type DeclarationListStatement struct {
	Statement
	declarationStatements []*ExpressionStatement
}

// NewDeclarationListStatement creates a new DeclarationListStatement from a list of DeclarationExpressions
func NewDeclarationListStatement(declarations ...*DeclarationExpression) *DeclarationListStatement {
	declarationStatements := make([]*ExpressionStatement, len(declarations))
	for i, decl := range declarations {
		stmt, err := NewExpressionStatement(decl)
		if err != nil {
			panic(err)
		}
		declarationStatements[i] = configureASTFromSource(stmt, decl)
	}
	return &DeclarationListStatement{declarationStatements: declarationStatements}
}

// GetDeclarationStatements returns the list of ExpressionStatements
func (d *DeclarationListStatement) GetDeclarationStatements() []*ExpressionStatement {
	declarationListStatementLabels := d.GetStatementLabels()

	for _, e := range d.declarationStatements {
		if declarationListStatementLabels != nil {
			// clear existing statement labels before setting labels
			if e.GetStatementLabels() != nil {
				e.ClearStatementLabels()
			}

			for label := declarationListStatementLabels.Front(); label != nil; label = label.Next() {
				e.AddStatementLabel(label.Value.(string))
			}
		}
	}

	return d.declarationStatements
}

// GetDeclarationExpressions returns the list of DeclarationExpressions
func (d *DeclarationListStatement) GetDeclarationExpressions() []*DeclarationExpression {
	declarations := make([]*DeclarationExpression, len(d.declarationStatements))
	for i, stmt := range d.declarationStatements {
		declarations[i] = stmt.GetExpression().(*DeclarationExpression)
	}
	return declarations
}
