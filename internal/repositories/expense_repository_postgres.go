package repositories

import (
	"database/sql"
	"money-buddy-backend/internal/models"
)

type expenseRepositoryPostgres struct {
	db *sql.DB
}

func NewExpenseRepositoryPostgres(db *sql.DB) ExpenseRepository {
	return &expenseRepositoryPostgres{db: db}
}

func (r *expenseRepositoryPostgres) CreateExpense(input models.CreateExpenseInput) (models.Expense, error) {
	var e models.Expense

	query := `
        INSERT INTO expenses (amount, category_id, memo, spent_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, amount, category_id, memo, spent_at
    `

	err := r.db.QueryRow(
		query,
		input.Amount,
		input.CategoryID,
		input.Memo,
		input.SpentAt,
	).Scan(
		&e.ID,
		&e.Amount,
		&e.CategoryID,
		&e.Memo,
		&e.SpentAt,
	)

	return e, err
}

func (r *expenseRepositoryPostgres) FindAll() ([]models.Expense, error) {
	rows, err := r.db.Query(`
        SELECT id, amount, category_id, memo, spent_at
        FROM expenses
        ORDER BY spent_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		rows.Scan(&e.ID, &e.Amount, &e.CategoryID, &e.Memo, &e.SpentAt)
		expenses = append(expenses, e)
	}

	return expenses, nil
}
