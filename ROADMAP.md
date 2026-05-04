# Roadmap

## 0.2 (current)

- One-shot agent loop: `openmelon -p "<intent>" --skill ...`
- LLM clients: Anthropic, OpenAI, OpenRouter (auto-detect via env)
- Image generation: OpenAI `/v1/images/generations`, OpenRouter chat-completions image models
- Token-streamed LLM output to stderr
- `--publish vbox` shells to `vbox-cli` to upload + post
- Provenance JSONL on every run

## 0.3

- Interactive REPL (`openmelon` with no args), bubbletea TUI
- Multi-candidate scene picker — when a skill emits multiple `scene_interpretation` candidates, pick before image gen
- `openmelon serve` — HTTP API for embedding into V-Box backend
- Skill files for Claude Code / Cursor / Codex (currently CLI-only via Bash)

## 0.4

- Long-term memory + labeling + review modules return as real implementations
- More image providers (Stability, Replicate)
- Skill catalog: more reference packages beyond food-street-realism + brand-logo

## 1.0

- Public Go API frozen for embedded use
- CLI flags + provenance schema frozen
- LTS policy + deprecation timeline
