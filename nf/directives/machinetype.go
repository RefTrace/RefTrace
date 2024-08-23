package directives

import (
	"reft-go/parser"
)

var _ Directive = (*MachineTypeDirective)(nil)

type MachineTypeDirective struct {
	MachineType string
}

func (a MachineTypeDirective) Type() DirectiveType { return MachineTypeDirectiveType }

func MakeMachineTypeDirective(mce *parser.MethodCallExpression) *MachineTypeDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MachineTypeDirective{MachineType: strValue}
			}
		}
	}
	return nil
}
