# Spec: Dependency Upgrades & Vulnerability Fix Verification (2026-05-31)

This spec defines the operational requirements and validation criteria to confirm successful upgrades of Server (Go/Buf) and Client (NPM) dependencies.

## 1. Goal

Upgrade system libraries while maintaining strict backwards compatibility and operational readiness.

---

## 2. Server (Go + Buf) Requirements

### A. Environment & Tools
- Go version must be locked to at least `1.24` in `go.mod`.
- Buf toolchain in `client/package.json` must be upgraded to modern version (`^1.70.0`).

### B. Compilation & Linter Contract
- `buf lint` must report zero linter issues in Protobuf definitions.
- `buf generate` must run without syntax errors and successfully generate files inside `server/gen/` and `client/gen/`.
- `go vet ./...` must compile clean.

### C. Runtime Contract
- `go test ./...` inside `server/` must return `PASS` for all existing unit tests.
- `go build` must output a working binary for Cloud Run deployment without errors.

---

## 3. Client (React/TypeScript) Requirements

### A. Security Audit Contract
- `npm audit` should report 0 high/critical vulnerabilities after fixes. Moderate vulnerabilities must be resolved where safe.

### B. Library Version Targets
- `@bufbuild/buf` upgraded to `^1.70.0`.
- Vite must be upgraded to resolve the development server security flaw (target version `^8.0.0` or latest available major recommended by npm audit).
- Dependencies must build cleanly without TypeScript transpilation or bundler failures.

### C. Build & Lint Contract
- `npm run format` runs formatting successfully.
- `npm run lint` finishes with `0` errors.
- `npm run build` completes successfully and produces optimized assets in the `client/dist/` directory.
- `npm run dev` boots up a development HTTP server listening on standard ports without script exceptions.

---

## 4. Error & Rollback Contingencies

- If package upgrades introduce compilation errors due to major breaking changes, we will:
  1. Identify the breaking API change (e.g., deprecated function calls).
  2. Pin to the highest non-breaking minor/major version if rewriting the code introduces excessive regressions.
  3. Ensure no local changes are left uncommitted on failure.
