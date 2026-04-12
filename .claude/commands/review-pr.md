# スキル: PRレビューコメント対応

GitHub Pull Request のレビューコメントを1件ずつコミットで対応し、各コメントに返信するスキルです。

## 基本フロー

```
1. PRの全レビューコメントを取得・分類
2. 各コメントを独立したコミットで対応
3. コミット後、そのコメントに対応内容とコミットハッシュを返信
4. lint/typecheck/test が通ることを確認
```

## ステップ詳細

### Step 1: コメント取得と分類

```bash
# PRのレビューコメント一覧を取得
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)
gh api repos/${REPO}/pulls/{PR番号}/comments \
  --jq '.[] | {id: .id, path: .path, body: .body[0:150]}'
```

コメントを以下の3種類に分類する：
- **対応する**: コードを修正して対応できるもの
- **PRで返信するだけ**: プロジェクトのルール・方針と矛盾するため修正しないもの
- **スキップ**: 既に対応済みのもの

### Step 2: 依存関係の確認

複数コメントをまとめざるを得ないケース（1コミットにまとめる）：
- SQLクエリ変更 → `sqlc generate` → リポジトリ修正（生成コードが絡む）
- 新しい型/関数の追加と、それを使う側の変更が同時に必要な場合
- pre-commit フックが途中状態でエラーになる場合

その他は原則 **1コメント = 1コミット** とする。

### Step 3: 対応とコミット

```bash
# コードを修正後
git add <変更ファイル>
git commit -m "#(issue番号) (対応内容を日本語で)"

# コミットハッシュを取得
HASH=$(git rev-parse --short HEAD)
```

### Step 4: コメントへの返信

```bash
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)
gh api repos/${REPO}/pulls/{PR番号}/comments/{COMMENT_ID}/replies \
  -f body="対応済みです（${HASH}）。

- (対応内容を箇条書きで記述)"
```

返信例：
```
対応済みです（b331112）。

- year/month の範囲チェック（year >= 1、month 1〜12）を Handler と UseCase の両層に追加
- 範囲外の場合は 400 Bad Request を返すように修正しました
```

### Step 5: 対応できないコメントへの返信

プロジェクトのルールや方針と矛盾するコメントには、コードを修正せず理由を返信する：

```bash
gh api repos/${REPO}/pulls/{PR番号}/comments/{COMMENT_ID}/replies \
  -f body="(矛盾する理由を説明するメッセージ)"
```

### Step 6: プッシュ

```bash
git push
```

## 注意事項

- 各コミット後に `go build ./...` または `npm run typecheck && npm run lint` でビルドが通ることを確認する
- pre-commit フックがある場合は自動でチェックされるが、フック通過後もビルドが通るか確認する
- 依存関係があるコメントはまとめてもよいが、その場合は返信に「XXとYYのコメントをまとめて対応しました」と明記する
- コメントへの返信は `対応済みです（{HASH}）。` + 対応内容の箇条書きを基本形とする
