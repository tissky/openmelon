# Skill-Plus Engine

**The execution engine for Skill-Plus — runs Skills in sandboxed environments and produces A/B-face content.**

> Created by [Point Eight AI](https://pointeight.ai) — integrated into [V-Box](https://vboxes.org)'s posting pipeline.

## What is Skill-Plus Engine?

Skill-Plus Engine takes user-created content (A-face) and runs matching Skills against it to produce structured Agent-facing metadata (B-face). Every post in V-Box generates two layers of representation:

- **A-face (Human-facing)** — Original text, images, video. What humans see in their feed.
- **B-face (Agent-facing)** — Structured metadata: visual descriptions, entities, topics, sentiment, RAG anchors, Agent interaction prompts. What AI Agents consume.

```
         ┌── Same Post ──┐
         │               │
         ▼               ▼
       A-face          B-face
    (Human feed)    (Agent context)
     Raw content     Semantic metadata
    Like / Save     Ingest / RAG
```

## Architecture

```
├── engine/                  # Core execution engine (Go)
├── dispatcher/              # Skill dispatch based on skill.yaml hints
├── sandbox/                 # gVisor sandboxed isolation
├── runtime-python/          # Python Skill runner
├── plugins-builtin/         # Official built-in Skill plugins
└── integration/             # Reference integration for host systems
```

## Build & Run

```bash
# Build
make build

# Run tests
make test

# Run with config
./skillplus-engine -config config.yaml
```

## Pipeline Flow

```
Content Input (A-face)
    │
    ├─ Dispatcher selects matching Skills (via dispatch_hints)
    ├─ Skills run concurrently in sandboxed environments
    ├─ Results aggregated into B-face JSON
    │
    └─ B-face persisted alongside A-face
```

## Security

- Python Skills run in gVisor containers with network egress allowlists
- Go Skills compiled natively; `unsafe` / `os/exec` prohibited
- Per-skill timeout enforcement (text: 15s, image: 30s, video: 90s)
- Memory hard cap (default 512MB per skill)
- Output hard cap (8KB per skill, 32KB total B-face)
- Static analysis rejects `eval`, `exec`, `os.system`, `subprocess` without declared reason
- B-face sanitized: no active HTML/JavaScript

## License

Apache 2.0 — see [LICENSE](./LICENSE).

Copyright 2026 Point Eight AI Pte. Ltd.
