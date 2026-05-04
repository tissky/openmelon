# CLAUDE.md — openmelon

Guidance for Claude Code (claude.ai/code) when working in this repo.

## What this is

`openmelon` is a content-creation agent CLI. Three usage modes:

1. **Standalone CLI** — `openmelon -p "<intent>"`
2. **Sub-agent via Skill** — invoked by Claude Code / Cursor / etc. via plain `bash openmelon -p ...`
3. **Embedded Go library** — `pkg/openmelon` for V-Box backend

> Repo: https://github.com/eight-acres-lab/openmelon
> Module: `github.com/eight-acres-lab/openmelon`

## Layout

```
cmd/openmelon/        # CLI entry, agent vs workflow dispatch, publish helper
internal/
  agent/              # one-shot orchestration: skill → LLM → image → artifact
  llm/                # pluggable LLM clients (Anthropic, OpenAI, OpenRouter) + SSE
  imagegen/           # pluggable image generators (OpenAI, OpenRouter)
  skillplus/          # subprocess wrapper to the `skillplus` console script
  artifacts/          # artifact write helper
  provenance/         # provenance JSONL append helper
  project/            # legacy 0.1 project.json loader
  workflow/           # legacy 0.1 workflow runner
  generation/         # legacy 0.1 shell generation provider
  version/            # build-time version variable
pkg/
  contracts/          # public Go types
  openmelon/          # public Go API for embedding
npm/                  # @e8s/openmelon Node distribution (downloads the binary)
examples/
  food-exploration/   # legacy 0.1 declarative example
  integrations/       # Skill files for Claude Code, Cursor
assets/               # logo + provenance for the logo
docs/                 # design notes, testing recipe
```

## Commands

```bash
go build -ldflags "-X github.com/eight-acres-lab/openmelon/internal/version.Version=$(git describe --tags --always)" -o ./openmelon ./cmd/openmelon
go test ./...
```

## Conventions

- **Module path**: `github.com/eight-acres-lab/openmelon`. Do not reintroduce `github.com/Jackyffight/openmelon`.
- **Provenance is mandatory.** Every artifact gets a JSONL line. Don't add code paths that skip it.
- **No vendor model defaults.** Code returns `ErrModelRequired` when no model id is passed; the user must specify `--llm-model` and `--image-model` explicitly.
- **Subprocess to skillplus.** Don't reimplement skill compilation in Go. Contract is JSON-in / JSON-out via `internal/skillplus`.
- **Streaming is opt-in via `Agent.StreamTo`.** Tests use `Complete` for determinism; `cmd/openmelon` sets `StreamTo = os.Stderr` in agent mode.

## Adding an LLM provider

Implement `llm.Client` (`Complete` + `Stream` + `Provider` + `Model`), register in `llm.New`. Reuse `internal/llm/sse.go` for SSE parsing.

## Adding an image provider

Implement `imagegen.Generator` (`Generate` + `Provider` + `Model`), register in `imagegen.New`. The two existing implementations show the two wire shapes (REST images endpoint vs chat-completions with `modalities`).

## Versioning

`internal/version/version.go` defaults to `"dev"`. Release builds override via `-ldflags`. The release script (`scripts/release.sh`) reads `git describe --tags`.
