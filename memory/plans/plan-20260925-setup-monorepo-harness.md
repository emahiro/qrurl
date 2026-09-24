# Plan: Setup Monorepo Lifecycle Hooks & AGENTS.md Harness (2026-09-25)

本計画は、`emahiro.dev#288` を参考に `qrurl` リポジトリ（Go server + React/TypeScript client のモノレポ構成）に自動検証ハーネス（Lifecycle Hooks）およびプロジェクト標準ガイドライン（`AGENTS.md`）を整備する手順をまとめたものです。

## Status: In Progress

---

## 1. 目的・背景

- **課題**:
  - `qrurl` は Go (server) と React/TS (client) のモノレポ構成をとっているが、ファイル編集後の自動検証・自動フォーマットの仕組みがない。
  - ルールファイルが古い `GEMINI.md`（短文のみ、かつ .gitignore 扱い）と `CLAUDE.md` に分かれており、最新の AI エージェント（Antigravity等）に対するプロジェクト全体のコンテキスト・制約が十分に集約されていない。
- **目標**:
  - `emahiro.dev#288` を踏襲し、`.agents/hooks.json` と `.agents/scripts/lint-hook.sh` を導入。
  - `server/` (Go) および `client/` (React/TS/Biome) の両方において、編集（`write_to_file` / `replace_file_content`）後に自動でフォーマット・静的解析を走らせる。
  - `GEMINI.md` から業界標準の `AGENTS.md` へ移行し、リポジトリ管理（Git tracked）とする。
  - 検証完了後、ブランチを作成してコミットし、PR を作成する。

---

## 2. アーキテクチャと詳細設計

### A. Lifecycle Hooks (`.agents/hooks.json` & `.agents/scripts/lint-hook.sh`)
- **トリガー**: `PostToolUse` (`write_to_file|replace_file_content`)
- **スクリプト入力**: JSON payload (`.toolCall.args.TargetFile`)
- **処理分岐**:
  - **Go (`server/` または `*.go`)**:
    - `gofmt -w <file>` で標準整形
    - `goimports -w -local github.com/emahiro/qrurl/server <file>`（利用可能な場合）
    - `cd server && go vet ./...` で型・静的解析チェック
  - **Client (`client/` または `*.ts|*.tsx|*.js|*.jsx|*.json|*.css`)**:
    - Biome CLI（PATH, npx, node_modules/.bin, ネイティブバイナリの順にフォールバック検出）
    - `biome check --write <file>` で自動修正・フォーマット実行
  - **Proto (`proto/` または `*.proto`)**:
    - 利用可能であれば `buf format -w <file>` を実行
- **スクリプト出力**: `{}` (exit 0)

### B. プロジェクトガイドライン (`AGENTS.md`)
- モノレポ全体の概要、サービスアーキテクチャ（LINE Bot + Web UI, Cloud Run + Firestore）
- Server / Client / Proto 各層の技術スタックと制約
- 標準開発コマンド
- タスク完了前の品質ゲート
- ライフサイクルフックと `work` スキルの運用方針

### C. クリーンアップ
- `GEMINI.md` の削除
- `.gitignore` の見直し（`GEMINI.md` のエントリ整理、`AGENTS.md` は追跡対象）

---

## 3. 作業タスク

- [x] 1. 設計計画の作成 (`memory/plans/plan-20260925-setup-monorepo-harness.md`)
- [x] 2. ライフサイクルフックの実装 (`.agents/hooks.json`, `.agents/scripts/lint-hook.sh`)
- [x] 3. フックスクリプトの実行権限付与 & 単体動作検証（Go / TS / 異常系）
- [x] 4. `AGENTS.md` の作成
- [x] 5. `GEMINI.md` の削除 & `.gitignore` の整理
- [x] 6. プロジェクト全体の最終検証（`go vet`, `go test`, `biome check`, `tsc`, `npm run build`）
- [ ] 7. Git ブランチ作成、コミット、PR 作成
