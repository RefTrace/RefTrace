package outputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Output = (*Eval)(nil)

type Eval struct {
	Command  string
	Emit     string
	Optional bool
	Topic    string
}

func (e *Eval) Attr(name string) (starlark.Value, error) {
	switch name {
	case "command":
		return starlark.String(e.Command), nil
	case "emit":
		return starlark.String(e.Emit), nil
	case "optional":
		return starlark.Bool(e.Optional), nil
	case "topic":
		return starlark.String(e.Topic), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("Eval has no attribute %q", name))
	}
}

func (e *Eval) AttrNames() []string {
	return []string{"command", "emit", "optional", "topic"}
}

func (e *Eval) String() string {
	return fmt.Sprintf("Eval(command=%q, emit=%q, optional=%v, topic=%q)",
		e.Command, e.Emit, e.Optional, e.Topic)
}

func (e *Eval) Type() string         { return "Eval" }
func (e *Eval) Freeze()              {} // No-op, as Eval is immutable
func (e *Eval) Truth() starlark.Bool { return starlark.Bool(e.Command != "") }
func (e *Eval) Hash() (uint32, error) {
	h := starlark.String(fmt.Sprintf("%s:%s:%v:%s",
		e.Command, e.Emit, e.Optional, e.Topic))
	return h.Hash()
}

func MakeEval(mce *parser.MethodCallExpression) (Output, error) {
	if mce.GetMethod().GetText() != "eval" {
		return nil, errors.New("invalid eval directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) < 1 || len(exprs) > 2 {
			return nil, errors.New("invalid eval directive: expected 1 to 2 arguments")
		}

		eval := &Eval{}

		for _, expr := range exprs {
			if ce, ok := expr.(*parser.ConstantExpression); ok {
				eval.Command = ce.GetText()
			}
			if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						valueExpr := entry.GetValueExpression()
						switch key.GetText() {
						case "emit":
							if value, ok := valueExpr.(*parser.VariableExpression); ok {
								eval.Emit = value.GetText()
							}
						case "optional":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									eval.Optional = boolVal
								}
							}
						case "topic":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								eval.Topic = value.GetText()
							}
						}
					}
				}
			}
		}

		if eval.Command != "" {
			return eval, nil
		}
	}
	return nil, errors.New("invalid eval directive")
}
