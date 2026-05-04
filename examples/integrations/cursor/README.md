# Cursor ↔ OpenMelon

Cursor doesn't have a Skills system the way Claude Code does, but it has **Project Rules** (`.cursor/rules/<name>.md`) that the in-editor LLM reads on every interaction in that project.

## Install

```bash
mkdir -p .cursor/rules
cp openmelon-create-content.md .cursor/rules/
```

The rules file applies project-scoped — Cursor will surface OpenMelon when the user is working in that repo. For workspace-wide availability, drop the same file into your home `.cursor/rules/` (Cursor merges global and project rules).

## Prerequisites

Same as Claude Code:

```bash
go build -o ~/bin/openmelon <openmelon-repo>/cmd/openmelon
pip install skillplus
npm i -g @e8s/vbox-cli              # only if you want to publish

export ANTHROPIC_API_KEY=...        # OR OPENAI_API_KEY
export OPENAI_API_KEY=...           # always, for image gen
export VBOX_API_KEY=...             # only when publishing
```

## Use it

In Cursor's chat (Cmd-L / Ctrl-L), ask:

> Generate a realistic V-Box post about eating beef noodles in an old neighborhood

Cursor reads the rules file, runs `openmelon -p "..." --skill skillplus:food-street-realism` in the integrated terminal, and reports back.

## Note on terminal access

Cursor needs terminal exec permission to run shell commands. If you've turned that off in Settings → Features → Composer → Allow terminal commands, the rules file becomes documentation rather than action. Re-enable it for OpenMelon to work.

## Rules-file format

Cursor accepts the same body as the Claude Code Skill — just drop the YAML frontmatter. See [`openmelon-create-content.md`](openmelon-create-content.md) in this directory.
