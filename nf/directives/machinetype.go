package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*MachineTypeDirective)(nil)

type MachineTypeDirective struct {
	MachineType string
}

func (a MachineTypeDirective) Type() DirectiveType { return MachineTypeDirectiveType }

func MakeMachineTypeDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid machine type directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MachineTypeDirective{MachineType: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid machine type directive")
}
