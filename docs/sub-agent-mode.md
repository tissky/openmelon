# Sub-Agent Mode

OpenMelon runs as a sub-agent inside any host that can shell out to a CLI:

```bash
openmelon -p "<intent>"
```

Headless `-p` reuses the same tool stack as the interactive TUI:

- Same project context (`<workdir>/.openmelon/project.json`, characters, references, materials).
- Same builtin tools (`list_characters`, `get_character`, `search`, `compile_skill`, `generate_image`, `save_artifact`, `bash`, `finish`).
- Same provider + model resolution (project credentials → global → env).
- Same bash permission policy (`project.json:settings.bash_permission_mode`).

The only difference: no UI to render an approval modal, so the bash tool is unavailable in `strict` mode. Set `/settings → trusted` (or `auto`) inside an interactive session first if your host needs to run shell from `-p`.

## What the host gets back

Stderr is the activity log (`[openmelon] project=… session=… llm=…`). Stdout, with `--json`, is a single JSON line summarizing the run: skill id + version, intent, generation prompt, image path + sha256, finished_at. The session dir under `<project>/.openmelon/sessions/<id>/` records the full conversation + generated artifacts; the host can post-process from there.

## Drop-in integrations

Skill files for Claude Code and Cursor are in [`examples/integrations/`](../examples/integrations/). Both are thin wrappers — they shell to `openmelon -p "$intent"` with sensible defaults.
