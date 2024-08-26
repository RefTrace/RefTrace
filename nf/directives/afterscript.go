package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*AfterScript)(nil)

func (a *AfterScript) String() string {
	return fmt.Sprintf("AfterScript(%q)", a.Script)
}

func (a *AfterScript) Type() string {
	return "afterscript"
}

func (a *AfterScript) Freeze() {
	// No mutable fields, so no action needed
}

func (a *AfterScript) Truth() starlark.Bool {
	return starlark.Bool(a.Script != "")
}

func (a *AfterScript) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(a.Script))
	return h.Sum32(), nil
}

type AfterScript struct {
	Script string
}

func MakeAfterScript(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &AfterScript{Script: strValue}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &AfterScript{Script: gstringExpr.GetText()}, nil
			}
		}
	}
	return nil, errors.New("invalid afterScript directive")
}
