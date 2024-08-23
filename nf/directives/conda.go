package directives

import (
	"reft-go/parser"
)

var _ Directive = (*Conda)(nil)

type Conda struct {
	Dependencies string
}

func (a Conda) Type() DirectiveType { return CondaType }

func MakeConda(mce *parser.MethodCallExpression) *Conda {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &Conda{Dependencies: strValue}
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &Conda{Dependencies: gstringExpr.GetText()}
			}
		}
	}
	return nil
}
