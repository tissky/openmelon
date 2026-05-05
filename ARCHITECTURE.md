# OpenMelon Architecture

OpenMelon is a content-production runtime. The core abstraction is a **project**: a durable, on-disk creative context the agent works inside, with persistent character / reference libraries and a session log.

## Operating modes

| Mode | Entry | When to use |
|---|---|---|
| **Interactive TUI** (primary) | `openmelon` (no args) | Human-in-the-loop content drafting. Multi-turn agent with slash commands, model pickers, bash approval, session resume. |
| **Headless one-shot** | `openmelon -p "<intent>"` | Scripts, sub-agent integration (Claude Code, Cursor), CI. Same tool stack, no TUI. |
| **Legacy declarative workflow** | `openmelon --project <path>` | Pre-0.3 staged pipelines (still supported for the `examples/food-exploration` style). |

## Project

A project owns durable creative context on disk under `<workdir>/.openmelon/`:

- `project.json` — name, description, persona, constraints, default models, settings
- `characters/` — registered people (portraits + description + tags) the agent can pull as reference images
- `references/` — named scenes, lighting, composition templates
- `materials/` — hash-addressed raw input pool
- `sessions/<ts>-<rnd>/` — per-run conversation log + generated images
- `artifacts/<slug>/<ts>/` — finalized outputs promoted via `save_artifact`
- `credentials.json` — per-project API keys (mode 0600), overrides global

A global registry under `~/.openmelon/` tracks projects and trusted directories.

## Tool-driven agent loop

Inside a project, the agent runs a classic ReAct-style loop:

```
                    ┌─────────────────────┐
                    │   user message      │
                    └──────────┬──────────┘
                               ▼
              ┌────────────────────────────────┐
              │  LLM (with tools registered)   │◀────────┐
              └────────────────┬───────────────┘         │
                               ▼                         │
                    finish?────yes────▶ done             │
                       │                                 │
                       no                                │
                       ▼                                 │
              dispatch tool (locally or with approval)   │
                       │                                 │
                       ▼                                 │
            tool result appended as assistant message ───┘
```

The model decides what to call. The tools are:

```
list_characters / get_character    pull people from your registry
list_references / get_reference    pull scenes, lighting, composition refs
search                             tag + grep across the project's libraries
read_file                          any file under the project workdir
compile_skill                      compile a skillplus package on demand
generate_image (refs[])            run the image model with optional anchors
save_artifact                      promote a session image to a final
bash (gated)                       inspect files, check outputs, run commands
finish                             end the loop with a summary + artifacts
```

## Skill-Plus integration

OpenMelon does not bake content "filters" into source. They live as [skillplus](https://github.com/eight-acres-lab/skillplus) packages — versioned, locale-aware, model-profile-aware bundles of system prompt + output schema. The agent's `compile_skill` tool shells out to the `skillplus` CLI; the user can pre-pick one in the TUI via `/skill` or have the model pick on its own.

A skill is a single function: take a user intent, produce a structured `generation_prompt` + `output_schema`. The image model paints from the prompt; the schema is validated.

## Bash + permission modes

The bash tool is gated by a four-tier policy controlled via `/settings`:

```
Tier 1  Trusted mode bypass    No prompt for any command.
Tier 2  Per-session allowlist  Binaries the user picked "always" this run.
Tier 3  Judge LLM              AUTO (read-only) / ASK / BLOCK (destructive).
Tier 4  User approval modal    Yes / Yes-always-for-<binary> / No.
```

Strict mode (default) only honors the judge's BLOCK; everything else asks. Auto mode runs read-only inspection silently. Trusted bypasses every gate — use only on throwaway projects.

## Sessions and provenance

Every TUI launch creates a new session dir. `messages.jsonl` records the full conversation (system + user + assistant + tool messages) as it happens; `meta.json` records project id + intent + timestamps + `resumed_from` if applicable. Generated images go into the session dir with their sha256 hash logged inline.

Resuming a session is `openmelon resume <id>` → loads the prior `messages.jsonl` into a fresh TUI's transcript, the model sees them as conversation context, a new session dir is opened to record the continuation. Old sessions are immutable.

In the legacy `--project` workflow mode, provenance is appended as JSONL lines to a single `provenance.jsonl` in the artifact dir.

## Streaming + tool-call wire format

LLM clients implement two interfaces:

- `llm.Client.{Complete, Stream}` — single-turn text completion, used by the legacy agent.
- `llm.ToolCaller.Chat` + `llm.StreamingToolCaller.StreamChat` — multi-turn message-list completion with tool calls.

Streaming via OpenAI-compatible endpoints uses `stream_options.include_usage=true`; the final chunk carries the usage block so the TUI can show running token counts. Tool-call deltas are accumulated by `tool_call_index` (vendors split `function.arguments` across many chunks) and surfaced as a single `[]ToolCall` per turn.

## Sub-agent role

`openmelon -p "..."` runs the same tool stack headless — same project context, same skills, same bash policy (driven by project settings). Drop-in Skill files for Claude Code and Cursor are in `examples/integrations/`. The bash tool requires `/settings` → `trusted` or `auto` for unattended use, since there's no UI to prompt for approval.
