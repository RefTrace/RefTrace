package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*Container)(nil)

type ContainerOptions struct {
	Options string
}

func (a ContainerOptions) Type() DirectiveType { return ContainerType }

func MakeContainerOptions(mce *parser.MethodCallExpression) (*ContainerOptions, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &ContainerOptions{Options: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &ContainerOptions{Options: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid containerOptions directive")
}
