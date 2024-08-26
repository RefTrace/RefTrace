package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*CacheDirective)(nil)

func (c *CacheDirective) String() string {
	return fmt.Sprintf("CacheDirective(Enabled: %t, Deep: %t, Lenient: %t)", c.Enabled, c.Deep, c.Lenient)
}

func (c *CacheDirective) Type() string {
	return "cache_directive"
}

func (c *CacheDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (c *CacheDirective) Truth() starlark.Bool {
	return starlark.Bool(c.Enabled)
}

func (c *CacheDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%t%t%t", c.Enabled, c.Deep, c.Lenient)))
	return h.Sum32(), nil
}

type CacheDirective struct {
	Enabled bool
	Deep    bool
	Lenient bool
}

func MakeCacheDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid cache directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolExpr, ok := constantExpr.GetValue().(bool); ok {
				if boolExpr {
					return &CacheDirective{Enabled: true}, nil
				} else {
					return &CacheDirective{Enabled: false}, nil
				}
			}
			if stringExpr, ok := constantExpr.GetValue().(string); ok {
				if stringExpr == "deep" {
					return &CacheDirective{Enabled: true, Deep: true}, nil
				}
				if stringExpr == "lenient" {
					return &CacheDirective{Enabled: true, Lenient: true}, nil
				}
			}
		}
	}
	return nil, errors.New("invalid cache directive")
}
