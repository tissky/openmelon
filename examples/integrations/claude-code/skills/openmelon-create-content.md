---
name: create-vbox-content
description: Generate a realistic V-Box-style post (image + structured output) from a one-sentence intent using OpenMelon. Use when the user wants to create social content, food-exploration posts, lifestyle photos, or anything described as "a real-feeling post about X". Do NOT use for code generation, debugging, or general Q&A.
---

# Create V-Box content via OpenMelon

You have access to `openmelon`, a content-creation agent CLI. In headless `-p` mode it runs the same tool stack as its interactive TUI: a tool-using LLM that pulls characters / references / search results from the project's `.openmelon/` directory, compiles a [skillplus](https://github.com/eight-acres-lab/skillplus) package when relevant, and generates an image with the user's configured image model.

## When to use this skill

Trigger this when the user asks for any of:

- a realistic social-media post (food, travel, lifestyle, "phone snapshot" feel)
- a 探店 / Xiaohongshu-style image
- "make me a post about X" / "generate content for V-Box about X"
- "create a fire-side / market / shop visit image"

Do **not** trigger this for: code, documentation, generic Q&A, raw image gen with no skill needed (the user can call the image API directly).

## Preconditions

- `openmelon` on PATH. Verify with `which openmelon`. If missing, tell the user to `npm i -g @e8s/openmelon`.
- `skillplus` on PATH. Verify with `which skillplus`. If missing: `npm i -g @e8s/skillplus`.
- A configured project. Check `<cwd>/.openmelon/project.json`. If missing, tell the user to `cd` into a project directory or run `openmelon` once to walk the first-run wizard.
- API key configured. Verify by `cat ~/.openmelon/credentials.json` exists. If not, tell the user to run `openmelon setup`.

## How to invoke

Single command. Pass the user's intent verbatim (Chinese or English both work):

```bash
openmelon -p "<the user's intent — keep it natural and complete>"
```

The agent decides which characters / references / skills to pull on its own. If the user explicitly names a skill, append `--skill <slug>` (bare slug, no `skillplus:` prefix).

## Optional flags

| Flag | When to use |
|---|---|
| `--llm-model <id>` | override the LLM model for this run |
| `--image-model <id>` | override the image model |
| `--image=false` | structuring only — useful for previewing the prompt before paying for image gen |
| `--json` | also print a structured summary to stdout (chain with other tools) |

The provider, default model, and bash permission mode are read from the project's `project.json` and `~/.openmelon/config.json`. Don't pass flags the user didn't ask for.

## Output you'll see

Activity log on stderr (turn-by-turn):

```
[openmelon] project=ai-talks session=20260506-... llm=openrouter/openai/gpt-5.5 image=openrouter/google/gemini-2.5-flash-image
[openmelon] intent: ...
[turn 1] reply (finish=tool_calls, tool_calls=2)
[turn 1] → list_characters({"query":"..."})
[turn 1] ← [{"slug":"...",...}]
[turn 1] → generate_image({"prompt":"...","reference_images":["/.../portrait.png"]})
[turn 1] ← {"path":".../draft-1.png","sha256":"..."}
...
```

After the command succeeds:

1. Tell the user the image path under `.openmelon/sessions/<id>/`.
2. Mention the resume id so they can continue: `openmelon resume <id>`.
3. If they look happy, offer to publish via the `publish-vbox-content` skill.

## Failure modes

- **No API key** — error says "run `openmelon setup`". Tell the user.
- **`skillplus` not found** — they haven't installed it. `npm i -g @e8s/skillplus`.
- **Bash unavailable in headless** — the agent tried to call bash but the project's bash mode is `strict`. Tell the user to run `openmelon` interactively, do `/settings → Auto-judge` (or Trusted), then retry.
- **Image generation timeout** — image models can take 30-90s. The default timeout is 5 minutes; transient TLS / 5xx are retried 3 times automatically. If still failing, the user's network may have a middlebox issue.
- **Content-policy block** — re-run with a less ambiguous intent. Don't force.

## Don't

- Don't run with `--publish vbox` from this skill — it's create-only. The publish skill is separate so the user can review the image first.
- Don't pass `--skill skillplus:<name>` — that prefix is the old format and the bare slug is what `skillplus` expects today.
- Don't summarize the activity log back to the user — show them the session dir + the image path. That's the relevant signal.
