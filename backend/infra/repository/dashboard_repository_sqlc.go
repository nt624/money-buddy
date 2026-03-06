package repository

import (
	"context"

	db "money-buddy-backend/db/generated"
	"money-buddy-backend/internal/services"
)

type dashboardRepositorySQLC struct {
	q *db.Queries
}

func NewDashboardRepositorySQLC(q *db.Queries) *dashboardRepositorySQLC {
	return &dashboardRepositorySQLC{q: q}
}

func (r *dashboardRepositorySQLC) GetMonthlySummary(ctx context.Context, userID string) (*services.MonthlySummary, error) {
	row, err := r.q.GetMonthlySummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &services.MonthlySummary{
		Income:     int64(row.Income),
		SavingGoal: int64(row.SavingGoal),
		FixedCosts: row.FixedCosts,
	}, nil
}

func (r *dashboardRepositorySQLC) GetMonthlyExpensesSummary(ctx context.Context, userID string) (*services.MonthlyExpensesSummary, error) {
	row, err := r.q.GetMonthlyExpensesSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &services.MonthlyExpensesSummary{
		ConfirmedExpenses: row.ConfirmedExpenses,
		PlannedExpenses:   row.PendingExpenses,
	}, nil
}
