# OpenMelon ↔ host-agent integrations

How to make an LLM-powered coding agent (Claude Code, Cursor, Codex, etc.) call OpenMelon as a sub-agent for content creation.

The pattern is the same across hosts: drop a small markdown file telling the host LLM **when** to delegate to OpenMelon and **how** to invoke the CLI. The host runs `bash openmelon -p "..."`. No daemon, no MCP, no protocol — just a CLI behind a documented intent boundary.

## Pick your host

| Host | Setup file | Walkthrough |
|---|---|---|
| **Claude Code** | `~/.claude/skills/<name>.md` | [`claude-code/`](claude-code/) |
| **Cursor** | `.cursor/rules/<name>.md` (project-scoped) | [`cursor/`](cursor/) |
| **Codex / generic** | system prompt / AGENTS.md / equivalent | adapt the Skill body — see notes below |

## Why Skills, not MCP

We considered MCP. For a fire-and-forget content-generation tool (one intent in, one image + post out), MCP's persistent-server model adds complexity that doesn't pay off — you get the same outcome by letting the host run a single `bash openmelon -p "..."`. The Skill file is the LLM-facing API contract; the binary is the implementation.

MCP would earn its complexity if OpenMelon grew long-lived stateful surfaces (live progress streaming back into the conversation, subscription to V-Box events, multi-tool catalogues with shared state). Until that happens, Skills win on simplicity.

## Adapting to a new host

Most agent CLIs accept some form of markdown-style instruction file. To wire up a new host:

1. Find the host's instruction-file location (e.g. `~/.cursorrules`, `~/.aider.conf.yml`, `AGENTS.md`).
2. Copy the body of [`claude-code/skills/openmelon-create-content.md`](claude-code/skills/openmelon-create-content.md), dropping the YAML frontmatter if the host doesn't use it.
3. Verify the host has bash exec permission (some agents disable it by default).
4. Test by asking the host: *"Use openmelon to create a Singapore food street post."*

If you build out a new host integration, please open a PR adding an entry above.
