package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*PenvDirective)(nil)

type PenvDirective struct {
	Environment string
}

func (p *PenvDirective) String() string {
	return fmt.Sprintf("PenvDirective(Environment: %q)", p.Environment)
}

func (p *PenvDirective) Type() string {
	return "penv_directive"
}

func (p *PenvDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (p *PenvDirective) Truth() starlark.Bool {
	return starlark.Bool(p.Environment != "")
}

func (p *PenvDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(p.Environment))
	return h.Sum32(), nil
}

func MakePenvDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid penv directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &PenvDirective{Environment: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid penv directive")
}
