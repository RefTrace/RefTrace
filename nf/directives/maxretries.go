package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*MaxRetriesDirective)(nil)

type MaxRetriesDirective struct {
	Num int
}

func (a MaxRetriesDirective) Type() DirectiveType { return MaxRetriesDirectiveType }

func MakeMaxRetriesDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max retries directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &MaxRetriesDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max retries directive")
}
