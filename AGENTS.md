# AGENTS.md - qrurl 開発ガイドライン

本リポジトリは **QR コード画像から URL を抽出する Web / LINE Bot サービス (`qrurl`)** です。
Go 製バックエンドサーバー（Cloud Run）と React / TypeScript 製フロントエンド（Vite）が同居する Monorepo 構成をとっています。

---

## 1. 主要技術スタック & アーキテクチャ

| 領域 | 技術 | 備考 |
|---|---|---|
| **Monorepo** | Go server + React/TypeScript client | 単一リポジトリでサーバー・クライアント・プロトコル定義を管理 |
| **Backend (`server/`)** | Go 1.24+, Connect-Go, Google Cloud Firestore, Google Cloud Run | LINE Messaging API, ZXing (QRデコード), Google Cloud Logging |
| **Frontend (`client/`)** | React 19, TypeScript 7, Vite 8, Tailwind CSS 4, Connect-Web | Firebase Hosting / 静的配信 |
| **API / Protocol (`proto/`)** | Protocol Buffers v3, Buf | Connect RPC 経由での通信 |
| **Database** | Google Cloud Firestore | トークン管理・セッション情報等の永続化 |
| **Lint & Format** | Server: `gofmt`, `goimports`, `go vet`<br>Client: Biome 2.5 (インデント: スペース2, クォート: ダブル) | Lifecycle Hooks (`.agents/hooks.json`) により自動実行 |
| **CI / CD** | GitHub Actions (`go.yml`, `nextjs.yml`, `deploy.yml`, `deploy-client.yml`) | Go vet/test, Biome, tsc, Vite build |

### ディレクトリ構成
```
/
├── proto/                    # Protocol Buffers 定義
│   ├── qrurl/v1/             # QRコード解析サービス定義
│   └── ping/v1/              # ヘルスチェックサービス定義
├── server/                   # Go バックエンドサーバー (Cloud Run)
│   ├── gen/                  # Buf により生成された Go コード
│   ├── handler/              # LINE Webhook ハンドラー
│   ├── infra/firestore/      # Firestore クライアント
│   ├── intercepter/          # Connect RPC インターセプター
│   ├── lib/                  # 共通ライブラリ (QRデコード, JWT, LINE, ロガー)
│   ├── middleware/           # HTTP ミドルウェア (CORS, ログ)
│   ├── model/                # ドメインモデル
│   ├── repository/           # リポジトリ層
│   ├── service/              # Connect RPC サービス実装
│   ├── main.go               # サーバーエントリーポイント
│   └── go.mod                # Go モジュール設定
├── client/                   # React / TypeScript フロントエンド
│   ├── gen/                  # Buf により生成された TypeScript コード
│   ├── src/                  # React アプリケーションコード (App.tsx, libs 等)
│   ├── biome.json            # Biome 設定 (スペース2, ダブルクォート)
│   ├── package.json          # クライアント依存設定
│   ├── tsconfig.json         # TypeScript 設定
│   └── vite.config.ts        # Vite 設定
├── .agents/
│   ├── hooks.json            # Antigravity ライフサイクルフック設定
│   └── scripts/lint-hook.sh  # ファイル編集後自動検証スクリプト (Go & Client)
├── memory/plans/             # 開発計画・設計ログ
├── buf.gen.yaml              # Buf コード生成設定
├── buf.yaml                  # Buf モジュール設定
└── firebase.json             # Firebase Hosting 設定
```

---

## 2. 最重要制約 (Critical Constraints)

エージェントは以下の制約を厳格に遵守すること：

1. **Monorepo の責務境界の維持**:
   - `server/` と `client/` はそれぞれ独立した依存管理（Go modules / npm）を行っています。
   - 両者間のインターフェース変更は必ず `proto/` を更新し、`buf generate` を通じて行います。
2. **基盤設定ファイルの保護**:
   - `client/biome.json`, `client/tsconfig*.json`, `server/go.mod`, `server/go.sum`, `buf.yaml`, `buf.gen.yaml`, `.github/workflows/*` はユーザーの明示的な指示なく変更しない。
3. **Go 実装原則**:
   - 標準の Go イディオムに従い、`error` の適切なハンドリングと構造化ログ（`lib/log`）を活用する。
   - 外部依存の注入（DI）により、テスタビリティを確保する。
4. **Client 実装原則**:
   - Biome の規約（スペース2、ダブルクォート、セミコロンあり）を遵守する。
   - Connect-Web による型安全な RPC クライアント通信を維持する。
5. **指示追従と最小限の実装**:
   - 要求されたスコープ外の過剰な抽象化や不要な機能追加は避ける（YAGNI, KISS）。

---

## 3. 開発コマンド

### Server (Go)
```bash
cd server
go run main.go               # サーバー起動 (要環境変数)
go vet ./...                 # 静的解析
go test ./...                # 単体テスト実行
go build -o qrurl main.go    # バイナリビルド
```

### Client (React / TypeScript)
```bash
cd client
npm run dev                  # Vite 開発サーバー起動
npm run lint                 # Biome lint (src/)
npm run format               # Biome フォーマット修正
npm run format:check         # Biome フォーマット検査
npm run check                # Biome lint & format 修正
npx tsc --noEmit             # TypeScript 型チェック
npm run build                # 型チェック & Vite 本番ビルド
```

### Protocol Buffers
```bash
# プロジェクトルートから実行
buf lint                     # Proto 定義の lint
buf generate                 # Go / TypeScript コードの再生成
```

---

## 4. 品質ゲート (タスク完了前)

※ ファイル編集時のフォーマットおよび検証は `.agents/hooks.json` の Hook により自動実行されます。
完了報告前に以下の全項目がエラーゼロであることを確認すること:

1. **Server (Go)**:
   - `cd server && go vet ./...` (静的解析エラーゼロ)
   - `cd server && go test ./...` (全テスト PASS)
2. **Client (React / TypeScript)**:
   - `cd client && npm run lint` (Biome エラーゼロ)
   - `cd client && npx tsc --noEmit` (TypeScript 型エラーゼロ)
   - `cd client && npm run build` (ビルド成功)
3. **Proto (変更時のみ)**:
   - `buf lint` (エラーゼロ)
   - `buf generate` 後の生成コードが最新であること

---

## 5. ワークフロー & スキルの活用

- **自律開発ワークフロー (`work`)**:
  - タスクの進行は `work` スキルを活用し、調査 → 計画 → TDD実装 → 検証 → コミットを自律的に進行する。
- **計画書・設計ログの記録**:
  - 中〜大規模な変更やアーキテクチャ見直しの際は、`memory/plans/plan-YYYYMMDD-{slug}.md` に計画書を作成・追跡する。
- **Lifecycle Hooks**:
  - `.agents/hooks.json` により、ファイル編集（`write_to_file`, `replace_file_content`）後に `.agents/scripts/lint-hook.sh` が自動起動し、Go および Client の自動フォーマット・検証を実行する。
