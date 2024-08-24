package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*TagDirective)(nil)

type TagDirective struct {
	Tag string
}

func (a TagDirective) Type() DirectiveType { return TagDirectiveType }

func MakeTagDirective(mce *parser.MethodCallExpression) (*TagDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &TagDirective{Tag: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &TagDirective{Tag: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid TagDirective directive")
}
