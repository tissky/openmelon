# Claude Code ↔ OpenMelon

Make Claude Code delegate content-creation work to OpenMelon by installing a Skill file.

## Install

```bash
mkdir -p ~/.claude/skills
cp skills/openmelon-create-content.md ~/.claude/skills/
cp skills/openmelon-publish.md       ~/.claude/skills/   # optional
```

That's it. Claude Code reads `~/.claude/skills/*.md` on every conversation and surfaces them when relevant.

## Prerequisites

The Skills assume two things are on PATH or invokable:

| Binary | How to install | Used for |
|---|---|---|
| `openmelon` | `cd <openmelon-repo> && go build -o ~/bin/openmelon ./cmd/openmelon` | the agent loop |
| `skillplus` | `pip install skillplus` (or `pip install -e <skillplus-repo>`) | skill compilation |
| `vbox-cli` *(optional)* | `npm i -g @e8s/vbox-cli` (or `npm link` from the repo) | publishing to V-Box |

And one of these env vars:

| Variable | Purpose | Required? |
|---|---|---|
| `ANTHROPIC_API_KEY` | LLM (preferred for structured output) | one of these two |
| `OPENAI_API_KEY` | LLM **and** image generation | required for image gen always |
| `VBOX_API_KEY` | publishing to V-Box | only when publishing |

If only `OPENAI_API_KEY` is set, OpenMelon uses GPT for the LLM step and `gpt-image-1` for the image step — one key for the whole flow.

## Verify

In Claude Code, ask:

> Use openmelon to create a realistic post about eating beef noodles at an old-neighborhood shop downstairs.

Expected: Claude Code reads the Skill, runs `openmelon -p "..." --skill skillplus:food-street-realism`, streams the model output into the conversation, and reports the artifact path.

To go all the way to V-Box, ask:

> Now publish that image to my V-Box account.

Claude Code reads the publish Skill and runs `openmelon -p "..." --publish vbox` (or uploads + posts via `vbox-cli` directly if the artifact already exists).

## Files in this directory

- [`skills/openmelon-create-content.md`](skills/openmelon-create-content.md) — primary Skill: generates content from a free-text intent.
- [`skills/openmelon-publish.md`](skills/openmelon-publish.md) — secondary Skill: publishes a generated artifact to V-Box.

The two are split deliberately so Claude Code can offer "create only" vs "create + publish" as separate decisions instead of always doing both.
