package usecase

import (
	"context"
	"errors"
	"strings"

	"pace-wallet-backend/internal/domain"
)


const CategoryNameMaxLen = 50

type createCategoryRepository interface {
	CreateCategory(ctx context.Context, userID string, name string) (domain.Category, error)
}

// CreateCategoryUseCase はカテゴリ作成ユースケースです。
type CreateCategoryUseCase struct {
	repo createCategoryRepository
}

// NewCreateCategoryUseCase は CreateCategoryUseCase の新しいインスタンスを生成します。
func NewCreateCategoryUseCase(repo createCategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{repo: repo}
}

// Execute はカテゴリ作成ユースケースを実行します。
func (uc *CreateCategoryUseCase) Execute(ctx context.Context, userID string, name string) (domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Category{}, &ValidationError{Message: "カテゴリ名を入力してください"}
	}
	if len([]rune(name)) > CategoryNameMaxLen {
		return domain.Category{}, &ValidationError{Message: "カテゴリ名は50文字以内で入力してください"}
	}

	cat, err := uc.repo.CreateCategory(ctx, userID, name)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateCategoryName) {
			return domain.Category{}, &ValidationError{Message: "同じ名前のカテゴリがすでに存在します"}
		}
		return domain.Category{}, &InternalError{Message: "internal error"}
	}
	return cat, nil
}
