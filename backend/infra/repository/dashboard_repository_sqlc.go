package repository

import (
	"context"

	db "pace-wallet-backend/db/generated"
	"pace-wallet-backend/internal/usecase"
)

type dashboardRepositorySQLC struct {
	q *db.Queries
}

func NewDashboardRepositorySQLC(q *db.Queries) *dashboardRepositorySQLC {
	return &dashboardRepositorySQLC{q: q}
}

func (r *dashboardRepositorySQLC) GetMonthlySummary(ctx context.Context, userID string) (*usecase.MonthlySummary, error) {
	row, err := r.q.GetMonthlySummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &usecase.MonthlySummary{
		Income:     int64(row.Income),
		SavingGoal: int64(row.SavingGoal),
		FixedCosts: row.FixedCosts,
	}, nil
}

func (r *dashboardRepositorySQLC) GetMonthlyExpensesSummary(ctx context.Context, userID string) (*usecase.MonthlyExpensesSummary, error) {
	row, err := r.q.GetMonthlyExpensesSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &usecase.MonthlyExpensesSummary{
		ConfirmedExpenses: row.ConfirmedExpenses,
		PlannedExpenses:   row.PendingExpenses,
	}, nil
}
