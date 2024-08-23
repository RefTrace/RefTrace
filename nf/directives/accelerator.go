package directives

import (
	"reft-go/parser"
)

var _ Directive = (*Accelerator)(nil)

type Accelerator struct {
	NumGPUs int
	GPUType string
}

func (a Accelerator) Type() DirectiveType { return AcceleratorType }

func MakeAccelerator(mce *parser.MethodCallExpression) *Accelerator {
	var numGPUs int = -1
	var gpuType string = ""
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		for _, expr := range exprs {
			if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if intValue, ok := value.(int); ok {
					numGPUs = intValue
				}
			}
			if mapExpr, ok := expr.(*parser.MapExpression); ok {
				entries := mapExpr.GetMapEntryExpressions()
				for _, entry := range entries {
					if entry.GetKeyExpression().GetText() == "type" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							gpuTypeVal := constantExpr.GetValue()
							if gpuTypeStr, ok := gpuTypeVal.(string); ok {
								gpuType = gpuTypeStr
							}
						}
					}
				}
			}
		}
	}
	if numGPUs != -1 {
		return &Accelerator{NumGPUs: numGPUs, GPUType: gpuType}
	}
	return nil
}
