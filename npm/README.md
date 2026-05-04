# @e8s/openmelon

```bash
npm install -g @e8s/openmelon
```

This package is a thin Node shim that downloads and spawns the [openmelon](https://github.com/eight-acres-lab/openmelon) Go binary on `npm install`. The binary is fetched from the matching GitHub Release for your platform (`darwin-arm64`, `darwin-x64`, `linux-arm64`, `linux-x64`) and verified against `SHASUMS256.txt`.

See [the main README](https://github.com/eight-acres-lab/openmelon#readme) for usage.

## Override the binary

```bash
# Skip the download (you'll provide a binary out-of-band)
OPENMELON_SKIP_DOWNLOAD=1 npm install -g @e8s/openmelon

# Or point at a local build
export OPENMELON_BIN=/path/to/your/openmelon
openmelon -p "..." --skill skillplus:food-street-realism
```

## License

[Apache 2.0](https://github.com/eight-acres-lab/openmelon/blob/main/LICENSE).
