package usecase

import "money-buddy-backend/internal/domain"

type listExpenseRepository interface {
	FindAll(userID string) ([]domain.Expense, error)
}

// ListExpensesUseCase は支出一覧取得ユースケースです。
type ListExpensesUseCase struct {
	repo listExpenseRepository
}

// NewListExpensesUseCase は ListExpensesUseCase の新しいインスタンスを生成します。
func NewListExpensesUseCase(repo listExpenseRepository) *ListExpensesUseCase {
	return &ListExpensesUseCase{repo: repo}
}

// Execute は支出一覧取得ユースケースを実行します。
func (uc *ListExpensesUseCase) Execute(userID string) ([]domain.Expense, error) {
	return uc.repo.FindAll(userID)
}
