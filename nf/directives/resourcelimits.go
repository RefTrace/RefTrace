package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*ResourceLimitsDirective)(nil)

type ResourceLimitsDirective struct {
	Cpus   *int
	Disk   *string
	Memory *string
	Time   *string
}

func (a ResourceLimitsDirective) Type() DirectiveType { return ResourceLimitsDirectiveType }

func MakeResourceLimitsDirective(mce *parser.MethodCallExpression) (Directive, error) {
	var cpus *int
	var disk *string
	var memory *string
	var time *string
	if args, ok := mce.GetArguments().(*parser.TupleExpression); ok {
		if len(args.GetExpressions()) != 1 {
			return nil, errors.New("invalid resource Limits directive")
		}
		expr := args.GetExpressions()[0]
		if namedArgListExpr, ok := expr.(*parser.NamedArgumentListExpression); ok {
			exprs := namedArgListExpr.GetMapEntryExpressions()
			for _, mapEntryExpr := range exprs {
				key := mapEntryExpr.GetKeyExpression().GetText()
				val := mapEntryExpr.GetValueExpression()
				if key == "cpus" {
					if constVal, ok := val.(*parser.ConstantExpression); ok {
						if intVal, ok := constVal.GetValue().(int); ok {
							cpus = &intVal
						}
					}
				}
				if key == "disk" {
					text := val.GetText()
					disk = &text
				}
				if key == "memory" {
					text := val.GetText()
					memory = &text
				}
				if key == "time" {
					text := val.GetText()
					time = &text
				}
			}
		}
	}
	return &ResourceLimitsDirective{
		Cpus:   cpus,
		Disk:   disk,
		Memory: memory,
		Time:   time,
	}, nil
}
