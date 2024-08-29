package nf

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"
	"sort"

	"go.starlark.net/starlark"
)

var _ parser.GroovyCodeVisitor = (*IncludeVisitor)(nil)

var _ starlark.Value = (*IncludedItem)(nil)
var _ starlark.HasAttrs = (*IncludedItem)(nil)

type IncludedItem struct {
	Name  string
	Alias string
}

// Implement starlark.Value interface
func (i *IncludedItem) String() string {
	if i.Alias != "" {
		return fmt.Sprintf("%s as %s", i.Name, i.Alias)
	}
	return i.Name
}
func (i *IncludedItem) Type() string         { return "IncludedItem" }
func (i *IncludedItem) Freeze()              {} // No-op
func (i *IncludedItem) Truth() starlark.Bool { return starlark.Bool(i.Name != "") }
func (i *IncludedItem) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(i.Name))
	h.Write([]byte(i.Alias))
	return h.Sum32(), nil
}

// Implement starlark.HasAttrs interface
func (i *IncludedItem) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(i.Name), nil
	case "alias":
		return starlark.String(i.Alias), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("IncludedItem has no attribute %q", name))
	}
}

func (i *IncludedItem) AttrNames() []string {
	return []string{"name", "alias"}
}

type IncludeStatement struct {
	Items      []IncludedItem
	ModulePath string
	LineNumber int
}

func (is IncludeStatement) String() string {
	return fmt.Sprintf("IncludeStatement(ModulePath: %q, Items: %v)", is.ModulePath, is.Items)
}

func (is IncludeStatement) Type() string         { return "IncludeStatement" }
func (is IncludeStatement) Freeze()              {} // No-op
func (is IncludeStatement) Truth() starlark.Bool { return starlark.Bool(len(is.Items) > 0) }

func (is IncludeStatement) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(is.ModulePath))
	for _, item := range is.Items {
		itemHash, err := item.Hash()
		if err != nil {
			return 0, err
		}
		h.Write([]byte(fmt.Sprintf("%d", itemHash)))
	}
	return h.Sum32(), nil
}

// Implement starlark.HasAttrs interface
func (is IncludeStatement) Attr(name string) (starlark.Value, error) {
	switch name {
	case "items":
		items := make([]starlark.Value, len(is.Items))
		for i, item := range is.Items {
			items[i] = &item
		}
		return starlark.NewList(items), nil
	case "module_path":
		return starlark.String(is.ModulePath), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("IncludeStatement has no attribute %q", name))
	}
}

func (is IncludeStatement) AttrNames() []string {
	return []string{"items", "module_path"}
}

type IncludeVisitor struct {
	includes []IncludeStatement
}

func NewIncludeVisitor() *IncludeVisitor {
	return &IncludeVisitor{
		includes: make([]IncludeStatement, 0),
	}
}

// Statements
func (v *IncludeVisitor) VisitBlockStatement(block *parser.BlockStatement) {
	for _, statement := range block.GetStatements() {
		v.VisitStatement(statement)
	}
}

func (v *IncludeVisitor) VisitForLoop(statement *parser.ForStatement) {
	v.VisitExpression(statement.GetCollectionExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *IncludeVisitor) VisitWhileLoop(statement *parser.WhileStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *IncludeVisitor) VisitDoWhileLoop(statement *parser.DoWhileStatement) {
	v.VisitStatement(statement.GetLoopBlock())
	v.VisitExpression(statement.GetBooleanExpression())
}

func (v *IncludeVisitor) VisitIfElse(statement *parser.IfStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetIfBlock())
	v.VisitStatement(statement.GetElseBlock())
}

func (v *IncludeVisitor) VisitExpressionStatement(statement *parser.ExpressionStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *IncludeVisitor) VisitReturnStatement(statement *parser.ReturnStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *IncludeVisitor) VisitAssertStatement(statement *parser.AssertStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitExpression(statement.GetMessageExpression())
}

func (v *IncludeVisitor) VisitTryCatchFinally(statement *parser.TryCatchStatement) {
	for _, resource := range statement.GetResourceStatements() {
		v.VisitStatement(resource)
	}
	v.VisitStatement(statement.GetTryStatement())
	for _, catchStatement := range statement.GetCatchStatements() {
		v.VisitStatement(catchStatement)
	}
	v.VisitStatement(statement.GetFinallyStatement())
}

func (v *IncludeVisitor) VisitSwitch(statement *parser.SwitchStatement) {
	v.VisitExpression(statement.GetExpression())
	for _, caseStatement := range statement.GetCaseStatements() {
		v.VisitStatement(caseStatement)
	}
	v.VisitStatement(statement.GetDefaultStatement())
}

func (v *IncludeVisitor) VisitCaseStatement(statement *parser.CaseStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *IncludeVisitor) VisitBreakStatement(statement *parser.BreakStatement) {}

func (v *IncludeVisitor) VisitContinueStatement(statement *parser.ContinueStatement) {}

func (v *IncludeVisitor) VisitThrowStatement(statement *parser.ThrowStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *IncludeVisitor) VisitSynchronizedStatement(statement *parser.SynchronizedStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *IncludeVisitor) VisitCatchStatement(statement *parser.CatchStatement) {
	v.VisitStatement(statement.GetCode())
}

func (v *IncludeVisitor) VisitEmptyStatement(statement *parser.EmptyStatement) {}

func (v *IncludeVisitor) VisitStatement(statement parser.Statement) {
	statement.Visit(v)
}

func (v *IncludeVisitor) getFromClosure(closure *parser.ClosureExpression) ([]IncludedItem, error) {
	var items []IncludedItem
	code, ok := closure.GetCode().(*parser.BlockStatement)
	if !ok {
		return items, errors.New("closure code is not a block statement")
	}
	for _, statement := range code.GetStatements() {
		exprStmt, ok := statement.(*parser.ExpressionStatement)
		if !ok {
			return items, errors.New("closure code statement is not an expression statement")
		}
		if varExpr, ok := exprStmt.GetExpression().(*parser.VariableExpression); ok {
			items = append(items, IncludedItem{Name: varExpr.GetName()})
		}
		if castExpr, ok := exprStmt.GetExpression().(*parser.CastExpression); ok {
			if nameExpr, ok := castExpr.GetExpression().(*parser.VariableExpression); ok {
				name := nameExpr.GetName()
				alias := nameExpr.BaseExpression.GetType().GetName()
				items = append(items, IncludedItem{Name: name, Alias: alias})
			}
		}
	}
	return items, nil
}

// Expressions
func (v *IncludeVisitor) VisitMethodCallExpression(call *parser.MethodCallExpression) {
	v.VisitExpression(call.GetObjectExpression())
	v.VisitExpression(call.GetMethod())
	v.VisitExpression(call.GetArguments())
	mce, ok := call.GetObjectExpression().(*parser.MethodCallExpression)
	if !ok {
		return
	}
	method, ok := mce.GetMethod().(*parser.ConstantExpression)
	if !ok {
		return
	}
	if method.GetText() != "include" {
		return
	}
	method, ok = call.GetMethod().(*parser.ConstantExpression)
	if !ok {
		return
	}
	if !(method.GetText() == "from") {
		return
	}
	args, ok := mce.GetArguments().(*parser.ArgumentListExpression)
	if !ok {
		return
	}
	if len(args.GetExpressions()) == 0 {
		return
	}
	closure, ok := args.GetExpressions()[0].(*parser.ClosureExpression)
	if !ok {
		return
	}
	items, err := v.getFromClosure(closure)
	if err != nil {
		return
	}
	args, ok = call.GetArguments().(*parser.ArgumentListExpression)
	if !ok {
		return
	}
	if len(args.GetExpressions()) != 1 {
		return
	}
	arg, ok := args.GetExpressions()[0].(*parser.ConstantExpression)
	if !ok {
		return
	}
	v.includes = append(v.includes, IncludeStatement{
		Items:      items,
		ModulePath: arg.GetText(),
		LineNumber: mce.GetLineNumber(),
	})
}

func (v *IncludeVisitor) VisitStaticMethodCallExpression(call *parser.StaticMethodCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *IncludeVisitor) VisitConstructorCallExpression(call *parser.ConstructorCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *IncludeVisitor) VisitTernaryExpression(expression *parser.TernaryExpression) {
	booleanExpr := expression.GetBooleanExpression()
	v.VisitExpression(booleanExpr)
	v.VisitExpression(expression.GetTrueExpression())
	v.VisitExpression(expression.GetFalseExpression())
}

func (v *IncludeVisitor) VisitShortTernaryExpression(expression *parser.ElvisOperatorExpression) {
	v.VisitTernaryExpression(expression.TernaryExpression)
}

func (v *IncludeVisitor) VisitBinaryExpression(expression *parser.BinaryExpression) {
	v.VisitExpression(expression.GetLeftExpression())
	v.VisitExpression(expression.GetRightExpression())
}

func (v *IncludeVisitor) VisitPrefixExpression(expression *parser.PrefixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitPostfixExpression(expression *parser.PostfixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitBooleanExpression(expression *parser.BooleanExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitClosureExpression(expression *parser.ClosureExpression) {
	if expression.IsParameterSpecified() {
		for _, parameter := range expression.GetParameters() {
			if parameter.HasInitialExpression() {
				v.VisitExpression(parameter.GetInitialExpression())
			}
		}
	}
	v.VisitStatement(expression.GetCode())
}

func (v *IncludeVisitor) VisitLambdaExpression(expression *parser.LambdaExpression) {
	v.VisitClosureExpression(expression.ClosureExpression)
}

func (v *IncludeVisitor) VisitTupleExpression(expression parser.ITupleExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *IncludeVisitor) VisitMapExpression(expression *parser.MapExpression) {
	entries := expression.GetMapEntryExpressions()
	exprs := make([]parser.Expression, len(entries))
	for i, entry := range entries {
		exprs[i] = entry
	}
	v.VisitListOfExpressions(exprs)
}

func (v *IncludeVisitor) VisitMapEntryExpression(expression *parser.MapEntryExpression) {
	v.VisitExpression(expression.GetKeyExpression())
	v.VisitExpression(expression.GetValueExpression())
}

func (v *IncludeVisitor) VisitListExpression(expression *parser.ListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *IncludeVisitor) VisitRangeExpression(expression *parser.RangeExpression) {
	v.VisitExpression(expression.GetFrom())
	v.VisitExpression(expression.GetTo())
}

func (v *IncludeVisitor) VisitPropertyExpression(expression *parser.PropertyExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *IncludeVisitor) VisitAttributeExpression(expression *parser.AttributeExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *IncludeVisitor) VisitFieldExpression(expression *parser.FieldExpression) {}

func (v *IncludeVisitor) VisitMethodPointerExpression(expression *parser.MethodPointerExpression) {
	v.VisitExpression(expression.GetExpression())
	v.VisitExpression(expression.GetMethodName())
}

func (v *IncludeVisitor) VisitMethodReferenceExpression(expression *parser.MethodReferenceExpression) {
	v.VisitMethodPointerExpression(expression.MethodPointerExpression)
}

func (v *IncludeVisitor) VisitConstantExpression(expression *parser.ConstantExpression) {}

func (v *IncludeVisitor) VisitClassExpression(expression *parser.ClassExpression) {}

func (v *IncludeVisitor) VisitVariableExpression(expression *parser.VariableExpression) {}

func (v *IncludeVisitor) VisitDeclarationExpression(expression *parser.DeclarationExpression) {
	v.VisitBinaryExpression(expression.BinaryExpression)
}

func (v *IncludeVisitor) VisitGStringExpression(expression *parser.GStringExpression) {
	v.VisitListOfExpressions(convertToExpressionSlice(expression.GetStrings()))
	v.VisitListOfExpressions(expression.GetValues())
}

func (v *IncludeVisitor) VisitArrayExpression(expression *parser.ArrayExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
	v.VisitListOfExpressions(expression.GetSizeExpression())
}

func (v *IncludeVisitor) VisitSpreadExpression(expression *parser.SpreadExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitSpreadMapExpression(expression *parser.SpreadMapExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitNotExpression(expression *parser.NotExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitUnaryMinusExpression(expression *parser.UnaryMinusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitUnaryPlusExpression(expression *parser.UnaryPlusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitBitwiseNegationExpression(expression *parser.BitwiseNegationExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitCastExpression(expression *parser.CastExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *IncludeVisitor) VisitArgumentlistExpression(expression *parser.ArgumentListExpression) {
	v.VisitTupleExpression(expression)
}

func (v *IncludeVisitor) VisitClosureListExpression(expression *parser.ClosureListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *IncludeVisitor) VisitEmptyExpression(expression *parser.EmptyExpression) {}

func (v *IncludeVisitor) VisitListOfExpressions(expressions []parser.Expression) {
	for _, expr := range expressions {
		v.VisitExpression(expr)
	}
}

func (v *IncludeVisitor) VisitExpression(expression parser.Expression) {
	expression.Visit(v)
}

// GetSortedIncludes returns a slice of IncludeStatement sorted by line number in ascending order
func (v *IncludeVisitor) GetSortedIncludes() []IncludeStatement {
	sortedIncludes := make([]IncludeStatement, 0, len(v.includes))
	for _, info := range v.includes {
		sortedIncludes = append(sortedIncludes, info)
	}

	sort.Slice(sortedIncludes, func(i, j int) bool {
		return sortedIncludes[i].LineNumber < sortedIncludes[j].LineNumber
	})

	return sortedIncludes
}
