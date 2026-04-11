package usecase

import (
	"context"
)

type reorderCategoriesRepository interface {
	UpdateCategorySortOrder(ctx context.Context, userID string, id int32, sortOrder int32) error
}

// ReorderCategoriesInput は並び替え入力DTOです。
type ReorderCategoriesInput struct {
	ID        int `json:"id"`
	SortOrder int `json:"sort_order"`
}

// ReorderCategoriesUseCase はカテゴリ並び替えユースケースです。
type ReorderCategoriesUseCase struct {
	repo reorderCategoriesRepository
}

// NewReorderCategoriesUseCase は ReorderCategoriesUseCase の新しいインスタンスを生成します。
func NewReorderCategoriesUseCase(repo reorderCategoriesRepository) *ReorderCategoriesUseCase {
	return &ReorderCategoriesUseCase{repo: repo}
}

// Execute はカテゴリ並び替えユースケースを実行します。
func (uc *ReorderCategoriesUseCase) Execute(ctx context.Context, userID string, items []ReorderCategoriesInput) error {
	for _, item := range items {
		if err := uc.repo.UpdateCategorySortOrder(ctx, userID, int32(item.ID), int32(item.SortOrder)); err != nil {
			return &InternalError{Message: "internal error"}
		}
	}
	return nil
}
