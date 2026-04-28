# OpenMelon

**A content-creation agent runtime for reproducible multimodal production.**

OpenMelon is an AI-era DaVinci: a professional runtime for planning, generating, reviewing, labeling, and iterating content across text, image, audio, and video.

It is not a general-purpose agent framework. It is built for content creation.

## Why OpenMelon

Modern generative models are powerful, but creative production still feels unstable:

- prompts drift between runs,
- visual style is hard to preserve,
- creator and project memory are not consistently applied,
- generated assets lack production provenance,
- feedback is hard to convert into better future outputs,
- content workflows are often one-shot prompts instead of durable pipelines.

OpenMelon treats content generation as a production process:

```text
Project -> Workflow -> Stage -> Skill-Plus Package -> Compiled Skill -> Generation -> Artifact -> Review -> Labels/Provenance -> Memory Update
```

## Core Responsibilities

OpenMelon manages:

- projects and creative briefs,
- long-term project memory,
- creator, character, and persona consistency,
- content vertical workflows,
- copy, image, audio, and video production stages,
- prompt, shot, and script enhancement,
- artifact labeling and provenance,
- evaluation and feedback loops,
- Skill-Plus package compilation and execution.

## Relationship to Skill-Plus

Skill-Plus was extracted from OpenMelon's production needs.

OpenMelon needs stable reusable content capabilities. Skill-Plus provides those capabilities as compilable packages.

```text
OpenMelon  = content-production runtime
Skill-Plus = compilable skill package standard and open ecosystem
```

OpenMelon compiles Skill-Plus packages into workflow-ready skills, executes them in the right stage, and attaches provenance to every generated artifact.

## Model Configuration

OpenMelon is model-driven. Different workflow stages can route to different models or providers:

```text
intent_planning          -> planner model
visual_concretization    -> Skill-Plus compiler / prompt director
image_generation         -> image generation model
audio_generation         -> audio model
visual_quality_review    -> reviewer model
```

See `config/openmelon.example.json` for the current configuration shape. The first runnable example uses a command provider for image generation so OpenMelon can produce a real image artifact while recording the model, command, prompt hash, labels, and provenance.

## Repository Layout

```text
├── docs/                 # Design docs for content-agent workflows
├── cmd/openmelon/        # CLI entrypoint
├── internal/             # Runtime implementation packages
├── pkg/                  # Public contracts
└── examples/             # End-to-end content workflow examples
```

## Status

OpenMelon is being rebuilt around the content-production runtime model. Current code is a skeleton for the new architecture.
