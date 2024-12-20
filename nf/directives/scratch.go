package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (s *ScratchDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(s.Line()),
		Directive: &pb.Directive_Scratch{
			Scratch: &pb.ScratchDirective{
				Enabled:   s.Enabled,
				Directory: s.Directory,
			},
		},
	}
}

var _ Directive = (*ScratchDirective)(nil)
var _ starlark.Value = (*ScratchDirective)(nil)
var _ starlark.HasAttrs = (*ScratchDirective)(nil)

func (s *ScratchDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "enabled":
		return starlark.Bool(s.Enabled), nil
	case "directory":
		return starlark.String(s.Directory), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("scratch directive has no attribute %q", name))
	}
}

func (s *ScratchDirective) AttrNames() []string {
	return []string{"enabled", "directory"}
}

type ScratchDirective struct {
	Enabled   bool
	Directory string
	line      int
}

func (s *ScratchDirective) Line() int {
	return s.line
}

func (s *ScratchDirective) String() string {
	return fmt.Sprintf("ScratchDirective(Enabled: %t, Directory: %q)", s.Enabled, s.Directory)
}

func (s *ScratchDirective) Type() string {
	return "scratch_directive"
}

func (s *ScratchDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (s *ScratchDirective) Truth() starlark.Bool {
	return starlark.Bool(s.Enabled)
}

func (s *ScratchDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%t", s.Enabled)))
	h.Write([]byte(s.Directory))
	return h.Sum32(), nil
}

func MakeScratchDirective(mce *parser.MethodCallExpression) (Directive, error) {
	enabled := false
	directory := ""
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Scratch directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				enabled = true
				directory = strValue
			}
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				enabled = boolValue
			}
		}
		if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
			enabled = true
			directory = gstringExpr.GetText()
		}
	}
	return &ScratchDirective{Enabled: enabled, Directory: directory, line: mce.GetLineNumber()}, nil
}
