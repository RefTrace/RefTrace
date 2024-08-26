package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ExecutorDirective)(nil)

func (e *ExecutorDirective) String() string {
	return fmt.Sprintf("ExecutorDirective(Executor: %q)", e.Executor)
}

func (e *ExecutorDirective) Type() string {
	return "executor_directive"
}

func (e *ExecutorDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (e *ExecutorDirective) Truth() starlark.Bool {
	return starlark.Bool(e.Executor != "")
}

func (e *ExecutorDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(e.Executor))
	return h.Sum32(), nil
}

var _ starlark.Value = (*ExecutorDirective)(nil)
var _ starlark.HasAttrs = (*ExecutorDirective)(nil)

func (e *ExecutorDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "executor":
		return starlark.String(e.Executor), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("executor directive has no attribute %q", name))
	}
}

func (e *ExecutorDirective) AttrNames() []string {
	return []string{"executor"}
}

type ExecutorDirective struct {
	Executor string
}

func MakeExecutorDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid executor directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &ExecutorDirective{Executor: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid executor directive")
}
