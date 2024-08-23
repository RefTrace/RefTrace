package directives

import (
	"reft-go/parser"
)

var _ Directive = (*Container)(nil)

type Container struct {
	Name string
}

func (a Container) Type() DirectiveType { return ContainerType }

func MakeContainer(mce *parser.MethodCallExpression) *Container {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &Container{Name: strValue}
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &Container{Name: gstringExpr.GetText()}
			}
		}
	}
	return nil
}
