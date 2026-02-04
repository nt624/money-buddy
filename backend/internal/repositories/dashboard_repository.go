package repositories

import "context"

// MonthlySummary は月次サマリー（収入・貯金目標・固定費）を表します。
type MonthlySummary struct {
	Income     int64
	SavingGoal int64
	FixedCosts int64
}

// MonthlyExpensesSummary は月次支出サマリー（確定支出・予定支出）を表します。
type MonthlyExpensesSummary struct {
	ConfirmedExpenses int64
	PlannedExpenses   int64
}

// DashboardRepository はダッシュボードリポジトリの振る舞いを表します。
type DashboardRepository interface {
	GetMonthlySummary(ctx context.Context, userID string) (*MonthlySummary, error)
	GetMonthlyExpensesSummary(ctx context.Context, userID string) (*MonthlyExpensesSummary, error)
}
