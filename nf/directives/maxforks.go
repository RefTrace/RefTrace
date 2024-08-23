package directives

import (
	"reft-go/parser"
)

var _ Directive = (*MaxForksDirective)(nil)

type MaxForksDirective struct {
	Num int
}

func (a MaxForksDirective) Type() DirectiveType { return MaxForksDirectiveType }

func MakeMaxForksDirective(mce *parser.MethodCallExpression) *MaxForksDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &MaxForksDirective{Num: intValue}
			}
		}
	}
	return nil
}
