package usecase

import (
	"context"
	"strings"

	"money-buddy-backend/internal/domain"
)

type createFixedCostRepository interface {
	CreateFixedCost(ctx context.Context, userID string, name string, amount int) (domain.FixedCost, error)
}

// CreateFixedCostUseCase は固定費作成ユースケースです。
type CreateFixedCostUseCase struct {
	repo createFixedCostRepository
}

// NewCreateFixedCostUseCase は CreateFixedCostUseCase の新しいインスタンスを生成します。
func NewCreateFixedCostUseCase(repo createFixedCostRepository) *CreateFixedCostUseCase {
	return &CreateFixedCostUseCase{repo: repo}
}

// Execute は固定費作成ユースケースを実行します。
func (uc *CreateFixedCostUseCase) Execute(ctx context.Context, userID string, name string, amount int) (domain.FixedCost, error) {
	name = strings.TrimSpace(name)

	if err := validateFixedCostInput(name, amount); err != nil {
		return domain.FixedCost{}, err
	}

	return uc.repo.CreateFixedCost(ctx, userID, name, amount)
}
