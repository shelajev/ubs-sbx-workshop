#!/usr/bin/env bash
# Download pinned app fixtures and the native adapter; never run app code on the host.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
. "$ROOT/scripts/common.sh"
. "$ROOT/scripts/versions.env"
[ $# -eq 0 ] || { [ $# -eq 1 ] && [ "$1" = --app-only ]; } \
  || die 'Usage: scripts/get-materials.sh [--app-only]'
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
release_url="https://github.com/shelajev/wad-sbx-workshop/releases/download/materials-v0.1.0"
for name in incident-triage-board.bundle SHA256SUMS; do
  curl --fail --location --silent --show-error --retry 3 "$release_url/$name" --output "$work/$name"
done
for name in incident-triage-board.bundle; do
  expected="$(awk -v n="$name" '$2 == n {print $1}' "$work/SHA256SUMS")"
  [ -n "$expected" ] || die "No published checksum for $name"
  actual="$(sha256_file "$work/$name")"
  [ "$actual" = "$expected" ] || die "Checksum mismatch: $name"
done
mkdir -p "$ROOT/.local"
if [ "${1:-}" != --app-only ]; then
  "$ROOT/scripts/download-mcp.sh"
fi
if [ -e "$ROOT/.local/app" ]; then
  info 'Keeping your existing .local/app; no application files changed.'
else
  git clone --config core.autocrlf=false --config core.eol=lf "$work/incident-triage-board.bundle" "$ROOT/.local/app"
  git -C "$ROOT/.local/app" remote remove origin
fi
for ref in app-00-starter app-01-warmup-solution app-02-feature-solution; do
  git -C "$ROOT/.local/app" rev-parse --verify "$ref^{commit}" >/dev/null
done
"$ROOT/scripts/prepare-app.sh"
ok 'Materials ready. Follow workshop/00-setup.md.'
