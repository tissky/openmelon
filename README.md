# OpenMelon

**Content-creation agent for the terminal — like Claude Code, but built for posts.**

OpenMelon is a content-production agent and runtime by [Point Eight AI](https://pointeight.ai). You can use it three ways:

1. **Standalone CLI** — `openmelon -p "write a Singapore food street post"` — talks to the model, picks the right [Skill-Plus](https://github.com/eight-acres-lab/skillplus) package, runs it, prints (or publishes) the result.
2. **Sub-agent / MCP server** — register OpenMelon as a Skill or MCP server in your existing agent (Claude Code, Cursor, Codex…) and have it delegate creation work to OpenMelon.
3. **Embedded Go library** — V-Box's own backend imports OpenMelon as the agent engine for content analysis and distribution. (The embedding contract is the `pkg/openmelon` Go package.)

OpenMelon is opinionated: it is built for content creation, not as a general-purpose agent framework.

## Status

**Pre-0.2 — current code is the workflow engine that became the foundation for the agent loop.** What runs today: load a `project.json`, pick a workflow, compile a Skill-Plus package via the reference compiler, execute stages, write artifacts + provenance JSONL. What does **not** run today: an interactive agent loop, MCP server mode, multi-vendor model clients, sub-agent delegation. See [`ROADMAP.md`](ROADMAP.md) — those land in 0.2.

## Try the food-exploration example today

You need [`skillplus`](https://github.com/eight-acres-lab/skillplus) installed (or its source tree adjacent to this one):

```bash
# from this repo's root
go run ./cmd/openmelon \
  --project examples/food-exploration/project.json \
  --compiler ../skillplus
```

That produces a JSON artifact under `artifacts/` plus a provenance JSONL line.

## What 0.2 will feel like

```bash
# install (after 0.2 release)
go install github.com/eight-acres-lab/openmelon/cmd/openmelon@latest

# one-shot
openmelon -p "Singapore 牛车水夜市的食物街快闪贴" \
  --skill skillplus:food-street-realism \
  --publish vbox

# interactive
openmelon

# MCP mode (for Claude Code / Cursor / etc.)
openmelon mcp

# HTTP mode (for V-Box backend)
openmelon serve --port 7842
```

The agent loop will route per-stage to a configurable model client (Anthropic / OpenAI / Google / OpenRouter — whatever the user has credentials for; OpenMelon does not pick for you).

## Architecture (today)

```
project.json                              ←── declarative workflow input
    │
    ▼
internal/project           internal/workflow
   load + validate    →    iterate stages
                                │
                                ▼
                       internal/skillplus      ←── shells out to `skillplus` compiler
                                │
                                ▼
                       internal/generation     ←── pluggable provider (today: command exec)
                                │
                                ▼
                       internal/artifacts      ←── write artifact
                       internal/provenance     ←── append JSONL provenance line
```

In 0.2, the agent loop sits in front of `project.json` (you don't have to write one) and the `generation` provider grows real model clients.

## Repository layout

```text
├── cmd/openmelon/        # CLI entrypoint (today: workflow runner; 0.2: agent loop)
├── internal/
│   ├── project/          # project.json loader + validation
│   ├── workflow/         # workflow / stage execution engine
│   ├── skillplus/        # subprocess adapter to the skillplus compiler
│   ├── generation/       # generation provider interface (CommandProvider today)
│   ├── artifacts/        # artifact write
│   └── provenance/       # provenance JSONL append
├── pkg/
│   ├── contracts/        # public Go types — the embedding contract
│   └── openmelon/        # public Go API for embedded use
├── config/               # example configs
├── examples/             # food-exploration end-to-end example
└── docs/                 # design notes
```

Modules that exist in the spec but are **deferred to 0.4**: `memory`, `labeling`, `review`, `roles`, `planner`. They were previously empty skeleton files; we deleted them rather than ship hollow placeholders. The 0.2 agent loop will use the simplest possible JSONL-on-disk substitute until those come back as real implementations.

## Where this fits in the e8s ecosystem

| Repo | Role |
|---|---|
| **[vbox-cli](https://github.com/eight-acres-lab/vbox-cli)** | V-Box terminal client — OpenMelon calls this as a builtin tool to publish |
| **[openmelon](https://github.com/eight-acres-lab/openmelon)** (this) | Content-creation agent — orchestrates skill compile + generation + publish |
| **[skillplus](https://github.com/eight-acres-lab/skillplus)** | Compilable skill packages — OpenMelon's "skills" come from here |

End-to-end story: OpenMelon receives a creation intent → picks a skillplus package → compiles it → runs the resulting stages with a model client → publishes the result via vbox-cli. Each piece is independently usable.

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`GOVERNANCE.md`](GOVERNANCE.md). RFC process for protocol/contract changes in [`RFC.md`](RFC.md). Code of Conduct in [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md). Security disclosures via [GitHub security advisories](https://github.com/eight-acres-lab/openmelon/security/advisories/new).

## License

[Apache 2.0](LICENSE).
