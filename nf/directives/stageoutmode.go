package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*StageOutModeDirective)(nil)
var _ starlark.Value = (*StageOutModeDirective)(nil)
var _ starlark.HasAttrs = (*StageOutModeDirective)(nil)

func (s *StageOutModeDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "mode":
		return starlark.String(s.Mode), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("stage_out_mode directive has no attribute %q", name))
	}
}

func (s *StageOutModeDirective) AttrNames() []string {
	return []string{"mode"}
}

type StageOutModeDirective struct {
	Mode string
}

func (s *StageOutModeDirective) String() string {
	return fmt.Sprintf("StageOutModeDirective(Mode: %q)", s.Mode)
}

func (s *StageOutModeDirective) Type() string {
	return "stage_out_mode_directive"
}

func (s *StageOutModeDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (s *StageOutModeDirective) Truth() starlark.Bool {
	return starlark.Bool(s.Mode != "")
}

func (s *StageOutModeDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(s.Mode))
	return h.Sum32(), nil
}

func MakeStageOutModeDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid StageOutMode directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &StageOutModeDirective{Mode: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid StageOutMode directive")
}
