package inputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Input = (*Stdin)(nil)

type Stdin struct {
	Var string
}

func (s *Stdin) Attr(name string) (starlark.Value, error) {
	switch name {
	case "var":
		return starlark.String(s.Var), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("stdin has no attribute %q", name))
	}
}

func (s *Stdin) AttrNames() []string {
	return []string{"var"}
}

// Implement other starlark.Value methods
func (s *Stdin) String() string       { return fmt.Sprintf("stdin(%s)", s.Var) }
func (s *Stdin) Type() string         { return "stdin" }
func (s *Stdin) Freeze()              {} // No-op, as Stdin is immutable
func (s *Stdin) Truth() starlark.Bool { return starlark.Bool(s.Var != "") }
func (s *Stdin) Hash() (uint32, error) {
	return starlark.String(s.Var).Hash()
}

func MakeStdin(mce *parser.MethodCallExpression) (Input, error) {
	if mce.GetMethod().GetText() != "stdin" {
		return nil, errors.New("invalid stdin directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid stdin directive")
		}
		if ve, ok := exprs[0].(*parser.VariableExpression); ok {
			return &Stdin{Var: ve.GetText()}, nil
		}
	}
	return nil, errors.New("invalid stdin directive")
}
