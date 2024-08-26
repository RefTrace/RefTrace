package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*QueueDirective)(nil)
var _ starlark.Value = (*QueueDirective)(nil)
var _ starlark.HasAttrs = (*QueueDirective)(nil)

func (q *QueueDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(q.Name), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("queue directive has no attribute %q", name))
	}
}

func (q *QueueDirective) AttrNames() []string {
	return []string{"name"}
}

type QueueDirective struct {
	Name string
}

func (q *QueueDirective) String() string {
	return fmt.Sprintf("QueueDirective(Name: %q)", q.Name)
}

func (q *QueueDirective) Type() string {
	return "queue_directive"
}

func (q *QueueDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (q *QueueDirective) Truth() starlark.Bool {
	return starlark.Bool(q.Name != "")
}

func (q *QueueDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(q.Name))
	return h.Sum32(), nil
}

func MakeQueueDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Queue directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &QueueDirective{Name: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid Queue directive")
}
