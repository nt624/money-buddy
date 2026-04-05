package repository

import (
	"context"
	"database/sql"
	"time"

	db "money-buddy-backend/db/generated"
	"money-buddy-backend/internal/domain"
	"money-buddy-backend/internal/usecase"
)

// sqlc-backed repository
type expenseRepositorySQLC struct {
	q *db.Queries
}

func NewExpenseRepositorySQLC(q *db.Queries) *expenseRepositorySQLC {
	return &expenseRepositorySQLC{q: q}
}

func (r *expenseRepositorySQLC) CreateExpense(userID string, input usecase.CreateExpenseInput) (domain.Expense, error) {
	var spentAt time.Time
	var err error

	spentAt, err = time.Parse(time.RFC3339, input.SpentAt)
	if err != nil {
		spentAt, err = time.Parse("2006-01-02", input.SpentAt)
		if err != nil {
			return domain.Expense{}, err
		}
		spentAt = time.Date(spentAt.Year(), spentAt.Month(), spentAt.Day(), 0, 0, 0, 0, time.UTC)
	}

	params := db.CreateExpenseParams{
		UserID:     userID,
		Amount:     int32(*input.Amount),
		CategoryID: int32(*input.CategoryID),
		Memo:       sql.NullString{String: input.Memo, Valid: input.Memo != ""},
		SpentAt:    spentAt,
		Status:     defaultStatus(input.Status),
	}

	id, err := r.q.CreateExpense(context.Background(), params)
	if err != nil {
		return domain.Expense{}, err
	}

	row, err := r.q.GetExpenseWithCategoryByID(context.Background(), db.GetExpenseWithCategoryByIDParams{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		return domain.Expense{}, err
	}

	return dbExpenseToModel(row), nil
}

func (r *expenseRepositorySQLC) FindAll(userID string) ([]domain.Expense, error) {
	items, err := r.q.ListExpenses(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	var out []domain.Expense
	for _, it := range items {
		out = append(out, dbListExpenseRowToModel(it))
	}

	return out, nil
}

func (r *expenseRepositorySQLC) FindByMonth(userID string, year, month int) ([]domain.Expense, error) {
	items, err := r.q.ListExpensesByMonth(context.Background(), db.ListExpensesByMonthParams{
		UserID: userID,
		Year:   int32(year),
		Month:  int32(month),
	})
	if err != nil {
		return nil, err
	}

	var out []domain.Expense
	for _, it := range items {
		out = append(out, dbListExpensesByMonthRowToModel(it))
	}
	return out, nil
}

func dbListExpensesByMonthRowToModel(e db.ListExpensesByMonthRow) domain.Expense {
	memo := ""
	if e.Memo.Valid {
		memo = e.Memo.String
	}
	return domain.Expense{
		ID:       int(e.ID),
		Amount:   int(e.Amount),
		Memo:     memo,
		SpentAt:  e.SpentAt.Format(time.RFC3339),
		Status:   e.Status,
		Category: domain.Category{ID: int(e.CategoryID), Name: e.CategoryName},
	}
}

func dbExpenseToModel(e db.GetExpenseWithCategoryByIDRow) domain.Expense {
	memo := ""
	if e.Memo.Valid {
		memo = e.Memo.String
	}

	return domain.Expense{
		ID:       int(e.ID),
		Amount:   int(e.Amount),
		Memo:     memo,
		SpentAt:  e.SpentAt.Format(time.RFC3339),
		Status:   e.Status,
		Category: domain.Category{ID: int(e.CategoryID), Name: e.CategoryName},
	}
}

func dbListExpenseRowToModel(e db.ListExpensesRow) domain.Expense {
	memo := ""
	if e.Memo.Valid {
		memo = e.Memo.String
	}

	return domain.Expense{
		ID:       int(e.ID),
		Amount:   int(e.Amount),
		Memo:     memo,
		SpentAt:  e.SpentAt.Format(time.RFC3339),
		Status:   e.Status,
		Category: domain.Category{ID: int(e.CategoryID), Name: e.CategoryName},
	}
}

func (r *expenseRepositorySQLC) GetExpenseByID(userID string, id int32) (domain.Expense, error) {
	row, err := r.q.GetExpenseWithCategoryByID(context.Background(), db.GetExpenseWithCategoryByIDParams{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		return domain.Expense{}, err
	}

	return dbExpenseToModel(row), nil
}

func (r *expenseRepositorySQLC) DeleteExpense(userID string, id int32) error {
	return r.q.DeleteExpense(context.Background(), db.DeleteExpenseParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *expenseRepositorySQLC) UpdateExpense(userID string, input usecase.UpdateExpenseInput) (domain.Expense, error) {
	var spentAt time.Time
	var err error

	spentAt, err = time.Parse(time.RFC3339, input.SpentAt)
	if err != nil {
		spentAt, err = time.Parse("2006-01-02", input.SpentAt)
		if err != nil {
			return domain.Expense{}, err
		}
		spentAt = time.Date(spentAt.Year(), spentAt.Month(), spentAt.Day(), 0, 0, 0, 0, time.UTC)
	}

	params := db.UpdateExpenseParams{
		ID:         int32(input.ID),
		Amount:     int32(*input.Amount),
		CategoryID: int32(*input.CategoryID),
		Memo:       sql.NullString{String: input.Memo, Valid: input.Memo != ""},
		SpentAt:    spentAt,
		Status:     defaultStatus(input.Status),
		UserID:     userID,
	}
	err = r.q.UpdateExpense(context.Background(), params)

	if err != nil {
		return domain.Expense{}, err
	}

	return r.GetExpenseByID(userID, int32(input.ID))
}
