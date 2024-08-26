package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*CpusDirective)(nil)

func (c *CpusDirective) String() string {
	return fmt.Sprintf("CpusDirective(%d)", c.Num)
}

func (c *CpusDirective) Type() string {
	return "cpus_directive"
}

func (c *CpusDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (c *CpusDirective) Truth() starlark.Bool {
	return starlark.Bool(c.Num > 0)
}

func (c *CpusDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%d", c.Num)))
	return h.Sum32(), nil
}

var _ starlark.Value = (*CpusDirective)(nil)
var _ starlark.HasAttrs = (*CpusDirective)(nil)

func (c *CpusDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "num":
		return starlark.MakeInt(c.Num), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("cpus directive has no attribute %q", name))
	}
}

func (c *CpusDirective) AttrNames() []string {
	return []string{"num"}
}

type CpusDirective struct {
	Num int
}

func MakeCpusDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid cpus directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &CpusDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid cpus directive")
}
