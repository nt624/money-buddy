# Money Buddy

## プロジェクト概要

**Money Buddy** は、「お金が貯まらない人が"貯まる生活"を続けられるように伴走する家計管理アプリ」です。

### コンセプト

「今の自分のお金の状況が一目でわかるアプリ」をコンセプトに、シンプルな入力とリアルタイムの残額表示により、貯まる行動を自然に習慣化することを目指しています。

### ターゲットユーザー

- お金の現状把握が苦手な方
- 計画立てが苦手な方
- 家計簿が続かない方
- 今いくら使えるかがすぐに分からない方

### 主な価値提供

アプリを開くたびに「あといくら使っていいか」が一目でわかり、無意識のうちに"貯まる生活"へ変わることを目指しています。

---

## 主な機能

### 実装済み機能 ✅

#### 1. 初期設定機能
- ユーザー登録（月収、貯金目標額の設定）
- 固定費の登録・管理
- 自由に使える変動費の自動計算（変動費 = 収入 - 固定費 - 貯金額）

#### 2. ダッシュボード機能
- **リアルタイム残額表示**
  - 今月の自由に使える残額を大きく表示
  - 残額計算ルール：`残額 = 変動費 - (確定支出 + 予定支出)`
  - 残額に応じた色分け表示（緑: 70%以上 / 黄: 30-70% / 赤: 30%未満）
- **月次サマリー**
  - 収入、貯金目標、固定費、変動費の一覧表示
  - 確定支出・予定支出の集計表示
- **レスポンシブデザイン**
  - モバイル、タブレット、デスクトップに最適化されたレイアウト

#### 3. 支出管理機能
- **支出の登録**
  - 金額、日付、カテゴリ、メモの入力
  - 支出ステータスの管理（予定 / 確定）
  - 予定支出：旅行・飲み会・イベントなど、これから発生する概算支出を事前登録
  - 確定支出：実際に使った金額
- **支出の更新**
  - 予定支出から確定支出への更新
  - 金額・カテゴリ・メモの編集
  - ステータス変更ルール：予定 → 確定は許可、確定 → 予定は禁止
  - カード型トグルUIで直感的なステータス切り替え
- **支出の削除**
- **支出一覧の表示**
  - ステータス別（確定/予定）の色分け表示
  - 編集・削除機能への簡単なアクセス

#### 4. カテゴリ管理
- デフォルトカテゴリの提供（食費、交通費、娯楽費など）
- カテゴリ一覧の取得

#### 5. UI/UX機能
- **ダークモード**
  - ライトモード/ダークモードの切り替え
  - システム設定への対応
  - ユーザー設定の永続化（localStorage）
- **デザインシステム**
  - Tailwind CSS 4によるモダンなデザイン
  - セマンティックカラートークンによる一貫したテーマ
  - アクセシビリティを考慮した色彩設計（WCAG AA準拠）

### 今後実装予定（MVP後）📋

- カテゴリ別支出サマリー
- 月別の比較グラフ・トレンド分析
- 無駄遣いアラート・通知機能
- 収支レポート（PDF生成など）
- カスタムカテゴリの作成・編集
- 固定費の編集・削除機能
- 今日までの利用ペースとの比較機能
- 認証機能（Firebase Auth）の実装
- モバイルアプリ（Expo / React Native）開発

---

## スクリーンショット

### 初期設定画面
![初期設定画面](docs/images/initial-setup.png)
_初回起動時に月収・貯金目標・固定費を設定_

### ダッシュボード（ライトモード）
![ダッシュボード - ライトモード](docs/images/dashboard.png)
_今月の残額と月次サマリーを一目で確認_

### ダッシュボード（ダークモード）
![ダッシュボード - ダークモード](docs/images/dashboard_dark.png)
_目に優しいダークモード対応_

### 支出入力・一覧画面
![支出入力・一覧画面](docs/images/expense-input.png)
_シンプルな入力フォームと見やすい支出一覧_

### モバイル表示
![モバイル表示](docs/images/dashboard_mobile.png)

_スマートフォンでも快適に操作可能_

---

## 使い方

### 1. 初期設定
初回起動時に以下の情報を入力します：
- **月収（手取り）**: 毎月の収入額
- **貯金目標額**: 毎月貯めたい金額
- **固定費**: 家賃、光熱費、通信費など毎月必ず発生する支出

これらを入力すると、自動的に**変動費（自由に使えるお金）**が計算されます。

### 2. ダッシュボードで残額確認
ホーム画面では以下の情報が表示されます：
- **今月あといくら使える？**: 大きく表示される残額（色で危険度を判断）
  - 🟢 緑: 余裕あり（70%以上残っている）
  - 🟡 黄: 注意（30-70%残っている）
  - 🔴 赤: 危険（30%未満、またはマイナス）
- **月次サマリー**: 収入、貯金目標、固定費、変動費
- **今月の支出**: 確定済みと予定の支出額

### 3. 支出の登録
「支出を追加」から以下を入力：
- **金額**: 使った（または使う予定の）金額
- **カテゴリ**: 食費、交通費、娯楽費など
- **日付**: 支出の日付
- **メモ**: 任意で詳細を記録
- **ステータス**: 確定（実際に使った）or 予定（これから使う）

**予定支出の活用例：**
- 来週の飲み会（3,000円）
- 月末の旅行（30,000円）
- 友人の誕生日プレゼント（5,000円）

事前に予定支出を登録しておくことで、「今使えるお金」がより正確に把握できます。

### 4. 支出の管理
- **編集**: 予定支出を確定に変更したり、金額を修正
- **削除**: 不要な支出を削除
- **一覧表示**: 全ての支出を確認、ステータス別に色分け表示

---

## 技術的な特徴

### アーキテクチャ
- **フロントエンド**: Next.js 16 App Router + React 19による高速なSPA
- **バックエンド**: Go + Ginによる高パフォーマンスなREST API
- **データベース**: PostgreSQLでリレーショナルデータを安全に管理
- **型安全性**: TypeScript（フロントエンド）+ sqlc（バックエンド）で完全な型安全を実現

### デザインシステム
- **Tailwind CSS 4**: 最新のユーティリティファーストCSS
- **カラートークン**: セマンティックな色定義で一貫性を確保
  - Primary（プライマリーアクション）
  - Secondary（セカンダリーアクション）
  - Success（成功・緑）
  - Warning（警告・黄）
  - Danger（危険・赤）
- **レスポンシブデザイン**: モバイルファースト設計
  - Mobile: ~640px
  - Tablet: 640px~1024px
  - Desktop: 1024px~

### UX設計の工夫
- **直感的な残額表示**: 数字だけでなく色でも危険度を伝達
- **予定支出機能**: 将来の支出を見越した残額計算
- **ダークモード**: 長時間の利用でも目に優しい
- **カード型トグルUI**: タップしやすい大きなボタンで操作性向上
- **リアルタイム更新**: 支出の追加・編集後すぐにダッシュボードに反映

---

## 技術スタック

### フロントエンド（Web）
- **フレームワーク**: Next.js 16 (App Router)
- **ライブラリ**: React 19
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS 4
- **Linter**: ESLint
- **テスト**: Jest

### バックエンド（API）
- **言語**: Go 1.25
- **Webフレームワーク**: Gin
- **データベース**: PostgreSQL
- **クエリビルダー**: sqlc
- **テストフレームワーク**: testify
- **API仕様**: OpenAPI 3.0

### データベース
- PostgreSQL
- 想定環境: Supabase / Railway / Render / Neon など

### 将来的な拡張
- **認証**: Firebase Auth
- **モバイル**: Expo (React Native)
- **インフラ**: AWS (ECS / RDS / S3)

---

## ディレクトリ構造

```
.
├── backend/                    # バックエンド（Go API）
│   ├── cmd/
│   │   └── server/            # エントリーポイント
│   │       └── main.go        # サーバー起動
│   ├── db/
│   │   ├── schema/            # データベーススキーマ定義
│   │   │   ├── users.sql
│   │   │   ├── fixed_costs.sql
│   │   │   ├── expenses.sql
│   │   │   └── categories.sql
│   │   ├── query/             # SQLクエリ定義（sqlc用）
│   │   │   ├── users.sql
│   │   │   ├── fixed_costs.sql
│   │   │   ├── expenses.sql
│   │   │   ├── categories.sql
│   │   │   └── dashboard.sql
│   │   ├── generated/         # sqlcによる自動生成コード
│   │   └── sqlc.yaml          # sqlc設定ファイル
│   ├── internal/
│   │   ├── db/                # データベース接続
│   │   ├── handlers/          # HTTPハンドラ（APIエンドポイント）
│   │   ├── models/            # ドメインモデル
│   │   ├── repositories/      # リポジトリインターフェース
│   │   └── services/          # ビジネスロジック
│   ├── infra/
│   │   ├── repository/        # リポジトリ実装
│   │   └── transaction/       # トランザクション管理
│   ├── openapi/
│   │   └── openapi.yaml       # API仕様書
│   ├── go.mod
│   └── go.sum
│
└── frontend/                   # フロントエンド（Next.js）
    ├── src/
    │   ├── app/               # Next.js App Router
    │   │   ├── page.tsx       # ホーム画面
    │   │   ├── layout.tsx     # レイアウト
    │   │   └── globals.css    # グローバルスタイル
    │   ├── components/        # Reactコンポーネント
    │   │   ├── ExpenseForm.tsx      # 支出入力フォーム
    │   │   └── ExpenseList.tsx      # 支出一覧表示
    │   ├── hooks/             # カスタムフック
    │   │   └── useExpenses.ts       # 支出管理フック
    │   └── lib/               # ユーティリティ・API
    │       ├── api/           # APIクライアント
    │       └── types/         # TypeScript型定義
    ├── public/                # 静的ファイル
    ├── package.json
    ├── tsconfig.json
    ├── next.config.ts
    └── tailwind.config.ts
```

---

## セットアップ / 実行方法

### 前提条件

- **Go**: 1.25以上
- **Node.js**: 20以上
- **PostgreSQL**: 14以上
- **sqlc**: sqlcのインストール（バックエンド開発に必要）

### データベースセットアップ

1. PostgreSQLデータベースを作成します：

```bash
createdb money_buddy
```

2. スキーマを適用します：

```bash
cd backend
psql -d money_buddy -f db/schema/users.sql
psql -d money_buddy -f db/schema/categories.sql
psql -d money_buddy -f db/schema/fixed_costs.sql
psql -d money_buddy -f db/schema/expenses.sql
```

3. 環境変数を設定します（必要に応じて `.env` ファイルを作成）：

```bash
export DATABASE_URL="postgres://ユーザー名:パスワード@localhost:5432/money_buddy?sslmode=disable"
```

### バックエンドの起動

```bash
cd backend

# 依存関係のインストール
go mod download

# サーバー起動（デフォルトポート: 8080）
go run cmd/server/main.go
```

### フロントエンドの起動

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバー起動（デフォルトポート: 3000）
npm run dev
```

ブラウザで `http://localhost:3000` にアクセスします。

---

## API仕様

API仕様は OpenAPI 3.0 形式で定義されています。詳細は [backend/openapi/openapi.yaml](backend/openapi/openapi.yaml) をご参照ください。

### ベースURL
```
http://localhost:8080
```

### エンドポイント一覧

#### 支出管理 (Expenses)
| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| POST | `/expenses` | 支出の登録 |
| GET | `/expenses` | 支出一覧の取得 |
| PUT | `/expenses/:id` | 支出の更新 |
| DELETE | `/expenses/:id` | 支出の削除 |

#### カテゴリ管理 (Categories)
| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| GET | `/categories` | カテゴリ一覧の取得 |

#### ユーザー管理 (Users)
| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| GET | `/user/me` | 現在のユーザー情報の取得 |

#### 初期設定 (Setup)
| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| POST | `/setup` | 初期設定の完了（ユーザー作成 + 固定費登録） |

#### ダッシュボード (Dashboard)
| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| GET | `/dashboard` | ダッシュボードデータの取得 |

### 主要なリクエスト/レスポンス例

#### 支出の登録 (POST /expenses)
**リクエスト:**
```json
{
  "amount": 3000,
  "category_id": 1,
  "memo": "ランチ",
  "spent_at": "2026-02-11",
  "status": "confirmed"
}
```

**レスポンス (201 Created):**
```json
{
  "expense": {
    "id": 1,
    "amount": 3000,
    "category": {
      "id": 1,
      "name": "食費"
    },
    "memo": "ランチ",
    "spent_at": "2026-02-11T00:00:00Z",
    "status": "confirmed"
  }
}
```

#### 支出一覧の取得 (GET /expenses)
**レスポンス (200 OK):**
```json
{
  "expenses": [
    {
      "id": 1,
      "amount": 3000,
      "category": {
        "id": 1,
        "name": "食費"
      },
      "memo": "ランチ",
      "spent_at": "2026-02-11T00:00:00Z",
      "status": "confirmed"
    }
  ]
}
```

#### ダッシュボードデータの取得 (GET /dashboard)
**レスポンス (200 OK):**
```json
{
  "income": 300000,
  "saving_goal": 50000,
  "fixed_costs": 150000,
  "variable_budget": 100000,
  "confirmed_expenses": 30000,
  "planned_expenses": 10000,
  "remaining": 60000
}
```

**フィールド説明:**
- `income`: 月収（手取り）
- `saving_goal`: 貯金目標額
- `fixed_costs`: 固定費の合計
- `variable_budget`: 変動費（自由に使えるお金）= income - saving_goal - fixed_costs
- `confirmed_expenses`: 確定済み支出の合計
- `planned_expenses`: 予定支出の合計
- `remaining`: 残額 = variable_budget - confirmed_expenses - planned_expenses

#### 初期設定 (POST /setup)
**リクエスト:**
```json
{
  "user": {
    "id": "user123",
    "income": 300000,
    "saving_goal": 50000
  },
  "fixed_costs": [
    {
      "name": "家賃",
      "amount": 80000
    },
    {
      "name": "光熱費",
      "amount": 15000
    }
  ]
}
```

**レスポンス (200 OK):**
```json
{
  "user": {
    "id": "user123",
    "income": 300000,
    "saving_goal": 50000,
    "created_at": "2026-02-11T10:00:00Z",
    "updated_at": "2026-02-11T10:00:00Z"
  },
  "fixed_costs": [
    {
      "id": 1,
      "user_id": "user123",
      "name": "家賃",
      "amount": 80000,
      "created_at": "2026-02-11T10:00:00Z",
      "updated_at": "2026-02-11T10:00:00Z"
    },
    {
      "id": 2,
      "user_id": "user123",
      "name": "光熱費",
      "amount": 15000,
      "created_at": "2026-02-11T10:00:00Z",
      "updated_at": "2026-02-11T10:00:00Z"
    }
  ]
}
```

### ビジネスルール

#### 支出ステータスの遷移ルール
- **許可**: `planned` (予定) → `confirmed` (確定)
- **禁止**: `confirmed` (確定) → `planned` (予定)

確定済みの支出を予定に戻すことはできません。これは実際に使ったお金を「使っていないことにする」ことを防ぐためです。

### エラーレスポンス

全てのエラーは以下の形式で返されます：

```json
{
  "error": "エラーメッセージ"
}
```

**HTTPステータスコード:**
- `400 Bad Request`: バリデーションエラー、無効なステータス遷移
- `404 Not Found`: リソースが見つからない（ユーザー未登録など）
- `422 Unprocessable Entity`: ビジネスロジックエラー
- `500 Internal Server Error`: サーバー内部エラー

---

## データモデル

### ER図

```mermaid
erDiagram
    Users ||--o{ FixedCosts : "has"
    Users ||--o{ Expenses : "has"
    Categories ||--o{ Expenses : "categorizes"

    Users {
        TEXT id PK "Firebase UID"
        INT income "月収（手取り）"
        INT saving_goal "月の貯金目標額"
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    FixedCosts {
        SERIAL id PK
        TEXT user_id FK "ユーザーID"
        TEXT name "固定費名"
        INT amount "金額"
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    Expenses {
        SERIAL id PK
        TEXT user_id FK "ユーザーID"
        INT amount "金額"
        INT category_id FK "カテゴリID"
        DATE spent_at "予定日/実施日"
        TEXT memo "メモ"
        TEXT status "planned/confirmed"
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    Categories {
        SERIAL id PK
        TEXT name "カテゴリ名"
        TIMESTAMP created_at
    }
```

### テーブル詳細

### Users（ユーザー）
| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | TEXT | Firebase UID（主キー） |
| income | INT | 月収（手取り） |
| saving_goal | INT | 月の貯金目標額 |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### FixedCosts（固定費）
| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | SERIAL | 主キー |
| user_id | TEXT | ユーザーID（外部キー） |
| name | TEXT | 固定費名（例: 家賃、光熱費） |
| amount | INT | 金額 |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### Expenses（支出）
| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | SERIAL | 主キー |
| user_id | TEXT | ユーザーID（外部キー） |
| amount | INT | 金額（概算 or 実額） |
| category_id | INT | カテゴリID |
| spent_at | DATE | 予定日 or 実施日 |
| memo | TEXT | メモ（任意） |
| status | TEXT | ステータス（planned / confirmed） |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### Categories（カテゴリ）
| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | SERIAL | 主キー |
| name | TEXT | カテゴリ名 |
| created_at | TIMESTAMP | 作成日時 |

---

## 開発

### バックエンド

#### テストの実行

```bash
cd backend
go test ./...
```

#### sqlcによるコード生成

スキーマやクエリを変更した場合は、sqlcで型安全なコードを再生成します：

```bash
cd backend
sqlc generate
```

### フロントエンド

#### テストの実行

```bash
cd frontend
npm test
```

#### Lintの実行

```bash
cd frontend
npm run lint
```

#### ビルド

```bash
cd frontend
npm run build
```

---

## プロジェクトの状態

本プロジェクトは**MVP（Minimum Viable Product）がほぼ完成**した段階です。

### ✅ 完成している機能
- 初期設定フロー（収入・貯金・固定費の設定）
- ダッシュボード（残額表示、月次サマリー、色分け表示）
- 支出の登録・更新・削除
- 予定支出と確定支出の管理
- カテゴリ管理
- ダークモード
- フルレスポンシブデザイン
- リアルタイム更新

### 🚧 開発中・今後の拡張
- 認証機能（Firebase Auth）
- カテゴリ別支出の詳細分析
- 月別トレンド分析
- モバイルアプリ版（React Native）

### 動作確認済み環境
- **ブラウザ**: Chrome, Safari, Firefox（最新版）
- **デバイス**: PC、タブレット、スマートフォン
- **OS**: macOS, Windows, iOS, Android

---

## 開発のこだわりポイント

### 1. ユーザー体験重視
- 「今いくら使えるか」を最優先で表示
- 予定支出機能で将来を見越した管理
- 色による直感的な危険度の可視化

### 2. 継続しやすい設計
- 最小限の入力項目（金額・日付・カテゴリのみ必須）
- ステータス切替のシンプルなUI
- モバイルでも快適な操作性

### 3. 技術的な堅牢性
- TypeScript + sqlcによる完全な型安全性
- PostgreSQLによる信頼性の高いデータ管理
- OpenAPI仕様書によるAPI設計の透明性
- テストコード完備（Go: testify, フロントエンド: Jest）

### 4. 拡張性
- クリーンアーキテクチャ（レイヤー分離）
- 依存性注入によるテスタビリティ
- 認証機能の追加を想定した設計
- モバイルアプリへの展開を見据えたAPI設計

---

## ライセンス

本プロジェクトのライセンスは未定です。

---

## 貢献

現在、本プロジェクトは個人開発中です。

---

## お問い合わせ

プロジェクトに関するご質問は、GitHubのIssuesよりお願いいたします。
