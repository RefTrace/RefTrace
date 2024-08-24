package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*ContainerOptions)(nil)

type ContainerOptions struct {
	Options string
}

func (a ContainerOptions) Type() DirectiveType { return ContainerOptionsType }

func MakeContainerOptions(mce *parser.MethodCallExpression) (Directive, error) {
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
