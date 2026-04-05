package usecase

import (
	"context"

	"pace-wallet-backend/internal/domain"
)

type listFixedCostsRepository interface {
	ListFixedCostsByUser(ctx context.Context, userID string) ([]domain.FixedCost, error)
}

// ListFixedCostsUseCase は固定費一覧取得ユースケースです。
type ListFixedCostsUseCase struct {
	repo listFixedCostsRepository
}

// NewListFixedCostsUseCase は ListFixedCostsUseCase の新しいインスタンスを生成します。
func NewListFixedCostsUseCase(repo listFixedCostsRepository) *ListFixedCostsUseCase {
	return &ListFixedCostsUseCase{repo: repo}
}

// Execute は固定費一覧取得ユースケースを実行します。
func (uc *ListFixedCostsUseCase) Execute(ctx context.Context, userID string) ([]domain.FixedCost, error) {
	return uc.repo.ListFixedCostsByUser(ctx, userID)
}
