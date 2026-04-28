# Governance

OpenMelon is governed as a content-creation agent/runtime project.

The repository contains workflow orchestration, project memory, multimodal artifact management, Skill-Plus integration, provenance, labeling, review, and runtime interfaces.

## Roles

| Role | Responsibilities |
| --- | --- |
| Contributor | Opens issues and pull requests for documentation, runtime behavior, workflow design, integrations, examples, or proposals. |
| Triager | Reproduces issues, applies labels, identifies duplicates, and routes work to reviewers. |
| Reviewer | Reviews pull requests for correctness, maintainability, tests, documentation, and security concerns. |
| Approver | Approves changes in owned areas before merge. |
| Maintainer | Merges pull requests, manages releases, maintains automation, and resolves process questions. |
| Workflow Approver | Approves changes that affect content workflow semantics, artifact contracts, provenance, labeling, or Skill-Plus integration. |
| Security Reviewer | Reviews changes involving runtime execution, network access, data retention, secrets, external APIs, filesystem access, and process execution. |
| Release Manager | Coordinates versioning, changelog entries, tags, release notes, compatibility updates, and release readiness. |

## Review Requirements

- All pull requests require passing checks once CI is enabled.
- All pull requests require at least one human reviewer approval before merge.
- Changes to owned paths require the relevant CODEOWNER review.
- Changes to workflow semantics, artifact contracts, provenance, labeling, or Skill-Plus integration require Workflow Approver review.
- Changes involving runtime execution, network access, data retention, secret handling, native-code dependencies, external APIs, or process execution require Security Reviewer review.
- Breaking changes require an accepted RFC before merge.

## Merge Policy

Maintainers should use squash merge by default. Breaking changes must be called out in the pull request and in release notes when applicable.

## Labels

The community process recognizes these labels: `type:docs`, `type:bugfix`, `type:workflow`, `type:memory`, `type:artifact`, `type:provenance`, `type:labeling`, `type:skillplus`, `type:security`, `type:dependency`, `breaking-change`, `needs-rfc`, `needs-security-review`, `needs-workflow-review`, and `do-not-merge`.
