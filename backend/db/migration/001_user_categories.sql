-- Migration: 001_user_categories
-- Purpose: カテゴリをユーザーごとに管理できるようにする
-- ※ このファイルを実行する前に必ずバックアップを取ること

-- ============================================================
-- Step 1: user_categories テーブル作成
-- ============================================================
CREATE TABLE IF NOT EXISTS user_categories (
  id SERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  name TEXT NOT NULL,
  category_type TEXT NOT NULL DEFAULT 'user',
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT user_categories_type_check CHECK (category_type IN ('default', 'user', 'other'))
);

CREATE INDEX IF NOT EXISTS idx_user_categories_user_id ON user_categories(user_id);

-- ============================================================
-- Step 2: 既存ユーザー × 既存カテゴリ から per-user カテゴリを生成
--   - 'その他' という名前のカテゴリは category_type='other' にする
--   - sort_order は旧 categories.id を流用（表示順の基準として使用）
-- ============================================================
INSERT INTO user_categories (user_id, name, category_type, sort_order, created_at, updated_at)
SELECT
  u.id AS user_id,
  c.name,
  CASE WHEN c.name = 'その他' THEN 'other' ELSE 'default' END AS category_type,
  c.id AS sort_order,
  now(),
  now()
FROM users u
CROSS JOIN categories c
ON CONFLICT DO NOTHING;

-- ============================================================
-- Step 3: expenses.category_id を user_categories.id に更新
--   マッピング: user_categories.sort_order = 旧 categories.id かつ user_id が一致
-- ============================================================
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS new_category_id INTEGER;

UPDATE expenses e
SET new_category_id = uc.id
FROM user_categories uc
WHERE uc.user_id = e.user_id
  AND uc.sort_order = e.category_id;

-- new_category_id が NULL のままの行がないことを確認してから次のステップへ
-- SELECT COUNT(*) FROM expenses WHERE new_category_id IS NULL;

-- ============================================================
-- Step 4: category_id 列を差し替え
-- ============================================================
ALTER TABLE expenses DROP COLUMN category_id;
ALTER TABLE expenses RENAME COLUMN new_category_id TO category_id;
ALTER TABLE expenses ALTER COLUMN category_id SET NOT NULL;

-- ============================================================
-- Step 5: (オプション) 旧 categories テーブルはバックアップ後に削除可能
--   DROP TABLE categories;
-- ============================================================
