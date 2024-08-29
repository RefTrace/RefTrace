package inputs

import (
	"errors"
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Input = (*Tuple)(nil)

type Tuple struct {
	Values []Input
}

func MakeTuple(mce *parser.MethodCallExpression) (Input, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		var values []Input
		for _, expr := range exprs {
			if mce, ok := expr.(*parser.MethodCallExpression); ok {
				methodName := mce.GetMethod().GetText()
				if methodName == "val" {
					if input, err := MakeVal(mce); err == nil {
						values = append(values, input)
						continue
					}
				}
				if methodName == "path" {
					if input, err := MakePath(mce); err == nil {
						values = append(values, input)
						continue
					}
				}
			}
		}
		return &Tuple{Values: values}, nil
	}
	return nil, errors.New("invalid tuple directive")
}

// Implement starlark.Value methods
func (t *Tuple) String() string {
	return fmt.Sprintf("tuple(%v)", t.Values)
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
	// This is a simple hash function and might not be suitable for all use cases
	var hash uint32
	for _, v := range t.Values {
		h, err := v.Hash()
		if err != nil {
			return 0, err
		}
		hash = hash*31 + h
	}
	return hash, nil
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
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("tuple has no attribute %q", name))
	}
}

func (t *Tuple) AttrNames() []string {
	return []string{"values"}
}
