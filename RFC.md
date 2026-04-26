# RFC Process

RFCs are required for changes that affect Skill-Plus Engine compatibility, runtime guarantees, sandbox guarantees, or community governance.

## Lifecycle

```text
draft -> review -> accepted -> implemented -> obsolete
```

## Changes That Require an RFC

- Changes to dispatcher matching semantics.
- Changes to registry manifest loading or validation semantics.
- Changes to runtime contract expectations.
- Changes to sandbox security model assumptions.
- Changes to B-face output contract expectations.
- Changes to timeout, memory, process, or output-cap guarantees.
- Adding a new runtime type.
- Breaking changes to integration contracts.
- Changes to license, CLA, or governance policy.

## Changes That Usually Do Not Require an RFC

- Typo fixes.
- Documentation clarifications that do not change semantics.
- Bug fixes that preserve compatibility.
- Internal refactors that do not change runtime, sandbox, dispatcher, registry, or B-face behavior.

## RFC Template

```markdown
# RFC XXXX: Title

## Summary

## Motivation

## Detailed Design

## Compatibility

## Security / Privacy Impact

## Alternatives Considered

## Migration Plan

## Rollback Plan
```

## Approval

An RFC is accepted when the relevant maintainers and approvers agree that the design is clear, compatible with project goals, and has an acceptable migration and rollback plan.
