package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"
	"strings"

	"go.starlark.net/starlark"
)

var _ Directive = (*Shell)(nil)
var _ starlark.Value = (*Shell)(nil)
var _ starlark.HasAttrs = (*Shell)(nil)

func (s *Shell) Attr(name string) (starlark.Value, error) {
	switch name {
	case "command":
		return starlark.String(s.Command), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("shell directive has no attribute %q", name))
	}
}

func (s *Shell) AttrNames() []string {
	return []string{"command"}
}

type Shell struct {
	Command string
}

func (s *Shell) String() string {
	return fmt.Sprintf("Shell(Command: %q)", s.Command)
}

func (s *Shell) Type() string {
	return "shell_directive"
}

func (s *Shell) Freeze() {
	// No mutable fields, so no action needed
}

func (s *Shell) Truth() starlark.Bool {
	return starlark.Bool(s.Command != "")
}

func (s *Shell) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(s.Command))
	return h.Sum32(), nil
}

func MakeShellDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		options := []string{}
		for _, expr := range exprs {
			if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
				text := constantExpr.GetText()
				options = append(options, text)
			}
		}
		joinedOptions := strings.Join(options, " ")
		return &Shell{Command: joinedOptions}, nil
	}
	return nil, errors.New("invalid Shell directive")
}
