# OpenMelon Architecture

OpenMelon is a content-production runtime.

Its core abstraction is not a chat message or a generic task. Its core abstraction is a production chain:

```text
Project -> Workflow -> Stage -> Compiled Skill -> Artifact -> Review -> Provenance -> Memory Update
```

## Project

A project stores the durable creative context:

- platform,
- target audience,
- content vertical,
- creator or character persona,
- style constraints,
- project memory,
- asset history,
- prior feedback.

## Workflow

A workflow is a staged production process for a content vertical.

Example stages:

1. intent planning,
2. angle selection,
3. copywriting,
4. visual concretization,
5. prompt generation,
6. image generation,
7. audio generation,
8. video planning,
9. review,
10. packaging.

## Skill-Plus Integration

OpenMelon does not treat skills as static text. It uses Skill-Plus packages.

At a workflow stage, OpenMelon can:

1. choose a Skill-Plus package,
2. compile it for a target and model profile,
3. inject project memory and runtime variables,
4. run the compiled skill,
5. validate output,
6. attach provenance to artifacts.

## Artifact

An artifact is a typed production output:

- copy draft,
- image prompt,
- generated image,
- audio script,
- voice asset,
- video shot list,
- review report.

Artifacts are not anonymous files. They carry labels and provenance.

## Provenance

OpenMelon records how an artifact was produced:

- workflow and stage,
- model used,
- Skill-Plus package,
- compiled skill target,
- runtime variables,
- raw prompt,
- generation parameters,
- evaluation result,
- user feedback.

This makes the production process itself a labeling process.

## Memory Update

After review, OpenMelon can write useful feedback back into project memory:

- which skill variants worked,
- which model profiles failed,
- which visual details improved realism,
- which platform-specific rules matter,
- which creator preferences should persist.

## Operating Modes

OpenMelon supports two intended modes:

- **Standalone mode** — a complete content-production agent/runtime.
- **Sub-agent mode** — a specialized content agent invoked by a broader orchestrator.

In both modes, OpenMelon owns content workflow structure and artifact provenance.
