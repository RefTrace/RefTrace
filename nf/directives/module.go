package directives

import (
	"reft-go/parser"
)

var _ Directive = (*ModuleDirective)(nil)

type ModuleDirective struct {
	Name string
}

func (a ModuleDirective) Type() DirectiveType { return ModuleDirectiveType }

func MakeModuleDirective(mce *parser.MethodCallExpression) *ModuleDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &ModuleDirective{Name: strValue}
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &ModuleDirective{Name: gstringExpr.GetText()}
			}
		}
	}
	return nil
}
