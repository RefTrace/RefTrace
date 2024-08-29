package nf

import (
	"fmt"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

type Module struct {
	Path      string
	Processes []Process
	Includes  []IncludeStatement
}

func BuildModule(filePath string) (*Module, error) {
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		return nil, err
	}

	includeVisitor := NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.includes

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	return &Module{
		Path:      filePath,
		Processes: processes,
		Includes:  includes,
	}, nil
}

func ConvertToStarlarkModule(m *Module) *StarlarkModule {
	starlarkProcesses := make([]*StarlarkProcess, len(m.Processes))
	for i, process := range m.Processes {
		starlarkProcesses[i] = ConvertToStarlarkProcess(process)
	}

	return &StarlarkModule{
		Path:      m.Path,
		Processes: starlarkProcesses,
		Includes:  m.Includes,
	}
}

var _ starlark.Value = (*StarlarkModule)(nil)
var _ starlark.HasAttrs = (*StarlarkModule)(nil)

type StarlarkModule struct {
	Path      string
	Processes []*StarlarkProcess
	Includes  []IncludeStatement
}

func (m *StarlarkModule) String() string {
	return fmt.Sprintf("Module(%s)", m.Path)
}

func (m *StarlarkModule) Type() string {
	return "module"
}

func (m *StarlarkModule) Freeze() {}

func (m *StarlarkModule) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (m *StarlarkModule) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: module")
}

func (m *StarlarkModule) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		return starlark.String(m.Path), nil
	case "processes":
		processes := make([]starlark.Value, len(m.Processes))
		for i, p := range m.Processes {
			processes[i] = p
		}
		return starlark.NewList(processes), nil
	case "includes":
		includes := make([]starlark.Value, len(m.Includes))
		for i, inc := range m.Includes {
			includes[i] = inc
		}
		return starlark.NewList(includes), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("module has no attribute %q", name))
	}
}

func (m *StarlarkModule) AttrNames() []string {
	return []string{"path", "processes", "includes"}
}
