package parser

import (
	"fmt"
	"strings"

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

// Function to create a CompilationFailedException
func createParsingFailedException(msg string, ctx antlr.ParserRuleContext) *SyntaxException {
	start := ctx.GetStart()
	stop := ctx.GetStop()

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
	moduleNode                   *ModuleNode
	classNodeList                []*ClassNode
	numberFormatError            *NumberFormatError
	sourceUnitName               string
	visitingAssertStatementCount int
}

// NewASTBuilder creates and initializes a new ASTBuilder instance
func NewASTBuilder(sourceUnitName string) *ASTBuilder {
	builder := &ASTBuilder{
		moduleNode:        NewModuleNode(), // Assuming you have a NewModuleNode function
		classNodeList:     make([]*ClassNode, 0),
		numberFormatError: nil,
		sourceUnitName:    sourceUnitName,
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

func (builder *ASTBuilder) VisitBlockStatement(ctx *BlockStatementContext) interface{} {
	panic("FOO")
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

func (v *ASTBuilder) VisitPackageDeclaration(ctx *PackageDeclarationContext) interface{} {
	return nil
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

func (v *ASTBuilder) VisitAssignmentExprAlt(ctx *AssignmentExprAltContext) interface{} {
	leftExpr := v.Visit(ctx.GetLeft()).(Expression)

	if varExpr, ok := leftExpr.(*VariableExpression); ok && v.isInsideParentheses(varExpr) {
		if varExpr.GetNodeMetaData(INSIDE_PARENTHESES_LEVEL).(int) > 1 {
			panic(createParsingFailedException("Nested parenthesis is not allowed in multiple assignment, e.g. ((a)) = b", ctx))
		}

		// Create a slice with a single Expression element
		expressions := []Expression{varExpr}
		tupleExpr := NewTupleExpressionWithExpressions(expressions...)

		return configureAST(
			NewBinaryExpression(
				configureAST(tupleExpr, ctx.GetLeft()),
				v.createGroovyToken(ctx.GetOp()),
				v.Visit(ctx.GetRight()).(Expression),
			),
			ctx,
		)
	}

	isValidLHS := func(expr Expression) bool {
		switch e := expr.(type) {
		case *VariableExpression:
			return !v.isInsideParentheses(e)
		case *PropertyExpression:
			return true
		case *BinaryExpression:
			return e.GetOperation().GetType() == LEFT_SQUARE_BRACKET
		default:
			return false
		}
	}

	if !isValidLHS(leftExpr) {
		panic(createParsingFailedException("The LHS of an assignment should be a variable or a field accessing expression", ctx))
	}

	return configureAST(
		NewBinaryExpression(
			leftExpr,
			v.createGroovyToken(ctx.GetOp()),
			v.Visit(ctx.GetRight()).(Expression),
		),
		ctx,
	)
}

func (builder *ASTBuilder) isInsideParentheses(nodeMetaDataHandler NodeMetaDataHandler) bool {
	insideParenLevel := nodeMetaDataHandler.GetNodeMetaData(INSIDE_PARENTHESES_LEVEL)
	if insideParenLevel == nil {
		return false
	}

	if level, ok := insideParenLevel.(int); ok {
		return level > 0
	}

	return false
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

	var newTokenType int
	switch tokenType {
	case GroovyParserRANGE_EXCLUSIVE_FULL, GroovyParserRANGE_EXCLUSIVE_LEFT, GroovyParserRANGE_EXCLUSIVE_RIGHT, GroovyParserRANGE_INCLUSIVE:
		newTokenType = RANGE_OPERATOR
	case GroovyParserSAFE_INDEX:
		newTokenType = LEFT_SQUARE_BRACKET
	default:
		newTokenType = Lookup(text, ANY)
	}

	return NewToken(
		newTokenType,
		text,
		token.GetLine(),
		token.GetColumn()+1,
	)
}

func (v *ASTBuilder) VisitPostfixExpression(ctx *PostfixExpressionContext) interface{} {
	pathExpr := v.VisitPathExpression(ctx.PathExpression().(*PathExpressionContext)).(Expression)

	if ctx.GetOp() != nil {
		postfixExpression := NewPostfixExpression(pathExpr, v.createGroovyToken(ctx.GetOp()))

		if v.visitingAssertStatementCount > 0 {
			// powerassert requires different column for values, so we have to copy the location of op
			return configureAST(postfixExpression, ctx)
		} else {
			return configureAST(postfixExpression, ctx)
		}
	}

	return configureAST(pathExpr, ctx)
}

func (v *ASTBuilder) VisitPathExpression(ctx *PathExpressionContext) interface{} {
	staticTerminalNode := ctx.STATIC()
	var primaryExpr Expression

	if staticTerminalNode != nil {
		primaryExpr = configureAST(NewVariableExpression(staticTerminalNode.GetText()), staticTerminalNode)
	} else {
		primaryExpr = v.Visit(ctx.Primary()).(Expression)
	}

	return v.createPathExpression(primaryExpr, ctx.AllPathElement())
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

func isTrue(obj NodeMetaDataHandler, key string) bool {
	value := obj.GetNodeMetaData(key)
	boolValue, ok := value.(bool)
	return ok && boolValue
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

func (v *ASTBuilder) VisitNonWildcardTypeArguments(ctx *NonWildcardTypeArgumentsContext) []*GenericsType {
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

func (v *ASTBuilder) VisitPathElement(ctx *PathElementContext) interface{} {
	baseExpr := ctx.GetNodeMetaData(PATH_EXPRESSION_BASE_EXPR).(Expression)
	if baseExpr == nil {
		panic("baseExpr is required!")
	}

	if ctx.NamePart() != nil {
		namePartExpr := v.VisitNamePart(ctx.NamePart().(*NamePartContext)).(Expression)
		genericsTypes := v.VisitNonWildcardTypeArguments(ctx.NonWildcardTypeArguments().(*NonWildcardTypeArgumentsContext))

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
		tuple := v.VisitIndexPropertyArgs(ctx.IndexPropertyArgs().(*IndexPropertyArgsContext)).(Tuple2)
		isSafeChain := isTrue(baseExpr, PATH_EXPRESSION_BASE_EXPR_SAFE_CHAIN)
		return configureAST(
			NewBinaryExpression(
				baseExpr,
				v.createGroovyToken(tuple.GetV1()),
				tuple.GetV2().(Expression),
				isSafeChain || ctx.IndexPropertyArgs().(*IndexPropertyArgsContext).SAFE_INDEX() != nil,
			),
			ctx,
		)
	} else if ctx.NamedPropertyArgs() != nil {
		mapEntryExpressionList := v.VisitNamedPropertyArgs(ctx.NamedPropertyArgs().(*NamedPropertyArgsContext)).([]MapEntryExpression)

		var right Expression
		mapEntryExpressionListSize := len(mapEntryExpressionList)
		if mapEntryExpressionListSize == 0 {
			right = configureAST(
				NewSpreadMapExpression(configureAST(NewMapExpression(), ctx.NamedPropertyArgs())),
				ctx.NamedPropertyArgs(),
			)
		} else if mapEntryExpressionListSize == 1 {
			firstKeyExpression := mapEntryExpressionList[0].GetKeyExpression()
			if _, ok := firstKeyExpression.(*SpreadMapExpression); ok {
				right = firstKeyExpression
			} else {
				listExpression := configureAST(
					NewListExpression(mapEntryExpressionList),
					ctx.NamedPropertyArgs(),
				)
				listExpression.SetWrapped(true)
				right = listExpression
			}
		} else {
			listExpression := configureAST(
				NewListExpression(mapEntryExpressionList),
				ctx.NamedPropertyArgs(),
			)
			listExpression.SetWrapped(true)
			right = listExpression
		}

		namedPropertyArgsContext := ctx.NamedPropertyArgs().(*NamedPropertyArgsContext)
		var token antlr.Token
		if namedPropertyArgsContext.LBRACK() == nil {
			token = namedPropertyArgsContext.SAFE_INDEX().GetSymbol()
		} else {
			token = namedPropertyArgsContext.LBRACK().GetSymbol()
		}
		return configureAST(
			NewBinaryExpression(baseExpr, v.createGroovyToken(token), right),
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
			return configureAST(v.createCallMethodCallExpression(attributeExpression, argumentsExpr, true), ctx)
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
				panic(createParsingFailedException("Primitive type literal: "+baseExprText+" cannot be used as a method name", ctx))
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
		closureExpression := v.VisitClosureOrLambdaExpression(ctx.ClosureOrLambdaExpression().(*ClosureOrLambdaExpressionContext)).(*ClosureExpression)

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
						configureAST(
							NewArgumentListExpression(
								configureAST(
									NewMapExpression(namedArgumentListExpression.GetMapEntryExpressions()),
									namedArgumentListExpression,
								),
								closureExpression,
							),
							tupleExpression,
						),
					)
				} else {
					methodCallExpression.SetArguments(
						configureAST(
							NewArgumentListExpression(closureExpression),
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
				configureAST(
					NewArgumentListExpression(closureExpression),
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
			configureAST(
				NewArgumentListExpression(closureExpression),
				closureExpression,
			),
		)

		return configureAST(methodCallExpression, ctx)
	}

	panic(createParsingFailedException("Unsupported path element: "+ctx.GetText(), ctx))
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
