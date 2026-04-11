package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"pace-wallet-backend/internal/domain"
)

type updateCategoryRepository interface {
	UpdateCategory(ctx context.Context, userID string, id int32, name string) (domain.Category, error)
}

// UpdateCategoryUseCase はカテゴリ更新ユースケースです。
type UpdateCategoryUseCase struct {
	repo updateCategoryRepository
}

// NewUpdateCategoryUseCase は UpdateCategoryUseCase の新しいインスタンスを生成します。
func NewUpdateCategoryUseCase(repo updateCategoryRepository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{repo: repo}
}

// Execute はカテゴリ更新ユースケースを実行します。
func (uc *UpdateCategoryUseCase) Execute(ctx context.Context, userID string, id int, name string) (domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Category{}, &ValidationError{Message: "カテゴリ名を入力してください"}
	}
	if len([]rune(name)) > CategoryNameMaxLen {
		return domain.Category{}, &ValidationError{Message: "カテゴリ名は50文字以内で入力してください"}
	}

	cat, err := uc.repo.UpdateCategory(ctx, userID, int32(id), name)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateCategoryName) {
			return domain.Category{}, &ValidationError{Message: "同じ名前のカテゴリがすでに存在します"}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, &NotFoundError{Message: "カテゴリが見つかりません"}
		}
		return domain.Category{}, &InternalError{Message: "internal error"}
	}
	return cat, nil
}
