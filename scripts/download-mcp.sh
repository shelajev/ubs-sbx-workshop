#!/usr/bin/env bash
# Download the native adapter and verify every file before making it available.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
. "$ROOT/scripts/common.sh"
. "$ROOT/scripts/versions.env"
asset="$(beans_mcp_asset)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
url="https://github.com/shelajev/ubs-sbx-workshop/releases/download/beans-mcp-v$BEANS_MCP_VERSION"
files=("$asset" LICENSE THIRD-PARTY-NOTICES.md release-manifest.json)
for name in "${files[@]}" SHA256SUMS; do
  curl -fsSL --retry 3 "$url/$name" -o "$work/$name"
done
for name in "${files[@]}"; do
  expected="$(awk -v n="$name" '$2 == n {print $1}' "$work/SHA256SUMS")"
  [ -n "$expected" ] && [ "$(sha256_file "$work/$name")" = "$expected" ] || die "Checksum mismatch: $name"
done
dest="$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION"
mkdir -p "$dest"
for name in "${files[@]}" SHA256SUMS; do copy_file "$work/$name" "$dest/$name"; done
ok "Downloaded and verified beans-mcp $BEANS_MCP_VERSION."
