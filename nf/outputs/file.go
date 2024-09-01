package outputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Output = (*File)(nil)

type File struct {
	Path     string
	Emit     string
	Optional bool
	Topic    string
}

func (v *File) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		return starlark.String(v.Path), nil
	case "emit":
		return starlark.String(v.Emit), nil
	case "optional":
		return starlark.Bool(v.Optional), nil
	case "topic":
		return starlark.String(v.Topic), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("File has no attribute %q", name))
	}
}

func (v *File) AttrNames() []string {
	return []string{"path", "emit", "optional", "topic"}
}

// Implement other starlark.Value methods
func (v *File) String() string {
	return fmt.Sprintf("File(path=%q, emit=%q, optional=%v, topic=%q)",
		v.Path, v.Emit, v.Optional, v.Topic)
}

func (v *File) Type() string         { return "File" }
func (v *File) Freeze()              {} // No-op, as File is immutable
func (v *File) Truth() starlark.Bool { return starlark.Bool(v.Path != "") }
func (v *File) Hash() (uint32, error) {
	// Include all fields in the hash calculation
	h := starlark.String(fmt.Sprintf("%s:%s:%v:%s",
		v.Path, v.Emit, v.Optional, v.Topic))
	return h.Hash()
}

func MakeFile(mce *parser.MethodCallExpression) (Output, error) {
	if mce.GetMethod().GetText() != "file" {
		return nil, errors.New("invalid file directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) < 1 || len(exprs) > 2 {
			return nil, errors.New("invalid file directive: expected 1 to 2 arguments")
		}

		file := &File{}

		for _, expr := range exprs {
			if ve, ok := expr.(*parser.VariableExpression); ok {
				file.Path = ve.GetText()
			}
			if ce, ok := expr.(*parser.ConstantExpression); ok {
				file.Path = ce.GetText()
			}
			if ge, ok := expr.(*parser.GStringExpression); ok {
				file.Path = ge.GetText()
			}
			if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						valueExpr := entry.GetValueExpression()
						switch key.GetText() {
						case "emit":
							if value, ok := valueExpr.(*parser.VariableExpression); ok {
								file.Emit = value.GetText()
							}
						case "optional":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								if boolVal, err := value.GetValue().(bool); err {
									file.Optional = boolVal
								}
							}
						case "topic":
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								file.Topic = value.GetText()
							}
						}
					}
				}
			}
		}

		if file.Path != "" {
			return file, nil
		}
	}
	return nil, errors.New("invalid file directive")
}
