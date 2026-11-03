#!/usr/bin/env bash
# harness-eval: reproducible agent harness health checks for qrurl monorepo.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EVAL_DIR="$REPO_ROOT/harness-eval"
BASELINE="$EVAL_DIR/baseline.json"
REPORT="$EVAL_DIR/REPORT.md"
RESULTS="$EVAL_DIR/.last-run.json"
COMPARE="${1:-}"

export PATH="/usr/local/go/bin:${PATH:-}"

RUN_ID="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
HOST="$(uname -s 2>/dev/null || echo unknown)"

# --- helpers ---
pass() { echo "pass"; }
fail() { echo "fail"; }

run_timed() {
  local name="$1"
  shift
  local start end elapsed_ms status output
  start="$(date +%s%3N)"
  set +e
  output="$("$@" 2>&1)"
  status=$?
  set -e
  end="$(date +%s%3N)"
  elapsed_ms=$((end - start))
  printf '%s' "$output" | tail -c 8000 >"$EVAL_DIR/.cmd-out.tmp"
  jq -n \
    --arg name "$name" \
    --arg cmd "$*" \
    --argjson status "$status" \
    --argjson elapsed_ms "$elapsed_ms" \
    --arg output "$(tr '\n' ' ' <"$EVAL_DIR/.cmd-out.tmp" | head -c 4000)" \
    '{name: $name, command: $cmd, exit_code: $status, elapsed_ms: $elapsed_ms, output_tail: $output}'
}

check_file_exists() {
  local id="$1" path="$2"
  if [[ -f "$REPO_ROOT/$path" || -d "$REPO_ROOT/$path" ]]; then
    jq -n --arg id "$id" --arg path "$path" '{id: $id, status: "pass", detail: $path}'
  else
    jq -n --arg id "$id" --arg path "$path" '{id: $id, status: "fail", detail: ("missing: " + $path)}'
  fi
}

check_json_valid() {
  local id="$1" path="$2"
  if ! [[ -f "$REPO_ROOT/$path" ]]; then
    jq -n --arg id "$id" --arg path "$path" '{id: $id, status: "fail", detail: ("missing: " + $path)}'
    return
  fi
  if jq empty "$REPO_ROOT/$path" 2>/dev/null; then
    jq -n --arg id "$id" --arg path "$path" '{id: $id, status: "pass", detail: $path}'
  else
    jq -n --arg id "$id" --arg path "$path" '{id: $id, status: "fail", detail: ("invalid JSON: " + $path)}'
  fi
}

# --- static checks ---
static_checks=()

DOC_FILES=(AGENTS.md)
for doc in "${DOC_FILES[@]}"; do
  static_checks+=("$(check_file_exists "doc:$doc" "$doc")")
done

PATHS_IN_AGENTS=(
  proto/
  server/
  client/
  .agents/hooks.json
  .agents/scripts/lint-hook.sh
  buf.yaml
  buf.gen.yaml
  firebase.json
)
for p in "${PATHS_IN_AGENTS[@]}"; do
  static_checks+=("$(check_file_exists "agents-path:$p" "$p")")
done

static_checks+=("$(check_json_valid "json:agents-hooks" ".agents/hooks.json")")

# Hook command path consistency
HOOK_CMD="$(jq -r '.. | objects | select(has("command")) | .command' "$REPO_ROOT/.agents/hooks.json" 2>/dev/null | head -1 || true)"
HOOK_RESOLVES="fail"
HOOK_DETAIL=""
if [[ -n "$HOOK_CMD" ]]; then
  if [[ "$HOOK_CMD" == *".agents/scripts/lint-hook.sh"* ]] && [[ -x "$REPO_ROOT/.agents/scripts/lint-hook.sh" ]]; then
    HOOK_RESOLVES="pass"
    HOOK_DETAIL="hook points to .agents/scripts/lint-hook.sh"
  elif [[ "$HOOK_CMD" == *"./scripts/lint-hook.sh"* ]]; then
    if [[ -f "$REPO_ROOT/.agents/scripts/lint-hook.sh" ]]; then
      HOOK_RESOLVES="fail"
      HOOK_DETAIL="hooks.json uses ./scripts/lint-hook.sh but script is at .agents/scripts/lint-hook.sh"
    else
      HOOK_RESOLVES="fail"
      HOOK_DETAIL="hook script not found for: $HOOK_CMD"
    fi
  elif [[ -f "$REPO_ROOT/.agents/scripts/lint-hook.sh" ]]; then
    HOOK_RESOLVES="pass"
    HOOK_DETAIL="hook resolves to .agents/scripts/lint-hook.sh"
  else
    HOOK_DETAIL="unrecognized hook command: $HOOK_CMD"
  fi
fi
static_checks+=("$(jq -n --arg s "$HOOK_RESOLVES" --arg d "$HOOK_DETAIL" '{id: "hook:lint-path", status: $s, detail: $d}')")

# Secret patterns in harness files (heuristic)
secret_hits=()
while IFS= read -r -d '' f; do
  if grep -qE '(AKIA[0-9A-Z]{16}|-----BEGIN (RSA |EC )?PRIVATE KEY-----|ghp_[A-Za-z0-9]{20,}|xox[baprs]-)' "$f" 2>/dev/null; then
    secret_hits+=("$f")
  fi
done < <(find "$REPO_ROOT/.agents" "$REPO_ROOT/.cursor" -type f \( -name '*.json' -o -name '*.sh' -o -name '*.md' \) -print0 2>/dev/null)

if ((${#secret_hits[@]} == 0)); then
  static_checks+=('{"id":"secrets:harness","status":"pass","detail":"no obvious secret patterns in .agents/.cursor"}')
else
  static_checks+=("$(jq -n --argjson files "$(printf '%s\n' "${secret_hits[@]}" | jq -R . | jq -s .)" '{id: "secrets:harness", status: "fail", detail: $files}')")
fi

# --- command checks (AGENTS.md quality gate) ---
command_checks=()

if command -v go >/dev/null 2>&1; then
  command_checks+=("$(run_timed "go-fix" bash -lc "cd \"$REPO_ROOT/server\" && go fix ./...")")
  command_checks+=("$(run_timed "go-vet" bash -lc "cd \"$REPO_ROOT/server\" && go vet ./...")")
  command_checks+=("$(run_timed "go-test" bash -lc "cd \"$REPO_ROOT/server\" && go test ./...")")
  command_checks+=("$(run_timed "go-build" bash -lc "cd \"$REPO_ROOT/server\" && go build -o /tmp/qrurl-harness main.go")")
else
  command_checks+=('{"name":"go-skipped","command":"go","exit_code":1,"elapsed_ms":0,"output_tail":"go not in PATH"}')
fi

if [[ -d "$REPO_ROOT/client/node_modules" ]]; then
  command_checks+=("$(run_timed "client-lint" bash -lc "cd \"$REPO_ROOT/client\" && npm run lint")")
  command_checks+=("$(run_timed "client-tsc" bash -lc "cd \"$REPO_ROOT/client\" && npx tsc --noEmit")")
  command_checks+=("$(run_timed "client-build" bash -lc "cd \"$REPO_ROOT/client\" && npm run build")")
else
  command_checks+=('{"name":"client-skipped","command":"npm","exit_code":1,"elapsed_ms":0,"output_tail":"client/node_modules missing; run npm ci in client/"}')
fi

if command -v buf >/dev/null 2>&1 || [[ -x "$REPO_ROOT/client/node_modules/.bin/buf" ]]; then
  BUF_BIN="buf"
  [[ -x "$REPO_ROOT/client/node_modules/.bin/buf" ]] && BUF_BIN="$REPO_ROOT/client/node_modules/.bin/buf"
  command_checks+=("$(run_timed "buf-lint" bash -lc "cd \"$REPO_ROOT\" && \"$BUF_BIN\" lint")")
else
  command_checks+=('{"name":"buf-lint-skipped","command":"buf lint","exit_code":0,"elapsed_ms":0,"output_tail":"buf not available; skipped"}')
fi

# Assemble JSON
static_json="$(printf '%s\n' "${static_checks[@]}" | jq -s '.')"
command_json="$(printf '%s\n' "${command_checks[@]}" | jq -s '.')"

pass_static="$(echo "$static_json" | jq '[.[] | select(.status=="pass")] | length')"
fail_static="$(echo "$static_json" | jq '[.[] | select(.status=="fail")] | length')"
pass_cmd="$(echo "$command_json" | jq '[.[] | select(.exit_code==0)] | length')"
fail_cmd="$(echo "$command_json" | jq '[.[] | select(.exit_code!=0)] | length')"

jq -n \
  --arg run_id "$RUN_ID" \
  --arg host "$HOST" \
  --arg repo "$REPO_ROOT" \
  --argjson static "$static_json" \
  --argjson commands "$command_json" \
  --argjson pass_static "$pass_static" \
  --argjson fail_static "$fail_static" \
  --argjson pass_commands "$pass_cmd" \
  --argjson fail_commands "$fail_cmd" \
  '{
    schema_version: 1,
    run_id: $run_id,
    host: $host,
    repo: $repo,
    summary: {
      static_pass: $pass_static,
      static_fail: $fail_static,
      command_pass: $pass_commands,
      command_fail: $fail_commands
    },
    static_checks: $static,
    command_checks: $commands
  }' >"$RESULTS"

# Compare to baseline
diff_section=""
if [[ -f "$BASELINE" ]]; then
  diff_section="$(jq -n \
    --slurpfile cur "$RESULTS" \
    --slurpfile base "$BASELINE" \
    'def ids($a): ($a.static_checks // []) | map(.id) + (($a.command_checks // []) | map(.name));
     def status_map($a):
       (($a.static_checks // []) | map({(.id): .status}) | add // {}) +
       (($a.command_checks // []) | map({(.name): (if .exit_code==0 then "pass" else "fail" end)}) | add // {});
     ($cur[0]) as $c | ($base[0]) as $b |
     {
       static_regressions: [($b.static_checks // [])[] | select(.status=="pass") | .id as $id | ($c.static_checks // [])[] | select(.id==$id and .status=="fail") | .id],
       static_improvements: [($b.static_checks // [])[] | select(.status=="fail") | .id as $id | ($c.static_checks // [])[] | select(.id==$id and .status=="pass") | .id],
       command_regressions: [($b.command_checks // [])[] | select(.exit_code==0) | .name as $n | ($c.command_checks // [])[] | select(.name==$n and .exit_code!=0) | .name],
       command_improvements: [($b.command_checks // [])[] | select(.exit_code!=0) | .name as $n | ($c.command_checks // [])[] | select(.name==$n and .exit_code==0) | .name],
       new_static_ids: [($c.static_checks // [])[] | .id] - [($b.static_checks // [])[] | .id],
       new_command_names: [($c.command_checks // [])[] | .name] - [($b.command_checks // [])[] | .name]
     }')"
fi

if [[ "$COMPARE" == "--write-baseline" ]]; then
  cp "$RESULTS" "$BASELINE"
  echo "Wrote baseline: $BASELINE"
fi

# REPORT.md
{
  echo "# Harness evaluation report"
  echo
  echo "- Run ID: \`$RUN_ID\`"
  echo "- Host: $HOST"
  echo "- Static: **${pass_static} pass** / **${fail_static} fail**"
  echo "- Commands: **${pass_cmd} pass** / **${fail_cmd} fail**"
  echo
  echo "## Static checks"
  echo
  echo "$static_json" | jq -r '.[] | "- [\(.status)] **\(.id)**: \(.detail)"'
  echo
  echo "## Command checks"
  echo
  echo "$command_json" | jq -r '.[] | "- [\((if .exit_code==0 then "pass" else "fail" end))] **\(.name)** (\(.elapsed_ms) ms): `\(.command)`"'
  if [[ -n "$diff_section" ]]; then
    echo
    echo "## Baseline comparison"
    echo
    echo '```json'
    echo "$diff_section" | jq .
    echo '```'
  else
    echo
    echo "## Baseline comparison"
    echo
    echo "No baseline.json yet. Run \`./harness-eval/run.sh --write-baseline\` after a green run to seed comparison."
  fi
} >"$REPORT"

cat "$REPORT"

if [[ "$fail_static" -gt 0 || "$fail_cmd" -gt 0 ]]; then
  exit 1
fi
