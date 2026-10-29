#!/usr/bin/env bash
# HOST: create (or refresh) the workshop's own Beans backlog under .local/.
#
# The backlog lives in the workshop control directory, which is never mounted into a
# sandbox. The MCP adapter reads this exact directory; nothing copies it anywhere else.
#
# Usage: scripts/backlog-init.sh [--force] [--disposable]
#   --force       replace seeded task files even if they were edited
#   --disposable  write the .workshop-disposable marker (allows result notes in the workshop backlog)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib/common.sh
. "$ROOT/scripts/common.sh"

FORCE=0
DISPOSABLE=0
for arg in "$@"; do
  case "$arg" in
    --force) FORCE=1 ;;
    --disposable) DISPOSABLE=1 ;;
    -h|--help) sed -n '2,12p' "$0"; exit 0 ;;
    *) die "unknown argument: $arg (see --help)" ;;
  esac
done

BACKLOG_DIR="$CONTROL_DIR/beans"
DATA_DIR="$BACKLOG_DIR/.beans"
CONFIG_FILE="$BACKLOG_DIR/.beans.yml"
SEED_DIR="$ROOT/backlog/seed"

require_cmd "$(beans_bin)" "Beans CLI" "scripts/install-beans.sh"

mkdir -p "$DATA_DIR"

if [ ! -f "$CONFIG_FILE" ]; then
  cat > "$CONFIG_FILE" <<'YML'
# Workshop backlog configuration. The MCP adapter is always started with an explicit
# --config pointing at this file, so Beans never searches upward from a caller's
# working directory and can never answer from an unrelated backlog.
beans:
    path: .beans
    prefix: task-
    id_length: 4
    default_status: todo
    default_type: task
YML
  info "wrote $CONFIG_FILE"
fi

copied=0
skipped=0
for seed in "$SEED_DIR"/*.md; do
  name="$(basename "$seed")"
  target="$DATA_DIR/$name"
  if [ -f "$target" ] && [ "$FORCE" -eq 0 ]; then
    if ! cmp -s "$seed" "$target"; then
      skipped=$((skipped + 1))
      warn "kept your edited $name (use --force to replace it)"
    fi
    continue
  fi
  copy_file "$seed" "$target" || die "could not write $target"
  copied=$((copied + 1))
done

copy_file "$SEED_DIR/SEED-VERSION" "$BACKLOG_DIR/SEED-VERSION"

if [ "$DISPOSABLE" -eq 1 ]; then
  printf 'This backlog is disposable workshop data. Created by scripts/backlog-init.sh --disposable.\n' \
    > "$DATA_DIR/.workshop-disposable"
  warn "disposable marker written: the result-note tool may write to this backlog"
else
  rm -f "$DATA_DIR/.workshop-disposable"
fi

info "backlog ready at $DATA_DIR (seed $(tr -d '\n' < "$BACKLOG_DIR/SEED-VERSION"), $copied written, $skipped preserved)"
"$(beans_bin)" --config "$CONFIG_FILE" --beans-path "$DATA_DIR" list
