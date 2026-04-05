# Pace Wallet - Project Guidelines for Claude Code

## Project Overview

**Pace Wallet**（ペースウォレット）は、「お金が貯まらない人が"貯まる生活"を続けられるように伴走する家計管理アプリ」です。

「今の自分のお金の状況が一目でわかるアプリ」をコンセプトに、シンプルな入力とリアルタイムの残額表示により、貯まる行動を自然に習慣化することを目指しています。

## Tech Stack

### Frontend
- **Framework**: Next.js 16 (App Router)
- **Library**: React 19
- **Language**: TypeScript
- **Styling**: Tailwind CSS 4
- **Linter**: ESLint
- **Testing**: Jest
- **Package manager**: npm
- **Deployment**: Vercel

### Backend
- **Language**: Go 1.25
- **Web framework**: Gin
- **Database**: PostgreSQL
- **Query builder**: sqlc
- **Auth**: Firebase Auth (Firebase Admin SDK)
- **Testing**: testify
- **API spec**: OpenAPI 3.0
- **Deployment**: Railway

### Infrastructure
- **Authentication**: Firebase Auth
- **Database hosting**: Neon / Supabase / Railway / Render

## Repository Structure

```
pace-wallet-app/
├── frontend/          # Next.js app
│   └── src/
│       ├── app/       # App Router pages & layouts
│       ├── components/
│       ├── contexts/
│       ├── hooks/
│       └── lib/       # API client, Firebase config, types
├── backend/           # Go API server
│   ├── cmd/server/    # Entry point
│   ├── db/            # Schema, queries, sqlc generated code
│   ├── infra/         # Repository & transaction implementations
│   ├── internal/      # Handlers, middleware, usecase, domain
│   └── openapi/       # OpenAPI spec
└── docs/              # Project documentation
```

## Rules

### Before creating or modifying code

**コードを作成・修正する前に、必ず対応する GitHub issue の番号をユーザーに確認してください。**

### Branch naming

- **必ず `main` ブランチから作成すること**
- Format: `(issue#)-(prefix)_description_in_english`
- Prefix must be one of:
  - `add` — 新機能追加
  - `fix` — バグ修正
  - `refactor` — リファクタリング
  - `update` — 既存機能の更新
  - `remove` — 機能・コードの削除
  - `docs` — ドキュメント
  - `test` — テスト追加・修正
  - `chore` — ビルド設定など雑務
- Words in description are separated by `_`
- Example: `42-add_expense_filter_by_category`

### Commit message format

- Format: `#(issue番号) (修正内容を現在時制の日本語で)`
- Example: `#42 カテゴリ別フィルター機能を追加する`
