package usecase

import (
	"context"
	"errors"

	"pace-wallet-backend/internal/domain"
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
	repo      reorderCategoriesRepository
	txManager TxManager
}

// NewReorderCategoriesUseCase は ReorderCategoriesUseCase の新しいインスタンスを生成します。
func NewReorderCategoriesUseCase(repo reorderCategoriesRepository, txManager TxManager) *ReorderCategoriesUseCase {
	return &ReorderCategoriesUseCase{repo: repo, txManager: txManager}
}

// Execute はカテゴリ並び替えユースケースを実行します。
// 複数の UPDATE を1トランザクションで実行し、部分更新を防ぎます。
func (uc *ReorderCategoriesUseCase) Execute(ctx context.Context, userID string, items []ReorderCategoriesInput) error {
	tx, err := uc.txManager.Begin(ctx)
	if err != nil {
		return &InternalError{Message: "internal error"}
	}
	txCtx := tx.Context(ctx)

	for _, item := range items {
		if err := uc.repo.UpdateCategorySortOrder(txCtx, userID, int32(item.ID), int32(item.SortOrder)); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, domain.ErrCategoryNotFound) {
				return &NotFoundError{Message: "カテゴリが見つかりません"}
			}
			return &InternalError{Message: "internal error"}
		}
	}

	if err := tx.Commit(); err != nil {
		return &InternalError{Message: "internal error"}
	}
	return nil
}
