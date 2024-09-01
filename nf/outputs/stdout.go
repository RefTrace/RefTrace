package outputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Output = (*Stdout)(nil)

type Stdout struct {
	Emit     string
	Optional bool
	Topic    string
}

func (s *Stdout) Attr(name string) (starlark.Value, error) {
	switch name {
	case "emit":
		return starlark.String(s.Emit), nil
	case "optional":
		return starlark.Bool(s.Optional), nil
	case "topic":
		return starlark.String(s.Topic), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("Stdout has no attribute %q", name))
	}
}

func (s *Stdout) AttrNames() []string {
	return []string{"emit", "optional", "topic"}
}

// Implement other starlark.Value methods
func (s *Stdout) String() string {
	return fmt.Sprintf("Stdout(emit=%q, optional=%v, topic=%q)",
		s.Emit, s.Optional, s.Topic)
}

func (s *Stdout) Type() string         { return "Stdout" }
func (s *Stdout) Freeze()              {} // No-op, as Stdout is immutable
func (s *Stdout) Truth() starlark.Bool { return starlark.Bool(true) }
func (s *Stdout) Hash() (uint32, error) {
	h := starlark.String(fmt.Sprintf("%s:%v:%s",
		s.Emit, s.Optional, s.Topic))
	return h.Hash()
}

func MakeStdout(mce *parser.MethodCallExpression) (Output, error) {
	if mce.GetMethod().GetText() != "stdout" {
		return nil, errors.New("invalid stdout directive")
	}
	if args, ok := mce.GetArguments().(*parser.TupleExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid stdout directive: expected 1 argument")
		}

		stdout := &Stdout{}

		expr := exprs[0]

		if me, ok := expr.(*parser.NamedArgumentListExpression); ok {
			entries := me.GetMapEntryExpressions()
			for _, entry := range entries {
				if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
					valueExpr := entry.GetValueExpression()
					switch key.GetText() {
					case "emit":
						if value, ok := valueExpr.(*parser.VariableExpression); ok {
							stdout.Emit = value.GetText()
						}
					case "optional":
						if value, ok := valueExpr.(*parser.ConstantExpression); ok {
							if boolVal, err := value.GetValue().(bool); err {
								stdout.Optional = boolVal
							}
						}
					case "topic":
						if value, ok := valueExpr.(*parser.ConstantExpression); ok {
							stdout.Topic = value.GetText()
						}
					}
				}
			}
		}

		return stdout, nil
	}
	return nil, errors.New("invalid stdout directive")
}
