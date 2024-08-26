package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ErrorStrategyDirective)(nil)

func (e *ErrorStrategyDirective) String() string {
	return fmt.Sprintf("ErrorStrategyDirective(Strategy: %q)", e.Strategy)
}

func (e *ErrorStrategyDirective) Type() string {
	return "error_strategy_directive"
}

func (e *ErrorStrategyDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (e *ErrorStrategyDirective) Truth() starlark.Bool {
	return starlark.Bool(e.Strategy != "")
}

func (e *ErrorStrategyDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(e.Strategy))
	return h.Sum32(), nil
}

type ErrorStrategyDirective struct {
	Strategy string
}

func MakeErrorStrategyDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid error strategy directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &ErrorStrategyDirective{Strategy: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid error strategy directive")
}
