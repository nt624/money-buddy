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
DATABASE_DSN=host=localhost port=5432 user=appuser password=password dbname=expense_db sslmode=disable

# CORS設定（開発環境）
ALLOWED_ORIGINS=http://localhost:3000

# Firebase Admin SDK
# ファイルパスまたはJSON文字列で設定
# 開発環境ではファイルパスを推奨
FIREBASE_CREDENTIALS_PATH=firebase-admin-key.json
# FIREBASE_CREDENTIALS_JSON={"type":"service_account",...}

# 環境
ENV=development
```

### 本番環境

本番環境では以下の環境変数を設定してください:

```env
# サーバー設定
PORT=8080

# データベース設定（本畮DBのDSN）
DATABASE_DSN=host=production-host port=5432 user=produser password=prodpass dbname=money_buddy sslmode=require

# CORS設定（本番フロントエンドのURL）
ALLOWED_ORIGINS=https://yourdomain.com

# Firebase Admin SDK
# ファイルパスまたはJSON文字列で設定（JSON文字列を優先）
# 本番環境ではセキュリティのためJSON文字列を推奨
FIREBASE_CREDENTIALS_PATH=/path/to/firebase-admin-key.json
# または
FIREBASE_CREDENTIALS_JSON={"type":"service_account","project_id":"your-project",...}

# 環境
ENV=production
```

## Firebase Admin SDKの設定

### 方法1: ファイルパス（開発環境推奨）

1. Firebase Consoleからサービスアカウントキーをダウンロード
2. `firebase-admin-key.json`として保存
3. `.gitignore`に追加されていることを確認（既に設定済み）
4. 環境変数 `FIREBASE_CREDENTIALS_PATH` にパスを設定

### 方法2: JSON文字列（本番環境推奨）

1. Firebase Consoleからサービスアカウントキーをダウンロード
2. JSONファイルの内容をそのまま環境変数 `FIREBASE_CREDENTIALS_JSON` に設定
3. クラウドサービスのシークレット管理機能を使用（AWS Secrets Manager、Cloud Secret Manager等）

**注意**: `FIREBASE_CREDENTIALS_JSON` が設定されている場合、`FIREBASE_CREDENTIALS_PATH` より優先されます

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
