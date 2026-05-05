# Testing OpenMelon end-to-end

Recipe for verifying the openmelon ↔ skillplus ↔ vbox-cli chain locally.

Three integration paths — pick whichever you have credentials for:

| Path | Audience | Requires |
|---|---|---|
| **A**. Interactive TUI | day-to-day creator workflow | API key for OpenRouter / OpenAI / Anthropic |
| **B**. Headless `-p` + V-Box publish | scripts, sub-agent integration | A's keys + `VBOX_API_KEY` + `vbox-cli` linked locally |
| **C**. Claude Code Skill | host-agent integration | A's keys + Claude Code installed |

---

## Setup (once)

### 1. Install openmelon + skillplus

```bash
npm install -g @e8s/openmelon @e8s/skillplus
```

Or build from source (clone the repo first):

```bash
cd path/to/openmelon
go install -ldflags "-X github.com/eight-acres-lab/openmelon/internal/version.Version=$(git describe --tags --always)" ./cmd/openmelon
cd path/to/skillplus && npm link
```

### 2. (Path B only) Link vbox-cli locally

```bash
cd path/to/vbox-cli
npm link
vbox-cli --version
```

### 3. First-run wizard

Just run `openmelon` in any directory. The wizard walks through trust → API key → LLM model → image model → project init. Credentials land in `~/.openmelon/credentials.json` (mode 0600), defaults in `~/.openmelon/config.json`.

---

## Path A — Interactive TUI

### A.1. The food-noodles demo

```bash
cd ~/work/test-melon
openmelon
```

Inside the TUI:

```
> /skill
[picker → choose food-street-realism]
(skill: food-street-realism) — applies to your next message

> Grab a bowl of beef noodles after work, write an authentic restaurant-visit post
```

Expected sequence (rendered live in viewport):

```
  ⏺ compile_skill({"skill":"food-street-realism","locale":"zh-CN"})
    ⎿ {"package":{"id":"food-street-realism",...},"compiled_prompt":"...",...}
  ⏺ generate_image({"prompt":"A handheld phone snapshot inside...",...})
    ⎿ {"path":".../draft-1.png","sha256":"abc123..."}

⠋ Streaming response · 0:24 · 1.4k in / 312 out · esc to cancel
```

### A.2. Verify the artifact

In the TUI:

```
> /settings
[switch to "Auto-judge" so bash auto-runs read-only commands]

> Open the image and tell me if it looks like a real phone photo
```

The agent calls `bash({"command":"file .openmelon/sessions/*/draft-1.png"})`, then `bash({"command":"open .openmelon/sessions/*/draft-1.png"})` — both auto-approved by the judge as read-only.

Or outside the TUI:

```bash
ls -lh .openmelon/sessions/*/draft-1.png
open .openmelon/sessions/*/draft-1.png
```

### A.3. Variations to try

```
> /model            # switch the LLM mid-session (e.g. claude-opus-4.7)
> /model-image      # switch the image model (e.g. gpt-5-image)
> /clear            # forget conversation history
> /save out.jsonl   # export the conversation
> /history          # print the message log
> /exit
```

### A.4. Resume a session

After exit, the shell prints:

```
session saved at /path/to/.openmelon/sessions/20260506-101203-a1b2c3d4
to resume:    openmelon resume 20260506-101203-a1b2c3d4
```

Then:

```bash
openmelon resume                            # list recent + pick id
openmelon resume 20260506-101203-a1b2c3d4   # load directly
```

The new TUI renders prior turns into the transcript and continues with full context.

---

## Path B — Headless `-p` + V-Box publish

`-p` runs the same tool stack without the TUI. Useful for scripts.

### B.1. One-shot create

```bash
cd ~/work/test-melon
openmelon -p "Grab a bowl of beef noodles after work, write an authentic restaurant-visit post"
```

Expected stderr (the activity log):

```
[openmelon] project=test-melon session=20260506-... llm=openrouter/openai/gpt-5.5 image=openrouter/google/gemini-2.5-flash-image
[openmelon] intent: Grab a bowl of beef noodles...
[turn 1] reply (finish=tool_calls, tool_calls=1)
[turn 1] → compile_skill(...)
[turn 1] ← {"package":{...},...}
...
```

### B.2. Headless + V-Box publish

```bash
openmelon -p "..." --publish vbox
```

Expected additional lines:

```
[openmelon] uploaded → fid=fid_xxx
[openmelon] published. vbox-cli response: {"status":"queued_for_review","content_id":"post_xxx"}
```

`queued_for_review` is **expected** — the post is gated by V-Box's Review Queue. Approve it in the V-Box app to make it live.

### B.3. Headless + bash

The bash tool in headless requires `/settings → trusted` (or `auto`) — no UI to render the approval modal. Set the mode interactively first:

```bash
openmelon                          # in TUI: /settings → Trusted, then /exit
openmelon -p "Inspect generated images and report sizes"   # bash now runs
```

---

## Path C — Claude Code Skill

### C.1. Install the Skills

```bash
mkdir -p ~/.claude/skills
cp path/to/openmelon/examples/integrations/claude-code/skills/openmelon-create-content.md ~/.claude/skills/
cp path/to/openmelon/examples/integrations/claude-code/skills/openmelon-publish.md       ~/.claude/skills/
```

### C.2. Verify openmelon is on PATH for Claude Code

Claude Code inherits PATH from the shell that launched it. After `npm install -g @e8s/openmelon`, `which openmelon` should resolve. If you `go install`-ed, ensure `$GOPATH/bin` is on PATH.

### C.3. Test from inside Claude Code

In a Claude Code conversation:

> Use openmelon to create a realistic post about eating beef noodles at an old-neighborhood shop downstairs.

Expected: Claude Code reads the Skill, recognizes the intent matches, runs `openmelon -p "..."` in its terminal, surfaces the streamed activity log, and reports the artifact path back into the conversation.

Then:

> Now publish that image to my V-Box account.

Expected: Claude Code runs `vbox-cli upload ...` then `vbox-cli post --media-fid ...`, surfaces the `queued_for_review` response.

---

## Smoke checks (no API keys)

```bash
# CLI works
openmelon help

# First-run wizard appears (you can Esc out at any step)
mkdir -p /tmp/smoke && cd /tmp/smoke && openmelon

# Sessions list works on existing project
openmelon resume

# Search a project's libraries (returns nothing on a fresh project, that's fine)
openmelon search "tag:character"

# Friendly error when no key is configured
openmelon -p "test"
# → "no API key for openrouter — run `openmelon setup` to configure"
```

---

## Known limitations (defer to later releases)

- Anthropic provider doesn't yet implement `ToolCaller.Chat` — picking Anthropic in the auth wizard works for legacy `-p --skill` mode but the TUI's tool-driven runtime needs OpenAI or OpenRouter. (0.4 will close this.)
- `compile_skill` shells to `skillplus` so the binary must be on PATH. `npm i -g @e8s/skillplus` handles this.
- No `edit_image` tool yet — to refine, ask the agent to call `generate_image` again with the prior generation passed as a reference. (0.4.)
