package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*TimeDirective)(nil)
var _ starlark.Value = (*TimeDirective)(nil)
var _ starlark.HasAttrs = (*TimeDirective)(nil)

func (t *TimeDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "duration":
		return starlark.String(t.Duration), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("time directive has no attribute %q", name))
	}
}

func (t *TimeDirective) AttrNames() []string {
	return []string{"duration"}
}

type TimeDirective struct {
	Duration string
}

func (t *TimeDirective) String() string {
	return fmt.Sprintf("TimeDirective(Duration: %q)", t.Duration)
}

func (t *TimeDirective) Type() string {
	return "time_directive"
}

func (t *TimeDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (t *TimeDirective) Truth() starlark.Bool {
	return starlark.Bool(t.Duration != "")
}

func (t *TimeDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(t.Duration))
	return h.Sum32(), nil
}

func MakeTimeDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Time directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &TimeDirective{Duration: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid Time directive")
}
