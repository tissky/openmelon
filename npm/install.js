#!/usr/bin/env node
// postinstall: download the openmelon binary for the current platform
// from this version's GitHub Release.
//
// Skipped when:
//   - OPENMELON_SKIP_DOWNLOAD=1 (CI / dev override)
//   - OPENMELON_BIN points at an existing binary (use that instead)
//   - the binary already exists at bin/openmelon-bin (re-run / cache)

import { existsSync, mkdirSync, chmodSync, createWriteStream } from "node:fs"
import { dirname, join } from "node:path"
import { fileURLToPath } from "node:url"
import { createHash } from "node:crypto"
import { readFileSync } from "node:fs"
import { get } from "node:https"

const __dirname = dirname(fileURLToPath(import.meta.url))
const BIN_DIR = join(__dirname, "bin")
const BIN_PATH = join(BIN_DIR, process.platform === "win32" ? "openmelon-bin.exe" : "openmelon-bin")

if (process.env.OPENMELON_SKIP_DOWNLOAD === "1") {
  console.log("[openmelon] OPENMELON_SKIP_DOWNLOAD=1 — skipping binary download")
  process.exit(0)
}

if (process.env.OPENMELON_BIN && existsSync(process.env.OPENMELON_BIN)) {
  console.log(`[openmelon] OPENMELON_BIN=${process.env.OPENMELON_BIN} — using that binary`)
  process.exit(0)
}

if (existsSync(BIN_PATH)) {
  console.log(`[openmelon] binary already present at ${BIN_PATH}`)
  process.exit(0)
}

// Map Node platform/arch → release artifact suffix.
const PLATFORM_MAP = {
  "darwin-arm64": "darwin-arm64",
  "darwin-x64": "darwin-amd64",
  "linux-arm64": "linux-arm64",
  "linux-x64": "linux-amd64",
}

const key = `${process.platform}-${process.arch}`
const suffix = PLATFORM_MAP[key]
if (!suffix) {
  console.error(`[openmelon] no prebuilt binary for ${key}.`)
  console.error(`[openmelon] supported: ${Object.keys(PLATFORM_MAP).join(", ")}`)
  console.error(`[openmelon] build from source: https://github.com/eight-acres-lab/openmelon`)
  process.exit(1)
}

const pkg = JSON.parse(readFileSync(join(__dirname, "package.json"), "utf8"))
const VERSION = `v${pkg.version}`
const BASE = `https://github.com/eight-acres-lab/openmelon/releases/download/${VERSION}`
const FILE = `openmelon-${VERSION}-${suffix}`
const URL = `${BASE}/${FILE}`
const SHASUMS_URL = `${BASE}/SHASUMS256.txt`

console.log(`[openmelon] downloading ${URL}`)

mkdirSync(BIN_DIR, { recursive: true })

await downloadAndVerify(URL, SHASUMS_URL, FILE, BIN_PATH)
chmodSync(BIN_PATH, 0o755)
console.log(`[openmelon] installed → ${BIN_PATH}`)

// --- helpers ---

function fetchToFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(dest)
    const req = get(url, { headers: { "user-agent": "openmelon-installer" } }, (res) => {
      if (res.statusCode === 301 || res.statusCode === 302) {
        // GitHub release downloads redirect to S3.
        file.close()
        fetchToFile(res.headers.location, dest).then(resolve).catch(reject)
        return
      }
      if (res.statusCode !== 200) {
        reject(new Error(`download ${url}: HTTP ${res.statusCode}`))
        return
      }
      res.pipe(file)
      file.on("finish", () => file.close(resolve))
    })
    req.on("error", reject)
  })
}

function fetchText(url) {
  return new Promise((resolve, reject) => {
    const req = get(url, { headers: { "user-agent": "openmelon-installer" } }, (res) => {
      if (res.statusCode === 301 || res.statusCode === 302) {
        fetchText(res.headers.location).then(resolve).catch(reject)
        return
      }
      if (res.statusCode !== 200) {
        reject(new Error(`fetch ${url}: HTTP ${res.statusCode}`))
        return
      }
      let data = ""
      res.setEncoding("utf8")
      res.on("data", (chunk) => (data += chunk))
      res.on("end", () => resolve(data))
    })
    req.on("error", reject)
  })
}

async function downloadAndVerify(url, shasumsURL, filename, dest) {
  let expected
  try {
    const shasums = await fetchText(shasumsURL)
    for (const line of shasums.split("\n")) {
      const [hash, name] = line.trim().split(/\s+/)
      if (name === filename) {
        expected = hash
        break
      }
    }
  } catch (err) {
    console.warn(`[openmelon] WARNING: could not fetch SHASUMS256.txt (${err.message}); skipping integrity check`)
  }

  await fetchToFile(url, dest)

  if (expected) {
    const actual = createHash("sha256").update(readFileSync(dest)).digest("hex")
    if (actual !== expected) {
      throw new Error(
        `[openmelon] sha256 mismatch for ${filename}: expected ${expected}, got ${actual}. ` +
          `Refusing to install a tampered binary.`,
      )
    }
    console.log(`[openmelon] sha256 verified: ${actual.slice(0, 16)}…`)
  }
}
