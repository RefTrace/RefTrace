package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ArrayDirective)(nil)
var _ starlark.Value = (*ArrayDirective)(nil)
var _ starlark.HasAttrs = (*ArrayDirective)(nil)

func (a *ArrayDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "size":
		return starlark.MakeInt(a.Size), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("array directive has no attribute %q", name))
	}
}

func (a *ArrayDirective) AttrNames() []string {
	return []string{"size"}
}

func (a *ArrayDirective) String() string {
	return fmt.Sprintf("ArrayDirective(Size: %d)", a.Size)
}

func (a *ArrayDirective) Type() string {
	return "array_directive"
}

func (a *ArrayDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (a *ArrayDirective) Truth() starlark.Bool {
	return starlark.Bool(a.Size > 0)
}

func (a *ArrayDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%d", a.Size)))
	return h.Sum32(), nil
}

type ArrayDirective struct {
	Size int
}

func MakeArrayDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid array directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &ArrayDirective{Size: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid array directive")
}
