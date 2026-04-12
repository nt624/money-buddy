-- name: UpsertMonthlySetting :one
INSERT INTO monthly_settings (user_id, year, month, income, saving_goal, updated_at)
VALUES ($1, $2, $3, $4, $5, now())
ON CONFLICT (user_id, year, month)
DO UPDATE SET income = $4, saving_goal = $5, updated_at = now()
RETURNING *;

-- name: GetMonthlySetting :one
SELECT * FROM monthly_settings
WHERE user_id = $1 AND year = $2 AND month = $3;

-- name: DeleteMonthlySetting :exec
DELETE FROM monthly_settings
WHERE user_id = $1 AND year = $2 AND month = $3;
