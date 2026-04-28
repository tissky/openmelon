# Contributing to OpenMelon

Thank you for your interest in contributing to OpenMelon.

OpenMelon is a content-creation-focused Agent/runtime. Contributions should preserve that focus: durable creative workflows, project memory, multimodal orchestration, Skill-Plus integration, and high-quality content production.

## Contributor License Agreement (CLA)

All contributors must sign the Point Eight Individual CLA, or Corporate CLA for company contributions, before a PR can be merged. The CLA Assistant bot will prompt you automatically on your first PR.

## How to Contribute

1. Fork this repository.
2. Create a topic branch from `main`.
3. Keep changes scoped to one clear purpose.
4. Add or update tests for behavior changes.
5. Run the local checks before opening a PR:

```bash
make test
```

If your change touches lint-sensitive code, also run:

```bash
make lint
```

## Areas of Contribution

Useful contribution areas include:

- content workflow orchestration,
- Skill-Plus dispatch and runtime integration,
- project memory and creative state management,
- multimodal tool coordination,
- image/copy/audio/video workflow stages,
- sandboxing and runtime safety,
- reference integrations and examples.

## Security-Sensitive Changes

Changes touching the following areas require additional review:

- sandbox behavior,
- runtime execution,
- network access or egress allowlists,
- filesystem persistence,
- secret handling,
- external API integrations,
- output sanitization,
- timeout, memory, process, or output caps,
- shell execution or process spawning.

See [SECURITY.md](./SECURITY.md) for details.

## Code Style

- Keep Go code formatted with `gofmt`.
- Prefer small, testable packages.
- Keep public interfaces explicit and documented.
- Avoid adding generic Agent features unless they support OpenMelon's content-creation mission.

## Pull Request Review

PRs should include:

- a concise summary,
- rationale for the change,
- tests or validation steps,
- notes about security or workflow implications.

## License

By contributing, you agree that your contributions will be licensed under Apache-2.0.
