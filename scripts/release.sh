#!/usr/bin/env bash
# Release a new openmelon version.
#
#   ./scripts/release.sh v0.2.0       # tag + build + GitHub release
#   ./scripts/release.sh v0.2.0 --dry-run
#
# Refuses to run with an unclean working tree. Builds darwin/linux ×
# amd64/arm64 binaries with the version baked in via -ldflags.

set -euo pipefail

VERSION="${1:-}"
DRY_RUN=""
[[ "${2:-}" == "--dry-run" ]] && DRY_RUN=1

if [[ -z "$VERSION" ]]; then
  echo "usage: $0 <version> [--dry-run]" >&2
  echo "       e.g. $0 v0.2.0" >&2
  exit 2
fi

if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-z0-9.]+)?$ ]]; then
  echo "version must look like v0.2.0 or v0.2.0-rc1, got: $VERSION" >&2
  exit 2
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "working tree is dirty; commit or stash first" >&2
  git status --short
  exit 1
fi

if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "tag $VERSION already exists" >&2
  exit 1
fi

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST="$ROOT/dist/$VERSION"
mkdir -p "$DIST"

LDFLAGS="-X github.com/eight-acres-lab/openmelon/internal/version.Version=$VERSION -s -w"

build_one() {
  local goos="$1" goarch="$2"
  local out="$DIST/openmelon-${VERSION}-${goos}-${goarch}"
  [[ "$goos" == "windows" ]] && out="$out.exe"
  echo "  → $goos/$goarch"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
    go build -ldflags "$LDFLAGS" -o "$out" ./cmd/openmelon
  (cd "$DIST" && shasum -a 256 "$(basename "$out")" >> SHASUMS256.txt)
}

echo "==> running tests"
go test ./... > /dev/null

echo "==> building binaries → $DIST"
rm -f "$DIST/SHASUMS256.txt"
build_one darwin amd64
build_one darwin arm64
build_one linux  amd64
build_one linux  arm64

echo "==> built artifacts:"
ls -lh "$DIST"

if [[ -n "$DRY_RUN" ]]; then
  echo "==> --dry-run; stopping before tag + release"
  exit 0
fi

echo "==> tagging $VERSION"
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

echo "==> creating GitHub release"
gh release create "$VERSION" \
  --title "$VERSION" \
  --generate-notes \
  "$DIST"/openmelon-* \
  "$DIST"/SHASUMS256.txt

echo "==> done. release: https://github.com/eight-acres-lab/openmelon/releases/tag/$VERSION"
