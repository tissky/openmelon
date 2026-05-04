# Quickstart: OpenMelon Phase 2 Runtime

**Date**: 2026-05-04  
**Branch**: `001-phase2-runtime-skeleton`

## Prerequisites

- Go 1.22+
- Python 3.9+（用于 Skill-Plus compiler）
- Skill-Plus 编译器：`pip install skillplus-compile` 或使用本地 reference 实现

## Build

```bash
cd openmelon
go build -o openmelon ./cmd/openmelon
```

## Setup Project File

将现有的 example JSON 拆分为两个文件：

**`examples/food-exploration/project.json`**（新建）:
```json
{
  "id": "demo-food-project",
  "name": "Street Food Realism Demo",
  "platform": "xiaohongshu",
  "audience": "urban food explorers",
  "persona": "ordinary after-work food explorer",
  "workflows": {
    "food_exploration": {
      "id": "food_exploration",
      "stages": [
        {
          "stage": "visual_concretization",
          "skillplus_package": "../../../skillplus/examples/food-street-realism.skillplus",
          "compile_target": "openmelon",
          "model_profile": "gpt-image-family",
          "locale": "zh-CN",
          "vars": { "realism_level": "high", "visual_evidence_density": "dense" }
        }
      ]
    }
  }
}
```

## Run Workflow (Prompt Generation Only)

```bash
./openmelon run \
  --project examples/food-exploration/project.json \
  --workflow food_exploration \
  --intent "下班后在老小区楼下吃一碗牛肉面，想发一条真实的探店帖" \
  --artifact-dir examples/food-exploration/artifacts \
  --compiler ../skillplus/compiler/reference
```

Expected output:
```
[openmelon] project: Street Food Realism Demo (demo-food-project)
[openmelon] workflow: food_exploration — 1 stage(s)
[openmelon] stage: visual_concretization — compiling skill...
[openmelon] compiled: food-street-realism@1.0.0 (target: openmelon)
[openmelon] artifact: image_prompt written → examples/food-exploration/artifacts/a3f2b1c4.image_prompt.txt
[openmelon] provenance appended → examples/food-exploration/artifacts/provenance.jsonl
[openmelon] done.
```

## Run Workflow (With Generation)

```bash
./openmelon run \
  --project examples/food-exploration/project.json \
  --workflow food_exploration \
  --intent "下班后在老小区楼下吃一碗牛肉面，想发一条真实的探店帖" \
  --artifact-dir examples/food-exploration/artifacts \
  --compiler ../skillplus/compiler/reference \
  --generate
```

（需要 `config/openmelon.example.json` 中配置的模型命令可以执行）

## Verify Output

```bash
# 查看生成的 image_prompt artifact
cat examples/food-exploration/artifacts/*.image_prompt.txt

# 查看 provenance 记录
cat examples/food-exploration/artifacts/provenance.jsonl | python3 -m json.tool
```

## Run Tests

```bash
# 全部单元测试
go test ./...

# 集成测试（需要 Python compiler 可用）
go test ./... -tags integration
```
