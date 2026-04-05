# Pace Wallet - Backend API

> Go + Gin + PostgreSQL + Firebase Auth による REST API

このディレクトリは Pace Wallet のバックエンド API です。

## 🚀 本番環境

- **デプロイ先**: Railway
- **URL**: https://money-buddy-production.up.railway.app
- **Health Check**: https://money-buddy-production.up.railway.app/health

## 技術スタック

- **言語**: Go 1.25
- **フレームワーク**: Gin
- **データベース**: PostgreSQL (Neon Serverless)
- **データベースドライバー**: pgx/v5
- **クエリビルダー**: sqlc（型安全なSQL生成）
- **認証**: Firebase Admin SDK（JWT検証）
- **API仕様**: OpenAPI 3.0
- **テスト**: testify
- **コンテナ**: Docker（Alpine Linux）

## 📁 プロジェクト構成

```
backend/
├── cmd/server/           # エントリーポイント
├── db/
│   ├── schema/          # DDL（テーブル定義）
│   ├── query/           # SQLクエリ（sqlc用）
│   └── generated/       # sqlc自動生成コード
├── internal/
│   ├── domain/         # ドメインエンティティ（エンティティのみ、入力DTOは含まない）
│   ├── usecase/        # ユースケース層（1ユースケース1ファイル）
│   ├── handlers/       # HTTPハンドラ層（1エンドポイント1ファイル）
│   ├── middleware/     # 認証ミドルウェア
│   └── db/             # DB接続
├── infra/
│   ├── auth/           # Firebase認証初期化（インフラ依存のためinfra層）
│   ├── repository/     # リポジトリ実装（sqlc）
│   └── transaction/    # トランザクション管理
├── openapi/
│   └── openapi.yaml    # OpenAPI 3.0仕様
├── Dockerfile          # 本番環境用コンテナ
├── .dockerignore
├── .env.example        # 環境変数テンプレート
└── ENV_SETUP.md        # セットアップガイド
```

## 🏗️ アーキテクチャ設計

### 基本方針：Goイディオムとクリーンアーキテクチャ

本プロジェクトは **Goイディオム**（コンシューマー側インターフェース）と **クリーンアーキテクチャ** の原則に従って設計しています。

### コンシューマー側インターフェース

クリーンアーキテクチャの原則に従い、インターフェースはプロデューサー側（`infra/repository`）ではなく、**コンシューマー側**（`internal/usecase`、`internal/handlers`）で定義しています。
各ユースケースは自身が必要とする最小限のメソッドのみを持つインターフェースを定義します（インターフェース分離原則）。

```go
// internal/usecase/create_expense.go
// このusecaseが必要とする最小インターフェースをusecase側で定義
type createExpenseRepository interface {
    CreateExpense(userID string, input CreateExpenseInput) (domain.Expense, error)
}
```

### 1ユースケース1ファイル（ISP達成）

`internal/usecase/` 配下の各ファイルは1つのユースケースのみを担当します。
各ファイルには入力DTO・最小リポジトリインターフェース・`Execute` メソッドを定義しています。

```go
type CreateExpenseUseCase struct { ... }
func (uc *CreateExpenseUseCase) Execute(userID string, input CreateExpenseInput) (domain.Expense, error)
```

### ハンドラー層にもISPを適用

`internal/handlers/` の各ファイルも、必要とするユースケースの最小インターフェースを自身で定義します。

```go
// internal/handlers/create_expense.go
type createExpenseUseCase interface {
    Execute(userID string, input usecase.CreateExpenseInput) (domain.Expense, error)
}
```

### Goのイディオム：具体型を返すコンストラクタ

`infra/repository/` のコンストラクタはインターフェース型ではなく **具体型**（`*xxxRepositorySQLC`）を返します。
これにより、依存関係の逆転をコンシューマー側のインターフェース定義で実現しつつ、Goの慣習に従ったコードになっています。

### インフラ依存のコードはinfra層へ

Firebase認証の初期化など、インフラ依存のコードは `internal/` ではなく `infra/` 配下に配置しています。

## 🛠️ ローカル開発環境のセットアップ

### 前提条件

- Go 1.25以上
- PostgreSQL 14以上
- sqlc（SQL変更時のみ必要）

### 1. 依存関係のインストール

```bash
go mod download
```

### 2. データベースのセットアップ

```bash
# データベース作成
createdb money_buddy

# スキーマ適用
psql -d money_buddy -f db/schema/users.sql
psql -d money_buddy -f db/schema/categories.sql
psql -d money_buddy -f db/schema/fixed_costs.sql
psql -d money_buddy -f db/schema/expenses.sql
```

### 3. 環境変数の設定

`.env` ファイルを作成：

```bash
# データベース（Pooled Connection推奨）
DATABASE_DSN=host=localhost port=5432 user=postgres password=yourpassword dbname=money_buddy sslmode=disable

# Firebase認証（開発環境）
FIREBASE_CREDENTIALS_PATH=./firebase-admin-key.json

# CORS設定
ALLOWED_ORIGINS=http://localhost:3000

# サーバー設定
PORT=8080
ENV=development
```

### 4. Firebase Admin SDKの設定

1. [Firebase Console](https://console.firebase.google.com/) でサービスアカウント鍵を生成
2. `firebase-admin-key.json` として保存

### 5. サーバー起動

```bash
go run cmd/server/main.go
```

サーバーが起動したら http://localhost:8080/health でヘルスチェック可能です。

## 🧪 テストの実行

```bash
# 全テスト実行
go test ./...

# カバレッジ付き
go test -cover ./...

# 特定のパッケージ
go test ./internal/usecase/...
```

## 🔧 sqlcによるコード生成

**重要**: 生成済みコードは `db/generated/` にコミット済みです。SQL変更時のみ再生成が必要です。

```bash
# sqlcのインストール（初回のみ）
go install github.com/sqlcdev/sqlc/cmd/sqlc@v1.30.0

# コード生成
cd db
sqlc generate
```

生成後は必ず差分を確認してコミットしてください。

## 🐳 Dockerビルド

```bash
# ビルド
docker build -t pace-wallet-backend .

# ローカル実行
docker run -p 8080:8080 \
  -e DATABASE_DSN="host=host.docker.internal port=5432 user=postgres password=yourpassword dbname=money_buddy sslmode=disable" \
  -e FIREBASE_CREDENTIALS_JSON='{"type":"service_account",...}' \
  -e ALLOWED_ORIGINS="http://localhost:3000" \
  pace-wallet-backend
```

## 📚 API仕様

詳細は [openapi/openapi.yaml](openapi/openapi.yaml) を参照してください。

### 主要エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/health` | ヘルスチェック（認証不要） |
| POST | `/setup` | 初期設定 |
| GET | `/user/me` | ユーザー情報取得 |
| PUT | `/user/me` | ユーザー情報更新 |
| GET | `/dashboard` | ダッシュボードデータ |
| GET/POST/PUT/DELETE | `/expenses` | 支出管理 |
| GET | `/categories` | カテゴリ一覧 |
| GET/POST/PUT/DELETE | `/fixed-costs` | 固定費管理 |

**認証**: 全エンドポイント（`/health`以外）は`Authorization: Bearer <Firebase ID Token>`が必要です。

## 🔒 セキュリティ

- Firebase Admin SDKによるJWT検証
- ユーザーIDはトークンから取得（改ざん不可）
- CORS設定による不正アクセス防止
- ログマスキング（トークン情報非表示）
- 非rootユーザーでのコンテナ実行

## 📦 デプロイ

### Railway（本番環境）

1. GitHubリポジトリを連携
2. 環境変数を設定（`.env.example`参照）
3. 自動デプロイが実行される

必要な環境変数：
- `DATABASE_DSN`（Neon Pooled Connection）
- `FIREBASE_CREDENTIALS_JSON`
- `ALLOWED_ORIGINS`（Vercelの本番URL）
- `ENV=production`
- `PORT=8080`

詳細は [ENV_SETUP.md](ENV_SETUP.md) を参照してください。

## 🐛 トラブルシューティング

### データベース接続エラー

- Neon Pooled Connection (`-pooler`)を使用しているか確認
- `DATABASE_DSN`のフォーマットが正しいか確認
- コネクションプール設定を確認（`internal/db/db.go`）

### Firebase認証エラー

- `FIREBASE_CREDENTIALS_PATH`または`FIREBASE_CREDENTIALS_JSON`が正しく設定されているか確認
- Firebase Consoleでサービスアカウント鍵を再生成

### CORS エラー

- `ALLOWED_ORIGINS`にフロントエンドのURLが含まれているか確認
- `https://`を含めているか確認

## 📖 関連ドキュメント

- [ルートREADME](../README.md) - プロジェクト全体の概要
- [ENV_SETUP.md](ENV_SETUP.md) - 環境変数の詳細設定
- [openapi.yaml](openapi/openapi.yaml) - API仕様書

## 🤝 コントリビューション

このプロジェクトは個人開発中です。

## 📝 ライセンス

未定
go run cmd/server/main.go
```

5. API の例（curl）

リクエスト例（`spent_at` は `YYYY-MM-DD` または RFC3339 を受け付けます）:

```bash
curl -X POST http://localhost:8080/expenses \
	-H "Content-Type: application/json" \
	-d '{
		"amount": 1500,
		"category_id": 2,
		"memo": "食費",
		"spent_at": "2025-01-03"
	}'
```

成功時のレスポンスは `{"expense": {...}}` です。

---

## 更新 API 例（PUT /expenses/:id）

- 経路: `PUT /expenses/:id`
- 仕様:
	- `status` は `planned` または `confirmed` のみ有効
	- 遷移ルール: `confirmed` → `planned` は禁止、`planned` → `confirmed` は許可
	- `spent_at` は `YYYY-MM-DD` または RFC3339 を受け付けます

リクエスト例:

```bash
curl -X PUT http://localhost:8080/expenses/42 \
	-H "Content-Type: application/json" \
	-d '{
		"amount": 700,
		"category_id": 5,
		"memo": "updated",
		"spent_at": "2025-07-01",
		"status": "confirmed"
	}'
```

成功レスポンス（200）例:

```json
{
	"expense": {
		"id": 42,
		"amount": 700,
		"memo": "updated",
		"spent_at": "2025-07-01",
		"status": "confirmed",
		"category": { "id": 5 }
	}
}
```

エラーレスポンス例:
- バリデーションエラー（400）: `{ "error": "amount must be greater than 0" }`
- ステータス遷移エラー（409）: `{ "error": "invalid status transition" }`
- 内部エラー（500）: `{ "error": "internal server error" }`

---

## 一覧 API 例（GET /expenses）

- 経路: `GET /expenses`
- レスポンスは `expenses` 配列を含むオブジェクト

リクエスト例:

```bash
curl -X GET http://localhost:8080/expenses
```

成功レスポンス（200）例:

```json
{
	"expenses": [
		{
			"id": 1,
			"amount": 1500,
			"memo": "食費",
			"spent_at": "2025-01-03T00:00:00Z",
			"status": "confirmed",
			"category": { "id": 2, "name": "food" }
		}
	]
}
```

---

## 作成 API 例（POST /expenses）

- 経路: `POST /expenses`
- `status` は省略可能（省略時は `confirmed` が適用）。有効値は `planned`/`confirmed`

リクエスト例:

```bash
curl -X POST http://localhost:8080/expenses \
	-H "Content-Type: application/json" \
	-d '{
		"amount": 1500,
		"category_id": 2,
		"memo": "食費",
		"spent_at": "2025-01-03",
		"status": "planned"
	}'
```

成功レスポンス（201）例:

```json
{
	"expense": {
		"id": 1,
		"amount": 1500,
		"memo": "食費",
		"spent_at": "2025-01-03",
		"status": "planned",
		"category": { "id": 2 }
	}
}
```

エラーレスポンス例:
- バリデーションエラー（400）: `{ "error": "amount must be greater than 0" }`
- 内部エラー（500）: `{ "error": "internal server error" }`

---

## 削除 API 例（DELETE /expenses/:id）

- 経路: `DELETE /expenses/:id`
- 成功時はボディなしで `204 No Content`

リクエスト例:

```bash
curl -X DELETE http://localhost:8080/expenses/123
```

レスポンス例:
- 成功（204）: ボディなし
- ID不正（400）: `{ "error": "invalid expense ID" }`
- バリデーションエラー（400）: `{ "error": "cannot delete planned expense" }`
- 内部エラー（500）: `{ "error": "internal server error" }`

---

## CI の推奨ステップ（例: GitHub Actions）

ワークフロー内に必ず `sqlc generate`（または生成済みの検証）を含めてください。例:

```yaml
- name: Install sqlc
	run: go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.30.0

- name: Generate sqlc code
	working-directory: ./db
	run: sqlc generate

- name: Check generated code is up to date
	run: git diff --exit-code db/generated || (echo "Generated files out of date" && exit 1)

- name: Build
	run: go build ./...
```

---

## コード構成（重要なファイル）
- `cmd/server/main.go` - サーバー起点。全ユースケースを初期化し `handlers.Register` に渡します。
- `internal/db/db.go` - DB 接続（DSN 設定）
- `internal/domain/` - ドメインエンティティ（エンティティのみ、入力DTOは含まない）
- `internal/usecase/` - ユースケース層（1ユースケース1ファイル）。各ファイルが入力DTOと最小インターフェースを定義します。
- `internal/handlers/` - Gin ハンドラ（1エンドポイント1ファイル）。全ルーティングは `register.go` に集約。
- `infra/repository/` - リポジトリ実装（sqlc）。コンストラクタは具体型を返します。
- `db/` - SQL & `sqlc` 設定
	- `db/generated/` - sqlc による生成物（このリポジトリに含める）

---

## 運用ガイドライン
- SQL やスキーマ変更時: `db/expenses.sql` や `schema/` を更新し、`sqlc generate` を実行、`db/generated` をコミット。
- 生成物の差分は PR で確認すること。大きな差分が出た場合は sqlc のバージョン差の可能性を疑う。
- 生成物を更新する際は、他の開発者に通知するか PR に生成手順を含めてください。

---