package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDashboardRepo struct {
	getMonthlySummaryFunc         func(ctx context.Context, userID string) (*MonthlySummary, error)
	getMonthlyExpensesSummaryFunc func(ctx context.Context, userID string) (*MonthlyExpensesSummary, error)
}

func (m *mockDashboardRepo) GetMonthlySummary(ctx context.Context, userID string) (*MonthlySummary, error) {
	if m.getMonthlySummaryFunc != nil {
		return m.getMonthlySummaryFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDashboardRepo) GetMonthlyExpensesSummary(ctx context.Context, userID string) (*MonthlyExpensesSummary, error) {
	if m.getMonthlyExpensesSummaryFunc != nil {
		return m.getMonthlyExpensesSummaryFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func TestGetDashboard_Success(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*MonthlySummary, error) {
			assert.Equal(t, "test-user", userID)
			return &MonthlySummary{
				Income:     300000,
				SavingGoal: 50000,
				FixedCosts: 100000,
			}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*MonthlyExpensesSummary, error) {
			assert.Equal(t, "test-user", userID)
			return &MonthlyExpensesSummary{
				ConfirmedExpenses: 80000,
				PlannedExpenses:   20000,
			}, nil
		},
	}

	uc := NewGetDashboardUseCase(repo)
	dashboard, err := uc.Execute(context.Background(), "test-user")

	require.NoError(t, err)
	require.NotNil(t, dashboard)

	assert.Equal(t, int64(300000), dashboard.Income)
	assert.Equal(t, int64(50000), dashboard.SavingGoal)
	assert.Equal(t, int64(100000), dashboard.FixedCosts)
	assert.Equal(t, int64(80000), dashboard.ConfirmedExpenses)
	assert.Equal(t, int64(20000), dashboard.PlannedExpenses)
	// 変動費 = 300000 - 100000 - 50000 = 150000
	assert.Equal(t, int64(150000), dashboard.VariableBudget)
	// 残額 = 150000 - (80000 + 20000) = 50000
	assert.Equal(t, int64(50000), dashboard.Remaining)
}

func TestGetDashboard_ZeroExpenses(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*MonthlySummary, error) {
			return &MonthlySummary{Income: 300000, SavingGoal: 50000, FixedCosts: 100000}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*MonthlyExpensesSummary, error) {
			return &MonthlyExpensesSummary{ConfirmedExpenses: 0, PlannedExpenses: 0}, nil
		},
	}

	uc := NewGetDashboardUseCase(repo)
	dashboard, err := uc.Execute(context.Background(), "test-user")

	require.NoError(t, err)
	assert.Equal(t, int64(150000), dashboard.VariableBudget)
	assert.Equal(t, int64(150000), dashboard.Remaining)
}

func TestGetDashboard_ZeroFixedCosts(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*MonthlySummary, error) {
			return &MonthlySummary{Income: 300000, SavingGoal: 50000, FixedCosts: 0}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*MonthlyExpensesSummary, error) {
			return &MonthlyExpensesSummary{ConfirmedExpenses: 100000, PlannedExpenses: 50000}, nil
		},
	}

	uc := NewGetDashboardUseCase(repo)
	dashboard, err := uc.Execute(context.Background(), "test-user")

	require.NoError(t, err)
	assert.Equal(t, int64(250000), dashboard.VariableBudget)
	assert.Equal(t, int64(100000), dashboard.Remaining)
}

func TestGetDashboard_UserNotFound(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*MonthlySummary, error) {
			return nil, sql.ErrNoRows
		},
	}

	uc := NewGetDashboardUseCase(repo)
	dashboard, err := uc.Execute(context.Background(), "non-existent-user")

	require.Error(t, err)
	require.Nil(t, dashboard)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestGetDashboard_ExpensesSummaryError(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*MonthlySummary, error) {
			return &MonthlySummary{Income: 300000, SavingGoal: 50000, FixedCosts: 100000}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*MonthlyExpensesSummary, error) {
			return nil, errors.New("database connection error")
		},
	}

	uc := NewGetDashboardUseCase(repo)
	dashboard, err := uc.Execute(context.Background(), "test-user")

	require.Error(t, err)
	require.Nil(t, dashboard)
	assert.Contains(t, err.Error(), "database connection error")
}

func TestGetDashboard_NegativeRemaining(t *testing.T) {
	repo := &mockDashboardRepo{
		getMonthlySummaryFunc: func(ctx context.Context, userID string) (*MonthlySummary, error) {
			return &MonthlySummary{Income: 300000, SavingGoal: 50000, FixedCosts: 100000}, nil
		},
		getMonthlyExpensesSummaryFunc: func(ctx context.Context, userID string) (*MonthlyExpensesSummary, error) {
			return &MonthlyExpensesSummary{ConfirmedExpenses: 120000, PlannedExpenses: 80000}, nil
		},
	}

	uc := NewGetDashboardUseCase(repo)
	dashboard, err := uc.Execute(context.Background(), "test-user")

	require.NoError(t, err)
	assert.Equal(t, int64(150000), dashboard.VariableBudget)
	assert.Equal(t, int64(-50000), dashboard.Remaining)
}
