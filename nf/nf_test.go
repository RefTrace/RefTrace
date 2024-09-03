package nf

import (
	"os"
	"path/filepath"
	"reft-go/parser"
	"testing"
)

func getTestDataDir() string {
	if dir := os.Getenv("NF_CORE_TEST_DATA"); dir != "" {
		return dir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get the home directory, fall back to "testdata"
		return "testdata"
	}

	return filepath.Join(homeDir, "reft-testdata")
}

func TestSarekEntireMain(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "sarek_entire_main.nf")
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
}

func TestSimpleWorkflow(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "simple_workflow.nf")
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
}

func TestSimpleProcess(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "simple_process.nf")
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
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "closure_directive.nf")
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
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "channel_frompath.nf")
	_, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}
	// cls := ast.GetClasses()[0]
	// TODO: fix this test
}

func TestCountProcesses(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-core")
	modules, err := ProcessDirectory(filePath)
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
	filePath := filepath.Join(getTestDataDir(), "nf-core", "airrflow")
	modules, err := ProcessDirectory(filePath)
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
	filePath := filepath.Join(getTestDataDir(), "nf-core", "airrflow/modules/local/enchantr/report_file_size.nf")
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
