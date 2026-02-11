# Money Buddy - バックエンド

## 環境変数の設定

### 開発環境

1. `.env.example`を`.env`にコピー:
```bash
cp .env.example .env
```

2. `.env`ファイルを編集して、必要な値を設定:

```env
# サーバー設定
PORT=8080

# データベース設定
DATABASE_URL=postgres://user:password@localhost:5432/money_buddy?sslmode=disable

# CORS設定（開発環境）
ALLOWED_ORIGINS=http://localhost:3000

# Firebase Admin SDK
FIREBASE_CREDENTIALS_PATH=firebase-admin-key.json

# 環境
ENV=development
```

### 本番環境

本番環境では以下の環境変数を設定してください:

```env
# サーバー設定
PORT=8080

# データベース設定（本番DBのURL）
DATABASE_URL=postgres://user:password@production-host:5432/money_buddy?sslmode=require

# CORS設定（本番フロントエンドのURL）
ALLOWED_ORIGINS=https://yourdomain.com

# Firebase Admin SDK（JSON文字列またはファイルパス）
FIREBASE_CREDENTIALS_PATH=/path/to/firebase-admin-key.json

# 環境
ENV=production
```

## Firebase Admin SDKの設定

1. Firebase Consoleからサービスアカウントキーをダウンロード
2. `firebase-admin-key.json`として保存
3. `.gitignore`に追加されていることを確認（既に設定済み）

## 起動方法

### 開発環境
```bash
go run cmd/server/main.go
```

### 本番環境
```bash
# ビルド
go build -o bin/server cmd/server/main.go

# 実行
./bin/server
```

## セキュリティ注意事項

- `firebase-admin-key.json`は**絶対にGitにコミットしないでください**
- 本番環境では必ず`ENV=production`を設定してください
- CORSの`ALLOWED_ORIGINS`は信頼できるドメインのみを設定してください
