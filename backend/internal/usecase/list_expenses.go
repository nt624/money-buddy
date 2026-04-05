package usecase

import (
	"time"

	"money-buddy-backend/internal/domain"
)

type MonthFilter struct {
	Year  int
	Month int
}

type listExpenseRepository interface {
	FindByMonth(userID string, year, month int) ([]domain.Expense, error)
}

// ListExpensesUseCase は支出一覧取得ユースケースです。
type ListExpensesUseCase struct {
	repo  listExpenseRepository
	nowFn func() time.Time
}

// NewListExpensesUseCase は ListExpensesUseCase の新しいインスタンスを生成します。
func NewListExpensesUseCase(repo listExpenseRepository) *ListExpensesUseCase {
	return &ListExpensesUseCase{repo: repo, nowFn: time.Now}
}

// Execute は支出一覧取得ユースケースを実行します。year/month が両方 0 の場合は当月を使用します。
func (uc *ListExpensesUseCase) Execute(userID string, filter MonthFilter) ([]domain.Expense, error) {
	year, month := filter.Year, filter.Month
	if year == 0 && month == 0 {
		now := uc.nowFn()
		year = now.Year()
		month = int(now.Month())
	}
	return uc.repo.FindByMonth(userID, year, month)
}
