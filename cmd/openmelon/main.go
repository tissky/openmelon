package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/eight-acres-lab/openmelon/internal/generation"
	"github.com/eight-acres-lab/openmelon/internal/project"
	"github.com/eight-acres-lab/openmelon/internal/skillplus"
	"github.com/eight-acres-lab/openmelon/internal/workflow"
)

func main() {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	projectFlag := fs.String("project", "examples/food-exploration/project.json", "Path to project.json")
	workflowFlag := fs.String("workflow", "", "Workflow ID to execute (default: first workflow)")
	intentFlag := fs.String("intent", "", "Free-text intent passed to the skill compiler")
	artifactDir := fs.String("artifact-dir", ".openmelon/artifacts", "Directory to write artifact files")
	compilerPath := fs.String("compiler", "../skillplus", "PYTHONPATH for the Skill-Plus Python compiler")
	generate := fs.Bool("generate", false, "Execute generation step (requires configured shell provider)")
	generateCmd := fs.String("generate-cmd", "", "Shell command for generation (stdin=prompt, stdout=content)")
	timeoutSec := fs.Int("timeout", 120, "Total execution timeout in seconds")
	provenancePath := fs.String("provenance", "", "Path to append-only provenance JSONL (default: <artifact-dir>/provenance.jsonl)")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutSec)*time.Second)
	defer cancel()

	if err := run(ctx, runOpts{
		projectPath:    *projectFlag,
		workflowID:     *workflowFlag,
		intent:         *intentFlag,
		artifactDir:    *artifactDir,
		compilerPath:   *compilerPath,
		doGenerate:     *generate,
		generateCmd:    *generateCmd,
		provenancePath: *provenancePath,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "[openmelon] error: %v\n", err)
		os.Exit(1)
	}
}

type runOpts struct {
	projectPath    string
	workflowID     string
	intent         string
	artifactDir    string
	compilerPath   string
	doGenerate     bool
	generateCmd    string
	provenancePath string
}

func run(ctx context.Context, opts runOpts) error {
	// 1. Load project
	fmt.Printf("[openmelon] loading project: %s\n", opts.projectPath)
	proj, err := project.Load(opts.projectPath)
	if err != nil {
		return fmt.Errorf("load project: %w", err)
	}
	fmt.Printf("[openmelon] project: %s (%s)\n", proj.Name, proj.Platform)

	// 2. Load workflows from project file
	workflows, err := workflow.LoadWorkflows(opts.projectPath)
	if err != nil {
		return fmt.Errorf("load workflows: %w", err)
	}

	// 3. Select workflow
	var wfDef *workflow.WorkflowDefinition
	if opts.workflowID != "" {
		var ok bool
		wfDef, ok = workflows[opts.workflowID]
		if !ok {
			return fmt.Errorf("workflow %q not found in project", opts.workflowID)
		}
	} else {
		for _, wf := range workflows {
			wfDef = wf
			break
		}
	}
	fmt.Printf("[openmelon] workflow: %s (%d stages)\n", wfDef.ID, len(wfDef.Stages))

	// 4. Set up compiler
	compiler := &skillplus.Compiler{CompilerPath: opts.compilerPath}

	// 5. Set up provider (optional)
	var provider generation.Provider
	if opts.doGenerate && opts.generateCmd != "" {
		provider = &generation.ShellProvider{Command: opts.generateCmd}
	}

	// 6. Resolve provenance path
	provPath := opts.provenancePath
	if provPath == "" {
		provPath = filepath.Join(opts.artifactDir, "provenance.jsonl")
	}

	// 7. Run engine
	engine := &workflow.Engine{}
	req := &workflow.RunRequest{
		Project:        proj,
		WorkflowDef:    wfDef,
		Intent:         opts.intent,
		ArtifactDir:    opts.artifactDir,
		CompilerPath:   opts.compilerPath,
		ProvenancePath: provPath,
		Compiler:       compiler,
		Provider:       provider,
		Generate:       opts.doGenerate,
	}

	fmt.Printf("[openmelon] running workflow stages...\n")
	results, err := engine.Run(ctx, req)
	if err != nil {
		return fmt.Errorf("engine run: %w", err)
	}

	for _, r := range results {
		fmt.Printf("[openmelon] stage %q → artifact %s written to %s\n",
			r.Stage, r.Artifact.ID, opts.artifactDir)
	}
	fmt.Printf("[openmelon] done. %d stage(s) completed.\n", len(results))
	return nil
}
