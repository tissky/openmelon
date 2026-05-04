# Implementation Plan: OpenMelon Phase 2 — Runtime Skeleton

**Branch**: `001-phase2-runtime-skeleton` | **Date**: 2026-05-04 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `specs/001-phase2-runtime-skeleton/spec.md`

## Summary

将 OpenMelon 从单 stage demo 演进为真正的工作流运行时引擎。核心交付：Project 文件加载、Skill-Plus Compiler 封装、generation.Provider 接口 + ShellProvider 实现、workflow.Engine 编排器、artifact 和 provenance 落盘持久化、`openmelon run` CLI 子命令。所有现有业务逻辑从 main.go 迁移到对应 internal 包。

## Technical Context

**Language/Version**: Go 1.22  
**Primary Dependencies**: 标准库（encoding/json, os/exec, crypto/sha256, flag, context）— 零外部依赖  
**Storage**: 本地文件系统（.txt artifacts + .jsonl provenance 追加日志）  
**Testing**: `go test ./...`，集成测试用 `-tags integration`  
**Target Platform**: macOS / Linux CLI  
**Project Type**: CLI tool + internal library packages  
**Performance Goals**: `openmelon run` 框架开销 < 1s（不含模型调用）  
**Constraints**: 零外部 Go 依赖；Python subprocess 不能 panic；单进程无并发  
**Scale/Scope**: 单用户本地工具，Phase 2 不考虑并发或远程执行

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

> **注意**: 项目尚无 constitution.md（仅有模板）。以下检查基于通用工程最佳实践。

| 检查项 | 状态 | 说明 |
|--------|------|------|
| 零外部依赖（Phase 2）| ✅ 通过 | 只用标准库，不引入新依赖 |
| 现有包不破坏性修改 | ✅ 通过 | 只追加新文件，不修改现有 struct |
| 单一入口点 | ✅ 通过 | cmd/openmelon/main.go 退化为 CLI 薄层 |
| 测试覆盖各 internal 包 | ✅ 通过 | 每个新包至少 1 个单元测试 |
| 无全局状态 | ✅ 通过 | Engine/Compiler/Provider 通过构造函数注入 |

**Post-Design Re-check**: data-model.md 设计通过所有检查项。无 constitution 违规。

## Project Structure

### Documentation (this feature)

```text
specs/001-phase2-runtime-skeleton/
├── plan.md              # 本文件
├── research.md          # Phase 0 研究结论
├── data-model.md        # Phase 1 实体模型
├── quickstart.md        # Phase 1 快速上手指南
├── contracts/
│   └── cli.md           # CLI 接口契约
├── checklists/
│   └── requirements.md  # Spec 质量检查清单
└── tasks.md             # Phase 2 任务列表（/speckit.tasks 生成）
```

### Source Code Layout

```text
cmd/
└── openmelon/
    └── main.go               # CLI 入口 — 仅 flag 解析 + engine.Run 调用（重构）

internal/
├── project/
│   ├── project.go            # Project struct（已有）
│   └── load.go               # Load(path) 函数（新增）
│
├── skillplus/
│   ├── compiled_skill.go     # CompiledSkill struct（已有）
│   └── compiler.go           # Compiler struct + Compile(ctx, req)（新增）
│
├── generation/
│   ├── request.go            # Request struct（已有）
│   ├── provider.go           # Provider interface + Trace + ProviderError（新增）
│   └── shell_provider.go     # ShellProvider.Generate(ctx, req)（新增）
│
├── artifacts/
│   ├── artifact.go           # Artifact struct + Type consts（已有）
│   └── writer.go             # Write(dir) + stableID()（新增）
│
├── provenance/
│   ├── provenance.go         # Record struct（已有，需扩展字段）
│   └── writer.go             # AppendRecord(path, rec)（新增）
│
└── workflow/
    ├── workflow.go           # Stage consts + Workflow struct（已有）
    ├── definition.go         # WorkflowDefinition / StageDefinition JSON struct（新增）
    └── engine.go             # Engine.Run(ctx, req)（新增）

config/
└── openmelon.example.json    # 更新：加入 workflow 定义示例

examples/food-exploration/
├── project.json              # 新增：从 beef-noodles.json 拆分出的项目文件
├── beef-noodles.json         # 保留：向后兼容
└── artifacts/                # 已有目录，产出放这里
```

**Structure Decision**: 标准 Go 项目布局（`cmd/` + `internal/`），每个关注点一个包，不引入多模块。

## Complexity Tracking

无 constitution 违规，此表不适用。
