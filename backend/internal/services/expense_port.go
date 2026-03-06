package services

import "money-buddy-backend/internal/models"

type ExpenseRepository interface {
	CreateExpense(userID string, input models.CreateExpenseInput) (models.Expense, error)
	FindAll(userID string) ([]models.Expense, error)
	GetExpenseByID(userID string, id int32) (models.Expense, error)
	DeleteExpense(userID string, id int32) error
	UpdateExpense(userID string, input models.UpdateExpenseInput) (models.Expense, error)
}
