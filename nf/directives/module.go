package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ModuleDirective)(nil)
var _ starlark.Value = (*ModuleDirective)(nil)
var _ starlark.HasAttrs = (*ModuleDirective)(nil)

func (m *ModuleDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(m.Name), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("module directive has no attribute %q", name))
	}
}

func (m *ModuleDirective) AttrNames() []string {
	return []string{"name"}
}

type ModuleDirective struct {
	Name string
	line int
}

func (m *ModuleDirective) Line() int {
	return m.line
}

func (m *ModuleDirective) String() string {
	return fmt.Sprintf("ModuleDirective(Name: %q)", m.Name)
}

func (m *ModuleDirective) Type() string {
	return "module_directive"
}

func (m *ModuleDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *ModuleDirective) Truth() starlark.Bool {
	return starlark.Bool(m.Name != "")
}

func (m *ModuleDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(m.Name))
	return h.Sum32(), nil
}

func MakeModuleDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &ModuleDirective{Name: strValue, line: mce.GetLineNumber()}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &ModuleDirective{Name: gstringExpr.GetText(), line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid module directive")
}
