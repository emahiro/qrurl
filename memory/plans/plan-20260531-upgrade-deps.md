# Plan: Go, Buf, and Client Library Upgrades (2026-05-31)

This plan outlines the steps to perform dependency upgrades and resolution of vulnerabilities for both the Go backend server and the React frontend client in separate Pull Requests.

## Status: Done

---

## 1. Goal

- **Server:**
  - `[x]` Upgrade Go version (to `1.25.0` automatically managed via toolchain to support modern library versions).
  - `[x]` Upgrade Go library dependencies in `server/go.mod` (Firestore, LINE Bot SDK, jwx, connect-go, API clients, protobuf, etc.).
  - `[x]` Upgrade `@bufbuild/buf` in `client/package.json` to get the latest compiler and protobuf configurations.
  - `[x]` Re-generate protobuf files via `buf generate` and confirm successful server build and test execution.

- **Client:**
  - `[x]` Fix security vulnerabilities via `npm audit fix` and `npm audit fix --force`.
  - `[x]` Upgrade key libraries (Vite to `8.0.14`, React, TypeScript, Prettier, ESLint).
  - `[x]` Verify that the client builds and the local development server starts without issues.

---

## 2. PR Separation Plan

To keep changes clean and easy to review, the work will be divided into two separate branches and PRs:

### PR 1: Server Upgrade (`upgrade-server-deps`)
- **Target Files:**
  - `server/go.mod`, `server/go.sum`
  - `client/package.json`, `client/package-lock.json` (just the `@bufbuild/buf` update)
  - `server/gen/` (regenerated Protobuf Go codes)
  - `client/gen/` (regenerated Protobuf TS codes)
- **Work Process:**
  1. `[x]` Create branch `upgrade-server-deps` from `main`.
  2. `[x]` Modify `server/go.mod` and trigger automatic Go toolchain upgrade to `1.25.0`.
  3. `[x]` Upgrade direct major dependencies individually to ensure maximum compatibility under Go 1.25.
  4. `[x]` Run `go mod tidy` to clean up.
  5. `[x]` Upgrade `@bufbuild/buf` inside `client/package.json` devDependencies.
  6. `[x]` Run `npm install` inside `client/` to apply buf package update.
  7. `[x]` Run `buf generate` from project root to regenerate Go and TS sources.
  8. `[x]` Run `go vet ./...` and `go test ./...` in `server/` to verify tests pass.
  9. `[x]` Commit and prepare for PR 1.

### PR 2: Client Upgrade (`upgrade-client-deps`)
- **Target Files:**
  - `client/package.json`, `client/package-lock.json`
- **Work Process:**
  1. `[x]` Create branch `upgrade-client-deps` from `main`.
  2. `[x]` Run `npm audit fix` inside the `client/` directory to fix all simple vulnerabilities.
  3. `[x]` Run `npm audit fix --force` to upgrade Vite (addressing high-risk esbuild/vite local dev dev-server vulnerabilities).
  4. `[x]` Verify updated Vite and nested dependencies resolved all 11 vulnerabilities (vulnerabilities count down to 0).
  5. `[x]` Run `npm run format` and `npm run lint` in `client/` to verify lint rules.
  6. `[x]` Run `npm run build` to confirm the production build completes successfully.
  7. `[x]` Start local development server (`npm run dev`) and test if it launches successfully.
  8. `[x]` Commit and prepare for PR 2.

---

## 3. Verification Plan

### Server Verification (PR 1):
- `[x]` `buf lint` (no errors)
- `[x]` `buf generate` (success)
- `[x]` `go vet ./...` (no compilation or static analysis issues)
- `[x]` `go test ./...` (all unit tests PASS)
- `[x]` `go build -o qrurl main.go` (build compiles successfully)

### Client Verification (PR 2):
- `[x]` `npm run format` (success)
- `[x]` `npm run lint` (no lint errors)
- `[x]` `npm run build` (production build compiles successfully)
- `[x]` Verify `npm run dev` launches and does not throw startup exceptions.

