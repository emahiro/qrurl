# Harness evaluation report

- Run ID: `2026-11-03T13:26:21Z`
- Host: Linux
- Static: **11 pass** / **1 fail**
- Commands: **8 pass** / **0 fail**

## Static checks

- [pass] **doc:AGENTS.md**: AGENTS.md
- [pass] **agents-path:proto/**: proto/
- [pass] **agents-path:server/**: server/
- [pass] **agents-path:client/**: client/
- [pass] **agents-path:.agents/hooks.json**: .agents/hooks.json
- [pass] **agents-path:.agents/scripts/lint-hook.sh**: .agents/scripts/lint-hook.sh
- [pass] **agents-path:buf.yaml**: buf.yaml
- [pass] **agents-path:buf.gen.yaml**: buf.gen.yaml
- [pass] **agents-path:firebase.json**: firebase.json
- [pass] **json:agents-hooks**: .agents/hooks.json
- [fail] **hook:lint-path**: hooks.json uses ./scripts/lint-hook.sh but script is at .agents/scripts/lint-hook.sh
- [pass] **secrets:harness**: no obvious secret patterns in .agents/.cursor

## Command checks

- [pass] **go-fix** (331 ms): `bash -lc cd "/workspace/server" && go fix ./...`
- [pass] **go-vet** (317 ms): `bash -lc cd "/workspace/server" && go vet ./...`
- [pass] **go-test** (320 ms): `bash -lc cd "/workspace/server" && go test ./...`
- [pass] **go-build** (276 ms): `bash -lc cd "/workspace/server" && go build -o /tmp/qrurl-harness main.go`
- [pass] **client-lint** (263 ms): `bash -lc cd "/workspace/client" && npm run lint`
- [pass] **client-tsc** (437 ms): `bash -lc cd "/workspace/client" && npx tsc --noEmit`
- [pass] **client-build** (720 ms): `bash -lc cd "/workspace/client" && npm run build`
- [pass] **buf-lint** (263 ms): `bash -lc cd "/workspace" && "/workspace/client/node_modules/.bin/buf" lint`

## Baseline comparison

No baseline.json yet. Run `./harness-eval/run.sh --write-baseline` after a green run to seed comparison.
