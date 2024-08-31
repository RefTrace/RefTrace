package inputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Input = (*Each)(nil)

type Each struct {
	Collection Input
}

func MakeEach(mce *parser.MethodCallExpression) (Input, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid each directive")
		}
		expr := exprs[0]
		if ve, ok := expr.(*parser.VariableExpression); ok {
			variableName := ve.GetText()
			valInput := &Val{Var: variableName}
			return &Each{Collection: valInput}, nil
		}
		if mce, ok := expr.(*parser.MethodCallExpression); ok {
			methodName := mce.GetMethod().GetText()
			if methodName == "path" {
				if input, err := MakePath(mce); err == nil {
					return &Each{Collection: input}, nil
				}
			}
		}
	}
	return nil, errors.New("invalid each directive")
}

// Implement starlark.Value methods
func (e *Each) String() string {
	return fmt.Sprintf("each(%v)", e.Collection)
}

func (e *Each) Type() string {
	return "each"
}

func (e *Each) Freeze() {
	e.Collection.Freeze()
}

func (e *Each) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (v *Each) Hash() (uint32, error) {
	return v.Collection.Hash()
}

// Implement starlark.HasAttrs methods
func (e *Each) Attr(name string) (starlark.Value, error) {
	switch name {
	case "collection":
		return e.Collection, nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("each has no attribute %q", name))
	}
}

func (e *Each) AttrNames() []string {
	return []string{"collection"}
}
