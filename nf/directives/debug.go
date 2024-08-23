package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*DebugDirective)(nil)

type DebugDirective struct {
	Enabled bool
}

func (a DebugDirective) Type() DirectiveType { return DebugDirectiveType }

func MakeDebugDirective(mce *parser.MethodCallExpression) (*DebugDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid debug directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				return &DebugDirective{Enabled: boolValue}, nil
			}
		}
	}
	return nil, errors.New("invalid debug directive")
}
