# @e8s/openmelon

A content-creation agent that lives in your terminal. Interactive TUI for image / multimodal content workflows, headless `-p` mode for scripting.

```bash
npm install -g @e8s/openmelon @e8s/skillplus
cd ~/work/your-project
openmelon          # first run: trust → API key → init project → REPL
```

This package is a thin Node shim that downloads and spawns the [openmelon](https://github.com/eight-acres-lab/openmelon) Go binary on `npm install`. The binary is fetched from the matching GitHub Release for your platform (`darwin-arm64`, `darwin-x64`, `linux-arm64`, `linux-x64`) and verified against `SHASUMS256.txt`.

See [the main README](https://github.com/eight-acres-lab/openmelon#readme) for the full workflow (TUI commands, slash palette, bash permission modes, sessions / resume, sub-agent integration).

## Override the binary

```bash
# Skip the download (you'll provide a binary out-of-band)
OPENMELON_SKIP_DOWNLOAD=1 npm install -g @e8s/openmelon

# Or point at a local build
export OPENMELON_BIN=/path/to/your/openmelon
openmelon
```

## License

[Apache 2.0](https://github.com/eight-acres-lab/openmelon/blob/main/LICENSE).
