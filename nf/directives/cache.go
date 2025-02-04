package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (c *CacheDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(c.Line()),
		Directive: &pb.Directive_Cache{
			Cache: &pb.CacheDirective{
				Enabled: c.Enabled,
				Deep:    c.Deep,
				Lenient: c.Lenient,
			},
		},
	}
}

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

var _ starlark.Value = (*CacheDirective)(nil)
var _ starlark.HasAttrs = (*CacheDirective)(nil)

func (c *CacheDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "enabled":
		return starlark.Bool(c.Enabled), nil
	case "deep":
		return starlark.Bool(c.Deep), nil
	case "lenient":
		return starlark.Bool(c.Lenient), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("cache directive has no attribute %q", name))
	}
}

func (c *CacheDirective) AttrNames() []string {
	return []string{"enabled", "deep", "lenient"}
}

type CacheDirective struct {
	Enabled bool
	Deep    bool
	Lenient bool
	line    int
}

func (c *CacheDirective) Line() int {
	return c.line
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
					return &CacheDirective{Enabled: true, line: mce.GetLineNumber()}, nil
				} else {
					return &CacheDirective{Enabled: false, line: mce.GetLineNumber()}, nil
				}
			}
			if stringExpr, ok := constantExpr.GetValue().(string); ok {
				if stringExpr == "deep" {
					return &CacheDirective{Enabled: true, Deep: true, line: mce.GetLineNumber()}, nil
				}
				if stringExpr == "lenient" {
					return &CacheDirective{Enabled: true, Lenient: true, line: mce.GetLineNumber()}, nil
				}
			}
		}
	}
	return nil, errors.New("invalid cache directive")
}
