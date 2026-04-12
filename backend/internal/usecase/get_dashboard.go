package usecase

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

// Dashboard はダッシュボード表示用のデータ構造です。
type Dashboard struct {
	Income            int64 // 月収
	SavingGoal        int64 // 貯金目標
	FixedCosts        int64 // 固定費合計
	VariableBudget    int64 // 変動費（自由に使える額）= 収入 - 固定費 - 貯金目標
	ConfirmedExpenses int64 // 確定支出
	PlannedExpenses   int64 // 予定支出
	Remaining         int64 // 残額 = 変動費 - (確定支出 + 予定支出)
}

type getDashboardRepository interface {
	GetMonthlySummary(ctx context.Context, userID string, year, month int) (*MonthlySummary, error)
	GetMonthlyExpensesSummary(ctx context.Context, userID string, year, month int) (*MonthlyExpensesSummary, error)
}

// GetDashboardUseCase はダッシュボード取得ユースケースです。
type GetDashboardUseCase struct {
	repo getDashboardRepository
}

// NewGetDashboardUseCase は GetDashboardUseCase の新しいインスタンスを生成します。
func NewGetDashboardUseCase(repo getDashboardRepository) *GetDashboardUseCase {
	return &GetDashboardUseCase{repo: repo}
}

// Execute はダッシュボード取得ユースケースを実行します。
func (uc *GetDashboardUseCase) Execute(ctx context.Context, userID string, year, month int) (*Dashboard, error) {
	summary, err := uc.repo.GetMonthlySummary(ctx, userID, year, month)
	if err != nil {
		return nil, err
	}

	expenses, err := uc.repo.GetMonthlyExpensesSummary(ctx, userID, year, month)
	if err != nil {
		return nil, err
	}

	variableBudget := summary.Income - summary.FixedCosts - summary.SavingGoal
	remaining := variableBudget - (expenses.ConfirmedExpenses + expenses.PlannedExpenses)

	return &Dashboard{
		Income:            summary.Income,
		SavingGoal:        summary.SavingGoal,
		FixedCosts:        summary.FixedCosts,
		VariableBudget:    variableBudget,
		ConfirmedExpenses: expenses.ConfirmedExpenses,
		PlannedExpenses:   expenses.PlannedExpenses,
		Remaining:         remaining,
	}, nil
}
