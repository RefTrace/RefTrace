package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*CacheDirective)(nil)

type CacheDirective struct {
	Enabled bool
	Deep    bool
	Lenient bool
}

func (a CacheDirective) Type() DirectiveType { return CacheDirectiveType }

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
