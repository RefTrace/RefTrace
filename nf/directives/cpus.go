package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*CpusDirective)(nil)

type CpusDirective struct {
	Num int
}

func (a CpusDirective) Type() DirectiveType { return CpusDirectiveType }

func MakeCpusDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid cpus directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &CpusDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid cpus directive")
}
