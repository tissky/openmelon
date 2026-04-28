# Security Policy

## Reporting a Vulnerability

Do not report security vulnerabilities through public GitHub issues, pull requests, or discussions.

Use the repository's private vulnerability reporting channel if available. If private vulnerability reporting is unavailable, contact a project maintainer privately and avoid public disclosure until the maintainers have had a reasonable opportunity to investigate and coordinate a fix.

## Response Expectations

Maintainers will acknowledge valid security reports as soon as practical, investigate the affected versions and components, coordinate remediation, and publish an advisory or release note when disclosure is appropriate.

## Security-Sensitive Areas

Changes touching the following areas require security review before merge:

- Runtime execution behavior.
- Network access and egress allowlists.
- File system persistence or temporary file cleanup.
- Secret handling, credentials, tokens, or authentication material.
- Dependencies with native code or unknown licensing.
- External API access from generation adapters or integrations.
- Project memory retention or user data persistence.
- Artifact provenance and label handling for private content.
- Any use of shell execution or process spawning.

## Supported Scope

This policy applies to the OpenMelon content-production runtime, workflow orchestration, project memory, artifact management, provenance, labeling, Skill-Plus integration, examples, and repository automation.

## Contributor Responsibility

Contributors are responsible for ensuring that their submissions do not include secrets, proprietary code, private data, or third-party confidential material. AI-assisted contributions must be reviewed by the contributor for correctness, licensing, and security before submission.
