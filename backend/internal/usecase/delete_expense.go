package usecase

import (
	"database/sql"
	"errors"

	"money-buddy-backend/internal/domain"
)

type deleteExpenseRepository interface {
	GetExpenseByID(userID string, id int32) (domain.Expense, error)
	DeleteExpense(userID string, id int32) error
}

// DeleteExpenseUseCase は支出削除ユースケースです。
type DeleteExpenseUseCase struct {
	repo deleteExpenseRepository
}

// NewDeleteExpenseUseCase は DeleteExpenseUseCase の新しいインスタンスを生成します。
func NewDeleteExpenseUseCase(repo deleteExpenseRepository) *DeleteExpenseUseCase {
	return &DeleteExpenseUseCase{repo: repo}
}

// Execute は支出削除ユースケースを実行します。
func (uc *DeleteExpenseUseCase) Execute(userID string, id int) error {
	expense, err := uc.repo.GetExpenseByID(userID, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &NotFoundError{Message: "支出が見つかりません"}
		}
		return &InternalError{Message: "internal error"}
	}
	if expense == (domain.Expense{}) {
		return &NotFoundError{Message: "支出が見つかりません"}
	}

	return uc.repo.DeleteExpense(userID, int32(id))
}
