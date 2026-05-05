# Roadmap

## 0.3 (current)

The "becomes a real interactive agent" release. From `-p` one-shot to a full-screen interactive TUI with the full creator workflow.

**Project / data model**
- `openmelon init`, per-project `.openmelon/` (project.json + characters/ + references/ + materials/ + sessions/ + artifacts/), global registry under `~/.openmelon/`.
- Trust list — every cwd must be confirmed once; subdirs auto-trust.
- Per-project + global API credentials at `credentials.json` (mode 0600). Project overrides global; both fall back to env. `openmelon project set-key` for per-project overrides.
- Auto-written `.openmelon/.gitignore` excludes `credentials.json` + `sessions/`.

**Content libraries**
- `openmelon character add|list|show|rm`, `openmelon reference add|list|show|rm`, hash-addressed `material add|list`.
- `openmelon search` — tag + grep across all libraries (operators: `tag:`, `kind:`, `-negative`, quoted phrases). Deliberately not vector.

**Runtime**
- Tool-using agent loop (`internal/runtime`) driven by `llm.ToolCaller`.
- Builtin tools: `list_characters`, `get_character`, `list_references`, `get_reference`, `search`, `read_file`, `compile_skill`, `generate_image` (with reference images), `save_artifact`, `bash`, `finish`.
- Streaming text via `llm.StreamingToolCaller.StreamChat` — text deltas fire as they arrive; tool-call deltas reassembled at end of turn. Usage tracked per turn.
- `RunInput.History` for multi-turn continuation.

**TUI** (`internal/tui`, bubbletea)
- Bordered-less single-line input that auto-grows to 10 lines. `›` prompt arrow.
- Slash command palette (`/`-prefix → floating list above input; Tab autocompletes, Enter executes highlighted).
- Top-of-screen header: `openmelon · <project> · <provider:llm-model> · img:<provider:img-model>`.
- Activity-aware spinner row: `⠋ Calling search · 0:12 · 1.2k in / 340 out · esc to cancel`.
- Bottom-anchored transcript (short content sits at viewport bottom, near input).
- Tool-call rendering: `⏺ name(args)` + dim `⎿ result` (`(no results)` for empty arrays).
- `/model` and `/model-image` selectors with curated top-10 / top-5 OpenRouter presets + Custom row. Hot-swaps `runtime.LLM` / image gen, persists to project.json.
- `/skill` picker — lists `skillplus list --json`, picked skill prepends a "compile_skill first" hint to the next message.
- `/settings` panel — bash permission mode (strict / auto-judge / trusted).
- Mouse wheel + PgUp/PgDn scroll.
- Ctrl+C ×2 quit; Esc cancels in-flight turn.

**Bash tool**
- 4-tier permission gate: trusted-mode bypass → per-session allowlist → judge LLM (AUTO/ASK/BLOCK) → user modal.
- Approval modal: Yes / Yes-always-for-`<binary>` / No.
- Judge LLM is the main agent LLM with a tight classifier prompt.
- Strict-mode default. Trusted mode bypasses every gate (no judge, no allowlist, no modal) — for throwaway projects only.

**Onboarding** (Codex-style, single alt-screen program)
- Trust → provider pick → API key (masked) → LLM model → image model → project init.
- Each step skipped if precondition is met. `openmelon setup` re-runs the auth wizard.

**Sessions / resume**
- Every TUI run records full transcript + meta + generated images under `sessions/<ts>-<rnd>/`.
- `openmelon resume [<id>]` lists recent or loads one. Resumed conversation renders into the new TUI's transcript and the model sees it as context. New session dir tagged with `resumed_from`.

**Headless `-p` consistency**
- `-p` reads API keys from the same project → global → env pipeline as the TUI.
- Wires the same tool stack incl. judge LLM + bash mode (no UI = no approval modal, so bash needs `/settings` → trusted/auto for headless use).

**Reliability**
- Image-gen HTTP: `DisableKeepAlives` + 3-attempt retry on transient errors (TLS bad-record-MAC, EOF, conn-reset, 5xx). Eliminates the most common transient failure mode.

## 0.4 (planned)

- Vision auto-describe on `character add` / `reference add` (single LLM call writes `.search` so the user doesn't have to).
- Anthropic ToolCaller parity (currently text-only).
- `/transcript` or Ctrl+O — full untruncated transcript view (no 240-char truncation on tool results).
- `edit_image` tool (round-trip a prior generation as a reference for refinement).
- `openmelon serve` — HTTP API surface for embedding into V-Box backend.
- Skill catalog inside TUI: `/skill add` / `/skill remove`, view installed locally vs bundled.

## 0.5

- More image providers (Stability, Replicate, Black Forest Labs direct).
- Markdown rendering in viewport (glamour).
- `@file` mention completion in input.
- Per-tool token cost tracking + per-session budget.

## 1.0

- Public Go API frozen for embedded use.
- CLI flags + project.json schema frozen.
- LTS policy + deprecation timeline.
