package nf

import "reft-go/parser"

type Module struct {
	Path      string
	Processes []Process
}

func BuildModule(filePath string) (*Module, error) {
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		return nil, err
	}
	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	return &Module{
		Path:      filePath,
		Processes: processes,
	}, nil
}
