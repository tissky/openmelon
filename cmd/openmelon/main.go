package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Example struct {
	Project          Project       `json:"project"`
	Intent           string        `json:"intent"`
	Workflow         string        `json:"workflow"`
	Stage            string        `json:"stage"`
	SkillPlusPackage string        `json:"skillplus_package"`
	Compile          CompileConfig `json:"compile"`
}

type Project struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Audience string `json:"audience"`
	Persona  string `json:"persona"`
}

type CompileConfig struct {
	Target       string            `json:"target"`
	ModelProfile string            `json:"model_profile"`
	Locale       string            `json:"locale"`
	Vars         map[string]string `json:"vars"`
}

type Config struct {
	Models  map[string]ModelConfig `json:"models"`
	Routing map[string]string      `json:"routing"`
}

type ModelConfig struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Role     string `json:"role"`
	Command  string `json:"command"`
}

type CompiledSkill struct {
	Target             string            `json:"target"`
	Package            PackageInfo       `json:"package"`
	CompiledPrompt     string            `json:"compiled_prompt"`
	RuntimeVars        map[string]string `json:"runtime_vars"`
	ModelProfile       string            `json:"model_profile"`
	Evaluation         Evaluation        `json:"evaluation"`
	ProvenanceTemplate map[string]any    `json:"provenance_template"`
	Lifecycle          map[string]any    `json:"lifecycle"`
	StageContract      map[string]any    `json:"stage_contract"`
	OutputSchema       map[string]any    `json:"output_schema"`
}

type PackageInfo struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`
}

type Evaluation struct {
	Checklist    []string `json:"checklist"`
	FailureModes []string `json:"failure_modes"`
}

type Output struct {
	Project       Project                `json:"project"`
	Workflow      string                 `json:"workflow"`
	Stage         string                 `json:"stage"`
	Intent        string                 `json:"intent"`
	ModelPlan     map[string]ModelConfig `json:"model_plan"`
	CompiledSkill CompiledSkill          `json:"compiled_skill"`
	Artifacts     []Artifact             `json:"artifacts"`
	Labels        map[string]string      `json:"labels"`
	Provenance    Provenance             `json:"provenance"`
	Review        Review                 `json:"review"`
}

type Artifact struct {
	ID      string            `json:"id"`
	Type    string            `json:"type"`
	Content string            `json:"content"`
	URI     string            `json:"uri,omitempty"`
	Model   string            `json:"model,omitempty"`
	Labels  map[string]string `json:"labels"`
}

type Provenance struct {
	ProjectID      string            `json:"project_id"`
	Workflow       string            `json:"workflow"`
	Stage          string            `json:"stage"`
	SkillPackage   PackageInfo       `json:"skill_package"`
	CompiledTarget string            `json:"compiled_target"`
	ModelProfile   string            `json:"model_profile"`
	RuntimeVars    map[string]string `json:"runtime_vars"`
	PromptHash     string            `json:"prompt_hash"`
	Generation     *GenerationTrace  `json:"generation,omitempty"`
	Source         map[string]any    `json:"source"`
}

type GenerationTrace struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Command  string `json:"command,omitempty"`
	Output   string `json:"output"`
}

type Review struct {
	Checklist    []string `json:"checklist"`
	FailureModes []string `json:"failure_modes"`
}

func main() {
	examplePath := flag.String("example", "examples/food-exploration/beef-noodles.json", "OpenMelon example input")
	compilerPath := flag.String("compiler", "../skillplus/compiler/reference", "Skill-Plus compiler PYTHONPATH")
	configPath := flag.String("config", "config/openmelon.example.json", "OpenMelon model config")
	artifactDir := flag.String("artifact-dir", "examples/food-exploration/artifacts", "Artifact output directory")
	generate := flag.Bool("generate", false, "Run configured generation model/tool and create final artifacts")
	flag.Parse()

	if err := run(*examplePath, *compilerPath, *configPath, *artifactDir, *generate); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(examplePath string, compilerPath string, configPath string, artifactDir string, generate bool) error {
	payload, err := os.ReadFile(examplePath)
	if err != nil {
		return fmt.Errorf("read example: %w", err)
	}

	var example Example
	if err := json.Unmarshal(payload, &example); err != nil {
		return fmt.Errorf("parse example: %w", err)
	}

	config, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	compiled, err := compileSkill(examplePath, compilerPath, example)
	if err != nil {
		return err
	}

	prompt := buildGenerationPrompt(example, compiled)
	promptLabels := labelsFor(example, compiled, "image_prompt")
	promptArtifact := Artifact{
		ID:      stableID(example.Project.ID + ":" + example.Workflow + ":" + example.Stage + ":" + compiled.Package.ID),
		Type:    "image_prompt",
		Content: prompt,
		Model:   modelForStage(config, "visual_concretization").Model,
		Labels:  promptLabels,
	}

	artifacts := []Artifact{promptArtifact}
	var generationTrace *GenerationTrace
	if generate {
		imageArtifact, trace, err := generateImageArtifact(config, artifactDir, prompt, example, compiled)
		if err != nil {
			return err
		}
		artifacts = append(artifacts, imageArtifact)
		generationTrace = &trace
	}

	output := Output{
		Project:       example.Project,
		Workflow:      example.Workflow,
		Stage:         example.Stage,
		Intent:        example.Intent,
		ModelPlan:     config.Models,
		CompiledSkill: compiled,
		Artifacts:     artifacts,
		Labels:        promptLabels,
		Provenance: Provenance{
			ProjectID:      example.Project.ID,
			Workflow:       example.Workflow,
			Stage:          example.Stage,
			SkillPackage:   compiled.Package,
			CompiledTarget: compiled.Target,
			ModelProfile:   compiled.ModelProfile,
			RuntimeVars:    compiled.RuntimeVars,
			PromptHash:     stableID(prompt),
			Generation:     generationTrace,
			Source:         compiled.ProvenanceTemplate,
		},
		Review: Review{
			Checklist:    compiled.Evaluation.Checklist,
			FailureModes: compiled.Evaluation.FailureModes,
		},
	}

	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("encode output: %w", err)
	}
	fmt.Println(string(encoded))
	return nil
}

func loadConfig(path string) (Config, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var config Config
	if err := json.Unmarshal(payload, &config); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return config, nil
}

func labelsFor(example Example, compiled CompiledSkill, artifactType string) map[string]string {
	return map[string]string{
		"content_vertical":  example.Workflow,
		"platform":          example.Project.Platform,
		"workflow_stage":    example.Stage,
		"skillplus_package": compiled.Package.ID,
		"skillplus_version": compiled.Package.Version,
		"model_profile":     compiled.ModelProfile,
		"artifact_type":     artifactType,
	}
}

func generateImageArtifact(config Config, artifactDir string, prompt string, example Example, compiled CompiledSkill) (Artifact, GenerationTrace, error) {
	modelKey := config.Routing["image_generation"]
	model := config.Models[modelKey]
	if model.Provider != "command" {
		return Artifact{}, GenerationTrace{}, fmt.Errorf("unsupported image generation provider: %s", model.Provider)
	}
	if model.Command == "" {
		return Artifact{}, GenerationTrace{}, fmt.Errorf("image generation command is required")
	}

	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		return Artifact{}, GenerationTrace{}, fmt.Errorf("create artifact dir: %w", err)
	}
	outputPath := filepath.Join(artifactDir, stableID(prompt)+".png")
	command := renderCommand(model.Command, map[string]string{
		"model":       model.Model,
		"output":      outputPath,
		"prompt_json": shellQuote(prompt),
	})
	cmd := exec.Command("sh", "-lc", command)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return Artifact{}, GenerationTrace{}, fmt.Errorf("run image generation: %w", err)
	}

	artifact := Artifact{
		ID:      stableID(outputPath),
		Type:    "image",
		Content: "generated image artifact",
		URI:     outputPath,
		Model:   model.Model,
		Labels:  labelsFor(example, compiled, "image"),
	}
	trace := GenerationTrace{
		Provider: model.Provider,
		Model:    model.Model,
		Command:  command,
		Output:   outputPath,
	}
	return artifact, trace, nil
}

func compileSkill(examplePath string, compilerPath string, example Example) (CompiledSkill, error) {
	packagePath := example.SkillPlusPackage
	if !filepath.IsAbs(packagePath) {
		packagePath = filepath.Clean(filepath.Join(filepath.Dir(examplePath), packagePath))
	}

	args := []string{"-m", "skillplus_compile", packagePath, "--target", example.Compile.Target, "--model-profile", example.Compile.ModelProfile, "--locale", example.Compile.Locale}
	for key, value := range example.Compile.Vars {
		args = append(args, "--var", key+"="+value)
	}
	cmd := exec.Command("python3", args...)
	cmd.Env = append(os.Environ(), "PYTHONPATH="+compilerPath)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return CompiledSkill{}, fmt.Errorf("compile skill: %w: %s", err, string(exitErr.Stderr))
		}
		return CompiledSkill{}, fmt.Errorf("compile skill: %w", err)
	}

	var compiled CompiledSkill
	if err := json.Unmarshal(out, &compiled); err != nil {
		return CompiledSkill{}, fmt.Errorf("parse compiled skill: %w", err)
	}
	return compiled, nil
}

func buildGenerationPrompt(example Example, compiled CompiledSkill) string {
	return strings.Join([]string{
		"手机随手拍的一张社媒探店照片，晚上九点半，在老小区楼下的一家普通牛肉面小店。",
		"一碗刚端上来的牛肉面放在有点发亮的塑料桌布上，拍摄者像是刚下班坐在靠墙小桌前，从胸口高度往下随手拍。",
		"画面里能看到一串钥匙、半瓶冰红茶、拆开的一次性筷子塑料膜、皱掉的纸巾、桌角贴着的点菜单。",
		"背景里有墙上褪色的菜单和后厨老板模糊的身影，旁边隐约有其他客人。",
		"牛肉面上有辣油、葱花、几片牛肉，碗边有一点汤汁溅出来，热气轻微往上冒。",
		"构图略微不正，碗不要完全居中，画面边缘可以裁掉一点饮料瓶，灯光是普通小店偏黄的顶灯，局部有一点过曝。",
		"整体像手机自动曝光的生活记录，不要商业摄影，不要电影感，不要精致摆盘，不要高级滤镜，不要文字水印。",
	}, "")
}

func modelForStage(config Config, stage string) ModelConfig {
	key := config.Routing[stage]
	return config.Models[key]
}

func renderCommand(command string, values map[string]string) string {
	for key, value := range values {
		command = strings.ReplaceAll(command, "{{"+key+"}}", value)
	}
	return command
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func stableID(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])[:16]
}
