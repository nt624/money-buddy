package services

import (
	"context"

	"money-buddy-backend/internal/repositories"
)

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

// DashboardService はダッシュボードサービスのインターフェースです。
type DashboardService interface {
	GetDashboard(ctx context.Context, userID string) (*Dashboard, error)
}

type dashboardService struct {
	repo repositories.DashboardRepository
}

// NewDashboardService は DashboardService の新しいインスタンスを作成します。
func NewDashboardService(repo repositories.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

// GetDashboard はダッシュボード表示用のデータを取得します。
func (s *dashboardService) GetDashboard(ctx context.Context, userID string) (*Dashboard, error) {
	// 月次サマリー（収入・貯金目標・固定費）を取得
	summary, err := s.repo.GetMonthlySummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 月次支出サマリー（確定支出・予定支出）を取得
	expenses, err := s.repo.GetMonthlyExpensesSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 変動費を計算: 収入 - 固定費 - 貯金目標
	variableBudget := summary.Income - summary.FixedCosts - summary.SavingGoal

	// 残額を計算: 変動費 - (確定支出 + 予定支出)
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
