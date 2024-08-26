package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*DebugDirective)(nil)

func (d *DebugDirective) String() string {
	return fmt.Sprintf("DebugDirective(Enabled: %t)", d.Enabled)
}

func (d *DebugDirective) Type() string {
	return "debug_directive"
}

func (d *DebugDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (d *DebugDirective) Truth() starlark.Bool {
	return starlark.Bool(d.Enabled)
}

func (d *DebugDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%t", d.Enabled)))
	return h.Sum32(), nil
}

var _ starlark.Value = (*DebugDirective)(nil)
var _ starlark.HasAttrs = (*DebugDirective)(nil)

func (d *DebugDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "enabled":
		return starlark.Bool(d.Enabled), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("debug directive has no attribute %q", name))
	}
}

func (d *DebugDirective) AttrNames() []string {
	return []string{"enabled"}
}

type DebugDirective struct {
	Enabled bool
}

func MakeDebugDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid debug directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				return &DebugDirective{Enabled: boolValue}, nil
			}
		}
	}
	return nil, errors.New("invalid debug directive")
}
