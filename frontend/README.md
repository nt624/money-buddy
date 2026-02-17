# Money Buddy - Frontend

> Next.js 16 + React 19 + TypeScript による家計管理アプリ

このディレクトリは Money Buddy のフロントエンドです。

## 🚀 本番環境

- **デプロイ先**: Vercel
- **URL**: https://money-buddy-app.vercel.app

## 技術スタック

- **フレームワーク**: Next.js 16 (App Router)
- **ライブラリ**: React 19
- **言語**: TypeScript 5
- **スタイリング**: Tailwind CSS 4
- **認証**: Firebase Authentication (SDK 9.x)
- **状態管理**: React Hooks（useState, useEffect, Context API）
- **HTTP通信**: Fetch API
- **Linter**: ESLint
- **テスト**: Jest

## 📁 プロジェクト構成

```
frontend/
├── src/
│   ├── app/                  # Next.js App Router
│   │   ├── page.tsx         # ダッシュボード
│   │   ├── layout.tsx       # ルートレイアウト
│   │   ├── login/           # ログイン画面
│   │   ├── settings/        # 設定画面
│   │   └── globals.css
│   ├── components/           # UIコンポーネント
│   │   ├── Dashboard.tsx
│   │   ├── ExpenseForm.tsx
│   │   ├── ExpenseList.tsx
│   │   ├── FixedCostForm.tsx
│   │   ├── InitialSetupForm.tsx
│   │   ├── Layout/          # レイアウト
│   │   │   ├── Header.tsx
│   │   │   ├── AuthGuard.tsx
│   │   │   └── Container.tsx
│   │   └── ui/              # 共通UIパーツ
│   ├── contexts/
│   │   └── AuthContext.tsx  # 認証状態管理
│   ├── hooks/               # カスタムフック
│   │   ├── useExpenses.ts
│   │   ├── useUser.ts
│   │   ├── useDashboard.ts
│   │   └── useTheme.ts     # ダークモード
│   └── lib/
│       ├── api/            # APIクライアント（7ファイル）
│       │   ├── client.ts   # 共通設定
│       │   ├── expenses.ts
│       │   ├── categories.ts
│       │   └── ...
│       ├── firebase/
│       │   └── config.ts   # Firebase初期化
│       └── types/          # TypeScript型定義
├── public/
├── .env.local.example      # 環境変数テンプレート
├── package.json
├── tsconfig.json
├── next.config.ts
├── tailwind.config.ts
└── jest.config.js
```

## 🛠️ ローカル開発環境のセットアップ

### 前提条件

- Node.js 20以上
- npm（Node.jsに同梱）

### 1. 依存関係のインストール

```bash
npm install
```

### 2. 環境変数の設定

`.env.local` ファイルを作成：

```bash
# バックエンドAPI（ローカル開発）
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# Firebase設定（Firebase Consoleから取得）
NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123456789:web:abc123
```

### 3. 開発サーバー起動

```bash
npm run dev
```

ブラウザで http://localhost:3000 にアクセス

## 📦 ビルド・本番起動

```bash
# プロダクションビルド
npm run build

# ビルド済みアプリの起動
npm run start
```

## 🧪 テスト・Lint

```bash
# Jest テスト実行
npm test

# ESLint チェック
npm run lint

# ESLint 自動修正
npm run lint:fix
```

## 🎨 主な機能

### 認証
- Firebase Authentication
- メール/パスワード認証
- Google OAuth
- 保護されたルート
- セッション期限切れの自動ログアウト

### ダッシュボード
- 今月の残額表示（色分け表示）
- 月次サマリー
- レスポンシブデザイン

### 支出管理
- 支出の登録・更新・削除
- 予定支出と確定支出の管理
- カテゴリ別表示
- リアルタイム更新

### UI/UX
- ダークモード対応
- モバイルファースト設計
- Tailwind CSS 4によるモダンなデザイン

## 🔧 開発のヒント

### API通信

すべてのAPI通信は `src/lib/api/` 配下のモジュールを使用：

```typescript
import { createExpense } from '@/lib/api/expenses'
import { getAuthHeaders } from '@/lib/api/client'

// 認証ヘッダーは自動的に付与される
const expense = await createExpense({ amount: 1000, ... })
```

### 認証状態の利用

```typescript
import { useAuth } from '@/contexts/AuthContext'

function MyComponent() {
  const { user, loading, signIn, signOut } = useAuth()
  
  if (loading) return <div>Loading...</div>
  if (!user) return <div>Not logged in</div>
  
  return <div>Hello, {user.email}</div>
}
```

### カスタムフックの活用

```typescript
import { useExpenses } from '@/hooks/useExpenses'
import { useDashboard } from '@/hooks/useDashboard'

function MyComponent() {
  const { expenses, createExpense, isLoading } = useExpenses()
  const { dashboard, refetch } = useDashboard()
  
  // ...
}
```

## 🐛 トラブルシューティング

### API接続エラー

- `NEXT_PUBLIC_API_BASE_URL` が正しく設定されているか確認
- バックエンドが起動しているか確認（http://localhost:8080/health）

### Firebase認証エラー

- Firebase設定（`NEXT_PUBLIC_FIREBASE_*`）が正しいか確認
- Firebase ConsoleでAuthorized domainsに`localhost`が追加されているか確認

### ビルドエラー

```bash
# node_modulesとキャッシュを削除して再インストール
rm -rf node_modules .next
npm install
npm run build
```

## 📚 関連ドキュメント

- [ルートREADME](../README.md) - プロジェクト全体の概要
- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Firebase Documentation](https://firebase.google.com/docs)

## 🚀 デプロイ

### Vercel（本番環境）

1. Vercelプロジェクトを作成
2. GitHubリポジトリを連携
3. 環境変数を設定：
   - `NEXT_PUBLIC_API_BASE_URL`（RailwayのURL）
   - `NEXT_PUBLIC_FIREBASE_*`（全7項目）
4. 自動デプロイが実行される

**重要**: `NEXT_PUBLIC_API_BASE_URL`には必ず`https://`を含めてください。

## 🤝 コントリビューション

このプロジェクトは個人開発中です。

## 📝 ライセンス

未定

フロントエンドは Fetch API を使って、バックエンドが提供する REST API と通信します。

## 開発環境のセットアップ

### 前提条件

- **Node.js**: バージョン 20 以上を推奨
- **npm**: Node.js に同梱

### インストール手順

1. **リポジトリをクローン**

```bash
git clone https://github.com/nt624/money-buddy-frontend.git
cd money-buddy-frontend
```

2. **依存関係をインストール**

```bash
npm install
```

3. **開発サーバーを起動**

```bash
npm run dev
```

4. **ブラウザでアクセス**

開発サーバーが起動したら、ブラウザで以下の URL にアクセスしてください：

```
http://localhost:3000
```

ページを編集すると、ブラウザが自動的にリロードされます。

### その他のコマンド

```bash
npm run build  # プロダクションビルド
npm run start  # ビルド済みアプリの起動
npm run lint   # ESLint によるコードチェック
npm run test   # Jest テストの実行
```

## 環境変数

現在の MVP では、API のベース URL はコード内にハードコードされていますが、将来的には環境変数で管理することを想定しています。

### .env.local の例

プロジェクトルートに `.env.local` ファイルを作成し、以下のように設定してください：

```env
# バックエンド API のベース URL
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# その他の環境変数（必要に応じて追加）
# NEXT_PUBLIC_ENVIRONMENT=development
```

**注意**: `.env.local` ファイルは `.gitignore` に含まれており、Git にコミットされません。

## ディレクトリ構成

```
money-buddy-frontend/
├── src/
│   ├── app/                    # Next.js App Router のページ
│   │   ├── page.tsx           # トップページ（支出入力画面）
│   │   ├── layout.tsx         # ルートレイアウト
│   │   └── globals.css        # グローバルスタイル
│   ├── components/             # 再利用可能な UI コンポーネント
│   │   ├── ExpenseForm.tsx    # 支出入力フォーム
│   │   └── ExpenseList.tsx    # 支出一覧表示
│   ├── hooks/                  # カスタム React フック
│   │   └── useExpenses.ts     # 支出データの取得・追加ロジック
│   └── lib/                    # ユーティリティ・API クライアント
│       ├── api/               # バックエンド API との通信
│       │   ├── expenses.ts    # 支出 API
│       │   └── categories.ts  # カテゴリ API
│       └── types/             # TypeScript 型定義
│           ├── expense.ts     # 支出関連の型
│           └── category.ts    # カテゴリ関連の型
├── public/                     # 静的ファイル
├── package.json               # プロジェクト設定と依存関係
├── tsconfig.json              # TypeScript 設定
├── next.config.ts             # Next.js 設定
└── README.md                  # このファイル
```

## 開発方針・設計メモ

### MVP 重視

まずは動くものを作り、フィードバックを得ることを優先しています。
完璧を目指すよりも、早くリリースして改善していくアプローチです。

### 状態管理はシンプルに

現在は React の `useState` と `useEffect` を使ったシンプルな状態管理です。
将来的にアプリが複雑になれば、状態管理ライブラリ（Zustand、Jotai など）の導入も検討します。

### UI より「動く・わかる」を優先

デザインの洗練度よりも、機能の実装と動作の確認を優先しています。
ユーザーが直感的に操作できる UI を目指しつつ、まずは機能を確実に動かすことに注力しています。

### コンポーネント設計

- **プレゼンテーショナルコンポーネント**: UI の表示のみを担当（ExpenseList）
- **コンテナコンポーネント**: ロジックとデータ取得を担当（page.tsx + useExpenses）
- コンポーネントは小さく、再利用可能に保つ

### API 通信

- Fetch API を直接使用（シンプルさを優先）
- エラーハンドリングは基本的なもののみ
- 将来的には SWR や React Query の導入も検討

## 今後やりたいこと

MVP の次のステップとして、以下の機能を検討しています：

### 短期的な改善

- **ローディング・エラーハンドリングの強化**
  - より詳細なエラーメッセージ
  - リトライ機能
  - トースト通知

- **環境変数の活用**
  - API_BASE_URL を `.env.local` から読み込む

- **テストの拡充**
  - コンポーネントのテスト追加
  - E2E テストの導入

### 中期的な機能追加

- **初期設定画面**
  - 月の予算設定
  - カテゴリのカスタマイズ

- **残額表示機能**
  - 使える金額の可視化

- **モバイル対応の強化**
  - レスポンシブデザインの改善
  - タッチ操作の最適化

### 長期的なビジョン

- 支出の編集・削除機能
- グラフやチャートでの可視化
- カテゴリ別の集計
- 月次レポート
- PWA 対応

---

## ライセンス

このプロジェクトはプライベートリポジトリです。
