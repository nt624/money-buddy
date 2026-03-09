package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"money-buddy-backend/internal/domain"
)

const (
	// BusinessMaxAmount は業務上の上限（個人向け家計簿の想定）
	BusinessMaxAmount = 1000000000
	MemoMaxLen        = 5000
)

type CreateExpenseInput struct {
	Amount     *int
	CategoryID *int
	Memo       string
	SpentAt    string
	Status     string
}

type createExpenseRepository interface {
	CreateExpense(userID string, input CreateExpenseInput) (domain.Expense, error)
}

type categoryExistsChecker interface {
	CategoryExists(ctx context.Context, id int32) (bool, error)
}

type CreateExpenseUseCase struct {
	repo     createExpenseRepository
	category categoryExistsChecker
}

func NewCreateExpenseUseCase(repo createExpenseRepository, category categoryExistsChecker) *CreateExpenseUseCase {
	return &CreateExpenseUseCase{repo: repo, category: category}
}

func (uc *CreateExpenseUseCase) Execute(userID string, input CreateExpenseInput) (domain.Expense, error) {
	if input.Amount == nil {
		return domain.Expense{}, &ValidationError{Message: "金額を入力してください"}
	}
	if *input.Amount <= 0 {
		return domain.Expense{}, &ValidationError{Message: "金額は1円以上で入力してください"}
	}
	if *input.Amount > BusinessMaxAmount {
		return domain.Expense{}, &ValidationError{Message: "金額は10億円以下で入力してください"}
	}

	if input.CategoryID == nil {
		return domain.Expense{}, &ValidationError{Message: "カテゴリを選択してください"}
	}
	if *input.CategoryID <= 0 {
		return domain.Expense{}, &ValidationError{Message: "有効なカテゴリを選択してください"}
	}

	if input.SpentAt == "" {
		return domain.Expense{}, &ValidationError{Message: "日付を入力してください"}
	}

	// 日付フォーマットの検証（RFC3339 をまず試し、失敗したら日付のみフォーマットを試す）
	var spentAt time.Time
	var err error
	spentAt, err = time.Parse(time.RFC3339, input.SpentAt)
	if err != nil {
		spentAt, err = time.Parse("2006-01-02", input.SpentAt)
		if err != nil {
			return domain.Expense{}, &ValidationError{Message: "日付の形式が正しくありません"}
		}
		spentAt = time.Date(spentAt.Year(), spentAt.Month(), spentAt.Day(), 0, 0, 0, 0, time.UTC)
	}
	if spentAt.IsZero() {
		return domain.Expense{}, &ValidationError{Message: "有効な日付を入力してください"}
	}

	if len(input.Memo) > MemoMaxLen {
		return domain.Expense{}, &ValidationError{Message: "メモは5000文字以内で入力してください"}
	}

	// Status は任意入力。指定された場合のみ正規化・検証する
	if input.Status != "" {
		if normalized, ok := domain.NormalizeStatus(input.Status); ok {
			input.Status = normalized
		} else {
			return domain.Expense{}, &ValidationError{Message: "ステータスは「予定」または「確定」を選択してください"}
		}
	}

	exists, err := uc.category.CategoryExists(context.Background(), int32(*input.CategoryID))
	if err != nil {
		return domain.Expense{}, &InternalError{Message: "internal error"}
	}
	if !exists {
		return domain.Expense{}, &ValidationError{Message: "カテゴリが存在しません"}
	}

	exp, err := uc.repo.CreateExpense(userID, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Expense{}, &NotFoundError{Message: "支出が見つかりません"}
		}
		lerr := strings.ToLower(err.Error())
		if strings.Contains(lerr, "foreign key") && (strings.Contains(lerr, "category") || strings.Contains(lerr, "category_id")) {
			return domain.Expense{}, &ValidationError{Message: "カテゴリが存在しません"}
		}
		return domain.Expense{}, &InternalError{Message: "internal error"}
	}

	return exp, nil
}
