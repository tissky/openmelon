#!/usr/bin/env node
// Shim: spawn the platform binary downloaded by install.js.

import { spawn } from "node:child_process"
import { existsSync } from "node:fs"
import { dirname, join } from "node:path"
import { fileURLToPath } from "node:url"

const __dirname = dirname(fileURLToPath(import.meta.url))
const FALLBACK = process.platform === "win32" ? "openmelon-bin.exe" : "openmelon-bin"
const BIN = process.env.OPENMELON_BIN || join(__dirname, FALLBACK)

if (!existsSync(BIN)) {
  console.error(`[openmelon] binary not found at ${BIN}`)
  console.error(`[openmelon] re-run \`npm install -g @e8s/openmelon\` to fetch it,`)
  console.error(`[openmelon] or set OPENMELON_BIN=/path/to/openmelon to point at a local build.`)
  process.exit(127)
}

const child = spawn(BIN, process.argv.slice(2), { stdio: "inherit" })
child.on("exit", (code, signal) => {
  if (signal) process.kill(process.pid, signal)
  process.exit(code ?? 1)
})
child.on("error", (err) => {
  console.error(`[openmelon] failed to spawn ${BIN}: ${err.message}`)
  process.exit(1)
})
