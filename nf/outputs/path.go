package outputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Output = (*Path)(nil)

type Path struct {
	Path          string
	Arity         string
	FollowLinks   bool
	Glob          bool
	Hidden        bool
	IncludeInputs bool
	MaxDepth      int
	PathType      string
	Emit          string
	Optional      bool
	Topic         string
}

func (v *Path) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		return starlark.String(v.Path), nil
	case "arity":
		return starlark.String(v.Arity), nil
	case "follow_links":
		return starlark.Bool(v.FollowLinks), nil
	case "glob":
		return starlark.Bool(v.Glob), nil
	case "hidden":
		return starlark.Bool(v.Hidden), nil
	case "include_inputs":
		return starlark.Bool(v.IncludeInputs), nil
	case "max_depth":
		return starlark.MakeInt(v.MaxDepth), nil
	case "path_type":
		return starlark.String(v.PathType), nil
	case "emit":
		return starlark.String(v.Emit), nil
	case "optional":
		return starlark.Bool(v.Optional), nil
	case "topic":
		return starlark.String(v.Topic), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("Path has no attribute %q", name))
	}
}

func (v *Path) AttrNames() []string {
	return []string{
		"path",
		"arity",
		"follow_links",
		"glob",
		"hidden",
		"include_inputs",
		"max_depth",
		"path_type",
		"emit",
		"optional",
		"topic",
	}
}

// Implement other starlark.Value methods
func (v *Path) String() string {
	return fmt.Sprintf("Path(path=%q, arity=%q, follow_links=%v, glob=%v, hidden=%v, include_inputs=%v, max_depth=%d, path_type=%q, emit=%q, optional=%v, topic=%q)",
		v.Path, v.Arity, v.FollowLinks, v.Glob, v.Hidden, v.IncludeInputs, v.MaxDepth, v.PathType, v.Emit, v.Optional, v.Topic)
}

func (v *Path) Type() string         { return "Path" }
func (v *Path) Freeze()              {} // No-op, as Path is immutable
func (v *Path) Truth() starlark.Bool { return starlark.Bool(v.Path != "") }
func (v *Path) Hash() (uint32, error) {
	// Include all fields in the hash calculation
	h := starlark.String(fmt.Sprintf("%s:%s:%v:%v:%v:%v:%d:%s:%s:%v:%s",
		v.Path, v.Arity, v.FollowLinks, v.Glob, v.Hidden, v.IncludeInputs, v.MaxDepth, v.PathType, v.Emit, v.Optional, v.Topic))
	return h.Hash()
}

func MakePath(mce *parser.MethodCallExpression) (Output, error) {
	if mce.GetMethod().GetText() != "path" {
		return nil, errors.New("invalid path directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) < 1 || len(exprs) > 2 {
			return nil, errors.New("invalid path directive: expected 1 to 3 arguments")
		}

		path := &Path{}

		for _, expr := range exprs {
			if ve, ok := expr.(*parser.VariableExpression); ok {
				path.Path = ve.GetText()
			}
			if ce, ok := expr.(*parser.ConstantExpression); ok {
				path.Path = ce.GetText()
			}
			if ge, ok := expr.(*parser.GStringExpression); ok {
				path.Path = ge.GetText()
			}
			if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						valueExpr := entry.GetValueExpression()
						switch key.GetText() {
						case "arity":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								path.Arity = value.GetText()
							}
						case "followLinks":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									path.FollowLinks = boolVal
								}
							}
						case "glob":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									path.Glob = boolVal
								}
							}
						case "hidden":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									path.Hidden = boolVal
								}
							}
						case "includeInputs":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									path.IncludeInputs = boolVal
								}
							}
						case "maxDepth":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if intVal, err := value.GetValue().(int); err {
									path.MaxDepth = intVal
								}
							}
						case "type":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								path.PathType = value.GetText()
							}
						case "emit":
							if value, ok := valueExpr.(*parser.VariableExpression); ok {
								path.Emit = value.GetText()
							}
						case "optional":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									path.Optional = boolVal
								}
							}
						case "topic":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								path.Topic = value.GetText()
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
