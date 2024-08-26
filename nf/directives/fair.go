package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*FairDirective)(nil)

func (f *FairDirective) String() string {
	return fmt.Sprintf("FairDirective(Enabled: %t)", f.Enabled)
}

func (f *FairDirective) Type() string {
	return "fair_directive"
}

func (f *FairDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (f *FairDirective) Truth() starlark.Bool {
	return starlark.Bool(f.Enabled)
}

func (f *FairDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%t", f.Enabled)))
	return h.Sum32(), nil
}

type FairDirective struct {
	Enabled bool
}

func MakeFairDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid fair directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				return &FairDirective{Enabled: boolValue}, nil
			}
		}
	}
	return nil, errors.New("invalid fair directive")
}
