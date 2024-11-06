package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*SpackDirective)(nil)
var _ starlark.Value = (*SpackDirective)(nil)
var _ starlark.HasAttrs = (*SpackDirective)(nil)

func (s *SpackDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "dependencies":
		return starlark.String(s.Dependencies), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("spack directive has no attribute %q", name))
	}
}

func (s *SpackDirective) AttrNames() []string {
	return []string{"dependencies"}
}

type SpackDirective struct {
	Dependencies string
	line         int
}

func (s *SpackDirective) Line() int {
	return s.line
}

func (s *SpackDirective) String() string {
	return fmt.Sprintf("SpackDirective(Dependencies: %q)", s.Dependencies)
}

func (s *SpackDirective) Type() string {
	return "spack_directive"
}

func (s *SpackDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (s *SpackDirective) Truth() starlark.Bool {
	return starlark.Bool(s.Dependencies != "")
}

func (s *SpackDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(s.Dependencies))
	return h.Sum32(), nil
}

func MakeSpackDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Spack directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &SpackDirective{Dependencies: strValue, line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid Spack directive")
}
