package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money-buddy-backend/internal/repositories"
)

// mockDashboardRepo は DashboardRepository のモック実装です
type mockDashboardRepo struct {
	getMonthlySummaryFunc         func(ctx context.Context, userID string) (*repositories.MonthlySummary, error)
	getMonthlyExpensesSummaryFunc func(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error)
}

func (m *mockDashboardRepo) GetMonthlySummary(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
	if m.getMonthlySummaryFunc != nil {
		return m.getMonthlySummaryFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDashboardRepo) GetMonthlyExpensesSummary(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
	if m.getMonthlyExpensesSummaryFunc != nil {
		return m.getMonthlyExpensesSummaryFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

// TestGetDashboard_Success は正常系のテストです
func TestGetDashboard_Success(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
			assert.Equal(t, "test-user", userID)
			return &repositories.MonthlySummary{
				Income:     300000,
				SavingGoal: 50000,
				FixedCosts: 100000,
			}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
			assert.Equal(t, "test-user", userID)
			return &repositories.MonthlyExpensesSummary{
				ConfirmedExpenses: 80000,
				PlannedExpenses:   20000,
			}, nil
		},
	}

	service := NewDashboardService(repo)
	dashboard, err := service.GetDashboard(context.Background(), "test-user")

	require.NoError(t, err)
	require.NotNil(t, dashboard)

	// 基本データの確認
	assert.Equal(t, int64(300000), dashboard.Income)
	assert.Equal(t, int64(50000), dashboard.SavingGoal)
	assert.Equal(t, int64(100000), dashboard.FixedCosts)
	assert.Equal(t, int64(80000), dashboard.ConfirmedExpenses)
	assert.Equal(t, int64(20000), dashboard.PlannedExpenses)

	// 計算結果の確認
	// 変動費 = 収入 - 固定費 - 貯金目標 = 300000 - 100000 - 50000 = 150000
	assert.Equal(t, int64(150000), dashboard.VariableBudget)
	// 残額 = 変動費 - (確定支出 + 予定支出) = 150000 - (80000 + 20000) = 50000
	assert.Equal(t, int64(50000), dashboard.Remaining)
}

// TestGetDashboard_ZeroExpenses は支出がゼロの場合のテストです
func TestGetDashboard_ZeroExpenses(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
			return &repositories.MonthlySummary{
				Income:     300000,
				SavingGoal: 50000,
				FixedCosts: 100000,
			}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
			return &repositories.MonthlyExpensesSummary{
				ConfirmedExpenses: 0,
				PlannedExpenses:   0,
			}, nil
		},
	}

	service := NewDashboardService(repo)
	dashboard, err := service.GetDashboard(context.Background(), "test-user")

	require.NoError(t, err)
	require.NotNil(t, dashboard)

	// 変動費 = 300000 - 100000 - 50000 = 150000
	assert.Equal(t, int64(150000), dashboard.VariableBudget)
	// 残額 = 変動費 - 0 = 150000
	assert.Equal(t, int64(150000), dashboard.Remaining)
}

// TestGetDashboard_ZeroFixedCosts は固定費がゼロの場合のテストです
func TestGetDashboard_ZeroFixedCosts(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
			return &repositories.MonthlySummary{
				Income:     300000,
				SavingGoal: 50000,
				FixedCosts: 0,
			}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
			return &repositories.MonthlyExpensesSummary{
				ConfirmedExpenses: 100000,
				PlannedExpenses:   50000,
			}, nil
		},
	}

	service := NewDashboardService(repo)
	dashboard, err := service.GetDashboard(context.Background(), "test-user")

	require.NoError(t, err)
	require.NotNil(t, dashboard)

	// 変動費 = 300000 - 0 - 50000 = 250000
	assert.Equal(t, int64(250000), dashboard.VariableBudget)
	// 残額 = 250000 - (100000 + 50000) = 100000
	assert.Equal(t, int64(100000), dashboard.Remaining)
}

// TestGetDashboard_UserNotFound はユーザーが存在しない場合のテストです
func TestGetDashboard_UserNotFound(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := NewDashboardService(repo)
	dashboard, err := service.GetDashboard(context.Background(), "non-existent-user")

	require.Error(t, err)
	require.Nil(t, dashboard)
	assert.Equal(t, sql.ErrNoRows, err)
}

// TestGetDashboard_ExpensesSummaryError は支出サマリー取得時のエラーをテストします
func TestGetDashboard_ExpensesSummaryError(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
			return &repositories.MonthlySummary{
				Income:     300000,
				SavingGoal: 50000,
				FixedCosts: 100000,
			}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := NewDashboardService(repo)
	dashboard, err := service.GetDashboard(context.Background(), "test-user")

	require.Error(t, err)
	require.Nil(t, dashboard)
	assert.Contains(t, err.Error(), "database connection error")
}

// TestGetDashboard_NegativeRemaining は残額がマイナスになる場合のテストです
func TestGetDashboard_NegativeRemaining(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
			return &repositories.MonthlySummary{
				Income:     300000,
				SavingGoal: 50000,
				FixedCosts: 100000,
			}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
			return &repositories.MonthlyExpensesSummary{
				ConfirmedExpenses: 120000,
				PlannedExpenses:   80000,
			}, nil
		},
	}

	service := NewDashboardService(repo)
	dashboard, err := service.GetDashboard(context.Background(), "test-user")

	require.NoError(t, err)
	require.NotNil(t, dashboard)

	// 変動費 = 300000 - 100000 - 50000 = 150000
	assert.Equal(t, int64(150000), dashboard.VariableBudget)
	// 残額 = 150000 - (120000 + 80000) = -50000 (マイナス)
	assert.Equal(t, int64(-50000), dashboard.Remaining)
}
