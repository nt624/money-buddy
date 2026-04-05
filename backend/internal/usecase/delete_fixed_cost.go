package usecase

import (
	"context"

	"pace-wallet-backend/internal/domain"
)

type deleteFixedCostRepository interface {
	ListFixedCostsByUser(ctx context.Context, userID string) ([]domain.FixedCost, error)
	DeleteFixedCost(ctx context.Context, id int32, userID string) error
}

// DeleteFixedCostUseCase は固定費削除ユースケースです。
type DeleteFixedCostUseCase struct {
	repo deleteFixedCostRepository
}

// NewDeleteFixedCostUseCase は DeleteFixedCostUseCase の新しいインスタンスを生成します。
func NewDeleteFixedCostUseCase(repo deleteFixedCostRepository) *DeleteFixedCostUseCase {
	return &DeleteFixedCostUseCase{repo: repo}
}

// Execute は固定費削除ユースケースを実行します。
func (uc *DeleteFixedCostUseCase) Execute(ctx context.Context, userID string, id int) error {
	fixedCosts, err := uc.repo.ListFixedCostsByUser(ctx, userID)
	if err != nil {
		return err
	}

	found := false
	for _, fc := range fixedCosts {
		if fc.ID == id {
			found = true
			break
		}
	}

	if !found {
		return &NotFoundError{Message: "固定費が見つかりません"}
	}

	return uc.repo.DeleteFixedCost(ctx, int32(id), userID)
}
