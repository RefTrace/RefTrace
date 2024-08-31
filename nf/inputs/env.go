package inputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Input = (*Env)(nil)

type Env struct {
	Var string
}

func (e *Env) Attr(name string) (starlark.Value, error) {
	switch name {
	case "var":
		return starlark.String(e.Var), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("env has no attribute %q", name))
	}
}

func (e *Env) AttrNames() []string {
	return []string{"var"}
}

// Implement other starlark.Value methods
func (e *Env) String() string       { return fmt.Sprintf("env(%s)", e.Var) }
func (e *Env) Type() string         { return "env" }
func (e *Env) Freeze()              {} // No-op, as Env is immutable
func (e *Env) Truth() starlark.Bool { return starlark.Bool(e.Var != "") }
func (e *Env) Hash() (uint32, error) {
	return starlark.String(e.Var).Hash()
}

func MakeEnv(mce *parser.MethodCallExpression) (Input, error) {
	if mce.GetMethod().GetText() != "env" {
		return nil, errors.New("invalid env directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid env directive")
		}
		if ve, ok := exprs[0].(*parser.VariableExpression); ok {
			return &Env{Var: ve.GetText()}, nil
		}
	}
	return nil, errors.New("invalid env directive")
}
