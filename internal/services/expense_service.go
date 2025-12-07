package services

import (
	"money-buddy-backend/internal/models"
	"money-buddy-backend/internal/repositories"
)

type ExpenseService interface {
	CreateExpense(input models.CreateExpenseInput) (models.Expense, error)
	ListExpenses() ([]models.Expense, error)
}

type expenseService struct {
	repo repositories.ExpenseRepository
}

func NewExpenseService(repo repositories.ExpenseRepository) ExpenseService {
	return &expenseService{repo: repo}
}

func (s *expenseService) CreateExpense(input models.CreateExpenseInput) (models.Expense, error) {
	// 今は特にビジネスロジックなし。将来:
	// - 上限チェック
	// - 不正カテゴリの排除
	// - 月内上限との突き合わせ
	return s.repo.CreateExpense(input)
}

func (s *expenseService) ListExpenses() ([]models.Expense, error) {
	return s.repo.FindAll()
}
