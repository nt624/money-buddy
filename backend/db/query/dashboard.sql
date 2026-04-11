-- name: GetMonthlySummary :one
SELECT
  COALESCE(ms.income, u.income)           AS income,
  COALESCE(ms.saving_goal, u.saving_goal) AS saving_goal,
  COALESCE(SUM(fc.amount), 0)::bigint     AS fixed_costs
FROM users u
LEFT JOIN fixed_costs fc ON fc.user_id = u.id
LEFT JOIN monthly_settings ms
  ON ms.user_id = u.id AND ms.year = $2 AND ms.month = $3
WHERE u.id = $1
GROUP BY u.id, ms.income, ms.saving_goal;

-- name: GetMonthlyExpensesSummary :one
SELECT
  COALESCE(SUM(CASE WHEN e.status = 'confirmed' THEN e.amount ELSE 0 END), 0)::bigint AS confirmed_expenses,
  COALESCE(SUM(CASE WHEN e.status = 'planned'   THEN e.amount ELSE 0 END), 0)::bigint AS pending_expenses
FROM expenses e
WHERE e.user_id = $1
  AND DATE_TRUNC('month', e.spent_at) = DATE_TRUNC('month', MAKE_DATE($2::int, $3::int, 1));