# CLI Contract: openmelon run

**Type**: Command-line interface  
**Date**: 2026-05-04  
**Branch**: `001-phase2-runtime-skeleton`

## Synopsis

```
openmelon run [flags]
```

## Flags

| Flag | Type | Default | Required | Description |
|------|------|---------|----------|-------------|
| `--project` | string | — | ✅ | Path to project JSON file |
| `--workflow` | string | — | ✅ | Workflow ID (must match a workflow in the project file) |
| `--intent` | string | — | ✅ | User intent text (passed to Skill-Plus as context) |
| `--artifact-dir` | string | `./artifacts` | | Output directory for artifacts and provenance |
| `--compiler` | string | `../skillplus/compiler/reference` | | Python compiler PYTHONPATH for Skill-Plus |
| `--generate` | bool | `false` | | Execute the generation provider (if false, only compile + build prompt) |
| `--timeout` | int | `120` | | Timeout in seconds for Python compiler subprocess calls |

## Project JSON Schema

```json
{
  "id": "demo-food-project",
  "name": "Street Food Realism Demo",
  "platform": "xiaohongshu",
  "audience": "urban food explorers",
  "persona": "ordinary after-work food explorer",
  "memory": {},
  "constraints": [],
  "workflows": {
    "food_exploration": {
      "id": "food_exploration",
      "name": "Food Exploration Workflow",
      "vertical": "food",
      "stages": [
        {
          "stage": "visual_concretization",
          "skillplus_package": "../../../skillplus/examples/food-street-realism.skillplus",
          "compile_target": "openmelon",
          "model_profile": "gpt-image-family",
          "locale": "zh-CN",
          "vars": {
            "realism_level": "high",
            "visual_evidence_density": "dense"
          }
        }
      ]
    }
  }
}
```

**Required fields**: `id`, `name`, `platform`  
**Workflow lookup**: `project.workflows[--workflow]` must exist, else error with list of available workflow IDs.

## Model Config JSON Schema

```json
{
  "models": {
    "gpt-image-family": {
      "provider": "shell",
      "model": "dall-e-3",
      "role": "image_generator",
      "command": "dalle3-cli generate --model dall-e-3"
    }
  }
}
```

**Config path**: `--config` flag (default: `config/openmelon.example.json`)  
**Lookup**: `config.models[stage.model_profile]` → must exist, else error.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success — all artifacts written |
| 1 | User error — bad flags, missing required file, invalid project JSON |
| 2 | Runtime error — Skill-Plus compiler failed, provider failed |

## Standard Output

On success:
```
[openmelon] project: Street Food Realism Demo (demo-food-project)
[openmelon] workflow: food_exploration — 1 stage(s)
[openmelon] stage: visual_concretization — compiling skill...
[openmelon] compiled: food-street-realism@1.0.0 (target: openmelon)
[openmelon] artifact: image_prompt written → artifacts/a3f2b1c4.image_prompt.txt
[openmelon] provenance appended → artifacts/provenance.jsonl
[openmelon] done.
```

On error:
```
[openmelon] error: stage "visual_concretization" — python3 not found in PATH
hint: install Python 3.9+ and ensure it is in PATH
```

## Standard Error

- Validation errors → stderr, exit 1
- Runtime failures → stderr with stage context, exit 2

## Artifact Output Structure

```
{artifact-dir}/
├── {artifact-id}.image_prompt.txt      # Artifact content
├── {artifact-id}.provenance.json       # Single artifact provenance snapshot
└── provenance.jsonl                    # Append-only provenance log (all runs)
```

**Artifact ID format**: `sha256(project_id:workflow_id:stage:package_id:intent_sha256)[:16]`

**provenance.jsonl line format**:
```json
{"artifact_id":"a3f2b1c4...","project_id":"demo-food-project","workflow_id":"food_exploration","stage":"visual_concretization","skill_package":"food-street-realism@1.0.0","compiled_target":"openmelon","model":"dall-e-3","prompt_hash":"9f3a...","timestamp":"2026-05-04T12:00:00Z","trace":{"provider_type":"shell","command":"dalle3-cli ...","duration_sec":3.2}}
```
