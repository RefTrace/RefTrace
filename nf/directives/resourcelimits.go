package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ResourceLimitsDirective)(nil)

func (r *ResourceLimitsDirective) String() string {
	return fmt.Sprintf("ResourceLimitsDirective(Cpus: %v, Disk: %v, Memory: %v, Time: %v)",
		ptrToString(r.Cpus), ptrToString(r.Disk), ptrToString(r.Memory), ptrToString(r.Time))
}

func (r *ResourceLimitsDirective) Type() string {
	return "resource_limits_directive"
}

func (r *ResourceLimitsDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (r *ResourceLimitsDirective) Truth() starlark.Bool {
	return starlark.Bool(r.Cpus != nil || r.Disk != nil || r.Memory != nil || r.Time != nil)
}

func (r *ResourceLimitsDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%v%v%v%v", ptrToString(r.Cpus), ptrToString(r.Disk), ptrToString(r.Memory), ptrToString(r.Time))))
	return h.Sum32(), nil
}

// Existing method
func (r ResourceLimitsDirective) DirectiveType() DirectiveType { return ResourceLimitsDirectiveType }

// Helper function to convert pointer to string
func ptrToString(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", v)
}

var _ starlark.Value = (*ResourceLimitsDirective)(nil)
var _ starlark.HasAttrs = (*ResourceLimitsDirective)(nil)

func (r *ResourceLimitsDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "cpus":
		if r.Cpus != nil {
			return starlark.MakeInt64(int64(*r.Cpus)), nil
		}
		return starlark.None, nil
	case "disk":
		if r.Disk != nil {
			return starlark.String(*r.Disk), nil
		}
		return starlark.None, nil
	case "memory":
		if r.Memory != nil {
			return starlark.String(*r.Memory), nil
		}
		return starlark.None, nil
	case "time":
		if r.Time != nil {
			return starlark.String(*r.Time), nil
		}
		return starlark.None, nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("resource_limits directive has no attribute %q", name))
	}
}

func (r *ResourceLimitsDirective) AttrNames() []string {
	return []string{"cpus", "disk", "memory", "time"}
}

type ResourceLimitsDirective struct {
	Cpus   *int
	Disk   *string
	Memory *string
	Time   *string
}

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
