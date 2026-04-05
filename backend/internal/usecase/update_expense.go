package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"pace-wallet-backend/internal/domain"
)

// UpdateExpenseInput は支出更新の入力DTOです。
type UpdateExpenseInput struct {
	ID         int
	Amount     *int
	CategoryID *int
	Memo       string
	SpentAt    string
	Status     string
}

type updateExpenseRepository interface {
	GetExpenseByID(userID string, id int32) (domain.Expense, error)
	UpdateExpense(userID string, input UpdateExpenseInput) (domain.Expense, error)
}

type updateCategoryExistsChecker interface {
	CategoryExists(ctx context.Context, id int32) (bool, error)
}

// UpdateExpenseUseCase は支出更新ユースケースです。
type UpdateExpenseUseCase struct {
	repo     updateExpenseRepository
	category updateCategoryExistsChecker
}

// NewUpdateExpenseUseCase は UpdateExpenseUseCase の新しいインスタンスを生成します。
func NewUpdateExpenseUseCase(repo updateExpenseRepository, category updateCategoryExistsChecker) *UpdateExpenseUseCase {
	return &UpdateExpenseUseCase{repo: repo, category: category}
}

// Execute は支出更新ユースケースを実行します。
func (uc *UpdateExpenseUseCase) Execute(userID string, input UpdateExpenseInput) (domain.Expense, error) {
	// 金額チェック
	if input.Amount == nil {
		return domain.Expense{}, &ValidationError{Message: "金額を入力してください"}
	}
	if *input.Amount <= 0 {
		return domain.Expense{}, &ValidationError{Message: "金額は1円以上で入力してください"}
	}
	if *input.Amount > BusinessMaxAmount {
		return domain.Expense{}, &ValidationError{Message: "金額は10億円以下で入力してください"}
	}

	// カテゴリID チェック
	if input.CategoryID == nil {
		return domain.Expense{}, &ValidationError{Message: "カテゴリを選択してください"}
	}
	if *input.CategoryID <= 0 {
		return domain.Expense{}, &ValidationError{Message: "有効なカテゴリを選択してください"}
	}

	// SpentAt の非空チェック
	if input.SpentAt == "" {
		return domain.Expense{}, &ValidationError{Message: "日付を入力してください"}
	}

	// 日付フォーマットの検証
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

	// Memo 長チェック
	if len(input.Memo) > MemoMaxLen {
		return domain.Expense{}, &ValidationError{Message: "メモは5000文字以内で入力してください"}
	}

	// 現在の状態を取得し、ステータス遷移のバリデーションを行う
	current, err := uc.repo.GetExpenseByID(userID, int32(input.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Expense{}, ErrInvalidStatusTransition
		}
		return domain.Expense{}, &InternalError{Message: "internal error"}
	}

	// カテゴリ存在チェック（現在のExpense取得後に実施）
	exists, err := uc.category.CategoryExists(context.Background(), int32(*input.CategoryID))
	if err != nil {
		return domain.Expense{}, &InternalError{Message: "internal error"}
	}
	if !exists {
		return domain.Expense{}, &ValidationError{Message: "カテゴリが存在しません"}
	}

	// 変更後ステータスの決定（未指定なら現状維持）
	desiredStatus := input.Status
	if desiredStatus == "" {
		desiredStatus = current.Status
	} else {
		if normalized, ok := domain.NormalizeStatus(desiredStatus); ok {
			desiredStatus = normalized
		} else {
			return domain.Expense{}, &ValidationError{Message: "ステータスは「予定」または「確定」を選択してください"}
		}
	}

	// 遷移ルール: confirmed → planned は禁止
	if strings.ToLower(current.Status) == "confirmed" && desiredStatus == "planned" {
		return domain.Expense{}, ErrInvalidStatusTransition
	}

	// リポジトリに渡す前に正規化済みステータスをセット
	input.Status = desiredStatus
	return uc.repo.UpdateExpense(userID, input)
}
