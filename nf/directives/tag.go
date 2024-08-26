package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*TagDirective)(nil)

type TagDirective struct {
	Tag string
}

func (t *TagDirective) String() string {
	return fmt.Sprintf("TagDirective(Tag: %q)", t.Tag)
}

func (t *TagDirective) Type() string {
	return "tag_directive"
}

func (t *TagDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (t *TagDirective) Truth() starlark.Bool {
	return starlark.Bool(t.Tag != "")
}

func (t *TagDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(t.Tag))
	return h.Sum32(), nil
}

func MakeTagDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &TagDirective{Tag: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &TagDirective{Tag: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid TagDirective directive")
}
