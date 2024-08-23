package directives

import (
	"reft-go/parser"
)

var _ Directive = (*Arch)(nil)

type Arch struct {
	Name   string
	Target string
}

func (a Arch) Type() DirectiveType { return ArchType }

func MakeArch(mce *parser.MethodCallExpression) *Arch {
	var name string = ""
	var target string = ""
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		for _, expr := range exprs {
			if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
				name = constantExpr.GetText()
			}
			if mapExpr, ok := expr.(*parser.MapExpression); ok {
				entries := mapExpr.GetMapEntryExpressions()
				for _, entry := range entries {
					if entry.GetKeyExpression().GetText() == "target" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							target = constantExpr.GetText()
						}
					}
				}
			}
		}
	}
	if name != "" {
		return &Arch{Name: name, Target: target}
	}
	return nil
}
