# Data Model: OpenMelon Phase 2 — Runtime Skeleton

**Date**: 2026-05-04  
**Branch**: `001-phase2-runtime-skeleton`

## Entity Overview

```
                        ┌─────────────┐
                        │   Project   │
                        │  (from file)│
                        └──────┬──────┘
                               │ 1
                               │
                       ┌───────┴────────┐
                       │    RunConfig   │
                       │ (workflow+flags)│
                       └───────┬────────┘
                               │ drives
                       ┌───────▼────────┐
                       │     Engine     │
                       │  (orchestrator)│
                       └──┬────────┬────┘
                          │        │
              ┌───────────▼──┐  ┌──▼────────────┐
              │  Compiler    │  │    Provider   │
              │ (Skill-Plus) │  │  (shell/http) │
              └───────┬──────┘  └──────┬────────┘
                      │                │
              ┌───────▼──────┐  ┌──────▼────────┐
              │ CompiledSkill│  │    Artifact   │
              └──────────────┘  └──────┬────────┘
                                       │
                               ┌───────▼────────┐
                               │   Provenance   │
                               │  (JSONL append)│
                               └────────────────┘
```

---

## Entities

### Project

存储创作项目的持久上下文，从 JSON 文件加载。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ID | string | ✅ | 唯一项目标识符，用于 artifact ID 哈希 |
| Name | string | ✅ | 人类可读项目名称 |
| Platform | string | ✅ | 目标发布平台（e.g. "xiaohongshu", "douyin"）|
| Audience | string | | 目标受众描述 |
| Persona | string | | 创作者人设描述 |
| Memory | map[string]string | | 项目级记忆键值对 |
| Constraints | []string | | 约束列表（e.g. "no food waste imagery"）|

**加载规则**: ID、Name、Platform 缺失 → 返回带字段名的错误；文件不存在 → wrap os.ErrNotExist。

**现有包**: `internal/project/project.go` — 已有 struct，**需添加 `Load(path string) (*Project, error)` 函数**

---

### WorkflowDefinition

Workflow 的可序列化定义，从 JSON 文件加载。（区别于 `internal/workflow/workflow.go` 的运行时 Workflow 结构体）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ID | string | ✅ | Workflow 唯一标识符 |
| Name | string | ✅ | 人类可读名称 |
| Vertical | string | | 内容垂类（e.g. "food", "travel"）|
| Stages | []StageDefinition | ✅ | 有序 stage 列表 |

**Phase 2**: Workflow 定义内嵌在 project.json 的 `workflows` 字段，或通过 `--workflow-file` 单独传入。

---

### StageDefinition

一个 workflow stage 的静态描述。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Stage | workflow.Stage | ✅ | stage 类型常量 |
| SkillPlusPackage | string | ✅ | .skillplus 包路径（相对或绝对）|
| CompileTarget | string | ✅ | 编译目标（e.g. "openmelon"）|
| ModelProfile | string | ✅ | 模型 profile key，对应 config 中的 models map |
| Vars | map[string]string | | 传递给 Skill-Plus compiler 的变量 |
| Locale | string | | 编译 locale，默认 "zh-CN" |

---

### CompileRequest / CompileResult

Compiler 的输入/输出边界。

**CompileRequest**:

| 字段 | 类型 | 说明 |
|------|------|------|
| PackagePath | string | .skillplus 包路径 |
| Target | string | 编译目标 |
| ModelProfile | string | 模型 profile |
| Locale | string | 语言 locale |
| Vars | map[string]string | 模板变量 |

**CompileResult** (wraps `internal/skillplus/CompiledSkill`):

| 字段 | 类型 | 说明 |
|------|------|------|
| PackageID | string | 包 ID |
| PackageVersion | string | 包版本 |
| Target | string | 编译目标 |
| ModelProfile | string | 模型 profile |
| RuntimeVars | map[string]string | 运行时变量（展开后）|
| Prompt | string | 编译后的完整 prompt |
| Evaluation | []string | 评估检查清单 |
| ProvenanceTemplate | map[string]any | provenance 源信息 |

**现有包**: `internal/skillplus/compiled_skill.go` — 已有 CompiledSkill struct，**需添加 `Compiler` struct 和 `Compile(ctx, req)` 方法**

---

### generation.Request / generation.Trace

| 字段 | 类型 | 说明 |
|------|------|------|
| ArtifactType | string | 期望输出的 artifact 类型 |
| Prompt | string | 完整 prompt |
| Model | string | 模型名称 |
| Params | map[string]string | 额外参数（temperature, steps 等）|

**Trace**（记录生产过程）:

| 字段 | 类型 | 说明 |
|------|------|------|
| ProviderType | string | "shell" / "http" |
| Model | string | 实际使用的模型名 |
| Command | string | shell 命令（仅 ShellProvider）|
| DurationSec | float64 | 耗时秒数 |

**现有包**: `internal/generation/request.go` — 已有 Request struct，**需添加 `Provider` 接口和 `ShellProvider` 实现**

---

### Artifact

一个生产输出单元，写入磁盘。

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | SHA256 stable ID（project:workflow:stage:package:intent_hash 的哈希）|
| Type | artifacts.Type | image_prompt / image / copy_draft 等 |
| Content | string | 文本内容 |
| Labels | map[string]string | 元数据标签（platform, audience 等）|
| Provenance | string | 对应的 Provenance record ID |

**磁盘格式**:
- `{artifact-dir}/{artifact-id}.{type}.txt` — 内容文件
- `{artifact-dir}/{artifact-id}.provenance.json` — 单个 artifact 的 provenance 快照

**现有包**: `internal/artifacts/artifact.go` — 已有 struct，**需添加 `Write(dir string) error` 方法**

---

### Provenance Record

追加到 `provenance.jsonl`，每行一条记录。

| 字段 | 类型 | 说明 |
|------|------|------|
| ArtifactID | string | 对应的 artifact ID |
| ProjectID | string | 项目 ID |
| WorkflowID | string | workflow ID |
| Stage | string | stage 名称 |
| SkillPackage | string | Skill-Plus 包路径 |
| CompiledTarget | string | 编译目标 |
| Model | string | 使用的模型 |
| PromptHash | string | prompt 的 SHA256（前 16 字节 hex）|
| GenerationParams | map[string]string | 生成参数 |
| EvaluationResult | string | 评估结果摘要 |
| Timestamp | string | RFC3339 时间戳 |
| Trace | *generation.Trace | 生产 trace（可选）|

**现有包**: `internal/provenance/provenance.go` — 已有 struct，**需添加 `AppendRecord(path string, rec *Record) error` 函数和 Timestamp/ArtifactID 字段**

---

### workflow.Engine

新增的运行时编排器（核心新增包）。

| 方法 | 签名 | 说明 |
|------|------|------|
| Run | `(ctx context.Context, req *RunRequest) (*RunResult, error)` | 执行整个 workflow |

**RunRequest**:

| 字段 | 类型 | 说明 |
|------|------|------|
| Project | *project.Project | 已加载的项目 |
| WorkflowDef | *WorkflowDefinition | Workflow 定义 |
| Intent | string | 用户意图文本 |
| ArtifactDir | string | artifact 输出目录 |
| CompilerPath | string | Python compiler PYTHONPATH |
| ProvenancePath | string | provenance.jsonl 路径 |
| Provider | generation.Provider | 已实例化的 provider |

**RunResult**:

| 字段 | 类型 | 说明 |
|------|------|------|
| Artifacts | []*artifacts.Artifact | 产出的 artifact 列表 |
| Provenance | []*provenance.Record | 对应的 provenance 记录 |

---

## State Transitions

### Engine Run State Machine

```
IDLE
  │  RunRequest
  ▼
COMPILING (foreach stage)
  │  CompileResult
  ▼
GENERATING
  │  content + Trace
  ▼
PERSISTING
  │  artifact file + provenance.jsonl append
  ▼
DONE
  
  (any step fails → ERROR,返回带 stage 名称的错误)
```

---

## New Files Required (Phase 2)

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/project/load.go` | **新增** | `Load(path)` 函数 |
| `internal/skillplus/compiler.go` | **新增** | `Compiler` struct + `Compile(ctx, req)` |
| `internal/generation/provider.go` | **新增** | `Provider` 接口 + `ShellProvider` + `Trace` |
| `internal/generation/shell_provider.go` | **新增** | `ShellProvider.Generate(ctx, req)` 实现 |
| `internal/artifacts/writer.go` | **新增** | `Write(dir)` + `stableID()` 函数 |
| `internal/provenance/writer.go` | **新增** | `AppendRecord(path, rec)` 函数 |
| `internal/workflow/engine.go` | **新增** | `Engine.Run(ctx, req)` |
| `internal/workflow/definition.go` | **新增** | `WorkflowDefinition` / `StageDefinition` JSON 结构 |
| `cmd/openmelon/main.go` | **重构** | 退化为 CLI 入口，业务逻辑迁移到 engine |
| `config/openmelon.example.json` | **更新** | 增加 workflow 定义和 project 路径示例 |
