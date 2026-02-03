package repository

import (
	"context"
	"database/sql"

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
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	// interface{} から int64 への変換
	fixedCosts := int64(0)
	if row.FixedCosts != nil {
		if fc, ok := row.FixedCosts.(int64); ok {
			fixedCosts = fc
		}
	}

	return &repositories.MonthlySummary{
		Income:     int64(row.Income),
		SavingGoal: int64(row.SavingGoal),
		FixedCosts: fixedCosts,
	}, nil
}

func (r *dashboardRepositorySQLC) GetMonthlyExpensesSummary(ctx context.Context, userID string) (*repositories.MonthlyExpensesSummary, error) {
	row, err := r.q.GetMonthlyExpensesSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	// interface{} から int64 への変換
	confirmedExpenses := int64(0)
	plannedExpenses := int64(0)

	if row.ConfirmedExpenses != nil {
		if ce, ok := row.ConfirmedExpenses.(int64); ok {
			confirmedExpenses = ce
		}
	}

	if row.PendingExpenses != nil {
		if pe, ok := row.PendingExpenses.(int64); ok {
			plannedExpenses = pe
		}
	}

	return &repositories.MonthlyExpensesSummary{
		ConfirmedExpenses: confirmedExpenses,
		PlannedExpenses:   plannedExpenses,
	}, nil
}
