package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*ExecutorDirective)(nil)

type ExecutorDirective struct {
	Executor string
}

func (a ExecutorDirective) Type() DirectiveType { return ExecutorDirectiveType }

func MakeExecutorDirective(mce *parser.MethodCallExpression) (*ExecutorDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid executor directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &ExecutorDirective{Executor: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid executor directive")
}
