package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*MaxSubmitAwaitDirective)(nil)

type MaxSubmitAwaitDirective struct {
	MaxSubmitAwait string
}

func (m *MaxSubmitAwaitDirective) String() string {
	return fmt.Sprintf("MaxSubmitAwaitDirective(MaxSubmitAwait: %q)", m.MaxSubmitAwait)
}

func (m *MaxSubmitAwaitDirective) Type() string {
	return "max_submit_await_directive"
}

func (m *MaxSubmitAwaitDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MaxSubmitAwaitDirective) Truth() starlark.Bool {
	return starlark.Bool(m.MaxSubmitAwait != "")
}

func (m *MaxSubmitAwaitDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(m.MaxSubmitAwait))
	return h.Sum32(), nil
}

func MakeMaxSubmitAwaitDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max submit await directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MaxSubmitAwaitDirective{MaxSubmitAwait: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max submit await directive")
}
