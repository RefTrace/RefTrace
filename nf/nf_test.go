package nf

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reft-go/parser"
	"sync"
	"sync/atomic"
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

func processFile(filePath string) (int, error) {
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		return 0, err
	}
	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	return len(processes), nil
}

func processDirectory(dir string) (int64, int64, error) {
	var totalFiles, totalProcesses int64
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
				nProcesses, err := processFile(path)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("error processing file %s: %v", path, err))
					mu.Unlock()
					return
				}
				atomic.AddInt64(&totalFiles, 1)
				atomic.AddInt64(&totalProcesses, int64(nProcesses))
			}(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return totalFiles, totalProcesses, err
	}

	if len(errors) > 0 {
		return totalFiles, totalProcesses, fmt.Errorf("encountered %d errors during processing: %v", len(errors), errors)
	}

	return totalFiles, totalProcesses, nil
}

func TestCountProcesses(t *testing.T) {
	filePath := filepath.Join("../parser/testdata", "nf-core", "sarek")
	nFiles, nProcesses, err := processDirectory(filePath)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}
	t.Logf("Processed %d files and found %d processes", nFiles, nProcesses)
}
