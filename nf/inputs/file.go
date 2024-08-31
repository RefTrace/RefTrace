package inputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Input = (*File)(nil)

type File struct {
	Path    string
	Arity   string
	StageAs string
}

func (v *File) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		return starlark.String(v.Path), nil
	case "arity":
		return starlark.String(v.Arity), nil
	case "stage_as":
		return starlark.String(v.StageAs), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("File has no attribute %q", name))
	}
}

func (v *File) AttrNames() []string {
	return []string{"path", "arity", "stage_as"}
}

// Implement other starlark.Value methods
func (v *File) String() string {
	return fmt.Sprintf("File(path=%s, arity=%s, stage_as=%s)", v.Path, v.Arity, v.StageAs)
}

func (v *File) Type() string         { return "File" }
func (v *File) Freeze()              {} // No-op, as File is immutable
func (v *File) Truth() starlark.Bool { return starlark.Bool(v.Path != "") }
func (v *File) Hash() (uint32, error) {
	return starlark.String(v.Path).Hash()
}

func MakeFile(mce *parser.MethodCallExpression) (Input, error) {
	if mce.GetMethod().GetText() != "file" {
		return nil, errors.New("invalid file directive")
	}
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) < 1 || len(exprs) > 2 {
			return nil, errors.New("invalid file directive: expected 1 to 3 arguments")
		}

		file := &File{}

		if len(exprs) == 1 {
			if ce, ok := exprs[0].(*parser.ConstantExpression); ok {
				file.Path = ce.GetText()
			} else {
				return nil, errors.New("invalid file argument")
			}
		}

		for _, expr := range exprs {
			if ve, ok := expr.(*parser.VariableExpression); ok {
				file.Path = ve.GetText()
			}
			if ce, ok := expr.(*parser.ConstantExpression); ok {
				file.Path = ce.GetText()
			}
			if me, ok := expr.(*parser.MapExpression); ok {
				entries := me.GetMapEntryExpressions()
				for _, entry := range entries {
					if key, ok := entry.GetKeyExpression().(*parser.ConstantExpression); ok {
						if key.GetText() == "arity" {
							valueExpr := entry.GetValueExpression()
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								file.Arity = value.GetText()
							}
						}
						if key.GetText() == "stageAs" {
							valueExpr := entry.GetValueExpression()
							if value, ok := valueExpr.(*parser.ConstantExpression); ok {
								file.StageAs = value.GetText()
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
