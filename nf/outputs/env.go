package outputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Output = (*Env)(nil)

type Env struct {
	Var      string
	Emit     string
	Optional bool
	Topic    string
}

func (e *Env) Attr(name string) (starlark.Value, error) {
	switch name {
	case "var":
		return starlark.String(e.Var), nil
	case "emit":
		return starlark.String(e.Emit), nil
	case "optional":
		return starlark.Bool(e.Optional), nil
	case "topic":
		return starlark.String(e.Topic), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("Env has no attribute %q", name))
	}
}

func (e *Env) AttrNames() []string {
	return []string{"var", "emit", "optional", "topic"}
}

// Implement other starlark.Value methods
func (e *Env) String() string {
	return fmt.Sprintf("Env(var=%q, emit=%q, optional=%v, topic=%q)",
		e.Var, e.Emit, e.Optional, e.Topic)
}

func (e *Env) Type() string         { return "Env" }
func (e *Env) Freeze()              {} // No-op, as Env is immutable
func (e *Env) Truth() starlark.Bool { return starlark.Bool(e.Var != "") }
func (e *Env) Hash() (uint32, error) {
	h := starlark.String(fmt.Sprintf("%s:%s:%v:%s",
		e.Var, e.Emit, e.Optional, e.Topic))
	return h.Hash()
}

func MakeEnv(mce *parser.MethodCallExpression) (Output, error) {
	if mce.GetMethod().GetText() != "env" {
		return nil, errors.New("invalid env directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) < 1 || len(exprs) > 2 {
			return nil, errors.New("invalid env directive: expected 1 to 2 arguments")
		}

		env := &Env{}

		for _, expr := range exprs {
			if ve, ok := expr.(*parser.VariableExpression); ok {
				env.Var = ve.GetText()
			}
			if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						valueExpr := entry.GetValueExpression()
						switch key.GetText() {
						case "emit":
							if value, ok := valueExpr.(*parser.VariableExpression); ok {
								env.Emit = value.GetText()
							}
						case "optional":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									env.Optional = boolVal
								}
							}
						case "topic":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								env.Topic = value.GetText()
							}
						}
					}
				}
			}
		}

		if env.Var != "" {
			return env, nil
		}
	}
	return nil, errors.New("invalid env directive")
}
