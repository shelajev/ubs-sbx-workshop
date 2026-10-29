#!/usr/bin/env bash
# Prepare the workshop's fixed working directory; preserve work when switching checkpoints.
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
[ $# -le 1 ] || { echo 'Usage: scripts/prepare-app.sh [checkpoint-tag]' >&2; exit 2; }
app="$root/sample-app"
ref="${1:-app-00-starter}"
if [ $# -eq 0 ] && [ -e "$app" ]; then
  echo 'Keeping your existing sample-app/.'
  exit 0
fi
git -C "$root/.local/app" rev-parse --verify "$ref^{commit}" >/dev/null
stage="$(mktemp -d "$root/.local/app-prepare.XXXXXX")"
trap 'rm -rf "$stage"' EXIT
git clone --quiet --no-hardlinks --config core.autocrlf=false --config core.eol=lf "$root/.local/app" "$stage/app"
git -C "$stage/app" switch --quiet -c workshop "$ref"
git -C "$stage/app" remote remove origin
git -C "$stage/app" apply --index "$root/scripts/patches/app-cli-paths.patch" \
  "$root/scripts/patches/app-browser-checks.patch"
git -C "$stage/app" -c user.name='Workshop setup' -c user.email=workshop@example.invalid \
  -c commit.gpgsign=false -c core.hooksPath=/dev/null commit --quiet \
  -m 'WORKSHOP-SETUP: fix CLI paths and browser checks'
if [ "$ref" = app-00-starter ]; then
  cat "$root/backlog/seed/task-101--warm-up-active-filter-count.md" > "$stage/app/WORKSHOP-TASK.md"
  printf '\n/WORKSHOP-TASK.md\n' >> "$stage/app/.git/info/exclude"
fi
if [ -e "$app" ]; then
  backup="$(mktemp -d "$root/.local/saved-app.XXXXXX")"
  mv "$app" "$backup/sample-app"
  printf 'Previous work saved in %s/sample-app\n' "$backup"
fi
mv "$stage/app" "$app"
printf 'This file stays outside the mounted application.\n' > "$root/host-only.txt"
printf 'Ready: sample-app/ at %s\n' "$ref"
