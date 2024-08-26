package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*EchoDirective)(nil)

func (e *EchoDirective) String() string {
	return fmt.Sprintf("EchoDirective(Enabled: %t)", e.Enabled)
}

func (e *EchoDirective) Type() string {
	return "echo_directive"
}

func (e *EchoDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (e *EchoDirective) Truth() starlark.Bool {
	return starlark.Bool(e.Enabled)
}

func (e *EchoDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%t", e.Enabled)))
	return h.Sum32(), nil
}

var _ starlark.Value = (*EchoDirective)(nil)
var _ starlark.HasAttrs = (*EchoDirective)(nil)

func (e *EchoDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "enabled":
		return starlark.Bool(e.Enabled), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("echo directive has no attribute %q", name))
	}
}

func (e *EchoDirective) AttrNames() []string {
	return []string{"enabled"}
}

type EchoDirective struct {
	Enabled bool
}

func MakeEchoDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid echo directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				return &EchoDirective{Enabled: boolValue}, nil
			}
		}
	}
	return nil, errors.New("invalid echo directive")
}
