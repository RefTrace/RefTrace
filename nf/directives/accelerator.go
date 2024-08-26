package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*Accelerator)(nil)
var _ starlark.Value = (*Accelerator)(nil)
var _ starlark.HasAttrs = (*Accelerator)(nil)

func (a *Accelerator) Attr(name string) (starlark.Value, error) {
	switch name {
	case "num_gpus":
		return starlark.MakeInt(a.NumGPUs), nil
	case "gpu_type":
		return starlark.String(a.GPUType), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("accelerator directive has no attribute %q", name))
	}
}

func (a *Accelerator) AttrNames() []string {
	return []string{"num_gpus", "gpu_type"}
}

func (a *Accelerator) String() string {
	return fmt.Sprintf("Accelerator(NumGPUs: %d, GPUType: %q)", a.NumGPUs, a.GPUType)
}

func (a *Accelerator) Type() string {
	return "accelerator"
}

func (a *Accelerator) Freeze() {
	// No mutable fields, so no action needed
}

func (a *Accelerator) Truth() starlark.Bool {
	return starlark.Bool(a.NumGPUs > 0)
}

func (a *Accelerator) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%d%s", a.NumGPUs, a.GPUType)))
	return h.Sum32(), nil
}

type Accelerator struct {
	NumGPUs int
	GPUType string
}

func MakeAccelerator(mce *parser.MethodCallExpression) (Directive, error) {
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
		return &Accelerator{NumGPUs: numGPUs, GPUType: gpuType}, nil
	}
	return nil, errors.New("invalid accelerator directive")
}
