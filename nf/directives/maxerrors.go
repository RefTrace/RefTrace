package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*MaxErrorsDirective)(nil)

type MaxErrorsDirective struct {
	Num int
}

func (a MaxErrorsDirective) Type() DirectiveType { return MaxErrorsDirectiveType }

func MakeMaxErrorsDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max errors directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &MaxErrorsDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max errors directive")
}
