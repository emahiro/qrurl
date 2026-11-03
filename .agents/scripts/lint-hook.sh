#!/usr/bin/env bash
set -e

# Antigravity PostToolUse Hook for Monorepo (Go Server & React/TypeScript Client)
# Input: JSON on stdin (with .toolCall.args.TargetFile)
# Output: {} on stdout

# Read stdin to payload
PAYLOAD=$(cat)

# Extract TargetFile from payload (with jq fallback to grep/sed)
TARGET_FILE=""
if command -v jq >/dev/null 2>&1; then
  TARGET_FILE=$(echo "$PAYLOAD" | jq -r '.toolCall.args.TargetFile // empty' 2>/dev/null || true)
fi
if [ -z "$TARGET_FILE" ]; then
  TARGET_FILE=$(echo "$PAYLOAD" | grep -o '"TargetFile":[ ]*"[^"]*"' | head -n 1 | sed 's/"TargetFile":[ ]*"//;s/"//' || true)
fi

# Determine repo root (parent of .agents)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Normalize TARGET_FILE to absolute path if needed
if [ -n "$TARGET_FILE" ]; then
  case "$TARGET_FILE" in
    /*) ;;
    *) TARGET_FILE="$REPO_ROOT/$TARGET_FILE" ;;
  esac
fi

# Run verification based on target file path/type
if [ -n "$TARGET_FILE" ] && [ -f "$TARGET_FILE" ]; then
  case "$TARGET_FILE" in
    *.go)
      # Server (Go) verification & modernize
      if [ -d "$REPO_ROOT/server" ] && command -v go >/dev/null 2>&1; then
        (cd "$REPO_ROOT/server" && go fix ./... >/dev/null 2>&1) || true
      fi
      if command -v gofmt >/dev/null 2>&1; then
        gofmt -w "$TARGET_FILE" >/dev/null 2>&1 || true
      fi
      if command -v goimports >/dev/null 2>&1; then
        MODULE_NAME=$(cd "$REPO_ROOT/server" 2>/dev/null && go list -m 2>/dev/null || echo "github.com/emahiro/qrurl/server")
        goimports -w -local "$MODULE_NAME" "$TARGET_FILE" >/dev/null 2>&1 || true
      fi
      if [ -d "$REPO_ROOT/server" ] && command -v go >/dev/null 2>&1; then
        (cd "$REPO_ROOT/server" && go vet ./... >/dev/null 2>&1) || true
      fi
      ;;
    *.ts|*.tsx|*.js|*.jsx|*.json|*.css)
      # Client (React/TypeScript) verification
      # Only run Biome if the file is inside client/ or matches client scope
      if [[ "$TARGET_FILE" == "$REPO_ROOT/client/"* ]] || [ -f "$REPO_ROOT/client/biome.json" ]; then
        (
          cd "$REPO_ROOT/client" 2>/dev/null || exit 0
          if command -v biome >/dev/null 2>&1; then
            biome check --write "$TARGET_FILE" >/dev/null 2>&1 || true
          elif command -v npx >/dev/null 2>&1; then
            npx @biomejs/biome check --write "$TARGET_FILE" >/dev/null 2>&1 || true
          elif [ -x "$REPO_ROOT/client/node_modules/.bin/biome" ]; then
            "$REPO_ROOT/client/node_modules/.bin/biome" check --write "$TARGET_FILE" >/dev/null 2>&1 || true
          else
            NATIVE_BIOME=$(find "$REPO_ROOT/client/node_modules/@biomejs" -type f -name biome 2>/dev/null | head -n 1)
            if [ -n "$NATIVE_BIOME" ] && [ -x "$NATIVE_BIOME" ]; then
              "$NATIVE_BIOME" check --write "$TARGET_FILE" >/dev/null 2>&1 || true
            fi
          fi
        ) || true
      fi
      ;;
    *.proto)
      # Protocol Buffers verification
      if command -v buf >/dev/null 2>&1; then
        (cd "$REPO_ROOT" && buf format -w "$TARGET_FILE" >/dev/null 2>&1) || true
      elif [ -x "$REPO_ROOT/client/node_modules/.bin/buf" ]; then
        (cd "$REPO_ROOT" && "$REPO_ROOT/client/node_modules/.bin/buf" format -w "$TARGET_FILE" >/dev/null 2>&1) || true
      fi
      ;;
  esac
fi

# Contract requires empty JSON object on stdout
echo "{}"
exit 0
