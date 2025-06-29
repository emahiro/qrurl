# qrurl

QR コードから URL を取り出すサービスです。
LINE の Bot と Web Browser に対応しています。
(Web ブラウザは開発中)

![qrurl の動作イメージ@LINE](imgs/スクリーンショット%202023-07-13%201.36.46.png)

友だち追加はこちら

![qrurl を友だち追加する](imgs/add-friend.png)

# Stacks

## Server

- Go
- Cloud Run
- LINE Messageing Platform

## Client

- React
- TypeScript
- Vite
- Tailwind CSS

# Development

## Server の起動

Server では Makefile を使用して環境変数を読み込みながら起動できます。

### 前提条件
- `.env.yaml` ファイルが `server/` ディレクトリに必要です
- 環境変数は YAML 形式で定義されている必要があります

### 起動方法

```bash
cd server
make run
```

### 手動起動（環境変数なし）

```bash
cd server
go run main.go
```

Makefile の `run` コマンドは `.env.yaml` から環境変数を読み込んで Go アプリケーションを起動します。

## Client の起動

```bash
cd client
npm install
npm run dev
```

# Contributer

@emahiro

