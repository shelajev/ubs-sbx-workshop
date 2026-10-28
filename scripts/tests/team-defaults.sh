#!/usr/bin/env bash
# Check subscription-only catch-up and the launcher's approval boundary.
set -euo pipefail
root="$(cd "$(dirname "$0")/../.." && pwd)"
scratch="$(mktemp -d)"
trap 'rm -rf "$scratch"' EXIT
fixture="$scratch/workshop with spaces"
mkdir -p "$fixture/scripts" "$fixture/chapters" "$scratch/bin"
cp "$root/scripts/use-chapter.sh" "$root/scripts/common.sh" "$fixture/scripts/"
cp -R "$root/chapters/support" "$fixture/chapters/"
printf 'coordinator\tclaude\tanthropic\tsonnet\ndeveloper\tclaude\tanthropic\tsonnet\nqa\tclaude\tanthropic\tsonnet\n' > "$scratch/expected-team"

for chapter in 04-team 05-mcp 06-human 07-factory; do
  cp -R "$root/chapters/$chapter" "$fixture/chapters/"
  diff -u "$scratch/expected-team" "$fixture/chapters/$chapter/team.tsv"
  if grep -q 'source: .*kits/pi' "$fixture/chapters/$chapter/sbxenv.yaml"; then
    echo "FAIL: $chapter requires the optional Pi kit" >&2
    exit 1
  fi
  bash "$fixture/scripts/use-chapter.sh" "$chapter"
  diff -u "$scratch/expected-team" "$fixture/factory/team.tsv"
  grep -F 'default: ../sample-app' "$fixture/factory/sbxenv.yaml"
done

printf 'coordinator claude anthropic opus\n' > "$fixture/factory/team.tsv"
bash "$fixture/scripts/use-chapter.sh" 04-team
test "$(cat "$fixture/factory/team.tsv")" = 'coordinator claude anthropic opus'

# Exercise real launcher argument assembly, with external control-plane calls recorded.
mkdir -p "$fixture/sample-app/.git" "$fixture/.local/chapters/bin"
cat > "$fixture/.local/chapters/bin/beans" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail
printf '{}\n'
STUB
cat > "$fixture/scripts/beans" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail
printf '{}\n'
STUB
cat > "$scratch/bin/sbx" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail
if [ "$1" = env ] && [ "$2" = create ]; then
  printf '%s\n' "$@" > "$WORKSHOP_LAUNCH_ARGS"
fi
STUB
chmod +x "$fixture/.local/chapters/bin/beans" "$fixture/scripts/beans" "$scratch/bin/sbx"
WORKSHOP_LAUNCH_ARGS="$scratch/create-args" PATH="$scratch/bin:$PATH" \
  bash "$fixture/chapters/support/launch" "$fixture/factory" wad-launch-test
printf '%s\n' env create "$fixture/factory/sbxenv.yaml" --env-arg name=wad-launch-test \
  --env-arg "app=$fixture/sample-app" > "$scratch/expected-args"
diff -u "$scratch/expected-args" "$scratch/create-args"
echo 'PASS: all-Claude catch-up, preserved choices and interactive launch approval.'
