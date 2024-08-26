package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*Arch)(nil)
var _ starlark.Value = (*Arch)(nil)
var _ starlark.HasAttrs = (*Arch)(nil)

func (a *Arch) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(a.Name), nil
	case "target":
		return starlark.String(a.Target), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("arch directive has no attribute %q", name))
	}
}

func (a *Arch) AttrNames() []string {
	return []string{"name", "target"}
}

func (a *Arch) String() string {
	return fmt.Sprintf("Arch(Name: %q, Target: %q)", a.Name, a.Target)
}

func (a *Arch) Type() string {
	return "arch"
}

func (a *Arch) Freeze() {
	// No mutable fields, so no action needed
}

func (a *Arch) Truth() starlark.Bool {
	return starlark.Bool(a.Name != "")
}

func (a *Arch) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(a.Name))
	h.Write([]byte(a.Target))
	return h.Sum32(), nil
}

type Arch struct {
	Name   string
	Target string
}

func MakeArch(mce *parser.MethodCallExpression) (Directive, error) {
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
		return &Arch{Name: name, Target: target}, nil
	}
	return nil, errors.New("invalid arch directive")
}
