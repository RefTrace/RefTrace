package directives

import (
	"reft-go/parser"
)

var _ Directive = (*CacheDirective)(nil)

type CacheDirective struct {
	Enabled bool
	Deep    bool
	Lenient bool
}

func (a CacheDirective) Type() DirectiveType { return CacheDirectiveType }

func MakeCacheDirective(mce *parser.MethodCallExpression) *CacheDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolExpr, ok := constantExpr.GetValue().(bool); ok {
				if boolExpr {
					return &CacheDirective{Enabled: true}
				} else {
					return &CacheDirective{Enabled: false}
				}
			}
			if stringExpr, ok := constantExpr.GetValue().(string); ok {
				if stringExpr == "deep" {
					return &CacheDirective{Enabled: true, Deep: true}
				}
				if stringExpr == "lenient" {
					return &CacheDirective{Enabled: true, Lenient: true}
				}
			}
		}
	}
	return nil
}
