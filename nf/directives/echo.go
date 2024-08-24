package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*EchoDirective)(nil)

type EchoDirective struct {
	Enabled bool
}

func (a EchoDirective) Type() DirectiveType { return EchoDirectiveType }

func MakeEchoDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid echo directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				return &EchoDirective{Enabled: boolValue}, nil
			}
		}
	}
	return nil, errors.New("invalid echo directive")
}
