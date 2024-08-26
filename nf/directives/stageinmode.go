package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*StageInModeDirective)(nil)
var _ starlark.Value = (*StageInModeDirective)(nil)
var _ starlark.HasAttrs = (*StageInModeDirective)(nil)

func (s *StageInModeDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "mode":
		return starlark.String(s.Mode), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("stage_in_mode directive has no attribute %q", name))
	}
}

func (s *StageInModeDirective) AttrNames() []string {
	return []string{"mode"}
}

type StageInModeDirective struct {
	Mode string
}

func (s *StageInModeDirective) String() string {
	return fmt.Sprintf("StageInModeDirective(Mode: %q)", s.Mode)
}

func (s *StageInModeDirective) Type() string {
	return "stage_in_mode_directive"
}

func (s *StageInModeDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (s *StageInModeDirective) Truth() starlark.Bool {
	return starlark.Bool(s.Mode != "")
}

func (s *StageInModeDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(s.Mode))
	return h.Sum32(), nil
}

func MakeStageInModeDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid StageInMode directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &StageInModeDirective{Mode: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid StageInMode directive")
}
