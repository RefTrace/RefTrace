package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*LabelDirective)(nil)
var _ starlark.Value = (*LabelDirective)(nil)
var _ starlark.HasAttrs = (*LabelDirective)(nil)

func (l *LabelDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "label":
		return starlark.String(l.Label), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("label directive has no attribute %q", name))
	}
}

func (l *LabelDirective) AttrNames() []string {
	return []string{"label"}
}

type LabelDirective struct {
	Label string
}

func (l *LabelDirective) String() string {
	return fmt.Sprintf("LabelDirective(Label: %q)", l.Label)
}

func (l *LabelDirective) Type() string {
	return "label_directive"
}

func (l *LabelDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (l *LabelDirective) Truth() starlark.Bool {
	return starlark.Bool(l.Label != "")
}

func (l *LabelDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(l.Label))
	return h.Sum32(), nil
}

func MakeLabelDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid label directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &LabelDirective{Label: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid label directive")
}
