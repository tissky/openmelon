---
name: create-vbox-content
description: Generate a realistic V-Box-style post (image + structured output) from a one-sentence intent, using OpenMelon + Skill-Plus packages. Use when the user wants to create social content, food-exploration posts, lifestyle photos, or anything described in terms of "a real-feeling post about X". Do NOT use for code generation, debugging, or general Q&A.
---

# Create V-Box content via OpenMelon

You have access to `openmelon`, a content-creation agent CLI. It compiles a [Skill-Plus](https://github.com/eight-acres-lab/skillplus) package (a structured "creative recipe"), sends it to an LLM along with the user's intent, parses the structured response, generates an image with OpenAI's image API, and writes everything plus a provenance line to `.openmelon/artifacts/`.

## When to use this skill

Trigger this when the user asks for any of:

- a realistic social-media post (food, travel, lifestyle, "phone snapshot" feel)
- a "real探店" / 探店 / Xiaohongshu-style image
- "make me a post about X" / "generate content for V-Box about X"
- "create a fire-side / market / shop visit image"

Do **not** trigger this for: code, documentation, generic Q&A, image generation that doesn't need a stabilized "skill-as-filter" — for raw image gen the user can call OpenAI's API directly.

## How to invoke

Single command. Pass the user's intent verbatim (Chinese or English both work):

```bash
openmelon -p "<the user's intent — keep it natural and complete>" \
  --skill skillplus:food-street-realism
```

Available skills (today):
- `skillplus:food-street-realism` — street-food / shop-visit / "real-feeling" posts. Optimized for `gpt-image-family` model profile.

If the user requests a different style, ask which skill to use; if no good fit exists, suggest using `food-street-realism` and noting that the skill catalogue is growing.

## Optional flags worth knowing

| Flag | When to use |
|---|---|
| `--llm openai` | force OpenAI for the LLM step (default `auto` picks Anthropic if both keys set) |
| `--llm-model <id>` | override the LLM model |
| `--image-model <id>` | override the image model (e.g. `dall-e-3` instead of `gpt-image-1`) |
| `--locale <code>` | locale for skill content (default `zh-CN`) |
| `--image=false` | structuring only — useful for previewing the prompt before paying for image gen |
| `--json` | also print a structured summary to stdout (useful when chaining with other tools) |

## Output you'll see

Streamed stderr while the LLM is generating:

```
[openmelon ...] skill=skillplus:food-street-realism llm=anthropic/claude-sonnet-4-6 image=openai/gpt-image-1
[openmelon] intent: ...
{"scene_interpretation":...,"generation_prompt":"...",...}
[openmelon] skill compiled: food-street-realism@0.1.0
[openmelon] generation prompt: ...
[openmelon] image: .openmelon/artifacts/food-street-realism-20260504-203045.png (sha256=abc123def456)
[openmelon] provenance: .openmelon/artifacts/provenance.jsonl
[openmelon] duration: 24.3s
```

After the command succeeds:

1. Tell the user the artifact path so they can open the PNG.
2. Quote 1-2 sentences of the `generation_prompt` so they see what the skill produced.
3. If they look happy, offer to publish via the `publish-vbox-content` skill.

## Failure modes

- **No API key** — error mentions `ANTHROPIC_API_KEY` or `OPENAI_API_KEY`. Tell the user which env var to set.
- **`skillplus` not found** — they haven't `pip install skillplus`. Tell them.
- **Skill not found** — message lists where it looked. Pass `--skill-root <dir>` if they keep packages somewhere unusual.
- **Image generation timeout** — `gpt-image-1` can take 30-90s. The default timeout is 5 minutes; if hit, retry once.
- **OpenAI returns content-policy block** — re-run with a less ambiguous intent. Don't force.

## Don't

- Don't run with `--publish vbox` from this skill — it's for create only. The publish skill is separate so the user can review the image first.
- Don't add ad-hoc flags the user didn't ask for. The defaults are tuned for the food-street-realism skill.
- Don't summarize the streamed JSON output back to the user — show them the generation_prompt and the image path; that's the relevant signal.
