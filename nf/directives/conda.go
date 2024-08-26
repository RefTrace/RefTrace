package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*Conda)(nil)

func (c *Conda) String() string {
	return fmt.Sprintf("Conda(%q)", c.Dependencies)
}

func (c *Conda) Type() string {
	return "conda"
}

func (c *Conda) Freeze() {
	// No mutable fields, so no action needed
}

func (c *Conda) Truth() starlark.Bool {
	return starlark.Bool(c.Dependencies != "")
}

func (c *Conda) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(c.Dependencies))
	return h.Sum32(), nil
}

var _ starlark.Value = (*Conda)(nil)
var _ starlark.HasAttrs = (*Conda)(nil)

func (c *Conda) Attr(name string) (starlark.Value, error) {
	switch name {
	case "dependencies":
		return starlark.String(c.Dependencies), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("conda directive has no attribute %q", name))
	}
}

func (c *Conda) AttrNames() []string {
	return []string{"dependencies"}
}

type Conda struct {
	Dependencies string
}

func MakeConda(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &Conda{Dependencies: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &Conda{Dependencies: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid conda directive")
}
