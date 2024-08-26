package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"
	"strings"

	"go.starlark.net/starlark"
)

var _ Directive = (*ClusterOptions)(nil)

func (c *ClusterOptions) String() string {
	return fmt.Sprintf("ClusterOptions(%q)", c.Options)
}

func (c *ClusterOptions) Type() string {
	return "cluster_options"
}

func (c *ClusterOptions) Freeze() {
	// No mutable fields, so no action needed
}

func (c *ClusterOptions) Truth() starlark.Bool {
	return starlark.Bool(c.Options != "")
}

func (c *ClusterOptions) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(c.Options))
	return h.Sum32(), nil
}

type ClusterOptions struct {
	Options string
}

func MakeClusterOptions(mce *parser.MethodCallExpression) (Directive, error) {
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
		return &ClusterOptions{Options: joinedOptions}, nil
	}
	return nil, errors.New("invalid clusterOptions directive")
}
