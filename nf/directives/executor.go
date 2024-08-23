package directives

import (
	"reft-go/parser"
)

var _ Directive = (*ExecutorDirective)(nil)

type ExecutorDirective struct {
	Executor string
}

func (a ExecutorDirective) Type() DirectiveType { return ExecutorDirectiveType }

func MakeExecutorDirective(mce *parser.MethodCallExpression) *ExecutorDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &ExecutorDirective{Executor: strValue}
			}
		}
	}
	return nil
}
