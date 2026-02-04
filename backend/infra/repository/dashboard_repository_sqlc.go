package repository

import (
	"context"

	db "money-buddy-backend/db/generated"
	"money-buddy-backend/internal/repositories"
)

type dashboardRepositorySQLC struct {
	q *db.Queries
}

func NewDashboardRepositorySQLC(q *db.Queries) repositories.DashboardRepository {
	return &dashboardRepositorySQLC{q: q}
}

func (r *dashboardRepositorySQLC) GetMonthlySummary(ctx context.Context, userID string) (*repositories.MonthlySummary, error) {
	row, err := r.q.GetMonthlySummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &repositories.MonthlySummary{
		Income:     int64(row.Income),
		SavingGoal: int64(row.SavingGoal),
		FixedCosts: row.FixedCosts,
	}, nil
}

func (r *dashboardRepositorySQLC) GetMonthlyExpensesSummary(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
	row, err := r.q.GetMonthlyExpensesSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &repositories.MonthlyExpensesSummary{
		ConfirmedExpenses: row.ConfirmedExpenses,
		PlannedExpenses:   row.PendingExpenses,
	}, nil
}
