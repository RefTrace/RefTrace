package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*TimeDirective)(nil)

type TimeDirective struct {
	Duration string
}

func (a TimeDirective) Type() DirectiveType { return TimeDirectiveType }

func MakeTimeDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Time directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &TimeDirective{Duration: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid Time directive")
}
