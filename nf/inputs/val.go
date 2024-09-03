package inputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Input = (*Val)(nil)

type Val struct {
	Var string
}

func (v *Val) Attr(name string) (starlark.Value, error) {
	switch name {
	case "var":
		return starlark.String(v.Var), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("val has no attribute %q", name))
	}
}

func (v *Val) AttrNames() []string {
	return []string{"var"}
}

// Implement other starlark.Value methods
func (v *Val) String() string       { return fmt.Sprintf("Val(%s)", v.Var) }
func (v *Val) Type() string         { return "val" }
func (v *Val) Freeze()              {} // No-op, as Val is immutable
func (v *Val) Truth() starlark.Bool { return starlark.Bool(v.Var != "") }
func (v *Val) Hash() (uint32, error) {
	return starlark.String(v.Var).Hash()
}

func MakeVal(mce *parser.MethodCallExpression) (Input, error) {
	if mce.GetMethod().GetText() != "val" {
		return nil, errors.New("invalid val directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid val directive")
		}
		if ve, ok := exprs[0].(*parser.VariableExpression); ok {
			return &Val{Var: ve.GetText()}, nil
		}
	}
	return nil, errors.New("invalid val directive")
}
