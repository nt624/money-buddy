# GitHub Copilot 指示書

## 使用言語
日本語でレビューしてください。

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

## 主な機能

### 実装済み機能 ✅

#### 1. 初期設定機能
- ユーザー登録（月収、貯金目標額の設定）
- 固定費の登録・管理
- 自由に使える変動費の自動計算（変動費 = 収入 - 固定費 - 貯金額）

#### 2. 支出管理機能
- **支出の登録**
  - 金額、日付、カテゴリ、メモの入力
  - 支出ステータスの管理（予定 / 確定）
  - 予定支出：旅行・飲み会・イベントなど、これから発生する概算支出を事前登録
  - 確定支出：実際に使った金額
- **支出の更新**
  - 予定支出から確定支出への更新
  - 金額・カテゴリ・メモの編集
  - ステータス変更ルール：予定 → 確定は許可、確定 → 予定は禁止
- **支出の削除**
- **支出一覧の表示**

#### 3. カテゴリ管理
- デフォルトカテゴリの提供
- カテゴリ一覧の取得

#### 4. 月次サマリー機能（バックエンド）
- 月の収入、貯金目標、固定費の合計取得
- 月の確定支出・予定支出の集計

### 未実装機能（MVP範囲内）🚧

#### リアルタイム残額表示（ホーム画面）
- 今月の自由に使える残額の表示
- 残額計算ルール：`残額 = 変動費 - (確定支出 + 予定支出)`
- 使った金額の累計表示
- 貯金の進捗表示
- 今日までの利用ペースとの比較

### 今後実装予定（MVP後）📋

- 危険度の可視化（色分けメーター、ゲージ）
- カテゴリ別支出サマリー
- 月別の比較グラフ
- 無駄遣いアラート
- 収支レポート（PDF生成など）
- 認証機能（Firebase Auth）の実装
- モバイルアプリ（Expo / React Native）開発

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

## データモデル

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

## プロジェクトの状態

本プロジェクトは現在開発中（MVP段階）です。

- **完成している部分**: 支出・固定費の登録、初期設定、カテゴリ管理、月次サマリー集計（バックエンド）
- **開発中の部分**: ホーム画面の残額表示UI、ダッシュボード機能
- **未実装の部分**: 認証機能、モバイルアプリ、分析機能、アラート機能
