package directives

import (
	"reft-go/parser"
)

var _ Directive = (*FairDirective)(nil)

type FairDirective struct {
	Enabled bool
}

func (a FairDirective) Type() DirectiveType { return FairDirectiveType }

func MakeFairDirective(mce *parser.MethodCallExpression) *FairDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				return &FairDirective{Enabled: boolValue}
			}
		}
	}
	return nil
}
