package nf

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reft-go/parser"
	"sync"
	"testing"
)

func TestSarekEntireMain(t *testing.T) {
	filePath := filepath.Join("../parser/testdata", "sarek_entire_main.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}
	paramVisitor := NewParamVisitor()
	paramVisitor.VisitBlockStatement(ast.StatementBlock)
	params := paramVisitor.GetSortedParams()
	if len(params) != 65 {
		t.Fatalf("Expected 65 params, got %d", len(params))
	}
	includeVisitor := NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.GetSortedIncludes()
	if len(includes) != 8 {
		t.Fatalf("Expected 8 includes, got %d", len(includes))
	}
	workflowVisitor := NewWorkflowVisitor()
	workflowVisitor.VisitBlockStatement(ast.StatementBlock)
	workflows := workflowVisitor.workflows
	if len(workflows) != 2 {
		t.Fatalf("Expected 2 workflows, got %d", len(workflows))
	}
	//stcVisitor := NewStcVisitor(cls)
	//stcVisitor.VisitBlockStatement(ast.StatementBlock)
}

func TestSimpleWorkflow(t *testing.T) {
	filePath := filepath.Join("./testdata", "simple_workflow.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	workflowVisitor := NewWorkflowVisitor()
	workflowVisitor.VisitBlockStatement(ast.StatementBlock)
	workflows := workflowVisitor.workflows
	if len(workflows) != 2 {
		t.Fatalf("Expected 2 workflows, got %d", len(workflows))
	}
	workflow := workflows[0]
	if len(workflow.Takes) != 2 {
		t.Fatalf("Expected 2 takes, got %d", len(workflow.Takes))
	}
	if len(workflow.Emits) != 2 {
		t.Fatalf("Expected 2 emits, got %d", len(workflow.Emits))
	}
	//stcVisitor := NewStcVisitor(cls)
	//stcVisitor.VisitBlockStatement(ast.StatementBlock)
}

func TestSimpleProcess(t *testing.T) {
	filePath := filepath.Join("./testdata", "simple_process.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	if len(processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(processes))
	}
	directives := processes[0].Directives
	if len(directives) != 51 {
		t.Fatalf("Expected 51 directives, got %d", len(directives))
	}
}

func TestClosureDirective(t *testing.T) {
	filePath := filepath.Join("./testdata", "closure_directive.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	if len(processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(processes))
	}
	directives := processes[0].Directives
	if len(directives) != 2 {
		t.Fatalf("Expected 2 directives, got %d", len(directives))
	}
}

func TestChannelFromPath(t *testing.T) {
	filePath := filepath.Join("./testdata", "channel_frompath.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}
	cls := ast.GetClasses()[0]
	stcVisitor := NewStcVisitor(cls)
	stcVisitor.VisitBlockStatement(ast.StatementBlock)
}

func processDirectory(dir string) ([]*Module, error) {
	var modules []*Module
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".nf" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				module, err := BuildModule(path)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("error processing file %s: %v", path, err))
					mu.Unlock()
					return
				}
				mu.Lock()
				modules = append(modules, module)
				mu.Unlock()
			}(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("encountered %d errors during processing: %v", len(errors), errors)
	}

	return modules, nil
}

func TestCountProcesses(t *testing.T) {
	filePath := filepath.Join("../parser/testdata", "nf-core")
	modules, err := processDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	totalProcesses := 0
	for _, module := range modules {
		totalProcesses += len(module.Processes)
	}

	t.Logf("Found %d modules with a total of %d processes", len(modules), totalProcesses)
}

func TestCountProcessesAirrflow(t *testing.T) {
	filePath := filepath.Join("../parser/testdata", "nf-core", "airrflow")
	modules, err := processDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	totalProcesses := 0
	for _, module := range modules {
		totalProcesses += len(module.Processes)
	}

	t.Logf("Found %d modules with a total of %d processes", len(modules), totalProcesses)
}

func TestIfStatementProcess(t *testing.T) {
	filePath := filepath.Join("../parser/testdata", "nf-core", "airrflow/modules/local/enchantr/report_file_size.nf")
	module, err := BuildModule(filePath)
	if err != nil {
		t.Fatalf("Failed to process file: %v", err)
	}
	if len(module.Processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(module.Processes))
	}
	process := module.Processes[0]
	if len(process.Directives) != 4 {
		t.Fatalf("Expected 4 directives, got %d", len(process.Directives))
	}
}
