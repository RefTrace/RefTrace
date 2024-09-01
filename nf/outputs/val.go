package outputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Output = (*Val)(nil)

type Val struct {
	Var      string
	Emit     string
	Optional bool
	Topic    string
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
func (v *Val) String() string       { return fmt.Sprintf("val(%s)", v.Var) }
func (v *Val) Type() string         { return "val" }
func (v *Val) Freeze()              {} // No-op, as Val is immutable
func (v *Val) Truth() starlark.Bool { return starlark.Bool(v.Var != "") }
func (v *Val) Hash() (uint32, error) {
	return starlark.String(v.Var).Hash()
}

func MakeVal(mce *parser.MethodCallExpression) (Output, error) {
	if mce.GetMethod().GetText() != "val" {
		return nil, errors.New("invalid val directive")
	}
	val := &Val{}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		for _, expr := range exprs {
			if ve, ok := expr.(*parser.VariableExpression); ok {
				val.Var = ve.GetText()
			}
			if ce, ok := expr.(*parser.ConstantExpression); ok {
				val.Var = ce.GetText()
			}
			if gse, ok := expr.(*parser.GStringExpression); ok {
				val.Var = gse.GetText()
			}
			if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						if key.GetText() == "emit" {
							valueExpr := entry.GetValueExpression()
							if value, ok := valueExpr.(*parser.VariableExpression); ok {
								val.Emit = value.GetText()
							}
						}
						if key.GetText() == "optional" {
							valueExpr := entry.GetValueExpression()
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolValue, ok := value.GetValue().(bool); ok {
									val.Optional = boolValue
								}
							}
						}
						if key.GetText() == "topic" {
							valueExpr := entry.GetValueExpression()
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								val.Topic = value.GetText()
							}
						}
					}
				}
			}
		}
	}
	if val.Var != "" {
		return val, nil
	}
	return nil, errors.New("invalid val directive")
}
