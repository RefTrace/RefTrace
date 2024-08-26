package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ResourceLabelsDirective)(nil)
var _ starlark.Value = (*ResourceLabelsDirective)(nil)
var _ starlark.HasAttrs = (*ResourceLabelsDirective)(nil)

func (r *ResourceLabelsDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "keys":
		starlarkKeys := make([]starlark.Value, len(r.Keys))
		for i, key := range r.Keys {
			starlarkKeys[i] = starlark.String(key)
		}
		return starlark.NewList(starlarkKeys), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("resource_labels directive has no attribute %q", name))
	}
}

func (r *ResourceLabelsDirective) AttrNames() []string {
	return []string{"keys"}
}

type ResourceLabelsDirective struct {
	Keys []string
}

func (r *ResourceLabelsDirective) String() string {
	return fmt.Sprintf("ResourceLabelsDirective(Keys: %q)", r.Keys)
}

func (r *ResourceLabelsDirective) Type() string {
	return "resource_labels_directive"
}

func (r *ResourceLabelsDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (r *ResourceLabelsDirective) Truth() starlark.Bool {
	return starlark.Bool(len(r.Keys) > 0)
}

func (r *ResourceLabelsDirective) Hash() (uint32, error) {
	h := fnv.New32()
	for _, key := range r.Keys {
		h.Write([]byte(key))
	}
	return h.Sum32(), nil
}

func MakeResourceLabelsDirective(mce *parser.MethodCallExpression) (Directive, error) {
	var keys []string = []string{}
	if args, ok := mce.GetArguments().(*parser.TupleExpression); ok {
		if len(args.GetExpressions()) != 1 {
			return nil, errors.New("invalid resource labels directive")
		}
		expr := args.GetExpressions()[0]
		if namedArgListExpr, ok := expr.(*parser.NamedArgumentListExpression); ok {
			exprs := namedArgListExpr.GetMapEntryExpressions()
			for _, mapEntryExpr := range exprs {
				key := mapEntryExpr.GetKeyExpression().GetText()
				keys = append(keys, key)
			}
		}
	}
	return &ResourceLabelsDirective{Keys: keys}, nil
}
