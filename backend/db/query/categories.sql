-- name: ListUserCategories :many
SELECT id, name, category_type, sort_order
FROM user_categories
WHERE user_id = $1
ORDER BY
  CASE category_type WHEN 'other' THEN 1 ELSE 0 END,
  sort_order ASC,
  created_at ASC;

-- name: UserCategoryExists :one
SELECT EXISTS (
  SELECT 1 FROM user_categories WHERE id = $1 AND user_id = $2
);

-- name: CreateUserCategory :one
INSERT INTO user_categories (user_id, name, category_type, sort_order)
VALUES (
  $1,
  $2,
  'user',
  COALESCE(
    (SELECT MAX(sort_order) + 1 FROM user_categories WHERE user_id = $1 AND category_type != 'other'),
    1
  )
)
RETURNING id, name, category_type, sort_order;

-- name: UpdateUserCategory :one
UPDATE user_categories
SET name = $3, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, name, category_type, sort_order;

-- name: UpdateUserCategorySortOrder :exec
UPDATE user_categories
SET sort_order = $3, updated_at = now()
WHERE id = $1 AND user_id = $2;

-- name: DeleteUserCategory :exec
DELETE FROM user_categories WHERE id = $1 AND user_id = $2;

-- name: CountExpensesByUserCategory :one
SELECT COUNT(*) FROM expenses WHERE user_id = $1 AND category_id = $2;

-- name: SeedDefaultCategoriesForUser :exec
INSERT INTO user_categories (user_id, name, category_type, sort_order)
VALUES
  ($1, '食費',       'default', 1),
  ($1, '交通費',     'default', 2),
  ($1, '日用品',     'default', 3),
  ($1, '趣味・娯楽', 'default', 4),
  ($1, '医療費',     'default', 5),
  ($1, '衣服',       'default', 6),
  ($1, 'その他',     'other',   0)
ON CONFLICT (user_id, name) DO NOTHING;
