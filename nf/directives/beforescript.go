package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*BeforeScript)(nil)

type BeforeScript struct {
	Script string
}

func (a BeforeScript) Type() DirectiveType { return BeforeScriptType }

func MakeBeforeScript(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &BeforeScript{Script: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &BeforeScript{Script: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid beforeScript directive")
}
