package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*ArrayDirective)(nil)

type ArrayDirective struct {
	Size int
}

func (a ArrayDirective) Type() DirectiveType { return ArrayDirectiveType }

func MakeArrayDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid array directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &ArrayDirective{Size: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid array directive")
}
