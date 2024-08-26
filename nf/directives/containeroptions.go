package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ContainerOptions)(nil)

func (c *ContainerOptions) String() string {
	return fmt.Sprintf("ContainerOptions(%q)", c.Options)
}

func (c *ContainerOptions) Type() string {
	return "container_options"
}

func (c *ContainerOptions) Freeze() {
	// No mutable fields, so no action needed
}

func (c *ContainerOptions) Truth() starlark.Bool {
	return starlark.Bool(c.Options != "")
}

func (c *ContainerOptions) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(c.Options))
	return h.Sum32(), nil
}

type ContainerOptions struct {
	Options string
}

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
