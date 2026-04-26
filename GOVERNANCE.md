# Governance

Skill-Plus Engine is governed as the reference execution engine for the Skill-Plus open standard. The repository contains execution, dispatching, registry loading, sandboxing, runtime behavior, built-in plugins, and host integration examples.

## Roles

| Role | Responsibilities |
| --- | --- |
| Contributor | Opens issues and pull requests for bug fixes, runtime improvements, sandbox changes, documentation, integration examples, or proposals. |
| Triager | Reproduces issues, applies labels, identifies duplicates, and routes work to reviewers. |
| Reviewer | Reviews pull requests for correctness, maintainability, tests, documentation, and security concerns. |
| Approver | Approves changes in owned areas before merge. |
| Maintainer | Merges pull requests, manages releases, maintains automation, and resolves process questions. |
| Spec Approver | Approves changes that affect dispatch semantics, registry behavior, runtime contract compatibility, and B-face output compatibility. |
| Security Reviewer | Reviews changes involving sandboxing, runtime execution, network access, data retention, secrets, external APIs, process execution, and malicious Skill risk. |
| Release Manager | Coordinates versioning, changelog entries, tags, release notes, compatibility matrix updates, and release readiness. |

## Review Requirements

- All pull requests require passing checks once CI is enabled.
- All pull requests require at least one human reviewer approval before merge.
- Changes to owned paths require the relevant CODEOWNER review.
- Changes to dispatcher semantics, registry validation, runtime contracts, or B-face output compatibility require Spec Approver review.
- Changes involving sandboxing, runtime execution, network access, data retention, secret handling, native-code dependencies, external APIs, output sanitization, or process execution require Security Reviewer review.
- Breaking changes require an accepted RFC before merge.

## Merge Policy

Maintainers should use squash merge by default. Breaking changes must be called out in the pull request and in release notes when applicable.

## Labels

The community process recognizes these labels: `type:docs`, `type:bugfix`, `type:engine`, `type:runtime`, `type:sandbox`, `type:security`, `type:dependency`, `breaking-change`, `needs-rfc`, `needs-security-review`, `needs-spec-review`, and `do-not-merge`.
