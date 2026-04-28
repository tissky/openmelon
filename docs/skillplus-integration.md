# Skill-Plus Integration

OpenMelon uses Skill-Plus packages as professional content capabilities.

At a workflow stage, OpenMelon should:

1. select a package by stage and role,
2. compile it with project memory and runtime variables,
3. route the compiled output to the configured model or provider,
4. validate output,
5. attach provenance to generated artifacts.

Skill-Plus gives OpenMelon a stable way to reduce prompt drift and make creative generation reproducible.

## Multi-Model Routing

OpenMelon is not a single-model wrapper. A workflow can configure different models for different roles:

| Stage | Role | Example model slot |
|---|---|---|
| `intent_planning` | planner | `planner` |
| `visual_concretization` | prompt director | `prompt_director` |
| `image_generation` | generator | `image_generator` |
| `visual_quality_review` | reviewer | `reviewer` |

This lets parsing, judgment, review, and generation use different models while still sharing one artifact/provenance trace.
