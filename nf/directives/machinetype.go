package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*MachineTypeDirective)(nil)
var _ starlark.Value = (*MachineTypeDirective)(nil)
var _ starlark.HasAttrs = (*MachineTypeDirective)(nil)

func (m *MachineTypeDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "machine_type":
		return starlark.String(m.MachineType), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("machine type directive has no attribute %q", name))
	}
}

func (m *MachineTypeDirective) AttrNames() []string {
	return []string{"machine_type"}
}

type MachineTypeDirective struct {
	MachineType string
}

func (m *MachineTypeDirective) String() string {
	return fmt.Sprintf("MachineTypeDirective(MachineType: %q)", m.MachineType)
}

func (m *MachineTypeDirective) Type() string {
	return "machine_type_directive"
}

func (m *MachineTypeDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MachineTypeDirective) Truth() starlark.Bool {
	return starlark.Bool(m.MachineType != "")
}

func (m *MachineTypeDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(m.MachineType))
	return h.Sum32(), nil
}

func MakeMachineTypeDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid machine type directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MachineTypeDirective{MachineType: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid machine type directive")
}
