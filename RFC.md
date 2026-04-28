# RFC Process

RFCs are required for changes that affect OpenMelon workflow semantics, artifact contracts, provenance, labeling, Skill-Plus integration, runtime guarantees, or community governance.

## Lifecycle

```text
draft -> review -> accepted -> implemented -> obsolete
```

## Changes That Require an RFC

- Changes to workflow stage semantics.
- Changes to project memory contracts.
- Changes to artifact or provenance schemas.
- Changes to content labeling semantics.
- Changes to Skill-Plus compilation or runtime integration.
- Adding a new runtime execution mode.
- Breaking changes to public contracts.
- Changes to license, CLA, or governance policy.

## Changes That Usually Do Not Require an RFC

- Typo fixes.
- Documentation clarifications that do not change semantics.
- Bug fixes that preserve compatibility.
- Internal refactors that do not change public workflow, artifact, provenance, labeling, or integration behavior.

## RFC Template

```markdown
# RFC XXXX: Title

## Summary

## Motivation

## Detailed Design

## Workflow / Artifact Impact

## Security / Privacy Impact

## Alternatives Considered

## Migration Plan

## Rollback Plan
```

## Approval

An RFC is accepted when the relevant maintainers and approvers agree that the design is clear, compatible with OpenMelon's content-production goals, and has an acceptable migration and rollback plan.
