#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
. "$ROOT/scripts/common.sh"
[ -d "$ROOT/.local/app/.git" ] || die 'Run scripts/get-materials.sh before this test.'
require_cmd node 'Node.js' 'install Node.js 22 or newer.'
work="$(mktemp -d "$ROOT/.local/app-paths-test.XXXXXX")"
cleanup() {
  case "$work" in "$ROOT"/.local/app-paths-test.*) rm -rf "$work" ;; esac
}
trap cleanup EXIT
fixture="$work/workshop path with spaces"
mkdir -p "$fixture/.local" "$fixture/scripts/patches" "$fixture/backlog/seed"
cp "$ROOT/scripts/prepare-app.sh" "$fixture/scripts/"
cp "$ROOT/scripts/patches/app-cli-paths.patch" "$fixture/scripts/patches/"
cp "$ROOT/scripts/patches/app-browser-checks.patch" "$fixture/scripts/patches/"
cp "$ROOT/backlog/seed/task-101--warm-up-active-filter-count.md" "$fixture/backlog/seed/"
git clone --quiet --no-hardlinks "$ROOT/.local/app" "$fixture/.local/app"

for ref in app-00-starter app-01-warmup-solution app-02-feature-solution; do
  GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=commit.gpgsign GIT_CONFIG_VALUE_0=true \
    bash "$fixture/scripts/prepare-app.sh" "$ref"
  app="$fixture/sample-app"
  [ "$(git -C "$app" rev-parse HEAD^)" = "$(git -C "$app" rev-parse "$ref^{commit}")" ]
  [ "$(git -C "$app" show -s --format='%an <%ae>' HEAD)" = 'Workshop setup <workshop@example.invalid>' ]
  [ -z "$(git -C "$app" status --porcelain)" ]
  if [ "$ref" = app-00-starter ]; then
    cmp "$ROOT/backlog/seed/task-101--warm-up-active-filter-count.md" "$app/WORKSHOP-TASK.md"
    grep -Fx '/WORKSHOP-TASK.md' "$app/.git/info/exclude" >/dev/null
  else
    [ ! -e "$app/WORKSHOP-TASK.md" ]
    ! grep -Fx '/WORKSHOP-TASK.md' "$app/.git/info/exclude" >/dev/null
  fi
  node --input-type=module - "$app" "$work" <<'NODE'
import assert from 'node:assert/strict';
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { spawnSync } from 'node:child_process';
import { pathToFileURL } from 'node:url';

for (const name of ['migrate', 'seed']) {
  const source = readFileSync(join(process.argv[2], 'src/db', `${name}.ts`), 'utf8');
  const imports = source.match(/^import .* from "node:url";$/m)?.[0];
  const guard = source.match(/^if \(process\.argv\[1\].*\{$/m)?.[0];
  assert.ok(imports && guard, `${name}: expected a URL import and CLI guard`);
  const fixture = join(process.argv[3], `${name} path #%.mjs`);
  writeFileSync(fixture, `${imports}\n${guard}\n  console.log('entered');\n}\n`);
  const direct = spawnSync(process.execPath, [fixture], { encoding: 'utf8' });
  assert.equal(direct.status, 0, direct.stderr);
  assert.equal(direct.stdout.trim(), 'entered', `${name}: CLI guard skipped the entry point`);
  const imported = spawnSync(process.execPath, [
    '--input-type=module', '-e', `await import(${JSON.stringify(pathToFileURL(fixture).href)})`,
  ], { encoding: 'utf8' });
  assert.equal(imported.status, 0, imported.stderr);
  assert.equal(imported.stdout, '', `${name}: importing the module ran its CLI`);
}
NODE
  head="$(git -C "$app" rev-parse HEAD)"
  printf 'learner work\n' > "$app/learner-work.txt"
  bash "$fixture/scripts/prepare-app.sh"
  [ "$(git -C "$app" rev-parse HEAD)" = "$head" ]
  [ "$(cat "$app/learner-work.txt")" = 'learner work' ]
done
backups=0
for backup in "$fixture"/.local/saved-app.*; do
  [ "$(cat "$backup/sample-app/learner-work.txt")" = 'learner work' ]
  backups=$((backups + 1))
done
[ "$backups" -eq 2 ]
printf 'Checkpoint setup, CLI paths, preserved work, and backups passed.\n'
