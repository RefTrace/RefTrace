package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*Container)(nil)

func (c *Container) String() string {
	return fmt.Sprintf("Container(%q)", c.Name)
}

func (c *Container) Type() string {
	return "container"
}

func (c *Container) Freeze() {
	// No mutable fields, so no action needed
}

func (c *Container) Truth() starlark.Bool {
	return starlark.Bool(c.Name != "")
}

func (c *Container) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(c.Name))
	return h.Sum32(), nil
}

type Container struct {
	Name string
}

func MakeContainer(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &Container{Name: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &Container{Name: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid container directive")
}
