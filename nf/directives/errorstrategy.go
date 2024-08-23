package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*ErrorStrategyDirective)(nil)

type ErrorStrategyDirective struct {
	Strategy string
}

func (a ErrorStrategyDirective) Type() DirectiveType { return ErrorStrategyDirectiveType }

func MakeErrorStrategyDirective(mce *parser.MethodCallExpression) (*ErrorStrategyDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid error strategy directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &ErrorStrategyDirective{Strategy: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid error strategy directive")
}
