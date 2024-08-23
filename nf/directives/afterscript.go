package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*AfterScript)(nil)

type AfterScript struct {
	Script string
}

func (a AfterScript) Type() DirectiveType { return AfterScriptType }

func MakeAfterScript(mce *parser.MethodCallExpression) (*AfterScript, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &AfterScript{Script: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &AfterScript{Script: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid afterScript directive")
}
