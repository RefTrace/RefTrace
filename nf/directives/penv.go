package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*PenvDirective)(nil)

type PenvDirective struct {
	Environment string
}

func (a PenvDirective) Type() DirectiveType { return PenvDirectiveType }

func MakePenvDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid penv directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &PenvDirective{Environment: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid penv directive")
}
