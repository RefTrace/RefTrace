package inputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Input = (*Path)(nil)

type Path struct {
	Path    string
	Arity   string
	StageAs string
}

func (v *Path) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		return starlark.String(v.Path), nil
	case "arity":
		return starlark.String(v.Arity), nil
	case "stage_as":
		return starlark.String(v.StageAs), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("Path has no attribute %q", name))
	}
}

func (v *Path) AttrNames() []string {
	return []string{"path", "arity", "stage_as"}
}

// Implement other starlark.Pathue methods
func (v *Path) String() string {
	return fmt.Sprintf("Path(path=%s, arity=%s, stage_as=%s)", v.Path, v.Arity, v.StageAs)
}

func (v *Path) Type() string         { return "Path" }
func (v *Path) Freeze()              {} // No-op, as Path is immutable
func (v *Path) Truth() starlark.Bool { return starlark.Bool(v.Path != "") }
func (v *Path) Hash() (uint32, error) {
	return starlark.String(v.Path).Hash()
}

func MakePath(mce *parser.MethodCallExpression) (Input, error) {
	if mce.GetMethod().GetText() != "path" {
		return nil, errors.New("invalid path directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) < 1 || len(exprs) > 2 {
			return nil, errors.New("invalid path directive: expected 1 to 3 arguments")
		}

		path := &Path{}

		if len(exprs) == 1 {
			if ce, ok := exprs[0].(*parser.ConstantExpression); ok {
				path.Path = ce.GetText()
			} else {
				return nil, errors.New("invalid path argument")
			}
		}

		for _, expr := range exprs {
			if ve, ok := expr.(*parser.VariableExpression); ok {
				path.Path = ve.GetText()
			}
			if ce, ok := expr.(*parser.ConstantExpression); ok {
				path.Path = ce.GetText()
			}
			if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						if key.GetText() == "arity" {
							valueExpr := entry.GetValueExpression()
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								path.Arity = value.GetText()
							}
						}
						if key.GetText() == "stageAs" {
							valueExpr := entry.GetValueExpression()
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								path.StageAs = value.GetText()
							}
						}
					}
				}
			}
		}

		if path.Path != "" {
			return path, nil
		}
	}
	return nil, errors.New("invalid path directive")
}
