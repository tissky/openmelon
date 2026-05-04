# Feature Specification: OpenMelon Phase 2 — Runtime Skeleton

**Feature Branch**: `001-phase2-runtime-skeleton`  
**Created**: 2026-05-04  
**Status**: Draft  
**Input**: 帮我记录所有里程碑。但是本需求先规划第一个里程碑

## Background

OpenMelon 的 Phase 0–1 已完成：定义了 Project / Workflow / Artifact / Provenance / Labeling / Memory / Roles 等核心数据结构，以及 Skill-Plus 集成接口。还有一个可运行的单 stage demo（`cmd/openmelon/main.go`），能完成 visual_concretization → image_prompt 的单次调用。

Phase 2 的目标是在这个基础上，建立一个真正可复用的运行时引擎，让任意 workflow 定义都能被加载和执行，并将 artifact 和 provenance 落盘持久化。

此 spec 覆盖 Phase 2 全部范围，并在 Milestone Map 里记录 Phase 3–5 供后续规划参考。

## Milestone Map（全项目）

| Milestone | 名称 | 状态 |
|-----------|------|------|
| Phase 0 | Foundation — 定位与核心概念 | 完成 |
| Phase 1 | Core Contracts — 数据结构与接口定义 | 完成 |
| Phase 2 | Runtime Skeleton — 可执行运行时 | 本 spec |
| Phase 3 | Content Workflows — 多工作流支持 | 待规划 |
| Phase 4 | Multimodal Production — 音频/视频/跨模态 | 待规划 |
| Phase 5 | Ecosystem — 社区包、子 Agent 模式 | 待规划 |

---

## User Scenarios & Testing

### User Story 1 — 完整工作流端到端执行 (Priority: P1)

作为内容创作者，我希望用一条命令就能从创意意图跑出一个 image_prompt artifact，并自动记录 provenance，而不需要手动拼接脚本或理解内部调用链。

**Why this priority**: 这是 OpenMelon 作为"运行时"而非"脚本集合"的最核心价值体现。

**Independent Test**: 执行 `./openmelon run --project examples/food-exploration/project.json --workflow food_exploration --intent "..."` 能产出 artifact 文件和 provenance 记录，不依赖任何其他 story。

**Acceptance Scenarios**:

1. **Given** 有效的 project.json 和 workflow 配置，**When** 执行 run 命令，**Then** 在 artifacts/ 目录下产出 `image_prompt_<id>.txt` 和对应的 `<id>.provenance.json`
2. **Given** 意图文本，**When** workflow 包含 visual_concretization stage，**Then** 该 stage 自动调用 Skill-Plus compiler 并将编译后的 prompt 注入生成请求
3. **Given** 模型配置指向一个 shell command provider，**When** 执行 generate，**Then** 命令被执行，输出作为 artifact content 保存
4. **Given** workflow 执行完成，**When** 查看 provenance.jsonl，**Then** 每个 artifact 都有完整的 provenance 记录（project_id, workflow, stage, skill_package, model, prompt_hash）

---

### User Story 2 — Project 从文件加载 (Priority: P1)

作为开发者，我希望 Project 信息从独立 JSON 文件加载，而不是内联在 example 里，这样同一个项目可以被多次不同 workflow 复用。

**Why this priority**: 现在 project 信息内联在 example JSON 里，无法复用。这是解耦的第一步。

**Independent Test**: 创建 project.json，执行 run 命令，系统能正确读取 project 并回显项目名称，即使没有 workflow 执行。

**Acceptance Scenarios**:

1. **Given** 一个合法的 project.json，**When** 系统加载它，**Then** Project 对象包含正确的 id、name、platform、audience、persona 字段
2. **Given** 缺少必填字段 id 的 project.json，**When** 加载，**Then** 返回明确错误，指出缺失字段
3. **Given** 文件路径不存在，**When** 加载，**Then** 返回文件不存在错误

---

### User Story 3 — Skill-Plus Adapter 正式化 (Priority: P1)

作为 workflow 引擎，我需要一个稳定的 Skill-Plus 编译接口，能通过包路径和编译参数调用 Python compiler，并将结果反序列化为 CompiledSkill 结构体。

**Why this priority**: 这是 workflow engine 的核心依赖。main.go 里已有实现，但未封装，无法测试，也无法在 engine 里复用。

**Independent Test**: 调用 `internal/skillplus.Compile()` 能返回正确的 CompiledSkill，可单独写单元测试 mock shell exec 验证。

**Acceptance Scenarios**:

1. **Given** 有效的 skillplus 包路径和编译配置，**When** 调用 Compile()，**Then** 返回包含 compiled_prompt、runtime_vars、evaluation 的 CompiledSkill
2. **Given** Python compiler 不可用，**When** 调用 Compile()，**Then** 返回明确错误，不 panic
3. **Given** 不支持的 target 类型，**When** 调用 Compile()，**Then** 错误消息包含可用的 target 列表

---

### User Story 4 — Model Provider 接口与 CommandProvider (Priority: P2)

作为 workflow 引擎，我需要一个统一的模型调用接口，当前先支持通过 shell command 调用外部工具，后续可扩展为 HTTP API。

**Why this priority**: 没有这个接口，workflow engine 无法结构化地调用任何模型，也无法添加第二个 provider。

**Independent Test**: 使用 CommandProvider 配置一个 echo 命令，调用 Generate(prompt) 能返回 prompt 文本作为输出，可单独单元测试。

**Acceptance Scenarios**:

1. **Given** config 中定义了 command provider，**When** workflow 调用 Generate()，**Then** shell command 被执行，stdout 作为 artifact 内容
2. **Given** shell command 返回非零退出码，**When** Generate()，**Then** 返回包含 stderr 内容的错误
3. **Given** 超时配置，**When** command 运行超过超时时间，**Then** 进程被终止并返回超时错误

---

### User Story 5 — Artifact 与 Provenance 落盘 (Priority: P2)

作为创作者，每次执行后我希望 artifact 内容和完整的生产 provenance 都能保存到磁盘，方便事后查看、复现和比较。

**Why this priority**: 没有落盘，OpenMelon 只是一次性脚本。落盘是使其成为生产运行时的关键。

**Independent Test**: 执行任意 workflow，检查 artifacts/ 目录和 provenance.jsonl 文件被正确创建，内容可读。

**Acceptance Scenarios**:

1. **Given** workflow 执行完成，**When** 检查 artifact 目录，**Then** 每个 artifact 有独立文件，文件名包含 artifact id 和类型
2. **Given** 多次执行同一 workflow，**When** 查看 provenance.jsonl，**Then** 每次执行追加一行，历史记录不被覆盖
3. **Given** artifact 目录不存在，**When** 首次执行，**Then** 目录自动创建

---

### Edge Cases

- 如果 Skill-Plus Python compiler 未安装，应给出可读的错误提示，指引用户如何安装
- 如果 workflow 定义的 stage 引用了不存在的 Skill-Plus 包，应在执行前检测并报错
- 如果同一个 artifact id 已存在于磁盘，应追加版本后缀而非静默覆盖
- 如果生成命令输出为空，artifact 记录为空内容，不报错，但在 provenance 中标注

## Requirements

### Functional Requirements

- **FR-001**: 系统必须能从独立的 JSON 文件加载 Project，并校验必填字段（id, name, platform）
- **FR-002**: 系统必须提供可测试的 `skillplus.Compile()` 函数，通过 subprocess 调用 Python 参考编译器
- **FR-003**: 系统必须定义 `generation.Provider` 接口，并实现 `CommandProvider`（shell exec）
- **FR-004**: 系统必须实现 `workflow.Engine`，能按 stage 顺序执行：compile skill → generate → create artifact
- **FR-005**: 系统必须将每个 artifact 写入磁盘（文本内容 + provenance JSON），并将 provenance 追加到 provenance.jsonl
- **FR-006**: CLI 必须支持 `run` 子命令，接受 `--project`, `--workflow`, `--intent`, `--artifact-dir` 参数
- **FR-007**: 系统必须有集成测试，使用 food-exploration 示例端到端验证整条链路
- **FR-008**: 现有 main.go 中的单 stage 逻辑必须被迁移到 workflow engine，main.go 退化为 CLI 入口

### Key Entities

- **Project**: 创作项目，包含 platform/audience/persona/constraints，从文件加载
- **Workflow**: 有序的 stage 列表，每个 stage 绑定一个 Skill-Plus 包和 role
- **CompiledSkill**: Skill-Plus compiler 输出，包含 compiled_prompt / runtime_vars / evaluation
- **Provider**: 模型调用接口，当前实现：CommandProvider（shell exec）
- **Artifact**: 生产输出，有 id / type / content / labels，写入磁盘
- **Provenance**: 生产溯源记录，追加到 provenance.jsonl

## Success Criteria

### Measurable Outcomes

- **SC-001**: 执行 `./openmelon run` 命令能在 10 秒内（不含模型调用耗时）完成并输出 artifact 文件
- **SC-002**: `internal/skillplus`、`internal/workflow`、`internal/generation` 三个包各有至少 1 个单元测试，`go test ./...` 全部通过
- **SC-003**: 集成测试覆盖 project load → skill compile → artifact write → provenance record 完整链路
- **SC-004**: `cmd/openmelon/main.go` 不再包含业务逻辑，仅为 CLI 参数解析与 engine 调用
- **SC-005**: 任何一个 stage 失败，运行时返回明确错误信息，包含 stage 名称和失败原因

## Assumptions

- Skill-Plus Python compiler 已安装在运行环境，路径通过 `--compiler` 参数传入（不做自动安装）
- 模型调用在 Phase 2 只需支持 CommandProvider（shell exec），不需要 HTTP API adapter（Phase 3 再加）
- Workflow 定义在 Phase 2 通过 JSON 文件描述，不需要 UI 或 DSL
- 持久化在 Phase 2 只需要本地文件系统，不需要数据库
- Phase 3（多 workflow）、Phase 4（音视频）、Phase 5（生态）超出本 spec 范围，不在此实现

## Out of Scope

- HTTP API 或 Web UI
- 数据库持久化
- 多轮 feedback 循环
- 音频/视频 workflow
- Sub-agent 模式
- Skill-Plus 包的自动发现或注册
