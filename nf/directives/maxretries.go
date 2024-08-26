package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*MaxRetriesDirective)(nil)

type MaxRetriesDirective struct {
	Num int
}

func (m *MaxRetriesDirective) String() string {
	return fmt.Sprintf("MaxRetriesDirective(Num: %d)", m.Num)
}

func (m *MaxRetriesDirective) Type() string {
	return "max_retries_directive"
}

func (m *MaxRetriesDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MaxRetriesDirective) Truth() starlark.Bool {
	return starlark.Bool(m.Num > 0)
}

func (m *MaxRetriesDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%d", m.Num)))
	return h.Sum32(), nil
}

func MakeMaxRetriesDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max retries directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &MaxRetriesDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max retries directive")
}
