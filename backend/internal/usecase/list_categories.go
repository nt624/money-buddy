package usecase

import (
	"context"

	"pace-wallet-backend/internal/domain"
)

type listCategoriesRepository interface {
	ListCategories(ctx context.Context, userID string) ([]domain.Category, error)
}

// ListCategoriesUseCase はカテゴリ一覧取得ユースケースです。
type ListCategoriesUseCase struct {
	repo listCategoriesRepository
}

// NewListCategoriesUseCase は ListCategoriesUseCase の新しいインスタンスを生成します。
func NewListCategoriesUseCase(repo listCategoriesRepository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{repo: repo}
}

// Execute はカテゴリ一覧取得ユースケースを実行します。
func (uc *ListCategoriesUseCase) Execute(ctx context.Context, userID string) ([]domain.Category, error) {
	return uc.repo.ListCategories(ctx, userID)
}
