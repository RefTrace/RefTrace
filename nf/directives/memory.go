package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*MemoryDirective)(nil)
var _ starlark.Value = (*MemoryDirective)(nil)
var _ starlark.HasAttrs = (*MemoryDirective)(nil)

func (m *MemoryDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "memory":
		return starlark.String(m.Memory), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("memory directive has no attribute %q", name))
	}
}

func (m *MemoryDirective) AttrNames() []string {
	return []string{"memory"}
}

type MemoryDirective struct {
	Memory string
}

func (m *MemoryDirective) String() string {
	return fmt.Sprintf("MemoryDirective(Memory: %q)", m.Memory)
}

func (m *MemoryDirective) Type() string {
	return "memory_directive"
}

func (m *MemoryDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MemoryDirective) Truth() starlark.Bool {
	return starlark.Bool(m.Memory != "")
}

func (m *MemoryDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(m.Memory))
	return h.Sum32(), nil
}

func MakeMemoryDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid memory directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MemoryDirective{Memory: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid memory directive")
}
