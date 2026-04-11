package usecase

import (
	"context"
)

type deleteCategoryRepository interface {
	CountExpensesByCategory(ctx context.Context, userID string, categoryID int32) (int64, error)
	DeleteCategory(ctx context.Context, userID string, id int32) error
}

// DeleteCategoryUseCase はカテゴリ削除ユースケースです。
type DeleteCategoryUseCase struct {
	repo deleteCategoryRepository
}

// NewDeleteCategoryUseCase は DeleteCategoryUseCase の新しいインスタンスを生成します。
func NewDeleteCategoryUseCase(repo deleteCategoryRepository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{repo: repo}
}

// Execute はカテゴリ削除ユースケースを実行します。
func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, userID string, id int) error {
	count, err := uc.repo.CountExpensesByCategory(ctx, userID, int32(id))
	if err != nil {
		return &InternalError{Message: "internal error"}
	}
	if count > 0 {
		return &ValidationError{Message: "このカテゴリは支出で使用されているため削除できません"}
	}

	return uc.repo.DeleteCategory(ctx, userID, int32(id))
}
