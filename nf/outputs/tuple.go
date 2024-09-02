package outputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Output = (*Tuple)(nil)

type Tuple struct {
	Values   []Output
	Emit     string
	Optional bool
	Topic    string
}

func MakeTuple(mce *parser.MethodCallExpression) (Output, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		var values []Output
		tuple := &Tuple{}

		for _, expr := range exprs {
			if mce, ok := expr.(*parser.MethodCallExpression); ok {
				methodName := mce.GetMethod().GetText()
				var output Output
				var err error

				switch methodName {
				case "val":
					output, err = MakeVal(mce)
				case "path":
					output, err = MakePath(mce)
				case "env":
					output, err = MakeEnv(mce)
				case "stdout":
					output, err = MakeStdout(mce)
				case "eval":
					output, err = MakeEval(mce)
				case "file":
					output, err = MakeFile(mce)
				}

				if err == nil && output != nil {
					values = append(values, output)
				}
			} else if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						valueExpr := entry.GetValueExpression()
						switch key.GetText() {
						case "emit":
							if value, ok := valueExpr.(*parser.VariableExpression); ok {
								tuple.Emit = value.GetText()
							}
						case "optional":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									tuple.Optional = boolVal
								}
							}
						case "topic":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								tuple.Topic = value.GetText()
							}
						}
					}
				}
			}
		}
		tuple.Values = values
		return tuple, nil
	}
	return nil, errors.New("invalid tuple directive")
}

// Implement starlark.Value methods
func (t *Tuple) String() string {
	return fmt.Sprintf("tuple(values=%v, emit=%q, optional=%v, topic=%q)",
		t.Values, t.Emit, t.Optional, t.Topic)
}

func (t *Tuple) Type() string {
	return "tuple"
}

func (t *Tuple) Freeze() {
	for _, v := range t.Values {
		v.Freeze()
	}
}

func (t *Tuple) Truth() starlark.Bool {
	return starlark.Bool(len(t.Values) > 0)
}

func (t *Tuple) Hash() (uint32, error) {
	var hash uint32
	for _, v := range t.Values {
		h, err := v.Hash()
		if err != nil {
			return 0, err
		}
		hash = hash*31 + h
	}
	h := starlark.String(fmt.Sprintf("%v:%s:%v:%s", hash, t.Emit, t.Optional, t.Topic))
	return h.Hash()
}

// Implement starlark.HasAttrs methods
func (t *Tuple) Attr(name string) (starlark.Value, error) {
	switch name {
	case "values":
		values := make([]starlark.Value, len(t.Values))
		for i, v := range t.Values {
			values[i] = v
		}
		return starlark.NewList(values), nil
	case "emit":
		return starlark.String(t.Emit), nil
	case "optional":
		return starlark.Bool(t.Optional), nil
	case "topic":
		return starlark.String(t.Topic), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("tuple has no attribute %q", name))
	}
}

func (t *Tuple) AttrNames() []string {
	return []string{"values", "emit", "optional", "topic"}
}
