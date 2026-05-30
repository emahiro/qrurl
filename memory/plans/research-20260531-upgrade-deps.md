# Research: Go, Buf, and Client Library Upgrades (2026-05-31)

This research document analyzes the current dependencies and drafts the upgrade path for the `qrurl` project's Server (Go + Buf) and Client (React + TypeScript).

## 1. Current Dependency Analysis

### A. Server (Go)
- **Go Version:** `1.24` / `toolchain go1.24.0`
- **Buf CLI:** `1.36.0` (Global), `1.25.0` (Local npm devDependency)
- **Key Go Dependencies (go.mod):**
  - `cloud.google.com/go/firestore` v1.11.0 -> Upgrade available: v1.16.0
  - `firebase.google.com/go` v3.13.0+incompatible -> Stable (Incompatible)
  - `github.com/bufbuild/connect-go` v1.9.0 -> Upgrade available: v1.10.0 (Connect-Go is deprecated, transitioned to `@connectrpc/connect`)
  - `github.com/lestrrat-go/jwx/v2` v2.0.11 -> Upgrade available: v2.1.6
  - `github.com/line/line-bot-sdk-go/v7` v7.20.0 -> Upgrade available: v7.21.0
  - `google.golang.org/api` v0.131.0 -> Upgrade available: v0.282.0
  - `google.golang.org/protobuf` v1.31.0 -> Upgrade available: v1.36.11

### B. Client (React/TypeScript)
- **Vite:** ^4.4.0 -> Upgrade available: v8.0.14
- **React:** ^18.2.0 -> Upgrade available: v19.2.6
- **TypeScript:** ^5.0.2 -> Upgrade available: v5.9.3 / v6.0.3
- **NPM Audit Vulnerabilities:**
  - 11 vulnerabilities (7 moderate, 4 high)
  - Major vulnerable packages: `esbuild <=0.24.2`, `vite <=6.4.1`, `postcss <8.5.10`, `rollup 3.x`
  - High risk vulnerability in `esbuild` / `vite` (development server access security issue).

---

## 2. Upgrade Risks & Strategies

### Risk 1: Buf Version Upgrades (Local vs Global)
- **Risk:** Discrepancy between global `buf` (1.36.0) and local devDependency `buf` (1.25.0) can lead to different generation outputs.
- **Strategy:** Update `@bufbuild/buf` in `client/package.json` to the latest `1.70.0`. This ensures buf generation matches modern Protobuf specs and plugins.

### Risk 2: connect-go deprecation
- **Risk:** `github.com/bufbuild/connect-go` is deprecated. Migrating to `connectrpc.com/connect` is a major breaking change requiring updates to imports in all Go files.
- **Strategy:**
  - Option A: Just upgrade `connect-go` to the latest `v1.10.0` (deprecated but compatible).
  - Option B: Migrate to `connectrpc.com/connect` and `connectrpc.com/cors` to align with the active community-supported packages.
  - *Recommendation:* We should first attempt Option A or ask the user, but standard library upgrades usually imply upgrading to the latest supported major version. We will outline both in the implementation plan.

### Risk 3: Vite Major Upgrade (v4 -> v8) & React v19 Upgrade
- **Risk:** Upgrading Vite from v4 to v8 and React from v18 to v19 are major updates.
- **Strategy:**
  - Standard `npm audit fix` will solve smaller dependency issues.
  - `npm audit fix --force` will upgrade Vite.
  - Run `npm run build` and `npm run lint` inside the client directory after upgrades to verify TS types, ESLint rules, and build configurations.

---

## 3. Scope of Work (PRs Separation)

As requested, we will split the work into two distinct Pull Requests (branches):

### PR 1: Server Dependency Upgrades (`upgrade-server-deps`)
- Upgrade Go toolchain to `1.24.4` (currently local `go version` is `1.24.4`).
- Upgrade Go dependencies in `server/go.mod` (Firestore, LINE Bot SDK, jwx, connect-go/connectrpc, protobuf).
- Upgrade `@bufbuild/buf` package in `client/package.json` devDependencies (since Buf is used to compile Go and TS definitions).
- Run `buf generate` and verify backend builds and tests (`go test ./...`).

### PR 2: Client Dependency Upgrades (`upgrade-client-deps`)
- Run `npm audit fix` inside the `client` directory.
- Perform necessary `--force` fixes to upgrade Vite.
- Upgrade other frontend library versions (React, TS, tailwindcss if applicable).
- Verify frontend runs successfully (`npm run dev`) and builds correctly (`npm run build`).

---

## 4. Verification Flow

### Server:
- `buf lint`
- `buf generate`
- `go vet ./...` (from `/server`)
- `go test ./...` (from `/server`)
- `go build` (from `/server`)

### Client:
- `npm run format` (from `/client`)
- `npm run lint` (from `/client`)
- `npm run build` (from `/client`)
- Check Vite dev server local launching (`npm run dev`)
