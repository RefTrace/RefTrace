package parser

import (
	"container/list"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
)

var _ GroovyParserVisitor = (*ASTBuilder)(nil)

const (
	QUESTION_STR     = "?"
	DOT_STR          = "."
	SUB_STR          = "-"
	ASSIGN_STR       = "="
	VALUE_STR        = "value"
	DOLLAR_STR       = "$"
	CALL_STR         = "call"
	THIS_STR         = "this"
	SUPER_STR        = "super"
	VOID_STR         = "void"
	SLASH_STR        = "/"
	SLASH_DOLLAR_STR = "/$"
	TDQ_STR          = "\"\"\""
	TSQ_STR          = "'''"
	SQ_STR           = "'"
	DQ_STR           = "\""
	DOLLAR_SLASH_STR = "$/"

	PACKAGE_INFO           = "package-info"
	PACKAGE_INFO_FILE_NAME = PACKAGE_INFO + ".groovy"

	CLASS_NAME                                = "CLASS_NAME"
	INSIDE_PARENTHESES_LEVEL                  = "INSIDE_PARENTHESES_LEVEL"
	IS_INSIDE_INSTANCEOF_EXPR                 = "IS_INSIDE_INSTANCEOF_EXPR"
	IS_SWITCH_DEFAULT                         = "IS_SWITCH_DEFAULT"
	IS_NUMERIC                                = "IS_NUMERIC"
	IS_STRING                                 = "IS_STRING"
	IS_INTERFACE_WITH_DEFAULT_METHODS         = "IS_INTERFACE_WITH_DEFAULT_METHODS"
	IS_INSIDE_CONDITIONAL_EXPRESSION          = "IS_INSIDE_CONDITIONAL_EXPRESSION"
	IS_COMMAND_EXPRESSION                     = "IS_COMMAND_EXPRESSION"
	IS_BUILT_IN_TYPE                          = "IS_BUILT_IN_TYPE"
	PATH_EXPRESSION_BASE_EXPR                 = "PATH_EXPRESSION_BASE_EXPR"
	PATH_EXPRESSION_BASE_EXPR_GENERICS_TYPES  = "PATH_EXPRESSION_BASE_EXPR_GENERICS_TYPES"
	PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN      = "PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN"
	CMD_EXPRESSION_BASE_EXPR                  = "CMD_EXPRESSION_BASE_EXPR"
	TYPE_DECLARATION_MODIFIERS                = "TYPE_DECLARATION_MODIFIERS"
	COMPACT_CONSTRUCTOR_DECLARATION_MODIFIERS = "COMPACT_CONSTRUCTOR_DECLARATION_MODIFIERS"
	CLASS_DECLARATION_CLASS_NODE              = "CLASS_DECLARATION_CLASS_NODE"
	VARIABLE_DECLARATION_VARIABLE_TYPE        = "VARIABLE_DECLARATION_VARIABLE_TYPE"
	ANONYMOUS_INNER_CLASS_SUPER_CLASS         = "ANONYMOUS_INNER_CLASS_SUPER_CLASS"
	INTEGER_LITERAL_TEXT                      = "INTEGER_LITERAL_TEXT"
	FLOATING_POINT_LITERAL_TEXT               = "FLOATING_POINT_LITERAL_TEXT"
	ENCLOSING_INSTANCE_EXPRESSION             = "ENCLOSING_INSTANCE_EXPRESSION"
	IS_YIELD_STATEMENT                        = "IS_YIELD_STATEMENT"
	PARAMETER_MODIFIER_MANAGER                = "PARAMETER_MODIFIER_MANAGER"
	PARAMETER_CONTEXT                         = "PARAMETER_CONTEXT"
	IS_RECORD_GENERATED                       = "IS_RECORD_GENERATED"
	RECORD_HEADER                             = "RECORD_HEADER"
	RECORD_TYPE_NAME                          = "groovy.transform.RecordType"
)

var (
	FOR_LOOP_DUMMY *Parameter = &Parameter{} // You might want to initialize this with appropriate values
)

var QUOTATION_MAP = map[string]string{
	DQ_STR:           DQ_STR,
	SQ_STR:           SQ_STR,
	TDQ_STR:          TDQ_STR,
	TSQ_STR:          TSQ_STR,
	SLASH_STR:        SLASH_STR,
	DOLLAR_SLASH_STR: SLASH_DOLLAR_STR,
}

type NumberFormatError struct {
	Context   antlr.ParserRuleContext
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

// SourcePosition is an interface that both antlr.ParserRuleContext and ASTNode should implement
type SourcePosition interface {
	GetStartLine() int
	GetStartColumn() int
	GetStopLine() int
	GetStopColumn() int
}

// Implement SourcePosition for antlr.ParserRuleContext
type parserRuleContextAdapter struct {
	antlr.ParserRuleContext
}

func (p parserRuleContextAdapter) GetStartLine() int {
	return p.GetStart().GetLine()
}

func (p parserRuleContextAdapter) GetStartColumn() int {
	return p.GetStart().GetColumn()
}

func (p parserRuleContextAdapter) GetStopLine() int {
	return p.GetStop().GetLine()
}

func (p parserRuleContextAdapter) GetStopColumn() int {
	return p.GetStop().GetColumn() + len(p.GetStop().GetText())
}

// Implement SourcePosition for antlr.Token
type tokenAdapter struct {
	antlr.Token
}

func (t tokenAdapter) GetStartLine() int {
	return t.GetLine()
}

func (t tokenAdapter) GetStartColumn() int {
	return t.GetColumn()
}

func (t tokenAdapter) GetStopLine() int {
	return t.GetLine() // Tokens typically represent a single line
}

func (t tokenAdapter) GetStopColumn() int {
	return t.GetColumn() + len(t.GetText()) - 1 // -1 because column is 0-based
}

// Implement SourcePosition for ASTNode
type astNodeAdapter struct {
	ASTNode
}

func (a astNodeAdapter) GetStartLine() int {
	return a.GetLineNumber()
}

func (a astNodeAdapter) GetStartColumn() int {
	return a.GetColumnNumber()
}

func (a astNodeAdapter) GetStopLine() int {
	return a.GetLastLineNumber()
}

func (a astNodeAdapter) GetStopColumn() int {
	return a.GetLastColumnNumber()
}

// Generic createParsingFailedException function
func createParsingFailedException[T SourcePosition](msg string, source T) *SyntaxException {
	return &SyntaxException{
		Message:           msg,
		StartLine:         source.GetStartLine(),
		StartCharPosition: source.GetStartColumn(),
		StopLine:          source.GetStopLine(),
		StopCharPosition:  source.GetStopColumn(),
	}
}

type ASTBuilder struct {
	BaseGroovyParserVisitor
	moduleNode                                *ModuleNode
	classNodeList                             []*ClassNode
	numberFormatError                         *NumberFormatError
	sourceUnitName                            string
	visitingAssertStatementCount              int
	visitingClosureCount                      int
	visitingLoopStatementCount                int
	visitingSwitchStatementCount              int
	visitingArrayInitializerCount             int
	switchExpressionRuleContextStack          *list.List
	switchExpressionVariableSeq               int
	classNodeStack                            *list.List
	anonymousInnerClassesDefinedInMethodStack *list.List
}

// NewASTBuilder creates and initializes a new ASTBuilder instance
func NewASTBuilder(sourceUnitName string) *ASTBuilder {
	builder := &ASTBuilder{
		moduleNode:                       NewModuleNode(), // Assuming you have a NewModuleNode function
		classNodeList:                    make([]*ClassNode, 0),
		numberFormatError:                nil,
		sourceUnitName:                   sourceUnitName,
		switchExpressionRuleContextStack: list.New(),
		classNodeStack:                   list.New(),
		anonymousInnerClassesDefinedInMethodStack: list.New(),
	}
	builder.BaseGroovyParserVisitor.VisitChildren = builder.VisitChildren
	return builder
}

func (builder *ASTBuilder) Visit(tree antlr.ParseTree) interface{} {
	if tree != nil {
		return tree.Accept(builder)
	}
	return nil
}

func (builder *ASTBuilder) VisitChildren(tree antlr.RuleNode) interface{} {
	for _, c := range tree.GetChildren() {
		v := c.(antlr.ParseTree)
		builder.Visit(v)
	}
	return nil
}

// pushAnonymousInnerClass adds a new list of InnerClassNode to the stack
func (v *ASTBuilder) pushAnonymousInnerClass(innerClassList *list.List) {
	v.anonymousInnerClassesDefinedInMethodStack.PushFront(innerClassList)
}

// popAnonymousInnerClass removes and returns the top list of InnerClassNode from the stack
func (v *ASTBuilder) popAnonymousInnerClass() *list.List {
	if v.anonymousInnerClassesDefinedInMethodStack.Len() == 0 {
		panic("empty anonymous inner class stack")
	}
	return v.anonymousInnerClassesDefinedInMethodStack.Remove(v.anonymousInnerClassesDefinedInMethodStack.Front()).(*list.List)
}

// peekAnonymousInnerClass returns the top list of InnerClassNode from the stack without removing it
func (v *ASTBuilder) peekAnonymousInnerClass() *list.List {
	if v.anonymousInnerClassesDefinedInMethodStack.Len() == 0 {
		panic("peek empty anonymous inner class stack")
	}
	return v.anonymousInnerClassesDefinedInMethodStack.Front().Value.(*list.List)
}

// addAnonymousInnerClass adds an InnerClassNode to the top list in the stack
func (v *ASTBuilder) addAnonymousInnerClass(innerClass *InnerClassNode) {
	if v.anonymousInnerClassesDefinedInMethodStack.Len() == 0 {
		panic("cannot add to empty anonymous inner class stack")
	}
	v.peekAnonymousInnerClass().PushBack(innerClass)
}

func (v *ASTBuilder) pushSwitchExpressionRuleContext(ctx antlr.ParserRuleContext) {
	v.switchExpressionRuleContextStack.PushFront(ctx)
}

func (v *ASTBuilder) popSwitchExpressionRuleContext() antlr.ParserRuleContext {
	if v.switchExpressionRuleContextStack.Len() == 0 {
		panic("empty rule context")
	}
	return v.switchExpressionRuleContextStack.Remove(v.switchExpressionRuleContextStack.Front()).(antlr.ParserRuleContext)
}

func (v *ASTBuilder) peekSwitchExpressionRuleContext() antlr.ParserRuleContext {
	if v.switchExpressionRuleContextStack.Len() == 0 {
		panic("peek empty rule context")
	}
	return v.switchExpressionRuleContextStack.Front().Value.(antlr.ParserRuleContext)
}

func (v *ASTBuilder) pushClassNode(classNode *ClassNode) {
	v.classNodeStack.PushFront(classNode)
}

func (v *ASTBuilder) popClassNode() *ClassNode {
	if v.classNodeStack.Len() == 0 {
		panic("empty class node stack")
	}
	return v.classNodeStack.Remove(v.classNodeStack.Front()).(*ClassNode)
}

func (v *ASTBuilder) peekClassNode() *ClassNode {
	if v.classNodeStack.Len() == 0 {
		panic("peek empty class node stack")
	}
	return v.classNodeStack.Front().Value.(*ClassNode)
}

func (builder *ASTBuilder) VisitCompilationUnit(ctx *CompilationUnitContext) interface{} {
	//builder.VisitPackageDeclaration(ctx.PackageDeclaration().(*PackageDeclarationContext))
	builder.Visit(ctx.PackageDeclaration())

	for _, node := range builder.VisitScriptStatements(ctx.ScriptStatements().(*ScriptStatementsContext)).([]ASTNode) {
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
		panic(createParsingFailedException(builder.numberFormatError.Exception.Error(), parserRuleContextAdapter{builder.numberFormatError.Context}))
	}

	return builder.moduleNode
}

func (builder *ASTBuilder) VisitScriptStatements(ctx *ScriptStatementsContext) interface{} {
	if ctx == nil {
		return []ASTNode{}
	}

	var nodes []ASTNode
	for _, stmt := range ctx.AllScriptStatement() {
		nodes = append(nodes, builder.Visit(stmt).(ASTNode))
	}

	return nodes
}

func (v *ASTBuilder) VisitPackageDeclaration(ctx *PackageDeclarationContext) *PackageNode {
	packageName := v.VisitQualifiedName(ctx.QualifiedName().(*QualifiedNameContext))
	v.moduleNode.SetPackageName(packageName + DOT_STR)

	packageNode := v.moduleNode.PackageNode
	annotations := v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext))

	packageNode.AddAnnotations(annotations)

	return configureAST(packageNode, ctx)
}

func (v *ASTBuilder) VisitImportDeclaration(ctx *ImportDeclarationContext) *ImportNode {
	annotations := v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext))

	hasStatic := ctx.STATIC() != nil
	hasStar := ctx.MUL() != nil
	hasAlias := ctx.alias != nil

	var importNode *ImportNode

	if hasStatic {
		if hasStar { // e.g. import static java.lang.Math.*
			qualifiedName := v.VisitQualifiedName(ctx.QualifiedName().(*QualifiedNameContext))
			importType := MakeFromString(qualifiedName)
			configureAST(importType, ctx.QualifiedName())

			v.moduleNode.AddStaticStarImportWithAnnotations(importType.GetText(), importType, annotations)
			var imports map[string]*ImportNode = v.moduleNode.GetStaticStarImports()
			importNode = lastMapValue(imports)
		} else { // e.g. import static java.lang.Math.pow
			identifierList := ctx.QualifiedName().(*QualifiedNameContext).AllQualifiedNameElement()
			identifierListSize := len(identifierList)

			qualifiedName := strings.Join(sliceMap(identifierList[:identifierListSize-1], func(e IQualifiedNameElementContext) string { return e.GetText() }), DOT_STR)
			importType := MakeFromString(qualifiedName)
			configureAST(importType, ctx.QualifiedName()) // qualifiedName() includes member name
			configureEndPosition(importType, identifierList[max(0, identifierListSize-2)].GetStop())

			memberName := identifierList[identifierListSize-1].GetText()
			simpleName := memberName
			if hasAlias {
				simpleName = ctx.alias.GetText()
			}

			v.moduleNode.AddStaticImport(importType, memberName, simpleName, annotations)
			importNode = lastMapValue(v.moduleNode.GetStaticImports())
		}
	} else {
		if hasStar { // e.g. import java.util.*
			qualifiedName := v.VisitQualifiedName(ctx.QualifiedName().(*QualifiedNameContext))
			v.moduleNode.AddStarImportWithAnnotations(qualifiedName+DOT_STR, annotations)
			importNode = last(v.moduleNode.GetStarImports())
		} else { // e.g. import java.util.Map
			qualifiedName := v.VisitQualifiedName(ctx.QualifiedName().(*QualifiedNameContext))
			importType := MakeFromString(qualifiedName)
			configureAST(importType, ctx.QualifiedName())

			simpleName := last(ctx.QualifiedName().(*QualifiedNameContext).AllQualifiedNameElement()).GetText()
			if hasAlias {
				simpleName = ctx.alias.GetText()
			}

			v.moduleNode.AddImportWithAnnotations(simpleName, importType, annotations)
			importNode = last(v.moduleNode.GetImports())
		}
	}

	return configureAST(importNode, ctx)
}

// Helper functions

func sliceMap[T any, R any](slice []T, f func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

func last[T any](slice []T) T {
	return slice[len(slice)-1]
}

func lastMapValue[K comparable, V any](m map[K]V) V {
	var last V
	for _, v := range m {
		last = v
	}
	return last
}

// makeAnnotationNode creates an AnnotationNode for the given annotation type
func makeAnnotationNode(annotationType reflect.Type) *AnnotationNode {
	classNode := Make(annotationType)
	node := NewAnnotationNode(classNode)
	// TODO: source offsets
	return node
}

// makeClassNode creates a ClassNode for the given class name
func makeClassNode(name string) *ClassNode {
	node := MakeFromString(name)
	// TODO: shared instances
	return node
}

func (v *ASTBuilder) VisitAssertStatement(ctx *AssertStatementContext) *AssertStatement {
	v.visitingAssertStatementCount++
	defer func() {
		v.visitingAssertStatementCount--
	}()

	conditionExpression := v.Visit(ctx.ce).(Expression)

	if binaryExpression, ok := conditionExpression.(*BinaryExpression); ok {
		if binaryExpression.GetOperation().GetType() == ASSIGN {
			panic(createParsingFailedException("Assignment expression is not allowed in the assert statement", astNodeAdapter{conditionExpression}))
		}
	}

	booleanExpression := configureAST(
		NewBooleanExpression(conditionExpression),
		ctx,
	)

	if ctx.me == nil {
		return configureAST(
			NewAssertStatement(booleanExpression),
			ctx,
		)
	}

	return configureAST(
		NewAssertStatementWithMessage(
			booleanExpression,
			v.Visit(ctx.me).(Expression),
		),
		ctx,
	)
}

func (v *ASTBuilder) VisitConditionalStatement(ctx *ConditionalStatementContext) Statement {
	if ctx.IfElseStatement() != nil {
		return configureAST(v.VisitIfElseStatement(ctx.IfElseStatement().(*IfElseStatementContext)), ctx)
	} else if ctx.SwitchStatement() != nil {
		return configureAST(v.VisitSwitchStatement(ctx.SwitchStatement().(*SwitchStatementContext)), ctx)
	}

	panic(createParsingFailedException("Unsupported conditional statement", parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitIfElseStatement(ctx *IfElseStatementContext) *IfStatement {
	conditionExpression := v.VisitExpressionInPar(ctx.ExpressionInPar().(*ExpressionInParContext))
	booleanExpression := configureAST(
		NewBooleanExpression(conditionExpression),
		ctx,
	)

	ifBlock := v.unpackStatement(v.Visit(ctx.tb).(Statement))
	var elseBlock Statement
	if ctx.ELSE() != nil {
		elseBlock = v.unpackStatement(v.Visit(ctx.fb).(Statement))
	} else {
		elseBlock = NewEmptyStatement()
	}

	return configureAST(NewIfStatement(booleanExpression, ifBlock, elseBlock), ctx)
}

func (v *ASTBuilder) VisitLoopStmtAlt(ctx *LoopStmtAltContext) Statement {
	v.pushSwitchExpressionRuleContext(ctx)
	v.visitingLoopStatementCount++
	defer func() {
		v.popSwitchExpressionRuleContext()
		v.visitingLoopStatementCount--
	}()

	return configureAST(v.Visit(ctx.LoopStatement()).(Statement), ctx)
}

func (v *ASTBuilder) VisitForStmtAlt(ctx *ForStmtAltContext) *ForStatement {
	controlTuple := v.VisitForControl(ctx.ForControl().(*ForControlContext))

	loopBlock := v.unpackStatement(v.Visit(ctx.Statement()).(Statement))

	var block Statement
	if loopBlock != nil {
		block = loopBlock
	} else {
		block = NewEmptyStatement()
	}

	return configureAST(
		NewForStatement(controlTuple.V1, controlTuple.V2, block),
		ctx,
	)
}

func (v *ASTBuilder) VisitForControl(ctx *ForControlContext) Tuple2[*Parameter, Expression] {
	if ctx.EnhancedForControl() != nil { // e.g. for(int i in 0..<10) {}
		return v.VisitEnhancedForControl(ctx.EnhancedForControl().(*EnhancedForControlContext))
	}

	if ctx.ClassicalForControl() != nil { // e.g. for(int i = 0; i < 10; i++) {}
		return v.VisitClassicalForControl(ctx.ClassicalForControl().(*ClassicalForControlContext))
	}

	panic(createParsingFailedException("Unsupported for control: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitForInit(ctx *ForInitContext) Expression {
	if ctx == nil {
		return EMPTY_EXPRESSION
	}

	if ctx.LocalVariableDeclaration() != nil {
		declarationListStatement := v.VisitLocalVariableDeclaration(ctx.LocalVariableDeclaration().(*LocalVariableDeclarationContext))
		declarationExpressions := declarationListStatement.GetDeclarationExpressions()

		if len(declarationExpressions) == 1 {
			return configureAST(declarationExpressions[0], ctx)
		} else {
			expressions := make([]Expression, len(declarationExpressions))
			for i, de := range declarationExpressions {
				expressions[i] = de
			}
			return configureAST(NewClosureListExpression(expressions), ctx)
		}
	}

	if ctx.ExpressionList() != nil {
		return v.translateExpressionList(ctx.ExpressionList().(*ExpressionListContext))
	}

	panic(createParsingFailedException("Unsupported for init: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitForUpdate(ctx *ForUpdateContext) Expression {
	if ctx == nil {
		return INSTANCE
	}

	return v.translateExpressionList(ctx.ExpressionList().(*ExpressionListContext))
}

func (v *ASTBuilder) translateExpressionList(ctx *ExpressionListContext) Expression {
	expressionList := v.VisitExpressionList(ctx)

	if len(expressionList) == 1 {
		return configureAST(expressionList[0], ctx)
	} else {
		return configureAST(NewClosureListExpression(expressionList), ctx)
	}
}

func (v *ASTBuilder) VisitEnhancedForControl(ctx *EnhancedForControlContext) Tuple2[*Parameter, Expression] {
	parameter := NewParameter(v.VisitType(ctx.Type_().(*TypeContext)), v.VisitVariableDeclaratorId(ctx.VariableDeclaratorId().(*VariableDeclaratorIdContext)).GetName())
	modifierManager := NewModifierManager(v, v.VisitVariableModifiersOpt(ctx.VariableModifiersOpt().(*VariableModifiersOptContext)))
	modifierManager.ProcessParameter(parameter)
	configureAST(parameter, ctx.VariableDeclaratorId())
	return NewTuple2(parameter, v.Visit(ctx.Expression()).(Expression))
}

func (v *ASTBuilder) VisitClassicalForControl(ctx *ClassicalForControlContext) Tuple2[*Parameter, Expression] {
	closureListExpression := NewEmptyClosureListExpression()

	closureListExpression.AddExpression(v.VisitForInit(ctx.ForInit().(*ForInitContext)))
	if ctx.Expression() != nil {
		closureListExpression.AddExpression(v.Visit(ctx.Expression()).(Expression))
	} else {
		closureListExpression.AddExpression(EMPTY_EXPRESSION)
	}
	closureListExpression.AddExpression(v.VisitForUpdate(ctx.ForUpdate().(*ForUpdateContext)))

	var foo Expression = closureListExpression

	return NewTuple2(FOR_LOOP_DUMMY, foo)
}

func (v *ASTBuilder) VisitWhileStmtAlt(ctx *WhileStmtAltContext) *WhileStatement {
	conditionAndBlock := v.createLoopConditionExpressionAndBlock(ctx.ExpressionInPar().(*ExpressionInParContext), ctx.Statement().(*StatementContext))

	var block Statement
	if conditionAndBlock.V2 != nil {
		block = conditionAndBlock.V2
	} else {
		block = NewEmptyStatement()
	}

	return configureAST(
		NewWhileStatement(conditionAndBlock.V1, block),
		parserRuleContextAdapter{ctx},
	)
}

func (v *ASTBuilder) VisitDoWhileStmtAlt(ctx *DoWhileStmtAltContext) *DoWhileStatement {
	conditionAndBlock := v.createLoopConditionExpressionAndBlock(ctx.ExpressionInPar().(*ExpressionInParContext), ctx.Statement().(*StatementContext))

	var block Statement
	if conditionAndBlock.V2 != nil {
		block = conditionAndBlock.V2
	} else {
		block = NewEmptyStatement()
	}

	return configureAST(
		NewDoWhileStatement(conditionAndBlock.V1, block),
		ctx,
	)
}

func (v *ASTBuilder) createLoopConditionExpressionAndBlock(eipc *ExpressionInParContext, sc *StatementContext) Tuple2[*BooleanExpression, Statement] {
	conditionExpression := v.VisitExpressionInPar(eipc)

	booleanExpression := configureASTFromSource(
		NewBooleanExpression(conditionExpression),
		conditionExpression,
	)

	loopBlock := v.unpackStatement(v.Visit(sc).(Statement))

	return NewTuple2(booleanExpression, loopBlock)
}

func (v *ASTBuilder) VisitTryCatchStatement(ctx *TryCatchStatementContext) Statement {
	resourcesExists := ctx.Resources() != nil
	catchExists := len(ctx.AllCatchClause()) > 0
	finallyExists := ctx.FinallyBlock() != nil

	if !(resourcesExists || catchExists || finallyExists) {
		panic(createParsingFailedException("Either a catch or finally clause or both is required for a try-catch-finally statement", parserRuleContextAdapter{ctx}))
	}

	tryCatchStatement := NewTryCatchStatement(
		v.Visit(ctx.Block()).(Statement),
		v.VisitFinallyBlock(ctx.FinallyBlock().(*FinallyBlockContext)),
	)

	if resourcesExists {
		for _, resource := range v.VisitResources(ctx.Resources().(*ResourcesContext)) {
			tryCatchStatement.AddResource(resource)
		}
	}

	for _, catchClause := range ctx.AllCatchClause() {
		for _, catchStmt := range v.VisitCatchClause(catchClause.(*CatchClauseContext)) {
			tryCatchStatement.AddCatch(catchStmt)
		}
	}

	return configureAST(tryCatchStatement, ctx)
}

func (v *ASTBuilder) VisitResources(ctx *ResourcesContext) []*ExpressionStatement {
	return v.VisitResourceList(ctx.ResourceList().(*ResourceListContext))
}

func (v *ASTBuilder) VisitResourceList(ctx *ResourceListContext) []*ExpressionStatement {
	var resources []*ExpressionStatement
	for _, resource := range ctx.AllResource() {
		resources = append(resources, v.VisitResource(resource.(*ResourceContext)))
	}
	return resources
}

func IsInstanceOf(obj interface{}, targetType interface{}) bool {
	return reflect.TypeOf(obj) == reflect.TypeOf(targetType).Elem()
}

func (v *ASTBuilder) VisitResource(ctx *ResourceContext) *ExpressionStatement {
	if ctx.LocalVariableDeclaration() != nil {
		declarationStatements := v.VisitLocalVariableDeclaration(ctx.LocalVariableDeclaration().(*LocalVariableDeclarationContext)).GetDeclarationStatements()

		if len(declarationStatements) > 1 {
			panic(createParsingFailedException("Multi resources can not be declared in one statement", parserRuleContextAdapter{ctx}))
		}

		return declarationStatements[0]
	} else if ctx.Expression() != nil {
		expression := v.Visit(ctx.Expression()).(Expression)
		isVariableDeclaration := false
		isVariableAccess := false

		if binaryExpr, ok := expression.(*BinaryExpression); ok {
			isVariableDeclaration = binaryExpr.GetOperation().GetType() == ASSIGN &&
				IsInstanceOf(binaryExpr.GetLeftExpression(), (*VariableExpression)(nil))
		}
		isVariableAccess = IsInstanceOf(expression, (*VariableExpression)(nil))

		if !(isVariableDeclaration || isVariableAccess) {
			panic(createParsingFailedException("Only variable declarations or variable access are allowed to declare resource", parserRuleContextAdapter{ctx}))
		}

		var assignmentExpression *BinaryExpression

		if isVariableDeclaration {
			assignmentExpression = expression.(*BinaryExpression)
		} else if isVariableAccess {
			// TODO: transform
			assignmentExpression = expression.(*BinaryExpression)
		} else {
			panic(createParsingFailedException("Unsupported resource declaration", parserRuleContextAdapter{ctx}))
		}

		variableExpr := NewVariableExpressionWithString(assignmentExpression.GetLeftExpression().GetText())
		configuredVarExpr := configureASTFromSource(variableExpr, assignmentExpression.GetLeftExpression())

		declExpr := NewDeclarationExpression(
			configuredVarExpr,
			assignmentExpression.GetOperation(),
			assignmentExpression.GetRightExpression(),
		)
		configuredDeclExpr := configureAST(declExpr, ctx)

		stmt, err := NewExpressionStatement(configuredDeclExpr)
		if err != nil {
			panic(err)
		}

		return configureAST(stmt, ctx)
	}

	panic(createParsingFailedException("Unsupported resource declaration: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitCatchClause(ctx *CatchClauseContext) []*CatchStatement {
	catchTypes := v.VisitCatchType(ctx.CatchType().(*CatchTypeContext))
	catchStatements := make([]*CatchStatement, 0, len(catchTypes))

	for _, e := range catchTypes {
		catchStatement := NewCatchStatement(
			NewParameter(e, v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))),
			v.VisitBlock(ctx.Block().(*BlockContext)),
		)
		catchStatements = append(catchStatements, configureAST(catchStatement, ctx))
	}

	return catchStatements
}

func (v *ASTBuilder) VisitCatchType(ctx *CatchTypeContext) []*ClassNode {
	if ctx == nil {
		return []*ClassNode{OBJECT_TYPE}
	}

	classNodes := make([]*ClassNode, 0, len(ctx.AllQualifiedClassName()))
	for _, qcn := range ctx.AllQualifiedClassName() {
		classNodes = append(classNodes, v.VisitQualifiedClassName(qcn.(*QualifiedClassNameContext)))
	}
	return classNodes
}

func (v *ASTBuilder) VisitFinallyBlock(ctx *FinallyBlockContext) Statement {
	if ctx == nil {
		return NewEmptyStatement()
	}

	return configureAST(
		v.createBlockStatement(v.Visit(ctx.Block()).(Statement)),
		ctx,
	)
}

func (v *ASTBuilder) VisitSwitchStatement(ctx *SwitchStatementContext) *SwitchStatement {
	v.pushSwitchExpressionRuleContext(ctx)
	v.visitingSwitchStatementCount++
	defer func() {
		v.popSwitchExpressionRuleContext()
		v.visitingSwitchStatementCount--
	}()

	var statementList []Statement
	for _, group := range ctx.AllSwitchBlockStatementGroup() {
		statementList = append(statementList, v.VisitSwitchBlockStatementGroup(group.(*SwitchBlockStatementGroupContext))...)
	}

	var caseStatementList []*CaseStatement
	var defaultStatementList []Statement

	for _, e := range statementList {
		if caseStmt, ok := e.(*CaseStatement); ok {
			caseStatementList = append(caseStatementList, caseStmt)
		} else if isTrue(e, IS_SWITCH_DEFAULT) {
			defaultStatementList = append(defaultStatementList, e)
		}
	}

	defaultStatementListSize := len(defaultStatementList)
	if defaultStatementListSize > 1 {
		panic(createParsingFailedException("a switch must only have one default branch", astNodeAdapter{defaultStatementList[0]}))
	}

	if defaultStatementListSize > 0 && IsInstanceOf(statementList[len(statementList)-1], (*CaseStatement)(nil)) {
		panic(createParsingFailedException("a default branch must only appear as the last branch of a switch", astNodeAdapter{defaultStatementList[0]}))
	}

	var defaultStatement Statement
	if defaultStatementListSize == 0 {
		defaultStatement = NewEmptyStatement()
	} else {
		defaultStatement = defaultStatementList[0]
	}

	return configureAST(
		NewSwitchStatementFull(
			v.VisitExpressionInPar(ctx.ExpressionInPar().(*ExpressionInParContext)),
			caseStatementList,
			defaultStatement,
		),
		ctx,
	)
}

func (v *ASTBuilder) VisitSwitchBlockStatementGroup(ctx *SwitchBlockStatementGroupContext) []Statement {
	labelCount := len(ctx.AllSwitchLabel())
	var firstLabelHolder []antlr.Token

	var statementList []Statement
	for i, label := range ctx.AllSwitchLabel() {
		tuple := v.VisitSwitchLabel(label.(*SwitchLabelContext))
		switch tuple.V1.GetTokenType() {
		case GroovyParserCASE:
			if len(statementList) == 0 {
				firstLabelHolder = append(firstLabelHolder, tuple.V1)
			}
			var blockStatements Statement
			if i == labelCount-1 {
				blockStatements = v.VisitBlockStatements(ctx.BlockStatements().(*BlockStatementsContext))
			} else {
				blockStatements = NewEmptyStatement()
			}
			statement := NewCaseStatement(tuple.V2, blockStatements)
			statementList = append(statementList, configureASTWithToken(statement, firstLabelHolder[0]))
		case GroovyParserDEFAULT:
			statement := v.VisitBlockStatements(ctx.BlockStatements().(*BlockStatementsContext))
			statement.SetNodeMetaData(IS_SWITCH_DEFAULT, true)
			statementList = append(statementList, statement)
		}
	}

	return statementList
}

func (v *ASTBuilder) VisitSwitchLabel(ctx *SwitchLabelContext) Tuple2[antlr.Token, Expression] {
	if ctx.CASE() != nil {
		return NewTuple2(ctx.CASE().GetSymbol(), v.Visit(ctx.Expression()).(Expression))
	} else if ctx.DEFAULT() != nil {
		var foo Expression = EMPTY_EXPRESSION
		return NewTuple2(ctx.DEFAULT().GetSymbol(), foo)
	}

	panic(createParsingFailedException("Unsupported switch label: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitSynchronizedStmtAlt(ctx *SynchronizedStmtAltContext) *SynchronizedStatement {
	return configureAST(
		NewSynchronizedStatement(
			v.VisitExpressionInPar(ctx.ExpressionInPar().(*ExpressionInParContext)),
			v.VisitBlock(ctx.Block().(*BlockContext)),
		),
		ctx,
	)
}

func (v *ASTBuilder) VisitReturnStmtAlt(ctx *ReturnStmtAltContext) *ReturnStatement {
	if _, ok := v.peekSwitchExpressionRuleContext().(*SwitchExpressionContext); ok {
		panic(createParsingFailedException("switch expression does not support `return`", parserRuleContextAdapter{ctx}))
	}

	var expr Expression
	if ctx.Expression() != nil {
		expr = v.Visit(ctx.Expression()).(Expression)
	} else {
		expr = EMPTY_EXPRESSION
	}

	return configureAST(NewReturnStatement(expr), ctx)
}

func (v *ASTBuilder) VisitThrowStmtAlt(ctx *ThrowStmtAltContext) *ThrowStatement {
	return configureAST(
		NewThrowStatement(v.Visit(ctx.Expression()).(Expression)),
		ctx,
	)
}

func (v *ASTBuilder) VisitLabeledStmtAlt(ctx *LabeledStmtAltContext) Statement {
	statement := v.Visit(ctx.Statement()).(Statement)
	statement.AddStatementLabel(v.VisitIdentifier(ctx.Identifier().(*IdentifierContext)))
	return statement
}

func (v *ASTBuilder) VisitBreakStatement(ctx *BreakStatementContext) *BreakStatement {
	if v.visitingLoopStatementCount == 0 && v.visitingSwitchStatementCount == 0 {
		panic(createParsingFailedException("break statement is only allowed inside loops or switches", parserRuleContextAdapter{ctx}))
	}

	if _, ok := v.peekSwitchExpressionRuleContext().(*SwitchExpressionContext); ok {
		panic(createParsingFailedException("switch expression does not support `break`", parserRuleContextAdapter{ctx}))
	}

	var label string
	if ctx.Identifier() != nil {
		label = v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))
	}

	return configureAST(NewBreakStatement(label), ctx)
}

func (v *ASTBuilder) VisitYieldStatement(ctx *YieldStatementContext) *ReturnStatement {
	returnStatement := NewReturnStatement(v.Visit(ctx.Expression()).(Expression))
	returnStatement.SetNodeMetaData(IS_YIELD_STATEMENT, true)
	return configureAST(returnStatement, ctx)
}

func (v *ASTBuilder) VisitYieldStmtAlt(ctx *YieldStmtAltContext) *ReturnStatement {
	return configureAST(v.VisitYieldStatement(ctx.YieldStatement().(*YieldStatementContext)), ctx)
}

func (v *ASTBuilder) VisitContinueStatement(ctx *ContinueStatementContext) *ContinueStatement {
	if v.visitingLoopStatementCount == 0 {
		panic(createParsingFailedException("continue statement is only allowed inside loops", parserRuleContextAdapter{ctx}))
	}

	if _, ok := v.peekSwitchExpressionRuleContext().(*SwitchExpressionContext); ok {
		panic(createParsingFailedException("switch expression does not support `continue`", parserRuleContextAdapter{ctx}))
	}

	var label string
	if ctx.Identifier() != nil {
		label = v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))
	}

	return configureAST(NewContinueStatement(label), ctx)
}

func (v *ASTBuilder) VisitSwitchExprAlt(ctx *SwitchExprAltContext) Expression {
	return configureAST(v.VisitSwitchExpression(ctx.SwitchExpression().(*SwitchExpressionContext)), ctx)
}

func (v *ASTBuilder) createDeclarationStatement(target Expression, init Expression) *DeclarationStatement {
	operator := NewToken(ASSIGN, "=", -1, -1)
	declExpr := NewDeclarationExpression(target, operator, init)
	return NewDeclarationStatement(declExpr)
}

// LocalVarX creates a new VariableExpression for a local variable with the given name.
func LocalVarX(name string) *VariableExpression {
	result := NewVariableExpressionWithString(name)
	result.SetAccessedVariable(result)
	return result
}

func (v *ASTBuilder) VisitSwitchExpression(ctx *SwitchExpressionContext) *MethodCallExpression {
	v.pushSwitchExpressionRuleContext(ctx)
	defer v.popSwitchExpressionRuleContext()

	v.validateSwitchExpressionLabels(ctx)
	var statementInfoList []Tuple3[[]Statement, bool, bool]
	for _, e := range ctx.AllSwitchBlockStatementExpressionGroup() {
		statementInfoList = append(statementInfoList, v.VisitSwitchBlockStatementExpressionGroup(e.(*SwitchBlockStatementExpressionGroupContext)))
	}

	if len(statementInfoList) == 0 {
		panic(createParsingFailedException("`case` or `default` branches are expected", parserRuleContextAdapter{ctx}))
	}

	isArrow := statementInfoList[0].V2
	if !isArrow && !slices.ContainsFunc(statementInfoList, func(e Tuple3[[]Statement, bool, bool]) bool {
		return e.V3
	}) {
		panic(createParsingFailedException("`yield` or `throw` is expected", parserRuleContextAdapter{ctx}))
	}

	var statementList []Statement
	for _, e := range statementInfoList {
		statementList = append(statementList, e.V1...)
	}

	var caseStatementList []*CaseStatement
	var defaultStatementList []Statement

	for _, e := range statementList {
		if caseStmt, ok := e.(*CaseStatement); ok {
			caseStatementList = append(caseStatementList, caseStmt)
		} else if isTrue(e, IS_SWITCH_DEFAULT) {
			defaultStatementList = append(defaultStatementList, e)
		}
	}

	defaultStatementListSize := len(defaultStatementList)
	if defaultStatementListSize > 1 {
		panic(createParsingFailedException("switch expression should have only one default case, which should appear at last", astNodeAdapter{defaultStatementList[0]}))
	}

	if defaultStatementListSize > 0 && IsInstanceOf(statementList[len(statementList)-1], (*CaseStatement)(nil)) {
		panic(createParsingFailedException("default case should appear at last", astNodeAdapter{defaultStatementList[0]}))
	}

	variableName := fmt.Sprintf("__$$sev%d", v.switchExpressionVariableSeq)
	v.switchExpressionVariableSeq++
	declarationStatement := v.createDeclarationStatement(
		LocalVarX(variableName),
		v.VisitExpressionInPar(ctx.ExpressionInPar().(*ExpressionInParContext)),
	)

	var defaultStatement Statement
	if defaultStatementListSize == 0 {
		defaultStatement = NewEmptyStatement()
	} else {
		defaultStatement = defaultStatementList[0]
	}

	switchStatement := configureAST(
		NewSwitchStatementFull(
			NewVariableExpressionWithString(variableName),
			caseStatementList,
			defaultStatement,
		),
		ctx,
	)

	closureExpression := configureAST(
		NewClosureExpression(nil, v.createBlockStatement(declarationStatement, switchStatement)),
		ctx,
	)

	callClosure := v.createMethodCallExpression(closureExpression, NewConstantExpression(CALL_STR))
	callClosure.SetImplicitThis(false)

	return configureAST(callClosure, ctx)
}

func (v *ASTBuilder) VisitSwitchBlockStatementExpressionGroup(ctx *SwitchBlockStatementExpressionGroupContext) Tuple3[[]Statement, bool, bool] {
	labelCnt := len(ctx.AllSwitchExpressionLabel())
	var firstLabelHolder []antlr.Token
	arrowCntHolder := 0

	var isArrowHolder bool
	var hasResultStmtHolder bool
	var result []Statement

	for i, e := range ctx.AllSwitchExpressionLabel() {
		tuple := v.VisitSwitchExpressionLabel(e.(*SwitchExpressionLabelContext))

		isArrow := tuple.V3 == GroovyParserARROW
		isArrowHolder = isArrow
		if isArrow {
			arrowCntHolder++
			if arrowCntHolder > 1 && len(firstLabelHolder) > 0 {
				panic(createParsingFailedException("`case ... ->` does not support falling through cases", tokenAdapter{firstLabelHolder[0]}))
			}
		}

		isLast := labelCnt-1 == i

		codeBlock := v.VisitBlockStatements(ctx.BlockStatements().(*BlockStatementsContext))
		statements := codeBlock.GetStatements()
		statementsCnt := len(statements)
		if statementsCnt == 0 {
			panic(createParsingFailedException("`yield` is expected", parserRuleContextAdapter{ctx.BlockStatements()}))
		}

		if isArrow && statementsCnt > 1 {
			panic(createParsingFailedException(fmt.Sprintf("Expect only 1 statement, but %d statements found", statementsCnt), parserRuleContextAdapter{ctx.BlockStatements()}))
		}

		if !isArrow {
			var hasYield, hasThrow bool
			codeBlock.Visit(&CodeVisitorSupport{
				VisitReturnStatementFunc: func(statement *ReturnStatement) {
					if isTrue(statement, IS_YIELD_STATEMENT) {
						hasYield = true
					}
				},
				VisitThrowStatementFunc: func(statement *ThrowStatement) {
					hasThrow = true
				},
			})

			if hasYield || hasThrow {
				hasResultStmtHolder = true
			}
		}

		exprOrBlockStatement := statements[0]
		if blockStatement, ok := exprOrBlockStatement.(*BlockStatement); ok {
			branchStatementList := blockStatement.GetStatements()
			if len(branchStatementList) == 1 {
				exprOrBlockStatement = branchStatementList[0]
			}
		}

		if _, ok := exprOrBlockStatement.(*ReturnStatement); !ok {
			if _, ok := exprOrBlockStatement.(*ThrowStatement); !ok {
				if isArrow {
					callClosure := v.createMethodCallExpression(
						configureASTFromSource(
							NewClosureExpression(nil, exprOrBlockStatement),
							exprOrBlockStatement,
						),
						NewConstantExpression(CALL_STR),
					)
					callClosure.SetImplicitThis(false)
					var resultExpr Expression
					if exprStmt, ok := exprOrBlockStatement.(*ExpressionStatement); ok {
						resultExpr = exprStmt.GetExpression()
					} else {
						resultExpr = callClosure
					}

					codeBlock = configureASTFromSource(
						v.createBlockStatement(configureASTFromSource(
							NewReturnStatement(resultExpr),
							exprOrBlockStatement,
						)),
						exprOrBlockStatement,
					)
				}
			}
		}

		switch tuple.V1.GetTokenType() {
		case GroovyParserCASE:
			if len(result) == 0 {
				firstLabelHolder = append(firstLabelHolder, tuple.V1)
			}
			for i, expr := range tuple.V2 {
				var stmt Statement
				if isLast && i == len(tuple.V2)-1 {
					stmt = codeBlock
				} else {
					stmt = NewEmptyStatement()
				}
				result = append(result,
					configureASTWithToken(
						NewCaseStatement(
							expr,
							stmt,
						),
						firstLabelHolder[0],
					),
				)
			}
		case GroovyParserDEFAULT:
			codeBlock.SetNodeMetaData(IS_SWITCH_DEFAULT, true)
			result = append(result, codeBlock)
		}
	}

	return NewTuple3(result, isArrowHolder, hasResultStmtHolder)
}

func (v *ASTBuilder) validateSwitchExpressionLabels(ctx *SwitchExpressionContext) {
	acMap := make(map[string][]*SwitchExpressionLabelContext)
	for _, group := range ctx.AllSwitchBlockStatementExpressionGroup() {
		for _, label := range group.(*SwitchBlockStatementExpressionGroupContext).AllSwitchExpressionLabel() {
			acText := label.(*SwitchExpressionLabelContext).GetAc().GetText()
			acMap[acText] = append(acMap[acText], label.(*SwitchExpressionLabelContext))
		}
	}

	if len(acMap) > 1 {
		var lastSelcList []*SwitchExpressionLabelContext
		for _, list := range acMap {
			lastSelcList = list
		}

		var keys []string
		for k := range acMap {
			keys = append(keys, k)
		}
		errorMsg := "`" + strings.Join(keys, "` and `") + "` cannot be used together"
		panic(createParsingFailedException(errorMsg, tokenAdapter{lastSelcList[0].GetAc()}))
	}
}

func (v *ASTBuilder) VisitSwitchExpressionLabel(ctx *SwitchExpressionLabelContext) Tuple3[antlr.Token, []Expression, int] {
	acType := ctx.GetAc().GetTokenType()
	if ctx.CASE() != nil {
		return NewTuple3(
			ctx.CASE().GetSymbol(),
			v.VisitExpressionList(ctx.ExpressionList().(*ExpressionListContext)),
			acType,
		)
	} else if ctx.DEFAULT() != nil {
		return NewTuple3(
			ctx.DEFAULT().GetSymbol(),
			[]Expression{EMPTY_EXPRESSION},
			acType,
		)
	}

	panic(createParsingFailedException("Unsupported switch expression label: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitTypeDeclaration(ctx *TypeDeclarationContext) *ClassNode {
	if ctx.ClassDeclaration() != nil { // e.g. class A {}
		ctx.ClassDeclaration().(*ClassDeclarationContext).PutNodeMetaData(TYPE_DECLARATION_MODIFIERS, v.VisitClassOrInterfaceModifiersOpt(ctx.ClassOrInterfaceModifiersOpt().(*ClassOrInterfaceModifiersOptContext)))
		return configureAST(v.VisitClassDeclaration(ctx.ClassDeclaration().(*ClassDeclarationContext)), ctx)
	}

	panic(createParsingFailedException("Unsupported type declaration: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitClassDeclaration(ctx *ClassDeclarationContext) *ClassNode {
	packageName := v.moduleNode.GetPackageName()
	if packageName == "" {
		packageName = ""
	}
	className := v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))
	if className == "var" {
		panic(createParsingFailedException("var cannot be used for type declarations", parserRuleContextAdapter{ctx.Identifier()}))
	}

	isAnnotation := ctx.AT() != nil
	if isAnnotation {
		if ctx.TypeParameters() != nil {
			panic(createParsingFailedException("annotation declaration cannot have type parameters", parserRuleContextAdapter{ctx.TypeParameters()}))
		}

		if ctx.EXTENDS() != nil {
			panic(createParsingFailedException("No extends clause allowed for annotation declaration", tokenAdapter{ctx.EXTENDS().GetSymbol()}))
		}

		if ctx.IMPLEMENTS() != nil {
			panic(createParsingFailedException("No implements clause allowed for annotation declaration", tokenAdapter{ctx.IMPLEMENTS().GetSymbol()}))
		}
	}

	isEnum := ctx.ENUM() != nil
	if isEnum {
		if ctx.TypeParameters() != nil {
			panic(createParsingFailedException("enum declaration cannot have type parameters", parserRuleContextAdapter{ctx.TypeParameters()}))
		}

		if ctx.EXTENDS() != nil {
			panic(createParsingFailedException("No extends clause allowed for enum declaration", tokenAdapter{ctx.EXTENDS().GetSymbol()}))
		}
	}

	isInterface := (ctx.INTERFACE() != nil && !isAnnotation)
	if isInterface {
		if ctx.IMPLEMENTS() != nil {
			panic(createParsingFailedException("No implements clause allowed for interface declaration", tokenAdapter{ctx.IMPLEMENTS().GetSymbol()}))
		}
	}

	modifierManager := NewModifierManager(v, ctx.GetNodeMetaData(TYPE_DECLARATION_MODIFIERS).([]*ModifierNode))

	finalModifier := modifierManager.Get(FINAL)
	sealedModifier := modifierManager.Get(SEALED)
	nonSealedModifier := modifierManager.Get(NON_SEALED)
	isFinal := finalModifier != nil
	isSealed := sealedModifier != nil
	isNonSealed := nonSealedModifier != nil

	isRecord := ctx.RECORD() != nil
	hasRecordHeader := ctx.FormalParameters() != nil
	if isRecord {
		if !hasRecordHeader {
			panic(createParsingFailedException("header declaration of record is expected", parserRuleContextAdapter{ctx.Identifier()}))
		}
		if ctx.EXTENDS() != nil {
			panic(createParsingFailedException("No extends clause allowed for record declaration", tokenAdapter{ctx.EXTENDS().GetSymbol()}))
		}
		if isSealed {
			panic(createParsingFailedException("`sealed` is not allowed for record declaration", astNodeAdapter{sealedModifier}))
		}
		if isNonSealed {
			panic(createParsingFailedException("`non-sealed` is not allowed for record declaration", astNodeAdapter{nonSealedModifier}))
		}
	} else {
		if hasRecordHeader {
			panic(createParsingFailedException("header declaration is only allowed for record declaration", parserRuleContextAdapter{ctx.FormalParameters()}))
		}
	}

	if isSealed && isNonSealed {
		panic(createParsingFailedException("type cannot be defined with both `sealed` and `non-sealed`", astNodeAdapter{nonSealedModifier}))
	}

	if isFinal && (isSealed || isNonSealed) {
		sealedStr := "sealed"
		if isNonSealed {
			sealedStr = "non-sealed"
		}
		panic(createParsingFailedException("type cannot be defined with both `"+sealedStr+"` and `final`", astNodeAdapter{finalModifier}))
	}

	if (isAnnotation || isEnum) && (isSealed || isNonSealed) {
		var mn *ModifierNode
		if isSealed {
			mn = sealedModifier
		} else {
			mn = nonSealedModifier
		}
		typeStr := "enum"
		if isAnnotation {
			typeStr = "annotation definition"
		}
		panic(createParsingFailedException("modifier `"+mn.GetText()+"` is not allowed for "+typeStr, astNodeAdapter{mn}))
	}

	hasPermits := ctx.PERMITS() != nil
	if !isSealed && hasPermits {
		panic(createParsingFailedException("only sealed type declarations should have `permits` clause", parserRuleContextAdapter{ctx}))
	}

	modifiers := modifierManager.GetClassModifiersOpValue()

	syntheticPublic := ((modifiers & ACC_SYNTHETIC) != 0)
	modifiers &= ^ACC_SYNTHETIC

	var classNode *ClassNode
	outerClass := v.peekClassNode()

	if isEnum {
		className := className
		if outerClass != nil {
			className = outerClass.GetName() + "$" + className
		} else {
			className = packageName + className
		}
		classNode = MakeEnumNode(
			className,
			modifiers,
			nil,
			outerClass,
		)
	} else if outerClass != nil {
		if outerClass.IsInterface() {
			modifiers |= ACC_STATIC
		}
		innerClassNode := NewInnerClassNode(
			outerClass,
			outerClass.GetName()+"$"+className,
			modifiers,
			OBJECT_TYPE.GetPlainNodeReference(),
		)
		classNode = innerClassNode.ClassNode
	} else {
		classNode = NewClassNode(
			packageName+className,
			modifiers,
			OBJECT_TYPE.GetPlainNodeReference(),
		)
	}

	configureAST(classNode, ctx)
	classNode.SetSyntheticPublic(syntheticPublic)
	classNode.SetGenericsTypes(v.VisitTypeParameters(ctx.TypeParameters().(*TypeParametersContext)))
	isInterfaceWithDefaultMethods := (isInterface && v.containsDefaultOrPrivateMethods(ctx))
	// TODO: handle this
	/*
		if isSealed {
			sealedAnnotationNode := makeAnnotationNode(Sealed)
			if ctx.ps != nil {
				permittedSubclassesListExpression := NewListExpression(v.VisitTypeList(ctx.Ps.(*TypeListContext)))
				sealedAnnotationNode.SetMember("permittedSubclasses", permittedSubclassesListExpression)
				configureAST(sealedAnnotationNode, ctx.PERMITS())
				sealedAnnotationNode.PutNodeMetaData("permits", true)
			}
			classNode.AddAnnotation(sealedAnnotationNode)
		} else if isNonSealed {
			classNode.AddAnnotation(makeAnnotationNode(NonSealed))
		}
		if ctx.TRAIT() != nil {
			classNode.AddAnnotation(makeAnnotationNode(Trait))
		}
	*/
	classNode.AddAnnotations(modifierManager.GetAnnotations())
	if isRecord && !slices.ContainsFunc(classNode.GetAnnotations(), func(a AnnotationNode) bool {
		return a.GetClassNode().GetName() == RECORD_TYPE_NAME
	}) {
		classNode.AddAnnotationNode(NewAnnotationNode(MakeWithoutCaching(RECORD_TYPE_NAME))) // TODO: makeAnnotationNode(RecordType)
	}

	if isInterfaceWithDefaultMethods {
		classNode.PutNodeMetaData(IS_INTERFACE_WITH_DEFAULT_METHODS, true)
	}
	classNode.PutNodeMetaData(CLASS_NAME, className)

	if ctx.CLASS() != nil || ctx.TRAIT() != nil {
		if ctx.scs != nil {
			scs := v.VisitTypeList(ctx.scs.(*TypeListContext))
			if len(scs) > 1 {
				panic(createParsingFailedException("Cannot extend multiple classes", tokenAdapter{ctx.EXTENDS().GetSymbol()}))
			}
			classNode.SetSuperClass(scs[0])
		}
		classNode.SetInterfaces(v.VisitTypeList(ctx.is.(*TypeListContext)))
		v.checkUsingGenerics(classNode)

	} else if isInterface {
		classNode.SetModifiers(classNode.GetModifiers() | ACC_INTERFACE | ACC_ABSTRACT)
		classNode.SetInterfaces(v.VisitTypeList(ctx.scs.(*TypeListContext)))
		v.checkUsingGenerics(classNode)
		v.hackMixins(classNode)

	} else if isEnum || isRecord {
		classNode.SetInterfaces(v.VisitTypeList(ctx.is.(*TypeListContext)))
		v.checkUsingGenerics(classNode)
		if isRecord {
			v.transformRecordHeaderToProperties(ctx, classNode)
		}

	} else if isAnnotation {
		classNode.SetModifiers(classNode.GetModifiers() | ACC_INTERFACE | ACC_ABSTRACT | ACC_ANNOTATION)
		classNode.AddInterface(ANNOTATION_NODE_TYPE)
		v.hackMixins(classNode)

	} else {
		panic(createParsingFailedException("Unsupported class declaration: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}

	v.pushClassNode(classNode)
	ctx.ClassBody().(*ClassBodyContext).PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
	v.VisitClassBody(ctx.ClassBody().(*ClassBodyContext))
	if isRecord {
		for _, field := range classNode.GetFields() {
			if !isTrue(field, IS_RECORD_GENERATED) && !field.IsStatic() {
				panic(createParsingFailedException("Instance field is not allowed in `record`", astNodeAdapter{field}))
			}
		}
	}
	v.popClassNode()

	// The first element in classNodeList determines what GCL#parseClass for
	// example will return. So we have to ensure it won't be an inner class.
	if outerClass == nil {
		v.addToClassNodeList(classNode)
	}
	//v.groovydocManager.Handle(classNode, ctx)

	return classNode
}

func (v *ASTBuilder) addToClassNodeList(classNode *ClassNode) {
	v.classNodeList = append(v.classNodeList, classNode) // GROOVY-11117: outer class first
	for _, innerClass := range classNode.innerClasses {
		v.addToClassNodeList(innerClass.ClassNode)
	}
}

func (v *ASTBuilder) checkUsingGenerics(classNode *ClassNode) {
	if !classNode.IsUsingGenerics() {
		if !classNode.IsEnum() && classNode.GetSuperClass().IsUsingGenerics() {
			classNode.SetUsingGenerics(true)
		} else if classNode.GetInterfaces() != nil {
			for _, interfaceNode := range classNode.GetInterfaces() {
				if interfaceNode.IsUsingGenerics() {
					classNode.SetUsingGenerics(true)
					break
				}
			}
		}
	}
}

func (v *ASTBuilder) transformRecordHeaderToProperties(ctx *ClassDeclarationContext, classNode *ClassNode) {
	parameters := v.VisitFormalParameters(ctx.FormalParameters().(*FormalParametersContext))
	classNode.PutNodeMetaData(RECORD_HEADER, parameters)

	for i, parameter := range parameters {
		parameterCtx := parameter.GetNodeMetaData(PARAMETER_CONTEXT).(*FormalParameterContext)
		parameterModifierManager := parameter.GetNodeMetaData(PARAMETER_MODIFIER_MANAGER).(*ModifierManager)
		propertyNode := v.declareProperty(parameterCtx.GroovyParserRuleContext, parameterModifierManager, parameter.GetType(), classNode, i,
			parameter, parameter.GetName(), parameter.GetModifiers()|ACC_FINAL, parameter.GetInitialExpression())
		propertyNode.GetField().PutNodeMetaData(IS_RECORD_GENERATED, true)
	}
}

func (v *ASTBuilder) containsDefaultOrPrivateMethods(ctx *ClassDeclarationContext) bool {
	var methodDeclarationContextList []*MethodDeclarationContext

	for _, bodyDecl := range ctx.ClassBody().AllClassBodyDeclaration() {
		memberDecl := bodyDecl.(*ClassBodyDeclarationContext).MemberDeclaration()
		if memberDecl == nil {
			continue
		}

		methodDecl := memberDecl.MethodDeclaration()
		if methodDecl == nil {
			continue
		}

		modifierManager := v.createModifierManager(methodDecl.(*MethodDeclarationContext))
		if modifierManager.ContainsAny(GroovyParserDEFAULT, GroovyParserPRIVATE) {
			methodDeclarationContextList = append(methodDeclarationContextList, methodDecl.(*MethodDeclarationContext))
		}
	}

	return len(methodDeclarationContextList) > 0
}

func (v *ASTBuilder) VisitClassBody(ctx *ClassBodyContext) interface{} {
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)
	if classNode == nil {
		panic("classNode should not be nil")
	}

	if ctx.EnumConstants() != nil {
		constants := ctx.EnumConstants().(*EnumConstantsContext)
		constants.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		v.VisitEnumConstants(ctx.EnumConstants().(*EnumConstantsContext))
	}

	for _, e := range ctx.AllClassBodyDeclaration() {
		foo := e.(*ClassBodyDeclarationContext)
		foo.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		v.VisitClassBodyDeclaration(e.(*ClassBodyDeclarationContext))
	}

	return nil
}

func (v *ASTBuilder) VisitEnumConstants(ctx *EnumConstantsContext) []*FieldNode {
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)
	if classNode == nil {
		panic("classNode should not be nil")
	}

	var fieldNodes []*FieldNode
	for _, e := range ctx.AllEnumConstant() {
		foo := e.(*EnumConstantContext)
		foo.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		fieldNodes = append(fieldNodes, v.VisitEnumConstant(e.(*EnumConstantContext)))
	}
	return fieldNodes
}

func (v *ASTBuilder) VisitEnumConstant(ctx *EnumConstantContext) *FieldNode {
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)
	if classNode == nil {
		panic("classNode should not be nil")
	}

	var anonymousInnerClassNode *InnerClassNode
	if ctx.AnonymousInnerClassDeclaration() != nil {
		foo := ctx.AnonymousInnerClassDeclaration().(*AnonymousInnerClassDeclarationContext)
		foo.PutNodeMetaData(ANONYMOUS_INNER_CLASS_SUPER_CLASS, classNode)
		anonymousInnerClassNode = v.VisitAnonymousInnerClassDeclaration(ctx.AnonymousInnerClassDeclaration().(*AnonymousInnerClassDeclarationContext))
	}

	enumConstant := AddEnumConstant(
		classNode,
		v.VisitIdentifier(ctx.Identifier().(*IdentifierContext)),
		v.createEnumConstantInitExpression(ctx.Arguments().(*ArgumentsContext), anonymousInnerClassNode),
	)

	enumConstant.AddAnnotations(v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext)))

	//v.groovydocManager.Handle(enumConstant, ctx)

	return configureAST(enumConstant, ctx)
}

func (v *ASTBuilder) createEnumConstantInitExpression(ctx *ArgumentsContext, anonymousInnerClassNode *InnerClassNode) Expression {
	if ctx == nil && anonymousInnerClassNode == nil {
		return nil
	}

	argumentListExpression := v.VisitArguments(ctx).(*TupleExpression)
	expressions := argumentListExpression.GetExpressions()

	if len(expressions) == 1 {
		expression := expressions[0]

		if namedArgListExpr, ok := expression.(*NamedArgumentListExpression); ok {
			mapEntryExpressionList := namedArgListExpr.GetMapEntryExpressions()
			expressions := make([]Expression, len(mapEntryExpressionList))
			for i, e := range mapEntryExpressionList {
				expressions[i] = e
			}
			listExpression := NewListExpressionWithExpressions(expressions)

			if anonymousInnerClassNode != nil {
				listExpression.AddExpression(
					configureASTFromSource(
						NewClassExpression(anonymousInnerClassNode.ClassNode),
						anonymousInnerClassNode,
					),
				)
			}

			if len(mapEntryExpressionList) > 1 {
				listExpression.SetWrapped(true)
			}

			return configureAST(listExpression, ctx)
		}

		if anonymousInnerClassNode == nil {
			if listExpr, ok := expression.(*ListExpression); ok {
				newListExpression := NewListExpression()
				newListExpression.AddExpression(listExpr)
				return configureAST(newListExpression, ctx)
			}
			return expression
		}

		listExpression := NewListExpression()

		if listExpr, ok := expression.(*ListExpression); ok {
			for _, expr := range listExpr.GetExpressions() {
				listExpression.AddExpression(expr)
			}
		} else {
			listExpression.AddExpression(expression)
		}

		listExpression.AddExpression(
			configureASTFromSource(
				NewClassExpression(anonymousInnerClassNode.ClassNode),
				anonymousInnerClassNode,
			),
		)

		return configureAST(listExpression, ctx)
	}

	listExpression := NewListExpressionWithExpressions(expressions)
	if anonymousInnerClassNode != nil {
		listExpression.AddExpression(
			configureASTFromSource(
				NewClassExpression(anonymousInnerClassNode.ClassNode),
				anonymousInnerClassNode,
			),
		)
	}

	if ctx != nil {
		listExpression.SetWrapped(true)
	}

	if ctx != nil {
		return configureAST(listExpression, ctx)
	}
	return configureASTFromSource(listExpression, anonymousInnerClassNode)
}

func (v *ASTBuilder) VisitClassBodyDeclaration(ctx *ClassBodyDeclarationContext) {
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)
	if ctx.MemberDeclaration() != nil {
		methodDecl := ctx.MemberDeclaration().(*MemberDeclarationContext)
		methodDecl.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		v.VisitMemberDeclaration(methodDecl)
	} else if ctx.Block() != nil {
		statement := v.VisitBlock(ctx.Block().(*BlockContext))
		if ctx.STATIC() != nil { // e.g. static { }
			classNode.AddStaticInitializerStatements([]Statement{statement}, false)
		} else { // e.g. { }
			classNode.AddObjectInitializerStatements(configureASTFromSource(v.createBlockStatement(statement), statement))
		}
	}
}

func (v *ASTBuilder) VisitMemberDeclaration(ctx *MemberDeclarationContext) {
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)
	if classNode == nil {
		panic("classNode should not be nil")
	}

	if ctx.MethodDeclaration() != nil {
		methodDecl := ctx.MethodDeclaration().(*MethodDeclarationContext)
		methodDecl.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		v.VisitMethodDeclaration(methodDecl)
	} else if ctx.FieldDeclaration() != nil {
		fieldDecl := ctx.FieldDeclaration().(*FieldDeclarationContext)
		fieldDecl.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		v.VisitFieldDeclaration(fieldDecl)
	} else if ctx.CompactConstructorDeclaration() != nil {
		compactConstructorDecl := ctx.CompactConstructorDeclaration().(*CompactConstructorDeclarationContext)
		compactConstructorDecl.PutNodeMetaData(COMPACT_CONSTRUCTOR_DECLARATION_MODIFIERS, v.VisitModifiersOpt(ctx.ModifiersOpt().(*ModifiersOptContext)))
		compactConstructorDecl.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		v.VisitCompactConstructorDeclaration(compactConstructorDecl)
	} else if ctx.ClassDeclaration() != nil {
		classDecl := ctx.ClassDeclaration().(*ClassDeclarationContext)
		classDecl.PutNodeMetaData(TYPE_DECLARATION_MODIFIERS, v.VisitModifiersOpt(ctx.ModifiersOpt().(*ModifiersOptContext)))
		classDecl.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
		v.VisitClassDeclaration(classDecl)
	}
}

func (v *ASTBuilder) VisitTypeParameters(ctx *TypeParametersContext) []*GenericsType {
	if ctx == nil {
		return nil
	}

	typeParameters := make([]*GenericsType, len(ctx.AllTypeParameter()))
	for i, tp := range ctx.AllTypeParameter() {
		typeParameters[i] = v.VisitTypeParameter(tp.(*TypeParameterContext))
	}
	return typeParameters
}

func (v *ASTBuilder) VisitTypeParameter(ctx *TypeParameterContext) *GenericsType {
	baseType := configureAST(MakeWithoutCaching(v.VisitClassName(ctx.ClassName().(*ClassNameContext))), ctx)
	baseType.AddTypeAnnotations(v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext)))
	genericsType := NewGenericsType(baseType, v.VisitTypeBound(ctx.TypeBound().(*TypeBoundContext)), nil)
	return configureAST(genericsType, ctx)
}

func (v *ASTBuilder) VisitTypeBound(ctx *TypeBoundContext) []*ClassNode {
	if ctx == nil {
		return nil
	}

	typeBounds := make([]*ClassNode, len(ctx.AllType_()))
	for i, t := range ctx.AllType_() {
		typeBounds[i] = v.VisitType(t.(*TypeContext))
	}
	return typeBounds
}

func (v *ASTBuilder) VisitFieldDeclaration(ctx *FieldDeclarationContext) {
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)
	if classNode == nil {
		panic("classNode should not be nil")
	}

	declaration := ctx.VariableDeclaration().(*VariableDeclarationContext)

	declaration.SetNodeMetaData(CLASS_DECLARATION_CLASS_NODE, classNode)
	v.VisitVariableDeclaration(ctx.VariableDeclaration().(*VariableDeclarationContext))
}

func (v *ASTBuilder) checkThisAndSuperConstructorCall(statement Statement) *ConstructorCallExpression {
	blockStatement, ok := statement.(*BlockStatement)
	if !ok { // method code must be a BlockStatement
		return nil
	}

	statementList := blockStatement.GetStatements()

	for i, s := range statementList {
		if exprStmt, ok := s.(*ExpressionStatement); ok {
			if constructorCall, ok := exprStmt.GetExpression().(*ConstructorCallExpression); ok && i != 0 {
				return constructorCall
			}
		}
	}

	return nil
}

func (v *ASTBuilder) createModifierManager(ctx *MethodDeclarationContext) *ModifierManager {
	var modifierNodeList []*ModifierNode

	if ctx.ModifiersOpt() != nil {
		modifierNodeList = v.VisitModifiersOpt(ctx.ModifiersOpt().(*ModifiersOptContext))
	}

	return NewModifierManager(v, modifierNodeList)
}

func (v *ASTBuilder) validateParametersOfMethodDeclaration(parameters []*Parameter, classNode *ClassNode) {
	if !classNode.IsInterface() {
		return
	}

	for _, parameter := range parameters {
		if parameter.HasInitialExpression() {
			panic(createParsingFailedException("Cannot specify default value for method parameter '"+parameter.GetName()+" = "+parameter.GetInitialExpression().GetText()+"' inside an interface", astNodeAdapter{parameter}))
		}
	}
}

func (v *ASTBuilder) VisitCompactConstructorDeclaration(ctx *CompactConstructorDeclarationContext) *MethodNode {
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)

	if !slices.ContainsFunc(classNode.GetAnnotations(), func(a AnnotationNode) bool {
		return a.GetClassNode().GetName() == RECORD_TYPE_NAME
	}) {
		panic(createParsingFailedException("Only record can have compact constructor", parserRuleContextAdapter{ctx}))
	}

	modifierManager := NewModifierManager(v, ctx.GetNodeMetaData(COMPACT_CONSTRUCTOR_DECLARATION_MODIFIERS).([]*ModifierNode))
	if modifierManager.ContainsAny(GroovyParserVAR) {
		panic(createParsingFailedException("var cannot be used for compact constructor declaration", parserRuleContextAdapter{ctx}))
	}

	methodName := v.VisitMethodName(ctx.MethodName().(*MethodNameContext))
	className := classNode.GetNodeMetaData(CLASS_NAME).(string)
	if methodName != className {
		panic(createParsingFailedException("Compact constructor should have the same name as record: "+className, parserRuleContextAdapter{ctx.MethodName()}))
	}

	header := classNode.GetNodeMetaData(RECORD_HEADER).([]*Parameter)
	code := v.VisitMethodBody(ctx.MethodBody().(*MethodBodyContext))
	code.Visit(&CodeVisitorSupport{
		VisitPropertyExpressionFunc: func(expression *PropertyExpression) {
			receiverText := expression.GetObjectExpression().GetText()
			propertyName := expression.GetPropertyAsString()
			if receiverText == THIS_STR && slices.ContainsFunc(header, func(p *Parameter) bool {
				return p.GetName() == propertyName
			}) {
				panic(createParsingFailedException("Cannot assign a value to final variable '"+propertyName+"'", astNodeAdapter{expression.GetProperty()}))
			}
		},
	})

	annos := classNode.GetAnnotationsOfType(TUPLE_TYPE)
	var tupleConstructor *AnnotationNode
	if len(annos) == 0 {
		tupleConstructor = makeAnnotationNode(reflect.TypeOf("tupletype"))
	} else {
		tupleConstructor = &annos[0]
	}
	tupleConstructor.SetMember("pre", NewClosureExpression(nil, code))
	if len(annos) == 0 {
		classNode.AddAnnotation(TUPLE_TYPE)
	}

	return nil
}

func (v *ASTBuilder) VisitMethodDeclaration(ctx *MethodDeclarationContext) *MethodNode {
	modifierManager := v.createModifierManager(ctx)

	if modifierManager.ContainsAny(GroovyParserVAR) {
		panic(createParsingFailedException("var cannot be used for method declarations", parserRuleContextAdapter{ctx}))
	}

	methodName := v.VisitMethodName(ctx.MethodName().(*MethodNameContext))
	returnType := v.VisitReturnType(ctx.ReturnType().(*ReturnTypeContext))
	parameters := v.VisitFormalParameters(ctx.FormalParameters().(*FormalParametersContext))
	exceptions := v.VisitQualifiedClassNameList(ctx.QualifiedClassNameList().(*QualifiedClassNameListContext))

	v.pushAnonymousInnerClass(list.New())
	code := v.VisitMethodBody(ctx.MethodBody().(*MethodBodyContext))
	anonymousInnerClassList := v.popAnonymousInnerClass()

	var methodNode *MethodNode
	// if classNode is not null, the method declaration is for class declaration
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)
	if classNode != nil {
		v.validateParametersOfMethodDeclaration(parameters, classNode)

		methodNode = v.createConstructorOrMethodNodeForClass(ctx, modifierManager, methodName, returnType, parameters, exceptions, code, classNode)
	} else { // script method declaration
		methodNode = v.createScriptMethodNode(modifierManager, methodName, returnType, parameters, exceptions, code)
	}

	for e := anonymousInnerClassList.Front(); e != nil; e = e.Next() {
		e.Value.(*InnerClassNode).SetEnclosingMethod(methodNode)
	}

	methodNode.SetGenericsTypes(v.VisitTypeParameters(ctx.TypeParameters().(*TypeParametersContext)))
	methodNode.SetSyntheticPublic(
		v.isSyntheticPublic(
			v.isAnnotationDeclaration(classNode),
			IsInstanceOf(classNode, (*EnumConstantClassNode)(nil)),
			ctx.ReturnType() != nil,
			modifierManager))

	if modifierManager.ContainsAny(STATIC) {
		for _, parameter := range methodNode.GetParameters() {
			parameter.SetInStaticContext(true)
		}

		methodNode.GetVariableScope().SetInStaticContext(true)
	}

	configureAST(methodNode, ctx)

	v.validateMethodDeclaration(ctx, methodNode, modifierManager, classNode)

	//v.groovydocManager.Handle(methodNode, ctx)

	return methodNode
}

func (v *ASTBuilder) validateMethodDeclaration(ctx *MethodDeclarationContext, methodNode *MethodNode, modifierManager *ModifierManager, classNode *ClassNode) {
	if ctx.t == 1 || ctx.t == 2 || ctx.t == 3 { // 1: normal method declaration; 2: abstract method declaration; 3: normal method declaration OR abstract method declaration
		if !(ctx.ModifiersOpt().Modifiers() != nil || ctx.ReturnType() != nil) {
			panic(createParsingFailedException("Modifiers or return type is required", parserRuleContextAdapter{ctx}))
		}
	}

	if ctx.t == 1 {
		if ctx.MethodBody() == nil {
			panic(createParsingFailedException("Method body is required", parserRuleContextAdapter{ctx}))
		}
	}

	if ctx.t == 2 {
		if ctx.MethodBody() != nil {
			panic(createParsingFailedException("Abstract method should not have method body", parserRuleContextAdapter{ctx}))
		}
	}

	isAbstractMethod := methodNode.IsAbstract()
	hasMethodBody := methodNode.code != nil && !IsInstanceOf(methodNode.code, (*ExpressionStatement)(nil))

	if ctx.ct == 9 { // script
		if isAbstractMethod || !hasMethodBody { // method should not be declared abstract in the script
			msg := fmt.Sprintf("You cannot define %s method[%s] %sin the script. Try %s%s%s",
				ternary(isAbstractMethod, "an abstract", "a"),
				methodNode.name,
				ternary(!hasMethodBody, "without method body ", ""),
				ternary(isAbstractMethod, "removing the 'abstract'", ""),
				ternary(isAbstractMethod && !hasMethodBody, " and", ""),
				ternary(!hasMethodBody, " adding a method body", ""))
			panic(createParsingFailedException(msg, astNodeAdapter{methodNode}))
		}
	} else {
		if ctx.ct == 4 { // trait
			if isAbstractMethod && hasMethodBody {
				panic(createParsingFailedException("Abstract method should not have method body", parserRuleContextAdapter{ctx}))
			}
		}

		if ctx.ct == 3 { // annotation
			if hasMethodBody {
				panic(createParsingFailedException("Annotation type element should not have body", parserRuleContextAdapter{ctx}))
			}
		}

		if !isAbstractMethod && !hasMethodBody { // non-abstract method without body in the non-script(e.g. class, enum, trait) is not allowed!
			panic(createParsingFailedException(
				fmt.Sprintf("You defined a method[%s] without a body. "+
					"Try adding a method body, or declare it abstract",
					methodNode.GetName()),
				astNodeAdapter{methodNode},
			))
		}

		isInterfaceOrAbstractClass := classNode != nil && classNode.IsAbstract() && !classNode.IsAnnotationDefinition()
		if isInterfaceOrAbstractClass && !modifierManager.ContainsAny(GroovyParserDEFAULT, GroovyParserPRIVATE) && isAbstractMethod && hasMethodBody {
			msg := fmt.Sprintf("You defined an abstract method[%s] with a body. Try removing the method body%s",
				methodNode.GetName(),
				ternary(classNode.IsInterface(), ", or declare it default or private", ""))
			panic(createParsingFailedException(msg, astNodeAdapter{methodNode}))
		}
	}

	modifierManager.Validate(methodNode)

	// TODO: add this
	/*
		if IsInstanceOf(methodNode, (*ConstructorNode)(nil)) {
			modifierManager.ValidateConstructor(methodNode.(*ConstructorNode))
		}
	*/
}

// Helper function for ternary operation
func ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func (v *ASTBuilder) createScriptMethodNode(modifierManager *ModifierManager, methodName string, returnType *ClassNode, parameters []*Parameter, exceptions []*ClassNode, code Statement) *MethodNode {
	var modifiers int
	if modifierManager.ContainsAny(GroovyParserPRIVATE) {
		modifiers = ACC_PRIVATE
	} else {
		modifiers = ACC_PUBLIC
	}

	methodNode := NewMethodNode(
		methodName,
		modifiers,
		returnType,
		parameters,
		exceptions,
		code,
	)
	modifierManager.ProcessMethodNode(methodNode)
	return methodNode
}

func (v *ASTBuilder) createConstructorOrMethodNodeForClass(ctx *MethodDeclarationContext, modifierManager *ModifierManager, methodName string, returnType *ClassNode, parameters []*Parameter, exceptions []*ClassNode, code Statement, classNode *ClassNode) *MethodNode {
	className := classNode.GetNodeMetaData(CLASS_NAME).(string)
	modifiers := modifierManager.GetClassMemberModifiersOpValue()

	hasReturnType := ctx.ReturnType() != nil
	hasMethodBody := ctx.MethodBody() != nil

	if !hasReturnType && hasMethodBody && methodName == className {
		return v.createConstructorNodeForClass(methodName, parameters, exceptions, code, classNode, modifiers).MethodNode
	} else {
		if !hasReturnType && hasMethodBody && modifierManager.GetModifierCount() == 0 {
			panic(createParsingFailedException("Invalid method declaration: "+methodName, parserRuleContextAdapter{ctx}))
		}
		return v.createMethodNodeForClass(ctx, modifierManager, methodName, returnType, parameters, exceptions, code, classNode, modifiers)
	}
}

func (v *ASTBuilder) createMethodNodeForClass(ctx *MethodDeclarationContext, modifierManager *ModifierManager, methodName string, returnType *ClassNode, parameters []*Parameter, exceptions []*ClassNode, code Statement, classNode *ClassNode, modifiers int) *MethodNode {
	if ctx.ElementValue() != nil { // the code of annotation method
		exprStmt, err := NewExpressionStatement(v.VisitElementValue(ctx.ElementValue().(*ElementValueContext)))
		if err != nil {
			panic(createParsingFailedException("Failed to create expression statement: "+err.Error(), parserRuleContextAdapter{ctx.ElementValue()}))
		}
		code = configureAST(exprStmt, ctx)
	}

	if !modifierManager.ContainsAny(STATIC) && classNode.IsInterface() && !(isTrue(classNode, IS_INTERFACE_WITH_DEFAULT_METHODS) && modifierManager.ContainsAny(GroovyParserDEFAULT, GroovyParserPRIVATE)) {
		modifiers |= ACC_ABSTRACT
	}
	methodNode := NewMethodNode(methodName, modifiers, returnType, parameters, exceptions, code)
	classNode.AddMethod(methodNode)

	methodNode.SetAnnotationDefault(ctx.ElementValue() != nil)
	return methodNode
}

func (v *ASTBuilder) createConstructorNodeForClass(methodName string, parameters []*Parameter, exceptions []*ClassNode, code Statement, classNode *ClassNode, modifiers int) *ConstructorNode {
	thisOrSuperConstructorCallExpression := v.checkThisAndSuperConstructorCall(code)
	if thisOrSuperConstructorCallExpression != nil {
		panic(createParsingFailedException(thisOrSuperConstructorCallExpression.GetText()+" should be the first statement in the constructor["+methodName+"]", astNodeAdapter{thisOrSuperConstructorCallExpression}))
	}

	return classNode.AddConstructorWithDetails(
		modifiers,
		parameters,
		exceptions,
		code,
	)
}

func (v *ASTBuilder) VisitMethodName(ctx *MethodNameContext) string {
	if ctx.Identifier() != nil {
		return v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))
	}

	if ctx.StringLiteral() != nil {
		return v.VisitStringLiteral(ctx.StringLiteral().(*StringLiteralContext)).GetText()
	}

	panic(createParsingFailedException("Unsupported method name: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitReturnType(ctx *ReturnTypeContext) *ClassNode {
	if ctx == nil {
		return dynamicType()
	}

	// TODO: handle this
	/*
		if ctx.StandardType() != nil {
			return v.VisitType(ctx.StandardType().(*TypeContext))
		}
	*/

	if ctx.VOID() != nil {
		if ctx.ct == 3 { // annotation
			panic(createParsingFailedException("annotation method cannot have void return type", parserRuleContextAdapter{ctx}))
		}

		return configureASTWithToken(VOID_TYPE.GetPlainNodeReference(), ctx.VOID().GetSymbol())
	}

	panic(createParsingFailedException("Unsupported return type: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitMethodBody(ctx *MethodBodyContext) Statement {
	if ctx == nil {
		return nil
	}

	return configureAST(v.VisitBlock(ctx.Block().(*BlockContext)), ctx)
}

func (v *ASTBuilder) VisitLocalVariableDeclaration(ctx *LocalVariableDeclarationContext) *DeclarationListStatement {
	return configureAST(v.VisitVariableDeclaration(ctx.VariableDeclaration().(*VariableDeclarationContext)), ctx)
}

func (v *ASTBuilder) createMultiAssignmentDeclarationListStatement(ctx *VariableDeclarationContext, modifierManager *ModifierManager) *DeclarationListStatement {
	elist := v.VisitTypeNamePairs(ctx.TypeNamePairs().(*TypeNamePairsContext))
	for _, e := range elist {
		modifierManager.ProcessVariableExpression(e.(*VariableExpression))
	}

	de := NewDeclarationExpression(
		configureAST(NewTupleExpressionWithExpressions(elist...), ctx.TypeNamePairs()),
		v.createGroovyTokenByType(ctx.ASSIGN().GetSymbol(), GroovyParserASSIGN),
		v.VisitVariableInitializer(ctx.VariableInitializer().(*VariableInitializerContext)),
	)

	configureAST(modifierManager.AttachAnnotations(de.AnnotatedNode), ctx)
	return configureAST(NewDeclarationListStatement(de), ctx)
}

func (v *ASTBuilder) VisitVariableDeclaration(ctx *VariableDeclarationContext) *DeclarationListStatement {
	var modifierManager *ModifierManager
	if ctx.Modifiers() != nil {
		modifierManager = NewModifierManager(v, v.VisitModifiers(ctx.Modifiers().(*ModifiersContext)))
	} else {
		modifierManager = NewModifierManager(v, []*ModifierNode{})
	}

	if ctx.TypeNamePairs() != nil { // e.g. def (int a, int b) = [1, 2]
		return v.createMultiAssignmentDeclarationListStatement(ctx, modifierManager)
	}

	variableType := v.VisitType(ctx.Type_().(*TypeContext))
	declarators := ctx.VariableDeclarators().(*VariableDeclaratorsContext)
	declarators.PutNodeMetaData(VARIABLE_DECLARATION_VARIABLE_TYPE, variableType)
	declarationExpressionList := v.VisitVariableDeclarators(ctx.VariableDeclarators().(*VariableDeclaratorsContext))

	// if classNode is not nil, the variable declaration is for class declaration. In other words, it is a field declaration
	classNode := ctx.GetNodeMetaData(CLASS_DECLARATION_CLASS_NODE).(*ClassNode)

	if classNode != nil {
		return v.createFieldDeclarationListStatement(ctx, modifierManager, variableType, declarationExpressionList, classNode)
	}

	size := len(declarationExpressionList)
	if size > 0 {
		for _, e := range declarationExpressionList {
			modifierManager.ProcessVariableExpression(e.GetVariableExpression())
			modifierManager.AttachAnnotations(e.AnnotatedNode)
		}

		declarationExpression := declarationExpressionList[0]
		if size == 1 {
			configureAST(declarationExpression, ctx)
		} else { // adjust start of first declaration
			declarationExpression.SetLineNumber(ctx.GetStart().GetLine())
			declarationExpression.SetColumnNumber(ctx.GetStart().GetColumn() + 1)
		}
	}

	return configureAST(NewDeclarationListStatement(declarationExpressionList...), ctx)
}

func (v *ASTBuilder) createFieldDeclarationListStatement(ctx *VariableDeclarationContext, modifierManager *ModifierManager, variableType *ClassNode, declarationExpressionList []*DeclarationExpression, classNode *ClassNode) *DeclarationListStatement {
	for i, declarationExpression := range declarationExpressionList {
		variableExpression := declarationExpression.GetLeftExpression().(*VariableExpression)

		fieldName := variableExpression.GetName()

		modifiers := modifierManager.GetClassMemberModifiersOpValue()

		var initialValue Expression
		if _, ok := declarationExpression.GetRightExpression().(*EmptyExpression); !ok {
			initialValue = declarationExpression.GetRightExpression()
		}
		defaultValue := v.findDefaultValueByType(variableType)

		if classNode.IsInterface() {
			if initialValue == nil {
				if defaultValue == nil {
					initialValue = nil
				} else {
					initialValue = NewConstantExpression(defaultValue)
				}
			}

			modifiers |= ACC_PUBLIC | ACC_STATIC | ACC_FINAL
		}

		if v.isFieldDeclaration(modifierManager, classNode) {
			v.declareField(ctx, modifierManager, variableType, classNode, i, variableExpression, fieldName, modifiers, initialValue)
		} else {
			v.declareProperty(ctx.GroovyParserRuleContext, modifierManager, variableType, classNode, i, variableExpression, fieldName, modifiers, initialValue)
		}
	}

	return nil
}

/*

type PropertyExpander struct {
	*Verifier
}

func NewPropertyExpander(cNode *ClassNode) *PropertyExpander {
	pe := &PropertyExpander{
		Verifier: NewVerifier(),
	}
	pe.SetClassNode(cNode)
	return pe
}

func (pe *PropertyExpander) CreateSetterBlock(propertyNode *PropertyNode, field *FieldNode) Statement {
	return NewExpressionStatement(
		NewBinaryExpression(
			NewVariableExpression(field),
			NewToken(ASSIGN, "="),
			NewVariableExpression(NewVariableExpression(VALUE_STR), field.GetType()),
		),
	)
}

func (pe *PropertyExpander) CreateGetterBlock(propertyNode *PropertyNode, field *FieldNode) Statement {
	return NewExpressionStatement(NewVariableExpression(field))
}

*/

func (v *ASTBuilder) declareProperty(ctx *GroovyParserRuleContext, modifierManager *ModifierManager, variableType *ClassNode, classNode *ClassNode, i int, startNode ASTNode, fieldName string, modifiers int, initialValue Expression) *PropertyNode {
	var propertyNode *PropertyNode
	fieldNode := classNode.GetDeclaredField(fieldName)

	if fieldNode != nil && !classNode.HasProperty(fieldName) {
		if fieldNode.HasInitialExpression() && initialValue != nil {
			panic(createParsingFailedException("The split property definition named '"+fieldName+"' must not have an initial value for both the field and the property", parserRuleContextAdapter{ctx}))
		}
		if !fieldNode.GetType().Equals(variableType) {
			panic(createParsingFailedException("The split property definition named '"+fieldName+"' must not have different types for the field and the property", parserRuleContextAdapter{ctx}))
		}
		classNode.RemoveField(fieldNode)

		propertyNode = NewPropertyNode(fieldNode, modifiers|ACC_PUBLIC, nil, nil)
		classNode.AddProperty(propertyNode)
		if initialValue != nil {
			fieldNode.SetInitialValueExpression(initialValue)
		}
		modifierManager.AttachAnnotations(propertyNode.AnnotatedNode)
		// TODO: implement this
		//propertyNode.AddAnnotationNode(makeAnnotationNode(CompileStatic))
		// expand properties early so AST transforms will be handled correctly
		// TODO: implement this
		//expander := NewPropertyExpander(classNode)
		//expander.VisitProperty(propertyNode)
	} else {
		fieldNode := NewFieldNode(fieldName, modifiers&^ACC_PUBLIC|ACC_PRIVATE, variableType, classNode, initialValue)
		propertyNode = NewPropertyNode(fieldNode, modifiers|ACC_PUBLIC, nil, nil)
		classNode.AddProperty(propertyNode)

		fieldNode = propertyNode.GetField()
		fieldNode.SetModifiers(modifiers&^ACC_PUBLIC | ACC_PRIVATE)
		fieldNode.SetSynthetic(!classNode.IsInterface())
		modifierManager.AttachAnnotations(fieldNode.AnnotatedNode)
		modifierManager.AttachAnnotations(propertyNode.AnnotatedNode)
		if i == 0 {
			configureAST(fieldNode, ctx)
		} else {
			configureASTFromSource(fieldNode, startNode)
		}
	}

	//v.groovydocManager.Handle(fieldNode, ctx)
	//v.groovydocManager.Handle(propertyNode, ctx)

	if i == 0 {
		configureAST(propertyNode, ctx)
	} else {
		configureASTFromSource(propertyNode, startNode)
	}
	return propertyNode
}

func (v *ASTBuilder) declareField(ctx *VariableDeclarationContext, modifierManager *ModifierManager, variableType *ClassNode, classNode *ClassNode, i int, variableExpression *VariableExpression, fieldName string, modifiers int, initialValue Expression) {
	var fieldNode *FieldNode
	propertyNode := classNode.GetProperty(fieldName)

	if propertyNode != nil && propertyNode.GetField().IsSynthetic() {
		if propertyNode.HasInitialExpression() && initialValue != nil {
			panic(createParsingFailedException("The split property definition named '"+fieldName+"' must not have an initial value for both the field and the property", parserRuleContextAdapter{ctx}))
		}
		if !propertyNode.GetType().Equals(variableType) {
			panic(createParsingFailedException("The split property definition named '"+fieldName+"' must not have different types for the field and the property", parserRuleContextAdapter{ctx}))
		}
		classNode.RemoveField(propertyNode.GetField())
		var initialExpr Expression
		if propertyNode.HasInitialExpression() {
			initialExpr = propertyNode.GetInitialExpression()
		} else {
			initialExpr = initialValue
		}
		fieldNode = NewFieldNode(fieldName, modifiers, variableType, classNode.Redirect(), initialExpr)
		propertyNode.SetField(fieldNode)
		// TODO: implement this
		// propertyNode.AddAnnotation(makeAnnotationNode(CompileStatic))
		classNode.AddField(fieldNode)
		// expand properties early so AST transforms will be handled correctly
		// TODO: implement this
		//expander := NewPropertyExpander(classNode)
		//expander.VisitProperty(propertyNode)
	} else {
		fieldNode = NewFieldNode(fieldName, modifiers, variableType, classNode, initialValue)
		classNode.AddField(fieldNode)
	}

	modifierManager.AttachAnnotations(fieldNode.AnnotatedNode)
	//v.groovydocManager.Handle(fieldNode, ctx)

	if i == 0 {
		configureAST(fieldNode, ctx)
	} else {
		configureASTFromSource(fieldNode, variableExpression)
	}
}

func (v *ASTBuilder) isFieldDeclaration(modifierManager *ModifierManager, classNode *ClassNode) bool {
	return classNode.IsInterface() || modifierManager.ContainsVisibilityModifier()
}

func (v *ASTBuilder) VisitTypeNamePairs(ctx *TypeNamePairsContext) []Expression {
	pairs := make([]Expression, 0, len(ctx.AllTypeNamePair()))
	for _, pair := range ctx.AllTypeNamePair() {
		pairs = append(pairs, v.VisitTypeNamePair(pair.(*TypeNamePairContext)))
	}
	return pairs
}

func (v *ASTBuilder) VisitTypeNamePair(ctx *TypeNamePairContext) *VariableExpression {
	return configureAST(
		NewVariableExpression(
			v.VisitVariableDeclaratorId(ctx.VariableDeclaratorId().(*VariableDeclaratorIdContext)).GetName(),
			v.VisitType(ctx.Type_().(*TypeContext)),
		),
		ctx,
	)
}

func (v *ASTBuilder) VisitVariableDeclarators(ctx *VariableDeclaratorsContext) []*DeclarationExpression {
	variableType := ctx.GetNodeMetaData(VARIABLE_DECLARATION_VARIABLE_TYPE).(*ClassNode)
	if variableType == nil {
		panic("variableType should not be nil")
	}

	declarationExpressions := make([]*DeclarationExpression, 0, len(ctx.AllVariableDeclarator()))
	for _, e := range ctx.AllVariableDeclarator() {
		variableDeclaratorContext := e.(*VariableDeclaratorContext)
		variableDeclaratorContext.PutNodeMetaData(VARIABLE_DECLARATION_VARIABLE_TYPE, variableType)
		declarationExpressions = append(declarationExpressions, v.VisitVariableDeclarator(variableDeclaratorContext))
	}
	return declarationExpressions
}

func (v *ASTBuilder) VisitVariableDeclarator(ctx *VariableDeclaratorContext) *DeclarationExpression {
	variableType := ctx.GetNodeMetaData(VARIABLE_DECLARATION_VARIABLE_TYPE).(*ClassNode)
	if variableType == nil {
		panic("variableType should not be nil")
	}

	var token *Token
	if ctx.ASSIGN() != nil {
		token = v.createGroovyTokenByType(ctx.ASSIGN().GetSymbol(), ASSIGN)
	} else {
		token = NewToken(ASSIGN, ASSIGN_STR, ctx.GetStart().GetLine(), 1)
	}

	return configureAST(
		NewDeclarationExpression(
			configureAST(
				NewVariableExpression(
					v.VisitVariableDeclaratorId(ctx.VariableDeclaratorId().(*VariableDeclaratorIdContext)).GetName(),
					variableType,
				),
				ctx.VariableDeclaratorId(),
			),
			token,
			v.VisitVariableInitializer(ctx.VariableInitializer().(*VariableInitializerContext)),
		),
		ctx,
	)
}

func (v *ASTBuilder) VisitVariableInitializer(ctx *VariableInitializerContext) Expression {
	if ctx == nil {
		return EMPTY_EXPRESSION
	}

	return configureAST(
		v.VisitEnhancedStatementExpression(ctx.EnhancedStatementExpression().(*EnhancedStatementExpressionContext)),
		ctx,
	)
}

func (v *ASTBuilder) VisitVariableInitializers(ctx *VariableInitializersContext) []Expression {
	if ctx == nil {
		return []Expression{}
	}

	initializers := make([]Expression, 0, len(ctx.AllVariableInitializer()))
	for _, initCtx := range ctx.AllVariableInitializer() {
		initializers = append(initializers, v.VisitVariableInitializer(initCtx.(*VariableInitializerContext)))
	}
	return initializers
}

func (v *ASTBuilder) VisitArrayInitializer(ctx *ArrayInitializerContext) []Expression {
	if ctx == nil {
		return []Expression{}
	}

	v.visitingArrayInitializerCount++
	defer func() {
		v.visitingArrayInitializerCount--
	}()

	return v.VisitVariableInitializers(ctx.VariableInitializers().(*VariableInitializersContext))
}

func (v *ASTBuilder) VisitBlock(ctx *BlockContext) Statement {
	if ctx == nil {
		return v.createBlockStatement()
	}

	return configureAST(
		v.VisitBlockStatementsOpt(ctx.BlockStatementsOpt().(*BlockStatementsOptContext)),
		ctx)
}

func (v *ASTBuilder) VisitCommandExprAlt(ctx *CommandExprAltContext) *ExpressionStatement {
	expr, err := NewExpressionStatement(v.VisitCommandExpression(ctx.CommandExpression().(*CommandExpressionContext)))
	if err != nil {
		panic(createParsingFailedException(err.Error(), parserRuleContextAdapter{ctx}))
	}
	return configureAST(expr, ctx)
}

func (v *ASTBuilder) getOriginalText(ctx antlr.ParserRuleContext) string {
	return ctx.GetStart().GetInputStream().GetText(ctx.GetStart().GetStart(), ctx.GetStop().GetStop())
}

func (v *ASTBuilder) VisitCommandExpression(ctx *CommandExpressionContext) Expression {
	hasArgumentList := ctx.ArgumentList() != nil
	hasCommandArgument := len(ctx.AllCommandArgument()) > 0

	if (hasArgumentList || hasCommandArgument) && v.visitingArrayInitializerCount > 0 {
		// To avoid ambiguities, command chain expression should not be used in array initializer
		// the old parser does not support either, so no breaking changes
		// SEE http://groovy.329449.n5.nabble.com/parrot-Command-expressions-in-array-initializer-tt5752273.html
		panic(createParsingFailedException("Command chain expression can not be used in array initializer", parserRuleContextAdapter{ctx}))
	}

	baseExpr := v.Visit(ctx.Expression()).(Expression)

	if (hasArgumentList || hasCommandArgument) && !v.isInsideParentheses(baseExpr) {
		if binaryExpr, ok := baseExpr.(*BinaryExpression); ok && binaryExpr.GetOperation().GetText() != "[" {
			panic(createParsingFailedException("Unexpected input: '"+v.getOriginalText(ctx.Expression())+"'", parserRuleContextAdapter{ctx.Expression()}))
		}
	}

	var methodCallExpression *MethodCallExpression

	if hasArgumentList {
		arguments := v.VisitEnhancedArgumentListInParArgs(ctx.ArgumentList().(*ArgumentListContext))

		if propertyExpr, ok := baseExpr.(*PropertyExpression); ok { // e.g. obj.a 1, 2
			methodCallExpression = configureAST(v.createMethodCallExpression(propertyExpr, arguments), ctx.Expression())
		} else if methodCallExpr, ok := baseExpr.(*MethodCallExpression); ok && !v.isInsideParentheses(baseExpr) { // e.g. m {} a, b  OR  m(...) a, b
			if arguments != nil {
				// The error should never be thrown.
				panic("When baseExpr is a instance of MethodCallExpression, which should follow NO argumentList")
			}
			methodCallExpression = methodCallExpr
		} else if !v.isInsideParentheses(baseExpr) &&
			(IsInstanceOf(baseExpr, (*VariableExpression)(nil)) || // e.g. m 1, 2
				IsInstanceOf(baseExpr, (*GStringExpression)(nil)) || // e.g. "$m" 1, 2
				(IsInstanceOf(baseExpr, (*ConstantExpression)(nil)) && isTrue(baseExpr, IS_STRING))) { // e.g. "m" 1, 2
			v.validateInvalidMethodDefinition(baseExpr, arguments)
			methodCallExpression = configureAST(v.createMethodCallExpression(baseExpr, arguments), ctx.Expression())
		} else { // e.g. a[x] b, new A() b, etc.
			methodCallExpression = configureAST(v.createCallMethodCallExpression(baseExpr, arguments), ctx.Expression())
		}

		methodCallExpression.SetNodeMetaData(IS_COMMAND_EXPRESSION, true)

		if !hasCommandArgument {
			return methodCallExpression
		}
	}

	if hasCommandArgument {
		baseExpr.PutNodeMetaData(IS_COMMAND_EXPRESSION, true)
	}

	result := methodCallExpression

	for _, cmdArgCtx := range ctx.AllCommandArgument() {
		commandArgumentContext := cmdArgCtx.(*CommandArgumentContext)
		commandArgumentContext.PutNodeMetaData(CMD_EXPRESSION_BASE_EXPR, result)
		result = v.VisitCommandArgument(commandArgumentContext).(*MethodCallExpression)
	}

	return configureAST(result, ctx)
}

func (v *ASTBuilder) validateInvalidMethodDefinition(baseExpr Expression, arguments Expression) {
	if variableExpr, ok := baseExpr.(*VariableExpression); ok {
		if v.isBuiltInType(baseExpr) || unicode.IsUpper(rune(variableExpr.GetText()[0])) {
			if argumentListExpr, ok := arguments.(*ArgumentListExpression); ok {
				expressionList := argumentListExpr.GetExpressions()
				if len(expressionList) == 1 {
					expression := expressionList[0]
					if methodCallExpr, ok := expression.(*MethodCallExpression); ok {
						methodCallArguments := methodCallExpr.GetArguments()

						// check the method call tails with a closure
						if argumentListExpr, ok := methodCallArguments.(*ArgumentListExpression); ok {
							methodCallArgumentExpressionList := argumentListExpr.GetExpressions()
							argumentCnt := len(methodCallArgumentExpressionList)
							if argumentCnt > 0 {
								lastArgumentExpression := methodCallArgumentExpressionList[argumentCnt-1]
								if closureExpr, ok := lastArgumentExpression.(*ClosureExpression); ok {
									if HasImplicitParameter(closureExpr) {
										panic(createParsingFailedException(
											"Method definition not expected here",
											sourcePosition{baseExpr.GetLineNumber(), baseExpr.GetColumnNumber(), baseExpr.GetLineNumber(), baseExpr.GetColumnNumber()},
										))
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (v *ASTBuilder) VisitCommandArgument(ctx *CommandArgumentContext) Expression {
	// e.g. x y a b     we call "x y" as the base expression
	baseExpr := ctx.GetNodeMetaData(CMD_EXPRESSION_BASE_EXPR).(Expression)

	primaryExpr := v.Visit(ctx.CommandPrimary()).(Expression)

	/*
		if ctx.EnhancedArgumentListInPar() != nil { // e.g. x y a b
			if _, ok := baseExpr.(*PropertyExpression); ok { // the branch should never reach, because a.b.c will be parsed as a path expression, not a method call
				panic(createParsingFailedException("Unsupported command argument: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
			}

			// the following code will process "a b" of "x y a b"
			methodCallExpression := NewMethodCallExpression(
				baseExpr,
				v.createConstantExpression(primaryExpr),
				v.VisitEnhancedArgumentListInPar(ctx.EnhancedArgumentListInPar().(*EnhancedArgumentListInParContext)),
			)
			methodCallExpression.SetImplicitThis(false)

			return configureAST(methodCallExpression, ctx)
		} else if len(ctx.AllPathElement()) > 0 { // e.g. x y a.b
			pathExpression := v.createPathExpression(
				configureASTFromSource(
					NewPropertyExpressionWithProperty(baseExpr, v.createConstantExpression(primaryExpr)),
					primaryExpr,
				),
				ctx.AllPathElement(),
			)

			return configureAST(pathExpression, ctx)
		}
	*/

	if len(ctx.AllPathElement()) > 0 { // e.g. x y a.b
		pathExpression := v.createPathExpression(
			configureASTFromSource(
				NewPropertyExpressionWithProperty(baseExpr, v.createConstantExpression(primaryExpr)),
				primaryExpr,
			),
			ctx.AllPathElement(),
		)

		return configureAST(pathExpression, ctx)
	}

	// e.g. x y a
	var propertyExpr Expression
	if _, ok := primaryExpr.(*VariableExpression); ok {
		propertyExpr = v.createConstantExpression(primaryExpr)
	} else {
		propertyExpr = primaryExpr
	}

	return configureASTFromSource(
		NewPropertyExpressionWithProperty(baseExpr, propertyExpr),
		primaryExpr,
	)
}

func (v *ASTBuilder) VisitCastParExpression(ctx *CastParExpressionContext) *ClassNode {
	return v.VisitType(ctx.Type_().(*TypeContext))
}

func (v *ASTBuilder) VisitParExpression(ctx *ParExpressionContext) Expression {
	expression := v.VisitExpressionInPar(ctx.ExpressionInPar().(*ExpressionInParContext))

	level := expression.GetNodeMetaData(INSIDE_PARENTHESES_LEVEL)
	if level == nil {
		level = new(int)
		expression.SetNodeMetaData(INSIDE_PARENTHESES_LEVEL, level)
	}
	*level.(*int)++

	return configureAST(expression, ctx)
}

func (v *ASTBuilder) VisitExpressionInPar(ctx *ExpressionInParContext) Expression {
	return v.VisitEnhancedStatementExpression(ctx.EnhancedStatementExpression().(*EnhancedStatementExpressionContext))
}

func (v *ASTBuilder) VisitEnhancedStatementExpression(ctx *EnhancedStatementExpressionContext) Expression {
	var expression Expression

	if ctx.StatementExpression() != nil {
		stmt := v.Visit(ctx.StatementExpression()).(Statement)
		exprStmt, ok := stmt.(*ExpressionStatement)
		if !ok {
			panic(createParsingFailedException("Expected ExpressionStatement", parserRuleContextAdapter{ctx}))
		}
		expression = exprStmt.GetExpression()
	} else if ctx.StandardLambdaExpression() != nil {
		expression = v.VisitStandardLambdaExpression(ctx.StandardLambdaExpression().(*StandardLambdaExpressionContext))
	} else {
		panic(createParsingFailedException("Unsupported enhanced statement expression: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}

	return configureAST(expression, ctx)
}

func (v *ASTBuilder) VisitPathExpression(ctx *PathExpressionContext) interface{} {
	staticTerminalNode := ctx.STATIC()
	var primaryExpr Expression

	if staticTerminalNode != nil {
		primaryExpr = NewVariableExpressionWithString(staticTerminalNode.GetText())
	} else {
		primaryExpr = v.Visit(ctx.Primary()).(Expression)
	}

	return v.createPathExpression(primaryExpr, ctx.AllPathElement())
}

func (v *ASTBuilder) VisitPathElement(ctx *PathElementContext) interface{} {
	baseExpr := ctx.GetNodeMetaData(PATH_EXPRESSION_BASE_EXPR).(Expression)
	if baseExpr == nil {
		panic("baseExpr is required!")
	}

	if ctx.NamePart() != nil {
		namePartExpr := v.VisitNamePart(ctx.NamePart().(*NamePartContext)).(Expression)
		genericsTypes := (v.VisitNonWildcardTypeArguments(ctx.NonWildcardTypeArguments().(*NonWildcardTypeArgumentsContext))).([]*GenericsType)

		if ctx.DOT() != nil {
			isSafeChain := isTrue(baseExpr, PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN)
			return v.createDotExpression(ctx, baseExpr, namePartExpr, genericsTypes, isSafeChain)
		} else if ctx.SAFE_DOT() != nil {
			return v.createDotExpression(ctx, baseExpr, namePartExpr, genericsTypes, true)
		} else if ctx.SAFE_CHAIN_DOT() != nil {
			expression := v.createDotExpression(ctx, baseExpr, namePartExpr, genericsTypes, true)
			expression.SetNodeMetaData(PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN, true)
			return expression
		} else if ctx.METHOD_POINTER() != nil {
			return configureAST(NewMethodPointerExpression(baseExpr, namePartExpr), ctx)
		} else if ctx.METHOD_REFERENCE() != nil {
			return configureAST(NewMethodReferenceExpression(baseExpr, namePartExpr), ctx)
		} else if ctx.SPREAD_DOT() != nil {
			if ctx.AT() != nil {
				attributeExpression := NewAttributeExpressionWithSafe(baseExpr, namePartExpr, true)
				attributeExpression.SetSpreadSafe(true)
				return configureAST(attributeExpression, ctx)
			} else {
				propertyExpression := NewPropertyExpressionWithSafe(baseExpr, namePartExpr, true)
				propertyExpression.SetNodeMetaData(PATH_EXPRESSION_BASE_EXPR_GENERICS_TYPES, genericsTypes)
				propertyExpression.SetSpreadSafe(true)
				return configureAST(propertyExpression, ctx)
			}
		}
	} else if ctx.Creator() != nil {
		creatorContext := ctx.Creator().(*CreatorContext)
		creatorContext.SetNodeMetaData(ENCLOSING_INSTANCE_EXPRESSION, baseExpr)
		return configureAST(v.VisitCreator(creatorContext), ctx)
	} else if ctx.IndexPropertyArgs() != nil {
		tuple := v.VisitIndexPropertyArgs(ctx.IndexPropertyArgs().(*IndexPropertyArgsContext))
		isSafeChain := isTrue(baseExpr, PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN)
		return configureAST(
			NewBinaryExpressionWithSafe(
				baseExpr,
				v.createGroovyToken(tuple.V1),
				tuple.V2.(Expression),
				isSafeChain || ctx.IndexPropertyArgs().(*IndexPropertyArgsContext).SAFE_INDEX() != nil,
			),
			ctx,
		)
	} else if ctx.NamedPropertyArgs() != nil {
		mapEntryExpressionList := v.VisitNamedPropertyArgs(ctx.NamedPropertyArgs().(*NamedPropertyArgsContext))

		expressions := make([]Expression, len(mapEntryExpressionList))
		for i, v := range mapEntryExpressionList {
			expressions[i] = &v
		}
		listExpression := configureAST(
			NewListExpressionWithExpressions(expressions),
			ctx.NamedPropertyArgs(),
		)

		namedPropertyArgsContext := ctx.NamedPropertyArgs().(*NamedPropertyArgsContext)
		var token antlr.Token
		if namedPropertyArgsContext.LBRACK() == nil {
			token = namedPropertyArgsContext.SAFE_INDEX().GetSymbol()
		} else {
			token = namedPropertyArgsContext.LBRACK().GetSymbol()
		}
		return configureAST(
			NewBinaryExpression(baseExpr, v.createGroovyToken(token), listExpression),
			ctx,
		)
	} else if ctx.Arguments() != nil {
		argumentsExpr := v.VisitArguments(ctx.Arguments().(*ArgumentsContext)).(Expression)
		configureAST(argumentsExpr, ctx)

		if v.isInsideParentheses(baseExpr) {
			return configureAST(v.createCallMethodCallExpression(baseExpr, argumentsExpr), ctx)
		}

		if attributeExpression, ok := baseExpr.(*AttributeExpression); ok {
			attributeExpression.SetSpreadSafe(false)
			return configureAST(v.createCallMethodCallExpressionWithImplicitThis(attributeExpression, argumentsExpr, true), ctx)
		}

		if propertyExpression, ok := baseExpr.(*PropertyExpression); ok {
			methodCallExpression := v.createMethodCallExpression(propertyExpression, argumentsExpr)
			return configureAST(methodCallExpression, ctx)
		}

		if variableExpression, ok := baseExpr.(*VariableExpression); ok {
			baseExprText := variableExpression.GetText()
			if baseExprText == VOID_STR {
				return configureAST(v.createCallMethodCallExpression(v.createConstantExpression(baseExpr), argumentsExpr), ctx)
			} else if isPrimitiveType(baseExprText) {
				panic(createParsingFailedException("Primitive type literal: "+baseExprText+" cannot be used as a method name", parserRuleContextAdapter{ctx}))
			}
		}

		if _, ok := baseExpr.(*VariableExpression); ok {
			// Handle VariableExpression
		} else if _, ok := baseExpr.(*GStringExpression); ok {
			// Handle GStringExpression
		} else if ce, ok := baseExpr.(*ConstantExpression); ok && isTrue(ce, IS_STRING) {
			// Handle ConstantExpression that is a string
		}

		baseExprText := baseExpr.GetText()
		if baseExprText == THIS_STR || baseExprText == SUPER_STR {
			if v.visitingClosureCount > 0 {
				return configureAST(
					NewMethodCallExpression(
						baseExpr,
						baseExprText,
						argumentsExpr,
					),
					ctx,
				)
			}

			var classNode *ClassNode
			if baseExprText == SUPER_STR {
				classNode = SUPER
			} else {
				classNode = THIS
			}
			return configureAST(
				NewConstructorCallExpression(classNode, argumentsExpr),
				ctx,
			)
		}

		methodCallExpression := v.createMethodCallExpression(baseExpr, argumentsExpr)
		return configureAST(methodCallExpression, ctx)
	} else if ctx.ClosureOrLambdaExpression() != nil {
		closureExpression := v.VisitClosureOrLambdaExpression(ctx.ClosureOrLambdaExpression().(*ClosureOrLambdaExpressionContext))

		if methodCallExpression, ok := baseExpr.(*MethodCallExpression); ok {
			argumentsExpression := methodCallExpression.GetArguments()

			if argumentListExpression, ok := argumentsExpression.(*ArgumentListExpression); ok {
				argumentListExpression.AddExpression(closureExpression)
				return configureAST(methodCallExpression, ctx)
			}

			if tupleExpression, ok := argumentsExpression.(*TupleExpression); ok {
				namedArgumentListExpression := tupleExpression.GetExpression(0).(*NamedArgumentListExpression)

				if len(tupleExpression.GetExpressions()) > 0 {
					methodCallExpression.SetArguments(
						configureASTFromSource(
							NewArgumentListExpressionFromSlice(
								configureASTFromSource(
									NewMapExpressionWithEntries(namedArgumentListExpression.GetMapEntryExpressions()),
									namedArgumentListExpression,
								),
								closureExpression,
							),
							tupleExpression,
						),
					)
				} else {
					methodCallExpression.SetArguments(
						configureASTFromSource(
							NewArgumentListExpressionFromSlice(closureExpression),
							tupleExpression,
						),
					)
				}

				return configureAST(methodCallExpression, ctx)
			}
		}

		if propertyExpression, ok := baseExpr.(*PropertyExpression); ok {
			methodCallExpression := v.createMethodCallExpression(
				propertyExpression,
				configureASTFromSource(
					NewArgumentListExpressionFromSlice(closureExpression),
					closureExpression,
				),
			)

			return configureAST(methodCallExpression, ctx)
		}

		if _, ok := baseExpr.(*VariableExpression); ok {
			// Handle VariableExpression
		} else if _, ok := baseExpr.(*GStringExpression); ok {
			// Handle GStringExpression
		} else if ce, ok := baseExpr.(*ConstantExpression); ok && isTrue(ce, IS_STRING) {
			// Handle ConstantExpression that is a string
		}

		methodCallExpression := v.createMethodCallExpression(
			baseExpr,
			configureASTFromSource(
				NewArgumentListExpressionFromSlice(closureExpression),
				closureExpression,
			),
		)

		return configureAST(methodCallExpression, ctx)
	}

	panic(createParsingFailedException("Unsupported path element: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) createDotExpression(ctx *PathElementContext, baseExpr Expression, namePartExpr Expression, genericsTypes []*GenericsType, safe bool) Expression {
	if ctx.AT() != nil { // e.g. obj.@a  OR  obj?.@a
		return configureAST(NewAttributeExpressionWithSafe(baseExpr, namePartExpr, safe), ctx)
	} else { // e.g. obj.p  OR  obj?.p
		propertyExpression := NewPropertyExpressionWithSafe(baseExpr, namePartExpr, safe)
		propertyExpression.SetNodeMetaData(PATH_EXPRESSION_BASE_EXPR_GENERICS_TYPES, genericsTypes)
		return configureAST(propertyExpression, ctx)
	}
}

func (v *ASTBuilder) createCallMethodCallExpression(baseExpr Expression, argumentsExpr Expression) *MethodCallExpression {
	return v.createCallMethodCallExpressionWithImplicitThis(baseExpr, argumentsExpr, false)
}

func (v *ASTBuilder) createCallMethodCallExpressionWithImplicitThis(baseExpr Expression, argumentsExpr Expression, implicitThis bool) *MethodCallExpression {
	methodCallExpression := NewMethodCallExpression(baseExpr, NewConstantExpression(CALL_STR), argumentsExpr)
	methodCallExpression.SetImplicitThis(implicitThis)
	return methodCallExpression
}

// []*GenericsType
func (v *ASTBuilder) VisitNonWildcardTypeArguments(ctx *NonWildcardTypeArgumentsContext) interface{} {
	if ctx == nil {
		return nil
	}

	typeList := v.VisitTypeList(ctx.TypeList().(*TypeListContext))
	genericsTypes := make([]*GenericsType, len(typeList))
	for i, t := range typeList {
		genericsTypes[i] = v.createGenericsType(t)
	}
	return genericsTypes
}

func (v *ASTBuilder) VisitTypeList(ctx *TypeListContext) []*ClassNode {
	if ctx == nil {
		return []*ClassNode{}
	}

	typeContexts := ctx.AllType_()
	classNodes := make([]*ClassNode, len(typeContexts))
	for i, typeCtx := range typeContexts {
		classNodes[i] = v.VisitType(typeCtx.(*TypeContext))
	}
	return classNodes
}

func (v *ASTBuilder) VisitArguments(ctx *ArgumentsContext) Expression {
	if ctx != nil && ctx.COMMA() != nil && ctx.EnhancedArgumentListInPar() == nil {
		panic(createParsingFailedException("Expression expected", tokenAdapter{ctx.COMMA().GetSymbol()}))
	}

	if ctx == nil || ctx.EnhancedArgumentListInPar() == nil {
		return NewArgumentListExpression()
	}

	return configureAST(v.VisitEnhancedArgumentListInPar(ctx.EnhancedArgumentListInPar().(*EnhancedArgumentListInParContext)), ctx)
}

func (v *ASTBuilder) VisitEnhancedArgumentListInParArgs(ctx *ArgumentListContext) Expression {
	if ctx == nil {
		return nil
	}

	var expressionList []Expression
	var mapEntryExpressionList []*MapEntryExpression

	for _, element := range ctx.AllArgumentListElement() {
		e := v.VisitEnhancedArgumentListElementArg(element.(*ArgumentListElementContext))

		if mapEntryExpr, ok := e.(*MapEntryExpression); ok {
			v.validateDuplicatedNamedParameter(mapEntryExpressionList, mapEntryExpr)
			mapEntryExpressionList = append(mapEntryExpressionList, mapEntryExpr)
		} else {
			expressionList = append(expressionList, e)
		}
	}

	if len(mapEntryExpressionList) == 0 { // e.g. arguments like  1, 2 OR  someArg, e -> e
		return configureAST(
			NewArgumentListExpressionFromSlice(expressionList...),
			ctx)
	}

	if len(expressionList) == 0 { // e.g. arguments like  x: 1, y: 2
		return configureAST(
			NewTupleExpressionWithExpressions(
				configureAST(
					NewNamedArgumentListExpressionWithEntries(mapEntryExpressionList),
					ctx)),
			ctx)
	}

	if len(mapEntryExpressionList) > 0 && len(expressionList) > 0 { // e.g. arguments like x: 1, 'a', y: 2, 'b', z: 3
		argumentListExpression := NewArgumentListExpressionFromSlice(expressionList...)
		argumentListExpression.PrependExpression(configureAST(NewMapExpressionWithEntries(mapEntryExpressionList), ctx))
		return configureAST(argumentListExpression, ctx)
	}

	panic(createParsingFailedException("Unsupported argument list: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitEnhancedArgumentListInPar(ctx *EnhancedArgumentListInParContext) Expression {
	if ctx == nil {
		return nil
	}

	var expressionList []Expression
	var mapEntryExpressionList []*MapEntryExpression

	for _, element := range ctx.AllEnhancedArgumentListElement() {
		e := v.VisitEnhancedArgumentListElement(element.(*EnhancedArgumentListElementContext))

		if mapEntryExpr, ok := e.(*MapEntryExpression); ok {
			v.validateDuplicatedNamedParameter(mapEntryExpressionList, mapEntryExpr)
			mapEntryExpressionList = append(mapEntryExpressionList, mapEntryExpr)
		} else {
			expressionList = append(expressionList, e)
		}
	}

	if len(mapEntryExpressionList) == 0 { // e.g. arguments like  1, 2 OR  someArg, e -> e
		return configureAST(
			NewArgumentListExpressionFromSlice(expressionList...),
			ctx)
	}

	if len(expressionList) == 0 { // e.g. arguments like  x: 1, y: 2
		return configureAST(
			NewTupleExpressionWithExpressions(
				configureAST(
					NewNamedArgumentListExpressionWithEntries(mapEntryExpressionList),
					ctx)),
			ctx)
	}

	if len(mapEntryExpressionList) > 0 && len(expressionList) > 0 { // e.g. arguments like x: 1, 'a', y: 2, 'b', z: 3
		argumentListExpression := NewArgumentListExpressionFromSlice(expressionList...)
		argumentListExpression.PrependExpression(configureAST(NewMapExpressionWithEntries(mapEntryExpressionList), ctx))
		return configureAST(argumentListExpression, ctx)
	}

	panic(createParsingFailedException("Unsupported argument list: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) validateDuplicatedNamedParameter(mapEntryExpressionList []*MapEntryExpression, mapEntryExpression *MapEntryExpression) {
	keyExpression := mapEntryExpression.GetKeyExpression()
	if keyExpression == nil || v.isInsideParentheses(keyExpression) {
		return
	}

	parameterName := keyExpression.GetText()
	isDuplicatedNamedParameter := false
	for _, m := range mapEntryExpressionList {
		if m.GetKeyExpression().GetText() == parameterName {
			isDuplicatedNamedParameter = true
			break
		}
	}
	if !isDuplicatedNamedParameter {
		return
	}

	panic(createParsingFailedException("Duplicated named parameter '"+parameterName+"' found", astNodeAdapter{mapEntryExpression}))
}

func (v *ASTBuilder) VisitEnhancedArgumentListElementArg(ctx *ArgumentListElementContext) Expression {
	if ctx.ExpressionListElement() != nil {
		return configureAST(v.VisitExpressionListElement(ctx.ExpressionListElement().(*ExpressionListElementContext)), ctx)
	}

	// TODO: implement this
	/*
		if ctx.StandardLambdaExpression() != nil {
			return configureAST(v.VisitStandardLambdaExpression(ctx.StandardLambdaExpression().(*StandardLambdaExpressionContext)), ctx)
		}

		if ctx.MapEntry() != nil {
			return configureAST(v.VisitMapEntry(ctx.MapEntry().(*MapEntryContext)), ctx)
		}
	*/

	panic(createParsingFailedException("Unsupported enhanced argument list element: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitEnhancedArgumentListElement(ctx *EnhancedArgumentListElementContext) Expression {
	if ctx.ExpressionListElement() != nil {
		return configureAST(v.VisitExpressionListElement(ctx.ExpressionListElement().(*ExpressionListElementContext)), ctx)
	}

	if ctx.StandardLambdaExpression() != nil {
		return configureAST(v.VisitStandardLambdaExpression(ctx.StandardLambdaExpression().(*StandardLambdaExpressionContext)), ctx)
	}

	// TODO: implement this

	/*
		if ctx.MapEntry() != nil {
			return configureAST(v.VisitMapEntry(ctx.MapEntry().(*MapEntryContext)), ctx)
		}
	*/

	panic(createParsingFailedException("Unsupported enhanced argument list element: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitStringLiteral(ctx *StringLiteralContext) *ConstantExpression {
	text := v.parseStringLiteral(ctx.StringLiteral().GetText())

	constantExpression := NewConstantExpression(text)
	constantExpression.PutNodeMetaData(IS_STRING, true)
	return configureAST(constantExpression, ctx)
}

func (v *ASTBuilder) parseStringLiteral(text string) string {
	slashyType := v.getSlashyType(text)
	startsWithSlash := false

	if strings.HasPrefix(text, TSQ_STR) || strings.HasPrefix(text, TDQ_STR) {
		text = RemoveCR(text) // remove CR in the multiline string

		text = TrimQuotations(text, 3)
	} else if strings.HasPrefix(text, SQ_STR) || strings.HasPrefix(text, DQ_STR) || (startsWithSlash == strings.HasPrefix(text, SLASH_STR)) {
		if startsWithSlash { // the slashy string can span rows, so we have to remove CR for it
			text = RemoveCR(text) // remove CR in the multiline string
		}

		text = TrimQuotations(text, 1)
	} else if strings.HasPrefix(text, DOLLAR_SLASH_STR) {
		text = RemoveCR(text)

		text = TrimQuotations(text, 2)
	}

	// handle escapes.
	return ReplaceEscapes(text, slashyType)
}

func (v *ASTBuilder) getSlashyType(text string) int {
	if strings.HasPrefix(text, SLASH_STR) {
		return SLASHY
	} else if strings.HasPrefix(text, DOLLAR_SLASH_STR) {
		return DOLLAR_SLASHY
	} else {
		return NONE_SLASHY
	}
}

func (v *ASTBuilder) VisitIndexPropertyArgs(ctx *IndexPropertyArgsContext) Tuple2[antlr.Token, Expression] {
	expressionList := v.VisitExpressionList(ctx.ExpressionList().(*ExpressionListContext))
	var token antlr.Token
	if ctx.LBRACK() == nil {
		token = ctx.SAFE_INDEX().GetSymbol()
	} else {
		token = ctx.LBRACK().GetSymbol()
	}

	if len(expressionList) == 1 {
		expr := expressionList[0]

		var indexExpr Expression
		if _, ok := expr.(*SpreadExpression); ok { // e.g. a[*[1, 2]]
			listExpression := NewListExpressionWithExpressions(expressionList)
			listExpression.SetWrapped(false)

			indexExpr = listExpression
		} else { // e.g. a[1]
			indexExpr = expr
		}

		return NewTuple2(token, indexExpr)
	}

	// e.g. a[1, 2]
	listExpression := NewListExpressionWithExpressions(expressionList)
	listExpression.SetWrapped(true)

	var expr Expression = configureAST(listExpression, ctx)

	return NewTuple2(token, expr)
}

func (v *ASTBuilder) VisitNamedPropertyArgs(ctx *NamedPropertyArgsContext) []MapEntryExpression {
	// TODO: implement this
	panic("FOO")
	//return v.VisitMapEntryList(ctx.MapEntryList().(*MapEntryListContext))
}

func (v *ASTBuilder) VisitNamePart(ctx *NamePartContext) Expression {
	if ctx.Identifier() != nil {
		return configureAST(NewConstantExpression(v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))), ctx)
	} else if ctx.StringLiteral() != nil {
		return configureAST(v.VisitStringLiteral(ctx.StringLiteral().(*StringLiteralContext)), ctx)
	} else if ctx.DynamicMemberName() != nil {
		return configureAST(v.VisitDynamicMemberName(ctx.DynamicMemberName().(*DynamicMemberNameContext)), ctx)
	} else if ctx.Keywords() != nil {
		return configureAST(NewConstantExpression(ctx.Keywords().GetText()), ctx)
	}

	panic(createParsingFailedException("Unsupported name part: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitDynamicMemberName(ctx *DynamicMemberNameContext) Expression {
	if ctx.ParExpression() != nil {
		return configureAST(v.VisitParExpression(ctx.ParExpression().(*ParExpressionContext)), ctx)
	} else if ctx.Gstring() != nil {
		return configureAST(v.VisitGstring(ctx.Gstring().(*GstringContext)), ctx)
	}

	panic(createParsingFailedException("Unsupported dynamic member name: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitPostfixExpression(ctx *PostfixExpressionContext) Expression {
	pathExpr := v.VisitPathExpression(ctx.PathExpression().(*PathExpressionContext)).(Expression)

	if ctx.GetOp() != nil {
		postfixExpression := NewPostfixExpression(pathExpr, v.createGroovyToken(ctx.GetOp()))

		if v.visitingAssertStatementCount > 0 {
			// powerassert requires different column for values, so we have to copy the location of op
			return configureASTWithToken(postfixExpression, ctx.GetOp())
		} else {
			return configureAST(postfixExpression, ctx)
		}
	}

	return configureAST(pathExpr, ctx)
}

func (v *ASTBuilder) VisitUnaryNotExprAlt(ctx *UnaryNotExprAltContext) Expression {
	if ctx.NOT() != nil {
		return configureAST(
			NewNotExpression(v.Visit(ctx.Expression()).(Expression)),
			ctx)
	}

	if ctx.BITNOT() != nil {
		return configureAST(
			NewBitwiseNegationExpression(v.Visit(ctx.Expression()).(Expression)),
			ctx)
	}

	panic(createParsingFailedException("Unsupported unary expression: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitCastExprAlt(ctx *CastExprAltContext) *CastExpression {
	expr := v.Visit(&ctx.ExpressionContext).(Expression)
	if varExpr, ok := expr.(*VariableExpression); ok && varExpr.IsSuperExpression() {
		createParsingFailedException("Cannot cast or coerce `super`", parserRuleContextAdapter{ctx}) // GROOVY-9391
	}
	cast := NewCastExpression(v.VisitCastParExpression(ctx.CastParExpression().(*CastParExpressionContext)), expr)
	return configureAST(cast, ctx)
}

func (v *ASTBuilder) VisitPowerExprAlt(ctx *PowerExprAltContext) *BinaryExpression {
	return v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext)
}

func (v *ASTBuilder) VisitUnaryAddExprAlt(ctx *UnaryAddExprAltContext) Expression {
	expression := v.Visit(ctx.Expression()).(Expression)
	switch ctx.op.GetTokenType() {
	case GroovyParserADD:
		if v.isNonStringConstantOutsideParentheses(expression) {
			return configureAST(expression, ctx)
		}
		return configureAST(NewUnaryPlusExpression(expression), ctx)

	case GroovyParserSUB:
		if v.isNonStringConstantOutsideParentheses(expression) {
			constantExpression := expression.(*ConstantExpression)
			integerLiteralText := constantExpression.GetNodeMetaData(INTEGER_LITERAL_TEXT)
			if integerLiteralText != nil {
				result := NewConstantExpression(ParseInteger(SUB_STR + integerLiteralText.(string)))
				/*
					if err != nil {
						panic(createParsingFailedException(err.Error(), ctx))
					}
				*/
				v.numberFormatError = nil // reset
				return configureAST(result, ctx)
			}

			floatingPointLiteralText := constantExpression.GetNodeMetaData(FLOATING_POINT_LITERAL_TEXT)
			if floatingPointLiteralText != nil {
				result := NewConstantExpression(ParseDecimal(SUB_STR + floatingPointLiteralText.(string)))
				/*
					if err != nil {
						panic(createParsingFailedException(err.Error(), ctx))
					}
				*/
				v.numberFormatError = nil // reset
				return configureAST(result, ctx)
			}

			panic("Failed to find the original number literal text: " + constantExpression.GetText())
		}
		return configureAST(NewUnaryMinusExpression(expression), ctx)

	case GroovyParserINC, GroovyParserDEC:
		return configureAST(NewPrefixExpression(v.createGroovyToken(ctx.op), expression), ctx)

	default:
		panic(createParsingFailedException("Unsupported unary operation: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}
}

func (v *ASTBuilder) isNonStringConstantOutsideParentheses(expression Expression) bool {
	if constantExpr, ok := expression.(*ConstantExpression); ok {
		_, isString := constantExpr.GetValue().(string)
		return !isString && !v.isInsideParentheses(expression)
	}
	return false
}

func (v *ASTBuilder) VisitMultiplicativeExprAlt(ctx *MultiplicativeExprAltContext) *BinaryExpression {
	return v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext)
}

func (v *ASTBuilder) VisitAdditiveExprAlt(ctx *AdditiveExprAltContext) interface{} {
	return v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext)
}

func (v *ASTBuilder) VisitShiftExprAlt(ctx *ShiftExprAltContext) Expression {
	left := v.Visit(ctx.left).(Expression)
	right := v.Visit(ctx.right).(Expression)

	if ctx.rangeOp != nil {
		return configureAST(NewRangeExpressionWithExclusive(left, right, strings.HasPrefix(ctx.rangeOp.GetText(), "<"), strings.HasSuffix(ctx.rangeOp.GetText(), "<")), ctx)
	}

	var op *Token
	var antlrToken antlr.Token

	if ctx.dlOp != nil {
		op = v.createGroovyTokenWithCardinality(ctx.dlOp, 2)
		antlrToken = ctx.dlOp
	} else if ctx.dgOp != nil {
		op = v.createGroovyTokenWithCardinality(ctx.dgOp, 2)
		antlrToken = ctx.dgOp
	} else if ctx.tgOp != nil {
		op = v.createGroovyTokenWithCardinality(ctx.tgOp, 3)
		antlrToken = ctx.tgOp
	} else {
		panic(createParsingFailedException("Unsupported shift expression: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}

	binaryExpression := NewBinaryExpression(left, op, right)
	if isTrue(ctx, IS_INSIDE_CONDITIONAL_EXPRESSION) {
		return configureASTWithToken(binaryExpression, antlrToken)
	}

	return configureAST(binaryExpression, ctx)
}

func (v *ASTBuilder) VisitRelationalExprAlt(ctx *RelationalExprAltContext) Expression {
	switch ctx.op.GetTokenType() {
	case GroovyParserAS:
		expr := v.Visit(ctx.left).(Expression)
		if varExpr, ok := expr.(*VariableExpression); ok && varExpr.IsSuperExpression() {
			createParsingFailedException("Cannot cast or coerce `super`", parserRuleContextAdapter{ctx}) // GROOVY-9391
		}
		cast := NewCastExpression(v.VisitType(ctx.Type_().(*TypeContext)), expr)
		return configureAST(cast, ctx)

	case GroovyParserINSTANCEOF, GroovyParserNOT_INSTANCEOF:
		ctx.Type_().(*TypeContext).PutNodeMetaData(IS_INSIDE_INSTANCEOF_EXPR, true)
		return configureAST(
			NewBinaryExpression(
				v.Visit(ctx.left).(Expression),
				v.createGroovyToken(ctx.op),
				configureAST(NewClassExpression(v.VisitType(ctx.Type_().(*TypeContext))), ctx.Type_()),
			),
			ctx,
		)

	case GroovyParserGT, GroovyParserGE, GroovyParserLT, GroovyParserLE, GroovyParserIN, GroovyParserNOT_IN:
		return v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext)

	default:
		panic(createParsingFailedException("Unsupported relational expression: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}
}

func (v *ASTBuilder) VisitEqualityExprAlt(ctx *EqualityExprAltContext) *BinaryExpression {
	return configureAST(
		v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext),
		ctx)
}

func (v *ASTBuilder) VisitRegexExprAlt(ctx *RegexExprAltContext) *BinaryExpression {
	return configureAST(
		v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext),
		ctx)
}

func (v *ASTBuilder) VisitAndExprAlt(ctx *AndExprAltContext) *BinaryExpression {
	return v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext)
}

func (v *ASTBuilder) VisitExclusiveOrExprAlt(ctx *ExclusiveOrExprAltContext) *BinaryExpression {
	return v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext)
}

func (v *ASTBuilder) VisitInclusiveOrExprAlt(ctx *InclusiveOrExprAltContext) *BinaryExpression {
	return v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext)
}

func (v *ASTBuilder) VisitLogicalAndExprAlt(ctx *LogicalAndExprAltContext) *BinaryExpression {
	return configureAST(
		v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext),
		ctx)
}

func (v *ASTBuilder) VisitLogicalOrExprAlt(ctx *LogicalOrExprAltContext) *BinaryExpression {
	return configureAST(
		v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext),
		ctx)
}

func (v *ASTBuilder) VisitImplicationExprAlt(ctx *ImplicationExprAltContext) *BinaryExpression {
	return configureAST(
		v.createBinaryExpression(ctx.left.(*ExpressionContext), ctx.right.(*ExpressionContext), ctx.op, &ctx.ExpressionContext),
		ctx)
}

func (v *ASTBuilder) VisitConditionalExprAlt(ctx *ConditionalExprAltContext) Expression {
	fbExpr := ctx.fb.(*ExpressionContext)
	fbExpr.PutNodeMetaData(IS_INSIDE_CONDITIONAL_EXPRESSION, true)

	if ctx.ELVIS() != nil { // e.g. a == 6 ?: 0
		return configureAST(
			NewElvisOperatorExpression(v.Visit(ctx.con).(Expression), v.Visit(ctx.fb).(Expression)),
			ctx)
	}

	tbExpr := ctx.tb.(*ExpressionContext)
	tbExpr.PutNodeMetaData(IS_INSIDE_CONDITIONAL_EXPRESSION, true)

	return configureAST(
		NewTernaryExpression(
			configureAST(NewBooleanExpression(v.Visit(ctx.con).(Expression)),
				ctx.con),
			v.Visit(ctx.tb).(Expression),
			v.Visit(ctx.fb).(Expression)),
		ctx)
}

func (v *ASTBuilder) VisitMultipleAssignmentExprAlt(ctx *MultipleAssignmentExprAltContext) *BinaryExpression {
	return configureAST(
		NewBinaryExpression(
			v.VisitVariableNames(ctx.left.(*VariableNamesContext)),
			v.createGroovyToken(ctx.op),
			v.Visit(ctx.right).(*ExpressionStatement).GetExpression()),
		ctx)
}

func (v *ASTBuilder) VisitAssignmentExprAlt(ctx *AssignmentExprAltContext) *BinaryExpression {
	leftExpr := v.Visit(ctx.left).(Expression)

	if _, ok := leftExpr.(*VariableExpression); ok && v.isInsideParentheses(leftExpr) {
		// it is a special multiple assignment whose variable count is only one, e.g. (a) = [1]
		insideParenthesesLevel := leftExpr.GetNodeMetaData(INSIDE_PARENTHESES_LEVEL).(int)
		if insideParenthesesLevel > 1 {
			panic(createParsingFailedException("Nested parenthesis is not allowed in multiple assignment, e.g. ((a)) = b", parserRuleContextAdapter{ctx}))
		}

		return configureAST(
			NewBinaryExpression(
				configureAST(NewTupleExpressionWithExpressions(leftExpr), ctx.left),
				v.createGroovyToken(ctx.op),
				v.Visit(ctx.right).(Expression)),
			ctx)
	}

	// the LHS expression should be a variable which is not inside any parentheses
	isValidLHS := false

	switch expr := leftExpr.(type) {
	case *VariableExpression:
		isValidLHS = !v.isInsideParentheses(leftExpr)
	case *PropertyExpression:
		isValidLHS = true
	case *BinaryExpression:
		isValidLHS = expr.GetOperation().GetType() == LEFT_SQUARE_BRACKET
	}

	if !isValidLHS {
		panic(createParsingFailedException("The LHS of an assignment should be a variable or a field accessing expression", parserRuleContextAdapter{ctx}))
	}

	return configureAST(
		NewBinaryExpression(
			leftExpr,
			v.createGroovyToken(ctx.op),
			v.Visit(ctx.right).(Expression)),
		ctx)
}

func (v *ASTBuilder) VisitIdentifierPrmrAlt(ctx *IdentifierPrmrAltContext) Expression {
	if ctx.TypeArguments() != nil {
		classNode := MakeFromString(ctx.Identifier().GetText())

		classNode.SetGenericsTypes(
			v.VisitTypeArguments(ctx.TypeArguments().(*TypeArgumentsContext)))

		return configureAST(NewClassExpression(classNode), ctx)
	}

	return configureAST(NewVariableExpressionWithString(v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))), ctx)
}

func (v *ASTBuilder) VisitNewPrmrAlt(ctx *NewPrmrAltContext) Expression {
	return configureAST(v.VisitCreator(ctx.Creator().(*CreatorContext)), ctx)
}

func (v *ASTBuilder) VisitThisPrmrAlt(ctx *ThisPrmrAltContext) *VariableExpression {
	return configureAST(NewVariableExpressionWithString(ctx.THIS().GetText()), ctx)
}

func (v *ASTBuilder) VisitSuperPrmrAlt(ctx *SuperPrmrAltContext) *VariableExpression {
	return configureAST(NewVariableExpressionWithString(ctx.SUPER().GetText()), ctx)
}

func (v *ASTBuilder) VisitCreator(ctx *CreatorContext) Expression {
	classNode := v.VisitCreatedName(ctx.CreatedName().(*CreatedNameContext))

	if ctx.Arguments() != nil { // create instance of class
		arguments := v.VisitArguments(ctx.Arguments().(*ArgumentsContext)).(Expression)
		enclosingInstanceExpression := ctx.GetNodeMetaData(ENCLOSING_INSTANCE_EXPRESSION).(Expression)

		if enclosingInstanceExpression != nil {
			if argumentListExpression, ok := arguments.(*ArgumentListExpression); ok {
				argumentListExpression.PrependExpression(enclosingInstanceExpression)
			} else if _, ok := arguments.(*TupleExpression); ok {
				panic(createParsingFailedException("Creating instance of non-static class does not support named parameters", astNodeAdapter{arguments}))
			} else if _, ok := arguments.(*NamedArgumentListExpression); ok {
				panic(createParsingFailedException("Unexpected arguments", parserRuleContextAdapter{ctx}))
			} else {
				panic(createParsingFailedException("Unsupported arguments", parserRuleContextAdapter{ctx})) // should never reach here
			}
			if constructorCallExpression, ok := enclosingInstanceExpression.(*ConstructorCallExpression); ok && !strings.Contains(classNode.GetName(), ".") {
				classNode.SetName(constructorCallExpression.GetType().GetName() + "." + classNode.GetName()) // GROOVY-8947
			}
		}

		if ctx.AnonymousInnerClassDeclaration() != nil {
			ctx.AnonymousInnerClassDeclaration().(*AnonymousInnerClassDeclarationContext).PutNodeMetaData(ANONYMOUS_INNER_CLASS_SUPER_CLASS, classNode)
			anonymousInnerClassNode := v.VisitAnonymousInnerClassDeclaration(ctx.AnonymousInnerClassDeclaration().(*AnonymousInnerClassDeclarationContext))

			anonymousInnerClassList := v.peekAnonymousInnerClass()
			if anonymousInnerClassList != nil { // if the anonymous class is created in a script, no anonymousInnerClassList is available.
				anonymousInnerClassList.PushBack(anonymousInnerClassNode)
			}

			constructorCallExpression := NewConstructorCallExpression(anonymousInnerClassNode.GetPlainNodeReference(), arguments)
			constructorCallExpression.SetUsingAnonymousInnerClass(true)

			return configureAST(constructorCallExpression, ctx)
		}

		constructorCallExpression := NewConstructorCallExpression(classNode, arguments)
		return configureAST(constructorCallExpression, ctx)
	}

	if len(ctx.AllDim()) > 0 { // create array
		var arrayExpression *ArrayExpression

		dimList := make([]Tuple3[Expression, []*AnnotationNode, antlr.TerminalNode], len(ctx.AllDim()))
		for i, dim := range ctx.AllDim() {
			dimList[i] = v.VisitDim(dim.(*DimContext))
		}

		var invalidDimLBrack antlr.TerminalNode
		var exprEmpty *bool
		emptyDimList := make([]Tuple3[Expression, []*AnnotationNode, antlr.TerminalNode], 0)
		dimWithExprList := make([]Tuple3[Expression, []*AnnotationNode, antlr.TerminalNode], 0)
		var latestDim Tuple3[Expression, []*AnnotationNode, antlr.TerminalNode]
		for _, dim := range dimList {
			if dim.V1 == nil {
				emptyDimList = append(emptyDimList, dim)
				trueVal := true
				exprEmpty = &trueVal
			} else {
				if exprEmpty != nil && *exprEmpty {
					invalidDimLBrack = latestDim.V3
				}

				dimWithExprList = append(dimWithExprList, dim)
				falseVal := false
				exprEmpty = &falseVal
			}

			latestDim = dim
		}

		if ctx.ArrayInitializer() != nil {
			if len(dimWithExprList) > 0 {
				panic(createParsingFailedException("dimension should be empty", tokenAdapter{dimWithExprList[0].V3.GetSymbol()}))
			}

			elementType := classNode
			for i := 0; i < len(emptyDimList)-1; i++ {
				elementType = v.createArrayType(elementType)
			}

			arrayExpression = NewArrayExpression(
				elementType,
				v.VisitArrayInitializer(ctx.ArrayInitializer().(*ArrayInitializerContext)),
				nil,
			)

		} else {
			if invalidDimLBrack != nil {
				panic(createParsingFailedException("dimension cannot be empty", parserRuleContextAdapter{ctx}))
			}

			if len(dimWithExprList) == 0 && len(emptyDimList) > 0 {
				panic(createParsingFailedException("dimensions cannot be all empty", tokenAdapter{emptyDimList[0].V3.GetSymbol()}))
			}

			var empties []Expression
			if len(emptyDimList) > 0 {
				empties = make([]Expression, len(emptyDimList))
				for i := range empties {
					empties[i] = EMPTY_EXPRESSION
				}
			} else {
				empties = []Expression{}
			}

			expressions := make([]Expression, len(dimWithExprList)+len(empties))
			for i, dim := range dimWithExprList {
				expressions[i] = dim.V1
			}
			copy(expressions[len(dimWithExprList):], empties)

			arrayExpression = NewArrayExpression(
				classNode,
				nil,
				expressions,
			)
		}

		annotations := make([][]*AnnotationNode, len(dimList))
		for i, dim := range dimList {
			annotations[i] = dim.V2
		}

		arrayExpression.SetType(
			v.createArrayTypeAnnotations(
				classNode,
				annotations,
				ctx,
			),
		)

		return configureAST(arrayExpression, ctx)
	}

	panic(createParsingFailedException("Unsupported creator: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitDim(ctx *DimContext) Tuple3[Expression, []*AnnotationNode, antlr.TerminalNode] {
	return NewTuple3(
		v.Visit(ctx.Expression()).(Expression),
		v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext)),
		ctx.LBRACK(),
	)
}

func nextAnonymousClassName(outerClass *ClassNode) string {
	anonymousClassCount := 0
	for _, innerClass := range outerClass.innerClasses {
		if innerClass.IsAnonymous() {
			anonymousClassCount++
		}
	}

	return outerClass.GetName() + "$" + strconv.Itoa(anonymousClassCount+1)
}

func (v *ASTBuilder) VisitAnonymousInnerClassDeclaration(ctx *AnonymousInnerClassDeclarationContext) *InnerClassNode {
	superClass := ctx.GetNodeMetaData(ANONYMOUS_INNER_CLASS_SUPER_CLASS).(*ClassNode)
	if superClass == nil {
		panic("superClass should not be nil")
	}

	var outerClass *ClassNode
	if v.classNodeStack.Len() > 0 {
		outerClass = v.peekClassNode()
	} else {
		outerClass = v.moduleNode.GetScriptClassDummy()
	}
	innerClassName := nextAnonymousClassName(outerClass)

	var anonymousInnerClass *InnerClassNode
	if ctx.t == 1 {
		anonymousInnerClass = NewEnumConstantClassNode(outerClass, innerClassName, superClass.GetPlainNodeReference()).InnerClassNode
		// and remove the final modifier from superClass to allow the sub class
		superClass.SetModifiers(superClass.GetModifiers() & ^ACC_FINAL)
	} else {
		anonymousInnerClass = NewInnerClassNode(outerClass, innerClassName, ACC_PUBLIC, superClass)
	}

	anonymousInnerClass.SetAnonymous(true)
	anonymousInnerClass.SetUsingGenerics(false)
	anonymousInnerClass.PutNodeMetaData(CLASS_NAME, innerClassName)
	configureAST(anonymousInnerClass, ctx)

	v.pushClassNode(anonymousInnerClass.ClassNode)
	classBody := ctx.ClassBody().(*ClassBodyContext)
	classBody.PutNodeMetaData(CLASS_DECLARATION_CLASS_NODE, anonymousInnerClass)
	v.VisitClassBody(classBody)
	v.popClassNode()

	if v.classNodeStack.Len() == 0 {
		v.addToClassNodeList(anonymousInnerClass.ClassNode)
	}

	return anonymousInnerClass
}

func (v *ASTBuilder) VisitCreatedName(ctx *CreatedNameContext) *ClassNode {
	var classNode *ClassNode

	if ctx.QualifiedClassName() != nil {
		classNode = v.VisitQualifiedClassName(ctx.QualifiedClassName().(*QualifiedClassNameContext))
		if ctx.TypeArgumentsOrDiamond() != nil {
			classNode.SetGenericsTypes(
				v.VisitTypeArgumentsOrDiamond(ctx.TypeArgumentsOrDiamond().(*TypeArgumentsOrDiamondContext)))
			configureAST(classNode, ctx)
		}
	} else if ctx.PrimitiveType() != nil {
		classNode = configureAST(v.VisitPrimitiveType(ctx.PrimitiveType().(*PrimitiveTypeContext)), ctx)
	}

	if classNode == nil {
		panic(createParsingFailedException("Unsupported created name: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}

	classNode.AddTypeAnnotations(v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext))) // GROOVY-11178

	return classNode
}

func (v *ASTBuilder) VisitMap(ctx *MapContext) *MapExpression {
	return configureAST(
		NewMapExpressionWithEntries(v.VisitMapEntryList(ctx.MapEntryList().(*MapEntryListContext))),
		ctx)
}

func (v *ASTBuilder) VisitMapEntryList(ctx *MapEntryListContext) []*MapEntryExpression {
	if ctx == nil {
		return []*MapEntryExpression{}
	}

	return v.createMapEntryList(ctx.AllMapEntry())
}

func (v *ASTBuilder) createMapEntryList(mapEntryContextList []IMapEntryContext) []*MapEntryExpression {
	if len(mapEntryContextList) == 0 {
		return []*MapEntryExpression{}
	}

	mapEntryList := make([]*MapEntryExpression, len(mapEntryContextList))
	for i, mapEntryContext := range mapEntryContextList {
		mapEntryList[i] = v.VisitMapEntry(mapEntryContext.(*MapEntryContext))
	}

	return mapEntryList
}

func (v *ASTBuilder) VisitMapEntry(ctx *MapEntryContext) *MapEntryExpression {
	var keyExpr Expression
	valueExpr := v.Visit(ctx.Expression()).(Expression)

	if ctx.MUL() != nil {
		keyExpr = configureAST(NewSpreadMapExpression(valueExpr), ctx)
	} else if ctx.MapEntryLabel() != nil {
		keyExpr = v.VisitMapEntryLabel(ctx.MapEntryLabel().(*MapEntryLabelContext))
	} else {
		panic(createParsingFailedException("Unsupported map entry: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}

	return configureAST(
		NewMapEntryExpression(keyExpr, valueExpr),
		ctx)
}

func (v *ASTBuilder) VisitMapEntryLabel(ctx *MapEntryLabelContext) Expression {
	if ctx.Keywords() != nil {
		return configureAST(v.VisitKeywords(ctx.Keywords().(*KeywordsContext)), ctx)
	} else if ctx.Primary() != nil {
		expression := v.Visit(ctx.Primary()).(Expression)

		// if the key is variable and not inside parentheses, convert it to a constant, e.g. [a:1, b:2]
		if varExpr, ok := expression.(*VariableExpression); ok && !v.isInsideParentheses(expression) {
			expression = configureASTFromSource(
				NewConstantExpression(varExpr.GetName()),
				expression)
		}

		return configureAST(expression, ctx)
	}

	panic(createParsingFailedException("Unsupported map entry label: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitKeywords(ctx *KeywordsContext) *ConstantExpression {
	return configureAST(NewConstantExpression(ctx.GetText()), ctx)
}

func (v *ASTBuilder) VisitBuiltInType(ctx *BuiltInTypeContext) *VariableExpression {
	var text string
	if ctx.VOID() != nil {
		text = ctx.VOID().GetText()
	} else if ctx.BuiltInPrimitiveType() != nil {
		text = ctx.BuiltInPrimitiveType().GetText()
	} else {
		panic(createParsingFailedException("Unsupported built-in type: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}

	variableExpression := NewVariableExpressionWithString(text)
	variableExpression.PutNodeMetaData(IS_BUILT_IN_TYPE, true)
	return configureAST(variableExpression, ctx)
}

func (v *ASTBuilder) VisitList(ctx *ListContext) *ListExpression {
	if ctx.COMMA() != nil && ctx.ExpressionList() == nil {
		panic(createParsingFailedException("Empty list constructor should not contain any comma(,)", tokenAdapter{ctx.COMMA().GetSymbol()}))
	}

	return configureAST(
		NewListExpressionWithExpressions(
			v.VisitExpressionList(ctx.ExpressionList().(*ExpressionListContext))),
		ctx)
}

func (v *ASTBuilder) VisitExpressionList(ctx *ExpressionListContext) []Expression {
	if ctx == nil {
		return []Expression{}
	}

	return v.createExpressionList(ctx.AllExpressionListElement())
}

func (v *ASTBuilder) createExpressionList(expressionListElementContextList []IExpressionListElementContext) []Expression {
	if len(expressionListElementContextList) == 0 {
		return []Expression{}
	}

	expressions := make([]Expression, len(expressionListElementContextList))
	for i, ctx := range expressionListElementContextList {
		expressions[i] = v.VisitExpressionListElement(ctx.(*ExpressionListElementContext))
	}
	return expressions
}

func (v *ASTBuilder) VisitExpressionListElement(ctx *ExpressionListElementContext) Expression {
	expression := v.Visit(ctx.Expression()).(Expression)

	v.validateExpressionListElement(ctx, expression)

	if ctx.MUL() != nil {
		if !ctx.canSpread {
			panic(createParsingFailedException("spread operator is not allowed here", tokenAdapter{ctx.MUL().GetSymbol()}))
		}

		return configureAST(NewSpreadExpression(expression), ctx)
	}

	return configureAST(expression, ctx)
}

func (v *ASTBuilder) validateExpressionListElement(ctx *ExpressionListElementContext, expression Expression) {
	if methodCallExpr, ok := expression.(*MethodCallExpression); ok && isTrue(expression, IS_COMMAND_EXPRESSION) {
		// statements like `foo(String a)` is invalid
		methodName := methodCallExpr.GetMethodAsString()
		if methodCallExpr.IsImplicitThis() && (unicode.IsUpper(rune(methodName[0])) || isPrimitiveType(methodName)) {
			panic(createParsingFailedException("Invalid method declaration", parserRuleContextAdapter{ctx}))
		}
	}
}

func (v *ASTBuilder) VisitIntegerLiteralAlt(ctx *IntegerLiteralAltContext) *ConstantExpression {
	text := ctx.IntegerLiteral().GetText()
	var num interface{}
	//var err error
	num = ParseInteger(text)
	// TODO: implement this
	/*
		if err != nil {
			v.numberFormatError = NewTuple2(ctx, err)
		}
	*/

	constantExpression := NewConstantExpression(num)
	constantExpression.PutNodeMetaData(INTEGER_LITERAL_TEXT, text)
	constantExpression.PutNodeMetaData(IS_NUMERIC, true)
	return configureAST(constantExpression, ctx)
}

func (v *ASTBuilder) VisitFloatingPointLiteralAlt(ctx *FloatingPointLiteralAltContext) *ConstantExpression {
	text := ctx.FloatingPointLiteral().GetText()
	var num interface{}
	//var err error
	num = ParseDecimal(text)
	// TODO: implement this
	/*
		if err != nil {
			v.numberFormatError = NewTuple2(ctx, err)
		}
	*/

	constantExpression := NewConstantExpression(num)
	constantExpression.PutNodeMetaData(FLOATING_POINT_LITERAL_TEXT, text)
	constantExpression.PutNodeMetaData(IS_NUMERIC, true)
	return configureAST(constantExpression, ctx)
}

func (v *ASTBuilder) VisitBooleanLiteralAlt(ctx *BooleanLiteralAltContext) *ConstantExpression {
	return configureAST(NewConstantExpression(ctx.BooleanLiteral().GetText() == "true"), ctx)
}

func (v *ASTBuilder) VisitNullLiteralAlt(ctx *NullLiteralAltContext) *ConstantExpression {
	return configureAST(NewConstantExpression(nil), ctx)
}

func (v *ASTBuilder) VisitGstring(ctx *GstringContext) *GStringExpression {
	stringLiteralList := make([]*ConstantExpression, 0)
	begin := ctx.GStringBegin().GetText()
	beginQuotation := beginQuotation(begin)
	stringLiteralList = append(stringLiteralList, configureASTWithToken(NewConstantExpression(v.parseGStringBegin(ctx, beginQuotation)), ctx.GStringBegin().GetSymbol()))

	for _, e := range ctx.AllGStringPart() {
		stringLiteralList = append(stringLiteralList, configureASTWithToken(NewConstantExpression(v.parseGStringPart(e, beginQuotation)), e.GetSymbol()))
	}

	stringLiteralList = append(stringLiteralList, configureASTWithToken(NewConstantExpression(v.parseGStringEnd(ctx, beginQuotation)), ctx.GStringEnd().GetSymbol()))

	var values []Expression
	for _, gstringValue := range ctx.AllGstringValue() {
		values = append(values, v.VisitGstringValue(gstringValue.(*GstringValueContext)))
	}

	verbatimText := strings.Builder{}
	verbatimText.Grow(len(ctx.GetText()))
	for i, n := 0, len(stringLiteralList); i < n; i++ {
		verbatimText.WriteString(stringLiteralList[i].GetValue().(string))

		if i == len(values) {
			continue
		}

		value := values[i]
		if value == nil {
			continue
		}

		isVariableExpression := IsInstanceOf(value, (*VariableExpression)(nil))
		verbatimText.WriteString(DOLLAR_STR)
		if !isVariableExpression {
			verbatimText.WriteString("{")
		}
		verbatimText.WriteString(value.GetText())
		if !isVariableExpression {
			verbatimText.WriteString("}")
		}
	}

	return configureAST(NewGStringExpressionWithValues(verbatimText.String(), stringLiteralList, values), ctx)
}

func hasArrow(e *GstringValueContext) bool {
	return e.Closure() != nil && e.Closure().ARROW() != nil
}

func (v *ASTBuilder) parseGStringEnd(ctx *GstringContext, beginQuotation string) string {
	text := strings.Builder{}
	text.WriteString(ctx.GStringEnd().GetText())
	text.WriteString(beginQuotation)

	return v.parseStringLiteral(text.String())
}

func (v *ASTBuilder) parseGStringPart(e antlr.TerminalNode, beginQuotation string) string {
	text := strings.Builder{}
	text.Grow(len(e.GetText()))
	text.WriteString(e.GetText()[:len(e.GetText())-1]) // remove the trailing $
	text.WriteString(beginQuotation)
	text.WriteString(QUOTATION_MAP[beginQuotation])

	return v.parseStringLiteral(text.String())
}

func (v *ASTBuilder) parseGStringBegin(ctx *GstringContext, beginQuotation string) string {
	text := strings.Builder{}
	text.Grow(len(ctx.GStringBegin().GetText()))
	text.WriteString(ctx.GStringBegin().GetText()[:len(ctx.GStringBegin().GetText())-1]) // remove the trailing $
	text.WriteString(QUOTATION_MAP[beginQuotation])

	return v.parseStringLiteral(text.String())
}

func beginQuotation(text string) string {
	if strings.HasPrefix(text, TDQ_STR) {
		return TDQ_STR
	} else if strings.HasPrefix(text, DQ_STR) {
		return DQ_STR
	} else if strings.HasPrefix(text, SLASH_STR) {
		return SLASH_STR
	} else if strings.HasPrefix(text, DOLLAR_SLASH_STR) {
		return DOLLAR_SLASH_STR
	} else {
		return string(text[0])
	}
}

func (v *ASTBuilder) VisitGstringValue(ctx *GstringValueContext) Expression {
	if ctx.GstringPath() != nil {
		return configureAST(v.VisitGstringPath(ctx.GstringPath().(*GstringPathContext)), ctx)
	}

	if ctx.Closure() != nil {
		closureExpression := v.VisitClosure(ctx.Closure().(*ClosureContext))
		if !hasArrow(ctx) {
			statementList := closureExpression.GetCode().(*BlockStatement).GetStatements()
			size := len(statementList)
			if size == 1 {
				statement := statementList[0]
				if expressionStatement, ok := statement.(*ExpressionStatement); ok {
					expression := expressionStatement.GetExpression()
					if _, ok := expression.(*DeclarationExpression); !ok {
						return expression
					}
				}
			} else if size == 0 { // e.g. "${}"
				return configureAST(NewConstantExpression(nil), ctx)
			}

			return configureAST(v.createCallMethodCallExpressionWithImplicitThis(closureExpression, NewArgumentListExpression(), true), ctx)
		}

		return configureAST(closureExpression, ctx)
	}

	panic(createParsingFailedException("Unsupported gstring value: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitGstringPath(ctx *GstringPathContext) Expression {
	variableExpression := NewVariableExpressionWithString(v.VisitIdentifier(ctx.Identifier().(*IdentifierContext)))

	if len(ctx.AllGStringPathPart()) > 0 {
		var propertyExpression Expression = configureAST(variableExpression, ctx.Identifier())

		for _, part := range ctx.AllGStringPathPart() {
			constantExpression := configureASTWithToken(
				NewConstantExpression(part.GetText()[1:]),
				part.GetSymbol(),
			)
			propertyExpression = configureASTWithToken(
				NewPropertyExpressionWithProperty(propertyExpression, constantExpression),
				part.GetSymbol(),
			)
		}

		return configureAST(propertyExpression, ctx)
	}

	return configureAST(variableExpression, ctx)
}

func (v *ASTBuilder) VisitStandardLambdaExpression(ctx *StandardLambdaExpressionContext) *LambdaExpression {
	v.pushSwitchExpressionRuleContext(ctx)
	defer v.popSwitchExpressionRuleContext()

	return configureAST(v.createLambda(ctx.StandardLambdaParameters(), ctx.LambdaBody()), ctx)
}

func (v *ASTBuilder) createLambda(standardLambdaParametersContext IStandardLambdaParametersContext, lambdaBodyContext ILambdaBodyContext) *LambdaExpression {
	return NewLambdaExpression(
		v.VisitStandardLambdaParameters(standardLambdaParametersContext.(*StandardLambdaParametersContext)),
		v.VisitLambdaBody(lambdaBodyContext.(*LambdaBodyContext)),
	)
}

func (v *ASTBuilder) VisitStandardLambdaParameters(ctx *StandardLambdaParametersContext) []*Parameter {
	if ctx.VariableDeclaratorId() != nil {
		variable := v.VisitVariableDeclaratorId(ctx.VariableDeclaratorId().(*VariableDeclaratorIdContext))
		parameter := NewParameter(dynamicType(), variable.GetName())
		configureASTFromSource(parameter, variable)
		return []*Parameter{parameter}
	}

	parameters := v.VisitFormalParameters(ctx.FormalParameters().(*FormalParametersContext))
	if len(parameters) > 0 {
		return parameters
	}
	return nil
}

func (v *ASTBuilder) VisitLambdaBody(ctx *LambdaBodyContext) Statement {
	if ctx.Block() != nil {
		return configureAST(v.VisitBlock(ctx.Block().(*BlockContext)), ctx)
	}
	return configureAST(v.Visit(ctx.StatementExpression()).(Statement), ctx)
}

func (v *ASTBuilder) VisitClosure(ctx *ClosureContext) *ClosureExpression {
	v.pushSwitchExpressionRuleContext(ctx)
	v.visitingClosureCount++
	defer func() {
		v.popSwitchExpressionRuleContext()
		v.visitingClosureCount--
	}()

	var parameters []*Parameter
	if ctx.FormalParameterList() != nil {
		parameters = v.VisitFormalParameterList(ctx.FormalParameterList().(*FormalParameterListContext))
	}

	code := v.VisitBlockStatementsOpt(ctx.BlockStatementsOpt().(*BlockStatementsOptContext))
	if ctx.ARROW() == nil {
		parameters = []*Parameter{}
		if code.IsEmpty() {
			configureAST(code, ctx)
		}
	}

	return configureAST(NewClosureExpression(parameters, code), ctx)
}

func (v *ASTBuilder) VisitFormalParameters(ctx *FormalParametersContext) []*Parameter {
	if ctx == nil {
		return []*Parameter{}
	}

	return v.VisitFormalParameterList(ctx.FormalParameterList().(*FormalParameterListContext))
}

func (v *ASTBuilder) VisitFormalParameterList(ctx *FormalParameterListContext) []*Parameter {
	if ctx == nil {
		return []*Parameter{}
	}

	var parameterList []*Parameter

	if ctx.ThisFormalParameter() != nil {
		parameterList = append(parameterList, v.VisitThisFormalParameter(ctx.ThisFormalParameter().(*ThisFormalParameterContext)))
	}

	formalParameterList := ctx.AllFormalParameter()
	if len(formalParameterList) > 0 {
		v.validateVarArgParameter(formalParameterList)

		for _, fp := range formalParameterList {
			parameterList = append(parameterList, v.VisitFormalParameter(fp.(*FormalParameterContext)))
		}
	}

	v.validateParameterList(parameterList)

	return parameterList
}

func (v *ASTBuilder) validateVarArgParameter(formalParameterList []IFormalParameterContext) {
	for i := 0; i < len(formalParameterList)-1; i++ {
		formalParameterContext := formalParameterList[i].(*FormalParameterContext)
		if formalParameterContext.ELLIPSIS() != nil {
			panic(createParsingFailedException("The var-arg parameter must be the last parameter", parserRuleContextAdapter{formalParameterContext}))
		}
	}
}

func (v *ASTBuilder) validateParameterList(parameterList []*Parameter) {
	for i := len(parameterList) - 1; i >= 0; i-- {
		parameter := parameterList[i]
		name := parameter.GetName()
		if name == "_" {
			continue // check this later
		}
		for _, otherParameter := range parameterList {
			if otherParameter == parameter {
				continue
			}
			if otherParameter.GetName() == name {
				panic(createParsingFailedException("Duplicated parameter '"+name+"' found.", astNodeAdapter{parameter}))
			}
		}
	}
}

func (v *ASTBuilder) VisitFormalParameter(ctx *FormalParameterContext) *Parameter {
	return v.processFormalParameter(ctx, ctx.VariableModifiersOpt().(*VariableModifiersOptContext), ctx.Type_().(*TypeContext), ctx.ELLIPSIS(), ctx.VariableDeclaratorId().(*VariableDeclaratorIdContext), ctx.Expression().(IExpressionContext))
}

func (v *ASTBuilder) VisitThisFormalParameter(ctx *ThisFormalParameterContext) *Parameter {
	return configureAST(NewParameter(v.VisitType(ctx.Type_().(*TypeContext)), THIS_STR), ctx)
}

func (v *ASTBuilder) VisitClassOrInterfaceModifiersOpt(ctx *ClassOrInterfaceModifiersOptContext) []*ModifierNode {
	if ctx.ClassOrInterfaceModifiers() != nil {
		return v.VisitClassOrInterfaceModifiers(ctx.ClassOrInterfaceModifiers().(*ClassOrInterfaceModifiersContext))
	}

	return []*ModifierNode{}
}

func (v *ASTBuilder) VisitClassOrInterfaceModifiers(ctx *ClassOrInterfaceModifiersContext) []*ModifierNode {
	modifiers := []*ModifierNode{}
	for _, modifier := range ctx.AllClassOrInterfaceModifier() {
		modifiers = append(modifiers, v.VisitClassOrInterfaceModifier(modifier.(*ClassOrInterfaceModifierContext)))
	}
	return modifiers
}

func (v *ASTBuilder) VisitClassOrInterfaceModifier(ctx *ClassOrInterfaceModifierContext) *ModifierNode {
	if ctx.Annotation() != nil {
		return configureAST(NewModifierNodeWithAnnotation(v.VisitAnnotation(ctx.Annotation().(*AnnotationContext)), ctx.GetText()), ctx)
	}

	if ctx.m != nil {
		return configureAST(NewModifierNodeWithText(ctx.m.GetTokenType(), ctx.GetText()), ctx)
	}

	panic(createParsingFailedException("Unsupported class or interface modifier: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitModifier(ctx *ModifierContext) *ModifierNode {
	if ctx.ClassOrInterfaceModifier() != nil {
		return configureAST(v.VisitClassOrInterfaceModifier(ctx.ClassOrInterfaceModifier().(*ClassOrInterfaceModifierContext)), ctx)
	}

	if ctx.m != nil {
		return configureAST(NewModifierNodeWithText(ctx.m.GetTokenType(), ctx.GetText()), ctx)
	}

	panic(createParsingFailedException("Unsupported modifier: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitModifiers(ctx *ModifiersContext) []*ModifierNode {
	modifiers := []*ModifierNode{}
	for _, modifier := range ctx.AllModifier() {
		modifiers = append(modifiers, v.VisitModifier(modifier.(*ModifierContext)))
	}
	return modifiers
}

func (v *ASTBuilder) VisitModifiersOpt(ctx *ModifiersOptContext) []*ModifierNode {
	if ctx.Modifiers() != nil {
		return v.VisitModifiers(ctx.Modifiers().(*ModifiersContext))
	}

	return []*ModifierNode{}
}

func (v *ASTBuilder) VisitVariableModifier(ctx *VariableModifierContext) *ModifierNode {
	if ctx.Annotation() != nil {
		return configureAST(NewModifierNodeWithAnnotation(v.VisitAnnotation(ctx.Annotation().(*AnnotationContext)), ctx.GetText()), ctx)
	}

	if ctx.m != nil {
		return configureAST(NewModifierNodeWithText(ctx.m.GetTokenType(), ctx.GetText()), ctx)
	}

	panic(createParsingFailedException("Unsupported variable modifier", parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitVariableModifiersOpt(ctx *VariableModifiersOptContext) []*ModifierNode {
	if ctx.VariableModifiers() != nil {
		return v.VisitVariableModifiers(ctx.VariableModifiers().(*VariableModifiersContext))
	}

	return []*ModifierNode{}
}

func (v *ASTBuilder) VisitVariableModifiers(ctx *VariableModifiersContext) []*ModifierNode {
	modifiers := make([]*ModifierNode, 0, len(ctx.AllVariableModifier()))
	for _, modifier := range ctx.AllVariableModifier() {
		modifiers = append(modifiers, v.VisitVariableModifier(modifier.(*VariableModifierContext)))
	}
	return modifiers
}

func (v *ASTBuilder) VisitEmptyDims(ctx *EmptyDimsContext) [][]*AnnotationNode {
	dimList := make([][]*AnnotationNode, 0, len(ctx.AllAnnotationsOpt()))
	for _, annotationsOpt := range ctx.AllAnnotationsOpt() {
		dimList = append(dimList, v.VisitAnnotationsOpt(annotationsOpt.(*AnnotationsOptContext)))
	}

	// Reverse the dimList
	for i, j := 0, len(dimList)-1; i < j; i, j = i+1, j-1 {
		dimList[i], dimList[j] = dimList[j], dimList[i]
	}

	return dimList
}

func (v *ASTBuilder) VisitEmptyDimsOpt(ctx *EmptyDimsOptContext) [][]*AnnotationNode {
	if ctx.EmptyDims() == nil {
		return [][]*AnnotationNode{}
	}

	return v.VisitEmptyDims(ctx.EmptyDims().(*EmptyDimsContext))
}

func (v *ASTBuilder) VisitType(ctx *TypeContext) *ClassNode {
	if ctx == nil {
		return dynamicType()
	}

	var classNode *ClassNode

	if ctx.GeneralClassOrInterfaceType() != nil {
		if isTrue(ctx, IS_INSIDE_INSTANCEOF_EXPR) {
			ctx.GeneralClassOrInterfaceType().(*GeneralClassOrInterfaceTypeContext).PutNodeMetaData(IS_INSIDE_INSTANCEOF_EXPR, true)
		}
		classNode = v.VisitClassOrInterfaceType(ctx.GeneralClassOrInterfaceType().(*GeneralClassOrInterfaceTypeContext))
	} else if ctx.PrimitiveType() != nil {
		classNode = v.VisitPrimitiveType(ctx.PrimitiveType().(*PrimitiveTypeContext))
	}

	if classNode == nil {
		if ctx.GetText() == VOID_STR {
			panic(createParsingFailedException("void is not allowed here", parserRuleContextAdapter{ctx}))
		}
		panic(createParsingFailedException("Unsupported type: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
	}

	classNode.AddTypeAnnotations(v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext)))

	dimList := v.VisitEmptyDimsOpt(ctx.EmptyDimsOpt().(*EmptyDimsOptContext))
	if len(dimList) > 0 {
		classNode = v.createArrayTypeWithAnnotations(classNode, dimList)
	}

	return configureAST(classNode, ctx)
}

func (v *ASTBuilder) VisitClassOrInterfaceType(ctx *GeneralClassOrInterfaceTypeContext) *ClassNode {
	var classNode *ClassNode
	if ctx.QualifiedClassName() != nil {
		if isTrue(ctx, IS_INSIDE_INSTANCEOF_EXPR) {
			ctx.QualifiedClassName().(*QualifiedClassNameContext).PutNodeMetaData(IS_INSIDE_INSTANCEOF_EXPR, true)
		}
		classNode = v.VisitQualifiedClassName(ctx.QualifiedClassName().(*QualifiedClassNameContext))
	} else {
		// TODO: implement this
		/*
			if isTrue(ctx, IS_INSIDE_INSTANCEOF_EXPR) {
				ctx.QualifiedStandardClassName().(*QualifiedStandardClassNameContext).PutNodeMetaData(IS_INSIDE_INSTANCEOF_EXPR, true)
			}
			classNode = v.VisitQualifiedStandardClassName(ctx.QualifiedStandardClassName().(*QualifiedStandardClassNameContext))
		*/
	}

	if ctx.TypeArguments() != nil {
		classNode.SetGenericsTypes(v.VisitTypeArguments(ctx.TypeArguments().(*TypeArgumentsContext)))
	}

	return configureAST(classNode, ctx)
}

func (v *ASTBuilder) VisitTypeArgumentsOrDiamond(ctx *TypeArgumentsOrDiamondContext) []*GenericsType {
	if ctx.TypeArguments() != nil {
		return v.VisitTypeArguments(ctx.TypeArguments().(*TypeArgumentsContext))
	}

	if ctx.LT() != nil { // e.g. <>
		return []*GenericsType{}
	}

	panic(createParsingFailedException("Unsupported type arguments or diamond: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitTypeArguments(ctx *TypeArgumentsContext) []*GenericsType {
	typeArguments := make([]*GenericsType, len(ctx.AllTypeArgument()))
	for i, typeArg := range ctx.AllTypeArgument() {
		typeArguments[i] = v.VisitTypeArgument(typeArg.(*TypeArgumentContext))
	}
	return typeArguments
}

func (v *ASTBuilder) VisitTypeArgument(ctx *TypeArgumentContext) *GenericsType {
	if ctx.QUESTION() != nil {
		baseType := configureASTWithToken(MakeWithoutCaching(QUESTION_STR), ctx.QUESTION().GetSymbol())
		baseType.AddTypeAnnotations(v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext)))

		if ctx.Type_() == nil {
			genericsType := NewGenericsTypeWithBasicType(baseType)
			genericsType.SetWildcard(true)
			return configureAST(genericsType, ctx)
		}

		var upperBounds []*ClassNode
		var lowerBound *ClassNode

		classNode := v.VisitType(ctx.Type_().(*TypeContext))
		if ctx.EXTENDS() != nil {
			upperBounds = []*ClassNode{classNode}
		} else if ctx.SUPER() != nil {
			lowerBound = classNode
		}

		genericsType := NewGenericsType(baseType, upperBounds, lowerBound)
		genericsType.SetWildcard(true)

		return configureAST(genericsType, ctx)
	} else if ctx.Type_() != nil {
		baseType := v.VisitType(ctx.Type_().(*TypeContext))
		return configureAST(v.createGenericsType(baseType), ctx)
	}

	panic(createParsingFailedException("Unsupported type argument: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitPrimitiveType(ctx *PrimitiveTypeContext) *ClassNode {
	return configureAST(MakeFromString(ctx.GetText()).GetPlainNodeReferenceHelper(false), ctx)
}

func (v *ASTBuilder) VisitVariableDeclaratorId(ctx *VariableDeclaratorIdContext) *VariableExpression {
	return configureAST(NewVariableExpressionWithString(v.VisitIdentifier(ctx.Identifier().(*IdentifierContext))), ctx)
}

func (v *ASTBuilder) VisitVariableNames(ctx *VariableNamesContext) *TupleExpression {
	var expressions []Expression
	for _, varDeclId := range ctx.AllVariableDeclaratorId() {
		expressions = append(expressions, v.VisitVariableDeclaratorId(varDeclId.(*VariableDeclaratorIdContext)))
	}
	return configureAST(NewTupleExpressionWithExpressions(expressions...), ctx)
}

func (v *ASTBuilder) VisitClosureOrLambdaExpression(ctx *ClosureOrLambdaExpressionContext) *ClosureExpression {
	if ctx.Closure() != nil {
		return configureAST(v.VisitClosure(ctx.Closure().(*ClosureContext)), ctx)
	} else if ctx.LambdaExpression() != nil {
		// TODO: implement this
		panic("LambdaExpression is not implemented")
		//return configureAST(v.VisitStandardLambdaExpression(ctx.LambdaExpression().(*LambdaExpressionContext)), ctx)
	}

	panic(createParsingFailedException("The node is not expected here"+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitBlockStatementsOpt(ctx *BlockStatementsOptContext) *BlockStatement {
	if ctx.BlockStatements() != nil {
		return configureAST(v.VisitBlockStatements(ctx.BlockStatements().(*BlockStatementsContext)), ctx)
	}

	return configureAST(v.createBlockStatement(), ctx)
}

func (v *ASTBuilder) VisitBlockStatements(ctx *BlockStatementsContext) *BlockStatement {
	var statements []Statement
	for _, stmt := range ctx.AllBlockStatement() {
		if s := v.VisitBlockStatement(stmt.(*BlockStatementContext)); s != nil {
			statements = append(statements, s)
		}
	}
	return configureAST(v.createBlockStatement(statements...), ctx)
}

func (v *ASTBuilder) VisitBlockStatement(ctx *BlockStatementContext) Statement {
	if ctx.LocalVariableDeclaration() != nil {
		return configureAST(v.VisitLocalVariableDeclaration(ctx.LocalVariableDeclaration().(*LocalVariableDeclarationContext)), ctx)
	}

	if ctx.Statement() != nil {
		astNode := v.Visit(ctx.Statement())

		if astNode == nil {
			return nil
		}

		switch node := astNode.(type) {
		case Statement:
			return node
		case *MethodNode:
			panic(createParsingFailedException("Method definition not expected here", parserRuleContextAdapter{ctx}))
		case *ImportNode:
			panic(createParsingFailedException("Import statement not expected here", parserRuleContextAdapter{ctx}))
		default:
			panic(createParsingFailedException(fmt.Sprintf("The statement(%T) not expected here", astNode), parserRuleContextAdapter{ctx}))
		}
	}

	panic(createParsingFailedException("Unsupported block statement: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitAnnotationsOpt(ctx *AnnotationsOptContext) []*AnnotationNode {
	if ctx == nil {
		return []*AnnotationNode{}
	}

	var annotations []*AnnotationNode
	for _, ann := range ctx.AllAnnotation() {
		annotations = append(annotations, v.VisitAnnotation(ann.(*AnnotationContext)))
	}
	return annotations
}

func (v *ASTBuilder) VisitAnnotation(ctx *AnnotationContext) *AnnotationNode {
	annotationName := v.VisitAnnotationName(ctx.AnnotationName().(*AnnotationNameContext))
	annotationNode := NewAnnotationNode(MakeFromString(annotationName))
	annotationElementValues := v.VisitElementValues(ctx.ElementValues().(*ElementValuesContext))

	for _, e := range annotationElementValues {
		annotationNode.AddMember(e.V1, e.V2)
	}
	configureAST(annotationNode.GetClassNode(), ctx.AnnotationName())
	return configureAST(annotationNode, ctx)
}

func (v *ASTBuilder) VisitElementValues(ctx *ElementValuesContext) []Tuple2[string, Expression] {
	if ctx == nil {
		return []Tuple2[string, Expression]{}
	}

	var annotationElementValues []Tuple2[string, Expression]

	if ctx.ElementValuePairs() != nil {
		for key, value := range v.VisitElementValuePairs(ctx.ElementValuePairs().(*ElementValuePairsContext)) {
			annotationElementValues = append(annotationElementValues, NewTuple2(key, value))
		}
	} else if ctx.ElementValue() != nil {
		annotationElementValues = append(annotationElementValues, NewTuple2(VALUE_STR, v.VisitElementValue(ctx.ElementValue().(*ElementValueContext))))
	}

	return annotationElementValues
}

func (v *ASTBuilder) VisitAnnotationName(ctx *AnnotationNameContext) string {
	return v.VisitQualifiedClassName(ctx.QualifiedClassName().(*QualifiedClassNameContext)).GetName()
}

func (v *ASTBuilder) VisitElementValuePairs(ctx *ElementValuePairsContext) map[string]Expression {
	elementValuePairs := make(map[string]Expression)
	for _, pair := range ctx.AllElementValuePair() {
		t := v.VisitElementValuePair(pair.(*ElementValuePairContext))
		elementValuePairs[t.V1] = t.V2
	}
	return elementValuePairs
}

func (v *ASTBuilder) VisitElementValuePair(ctx *ElementValuePairContext) Tuple2[string, Expression] {
	return NewTuple2(ctx.ElementValuePairName().GetText(), v.VisitElementValue(ctx.ElementValue().(*ElementValueContext)))
}

func (v *ASTBuilder) VisitElementValue(ctx *ElementValueContext) Expression {
	if ctx.Expression() != nil {
		return configureAST(v.Visit(ctx.Expression()).(Expression), ctx)
	}

	if ctx.Annotation() != nil {
		return configureAST(NewAnnotationConstantExpression(v.VisitAnnotation(ctx.Annotation().(*AnnotationContext))), ctx)
	}

	if ctx.ElementValueArrayInitializer() != nil {
		return configureAST(v.VisitElementValueArrayInitializer(ctx.ElementValueArrayInitializer().(*ElementValueArrayInitializerContext)), ctx)
	}

	panic(createParsingFailedException("Unsupported element value: "+ctx.GetText(), parserRuleContextAdapter{ctx}))
}

func (v *ASTBuilder) VisitElementValueArrayInitializer(ctx *ElementValueArrayInitializerContext) *ListExpression {
	var elementValues []Expression
	for _, elementValue := range ctx.AllElementValue() {
		elementValues = append(elementValues, v.VisitElementValue(elementValue.(*ElementValueContext)))
	}
	return configureAST(NewListExpressionWithExpressions(elementValues), ctx)
}

func (v *ASTBuilder) VisitClassName(ctx *ClassNameContext) string {
	return ctx.GetText()
}

func (v *ASTBuilder) VisitIdentifier(ctx *IdentifierContext) string {
	return ctx.GetText()
}

func (v *ASTBuilder) VisitQualifiedName(ctx *QualifiedNameContext) string {
	var elements []string
	for _, element := range ctx.AllQualifiedNameElement() {
		elements = append(elements, element.GetText())
	}
	return strings.Join(elements, DOT_STR)
}

func (v *ASTBuilder) VisitAnnotatedQualifiedClassName(ctx *AnnotatedQualifiedClassNameContext) *ClassNode {
	classNode := v.VisitQualifiedClassName(ctx.QualifiedClassName().(*QualifiedClassNameContext))

	classNode.AddTypeAnnotations(v.VisitAnnotationsOpt(ctx.AnnotationsOpt().(*AnnotationsOptContext)))

	return classNode
}

func (v *ASTBuilder) VisitQualifiedClassNameList(ctx *QualifiedClassNameListContext) []*ClassNode {
	if ctx == nil {
		return []*ClassNode{}
	}

	var classNodes []*ClassNode
	for _, annotatedQualifiedClassName := range ctx.AllAnnotatedQualifiedClassName() {
		classNodes = append(classNodes, v.VisitAnnotatedQualifiedClassName(annotatedQualifiedClassName.(*AnnotatedQualifiedClassNameContext)))
	}
	return classNodes
}

func (v *ASTBuilder) VisitQualifiedClassName(ctx *QualifiedClassNameContext) *ClassNode {
	return v.createClassNode(ctx.GroovyParserRuleContext)
}

func (v *ASTBuilder) VisitQualifiedStandardClassName(ctx *QualifiedStandardClassNameContext) *ClassNode {
	return v.createClassNode(ctx.GroovyParserRuleContext)
}

func (v *ASTBuilder) createArrayTypeWithAnnotations(elementType *ClassNode, dimAnnotationsList [][]*AnnotationNode) *ClassNode {
	arrayType := elementType
	for i := len(dimAnnotationsList) - 1; i >= 0; i-- {
		arrayType = v.createArrayType(arrayType)
		arrayType.AddAnnotations(dimAnnotationsList[i])
	}
	return arrayType
}

func (v *ASTBuilder) createArrayType(elementType *ClassNode) *ClassNode {
	if IsPrimitiveVoid(elementType) {
		panic(createParsingFailedException("void[] is an invalid type", astNodeAdapter{elementType}))
	}
	return elementType.MakeArray()
}

func (v *ASTBuilder) createClassNode(ctx *GroovyParserRuleContext) *ClassNode {
	result := MakeFromString(ctx.GetText())
	if isTrue(ctx, IS_INSIDE_INSTANCEOF_EXPR) {
		// type in the "instanceof" expression shouldn't have redirect
	} else {
		result = v.proxyClassNode(result)
	}
	return configureAST(result, ctx)
}

func (v *ASTBuilder) proxyClassNode(classNode *ClassNode) *ClassNode {
	if !classNode.IsUsingGenerics() {
		return classNode
	}

	cn := MakeWithoutCaching(classNode.GetName())
	cn.SetRedirect(classNode)
	return cn
}

func (v *ASTBuilder) createMethodCallExpression(baseExpr Expression, arguments Expression) *MethodCallExpression {
	var methodCallExpression *MethodCallExpression

	if propertyExpression, ok := baseExpr.(*PropertyExpression); ok {
		// Case for property expressions
		methodCallExpression = NewMethodCallExpression(
			propertyExpression.GetObjectExpression(),
			propertyExpression.GetProperty(),
			arguments,
		)

		methodCallExpression.SetImplicitThis(false)
		methodCallExpression.SetSafe(propertyExpression.IsSafe())
		methodCallExpression.SetSpreadSafe(propertyExpression.IsSpreadSafe())

		// method call obj*.m(): "safe"(false) and "spreadSafe"(true)
		// property access obj*.p: "safe"(true) and "spreadSafe"(true)
		// so we have to reset safe here.
		if propertyExpression.IsSpreadSafe() {
			methodCallExpression.SetSafe(false)
		}

		// if the generics types metadata is not empty, it is a generic method call, e.g. obj.<Integer>a(1, 2)
		methodCallExpression.SetGenericsTypes(
			propertyExpression.GetNodeMetaData(PATH_EXPRESSION_BASE_EXPR_GENERICS_TYPES).([]*GenericsType),
		)
	} else {
		// Case for other expressions (e.g., m(1, 2) or m 1, 2)
		thisExpr := NewVariableExpressionWithString("this")
		thisExpr.SetColumnNumber(baseExpr.GetColumnNumber())
		thisExpr.SetLineNumber(baseExpr.GetLineNumber())

		var method Expression
		if _, ok := baseExpr.(*VariableExpression); ok {
			method = v.createConstantExpression(baseExpr)
		} else {
			method = baseExpr
		}

		methodCallExpression = NewMethodCallExpression(thisExpr, method, arguments)
	}

	return methodCallExpression
}

func (v *ASTBuilder) processFormalParameter(ctx *FormalParameterContext, variableModifiersOptContext *VariableModifiersOptContext, typeContext *TypeContext, ellipsis antlr.TerminalNode, variableDeclaratorIdContext *VariableDeclaratorIdContext, expressionContext IExpressionContext) *Parameter {
	classNode := v.VisitType(typeContext)

	if ellipsis != nil {
		classNode = v.createArrayType(classNode)
		if typeContext == nil {
			configureASTWithToken(classNode, ellipsis.GetSymbol())
		} else {
			configureASTWithInitialStop(classNode, typeContext, configureASTWithToken(NewConstantExpression("..."), ellipsis.GetSymbol()))
		}
	}

	modifierManager := NewModifierManager(v, v.VisitVariableModifiersOpt(variableModifiersOptContext))
	parameter := modifierManager.ProcessParameter(
		configureAST(
			NewParameter(
				classNode,
				v.VisitVariableDeclaratorId(variableDeclaratorIdContext).GetName(),
			),
			ctx,
		),
	)
	parameter.PutNodeMetaData(PARAMETER_MODIFIER_MANAGER, modifierManager)
	parameter.PutNodeMetaData(PARAMETER_CONTEXT, ctx)

	if expressionContext != nil {
		parameter.SetInitialExpression(v.Visit(expressionContext).(Expression))
	}

	return parameter
}

func (v *ASTBuilder) createPathExpression(primaryExpr Expression, pathElementContextList []IPathElementContext) Expression {
	result := primaryExpr

	for _, e := range pathElementContextList {
		pathElementContext := e.(*PathElementContext)
		pathElementContext.SetNodeMetaData(PATH_EXPRESSION_BASE_EXPR, result)
		expression := v.VisitPathElement(pathElementContext).(Expression)

		if isTrue(result, PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN) {
			expression.SetNodeMetaData(PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN, true)
		}

		result = expression
	}

	return result
}

func (v *ASTBuilder) createGenericsType(classNode *ClassNode) *GenericsType {
	genericsType := NewGenericsTypeWithBasicType(classNode)
	return configureASTFromSource(genericsType, classNode)
}

func (v *ASTBuilder) createConstantExpression(expression Expression) *ConstantExpression {
	if constantExpr, ok := expression.(*ConstantExpression); ok {
		return constantExpr
	}

	return configureASTFromSource(NewConstantExpression(expression.GetText()), expression)
}

func (v *ASTBuilder) createBinaryExpressionHelper(left, right *ExpressionContext, op antlr.Token) *BinaryExpression {
	return NewBinaryExpression(
		v.Visit(left).(Expression),
		v.createGroovyToken(op),
		v.Visit(right).(Expression),
	)
}

func (v *ASTBuilder) createBinaryExpression(left, right *ExpressionContext, op antlr.Token, ctx *ExpressionContext) *BinaryExpression {
	binaryExpression := v.createBinaryExpressionHelper(left, right, op)

	if ctx != nil {
		if isTrue(ctx, IS_INSIDE_CONDITIONAL_EXPRESSION) {
			return configureASTWithToken(binaryExpression, op)
		}
		return configureAST(binaryExpression, ctx)
	}

	return binaryExpression
}

func (v *ASTBuilder) unpackStatement(statement Statement) Statement {
	if declarationListStatement, ok := statement.(*DeclarationListStatement); ok {
		expressionStatementList := declarationListStatement.GetDeclarationStatements()

		if len(expressionStatementList) == 1 {
			return expressionStatementList[0]
		}

		return configureASTFromSource(v.createBlockStatement(statement), statement) // if DeclarationListStatement contains more than 1 declarations, maybe it's better to create a block to hold them
	}

	return statement
}

func (v *ASTBuilder) createBlockStatement(statements ...Statement) *BlockStatement {
	return v.createBlockStatementFromList(statements)
}

func (v *ASTBuilder) createBlockStatementFromList(statementList []Statement) *BlockStatement {
	return v.appendStatementsToBlockStatement(NewBlockStatement(), statementList...)
}

func (v *ASTBuilder) appendStatementsToBlockStatement(bs *BlockStatement, statements ...Statement) *BlockStatement {
	return v.appendStatementsToBlockStatementFromList(bs, statements)
}

func (v *ASTBuilder) appendStatementsToBlockStatementFromList(bs *BlockStatement, statementList []Statement) *BlockStatement {
	for _, e := range statementList {
		if declarationListStmt, ok := e.(*DeclarationListStatement); ok {
			for _, decl := range declarationListStmt.GetDeclarationStatements() {
				bs.AddStatement(decl)
			}
		} else {
			bs.AddStatement(e)
		}
	}
	return bs
}

func (v *ASTBuilder) isAnnotationDeclaration(classNode *ClassNode) bool {
	return classNode != nil && classNode.IsAnnotationDefinition()
}

func (v *ASTBuilder) isSyntheticPublic(isAnnotationDeclaration, isAnonymousInnerEnumDeclaration, hasReturnType bool, modifierManager *ModifierManager) bool {
	if modifierManager.ContainsVisibilityModifier() {
		return false
	}

	if isAnnotationDeclaration {
		return true
	}

	if hasReturnType && (modifierManager.ContainsAny(GroovyParserDEF, GroovyParserVAR)) {
		return true
	}

	if !hasReturnType || modifierManager.ContainsNonVisibilityModifier() || modifierManager.ContainsAnnotations() {
		return true
	}

	return isAnonymousInnerEnumDeclaration
}

// the mixins of interface and annotation should be nil
func (v *ASTBuilder) hackMixins(classNode *ClassNode) {
	classNode.SetMixins(nil)
}

var TYPE_DEFAULT_VALUE_MAP = map[*ClassNode]interface{}{
	INT_TYPE:     0,
	LONG_TYPE:    int64(0),
	DOUBLE_TYPE:  0.0,
	FLOAT_TYPE:   float32(0.0),
	SHORT_TYPE:   int16(0),
	BYTE_TYPE:    int8(0),
	CHAR_TYPE:    rune(0),
	BOOLEAN_TYPE: false,
}

func (v *ASTBuilder) findDefaultValueByType(t *ClassNode) interface{} {
	return TYPE_DEFAULT_VALUE_MAP[t]
}

func (v *ASTBuilder) isPackageInfoDeclaration() bool {
	name := v.sourceUnitName
	return name != "" && strings.HasSuffix(name, PACKAGE_INFO_FILE_NAME)
}

func (v *ASTBuilder) isBlankScript() bool {
	return v.moduleNode.GetStatementBlock().IsEmpty() && len(v.moduleNode.GetMethods()) == 0 && len(v.moduleNode.GetClasses()) == 0
}

func (v *ASTBuilder) isInsideParentheses(nodeMetaDataHandler NodeMetaDataHandler) bool {
	insideParenLevel := nodeMetaDataHandler.GetNodeMetaData(INSIDE_PARENTHESES_LEVEL)
	if insideParenLevel == nil {
		return false
	}
	return insideParenLevel.(int) > 0
}

func (v *ASTBuilder) isBuiltInType(expression Expression) bool {
	if variableExpr, ok := expression.(*VariableExpression); ok {
		return isTrue(variableExpr, IS_BUILT_IN_TYPE)
	}
	return false
}

func (v *ASTBuilder) createGroovyTokenByType(token antlr.Token, tokenType int) *Token {
	if token == nil {
		panic("token should not be nil")
	}
	return NewToken(tokenType, token.GetText(), token.GetLine(), token.GetColumn())
}

func (v *ASTBuilder) createGroovyToken(token antlr.Token) *Token {
	return v.createGroovyTokenWithCardinality(token, 1)
}

func (v *ASTBuilder) createGroovyTokenWithCardinality(token antlr.Token, cardinality int) *Token {
	tokenText := token.GetText()
	tokenType := token.GetTokenType()
	text := tokenText
	if cardinality != 1 {
		text = strings.Repeat(tokenText, cardinality)
	}

	var finalTokenType int
	switch tokenType {
	case GroovyParserRANGE_EXCLUSIVE_FULL, GroovyParserRANGE_EXCLUSIVE_LEFT, GroovyParserRANGE_EXCLUSIVE_RIGHT, GroovyParserRANGE_INCLUSIVE:
		finalTokenType = RANGE_OPERATOR
	case GroovyParserSAFE_INDEX:
		finalTokenType = LEFT_SQUARE_BRACKET
	default:
		finalTokenType = Lookup(text, ANY)
	}

	return NewToken(finalTokenType, text, token.GetLine(), token.GetColumn()+1)
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

func isTrue(obj NodeMetaDataHandler, key string) bool {
	value := obj.GetNodeMetaData(key)
	boolValue, ok := value.(bool)
	return ok && boolValue
}

func isPrimitiveType(name string) bool {
	switch name {
	case "int", "void", "boolean", "byte", "char", "short", "double", "float", "long":
		return true
	default:
		return false
	}
}

// Helper function to check if a ClassNode represents void
func isPrimitiveVoid(classNode *ClassNode) bool {
	// Assuming ClassHelper.isPrimitiveVoid is implemented similarly in Go
	return classNode == VOID_TYPE
}

func (v *ASTBuilder) createArrayTypeAnnotations(elementType *ClassNode, dimAnnotationsList [][]*AnnotationNode, ctx antlr.ParserRuleContext) *ClassNode {
	arrayType := elementType

	for i := len(dimAnnotationsList) - 1; i >= 0; i-- {
		arrayType = v.createArrayType(arrayType)
		arrayType.AddAnnotations(dimAnnotationsList[i])
	}

	return arrayType
}
