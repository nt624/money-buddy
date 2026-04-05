package usecase

import (
	"context"
	"strings"

	"pace-wallet-backend/internal/domain"
)

type updateFixedCostRepository interface {
	UpdateFixedCost(ctx context.Context, id int32, userID string, name string, amount int) error
	ListFixedCostsByUser(ctx context.Context, userID string) ([]domain.FixedCost, error)
}

// UpdateFixedCostUseCase は固定費更新ユースケースです。
type UpdateFixedCostUseCase struct {
	repo updateFixedCostRepository
}

// NewUpdateFixedCostUseCase は UpdateFixedCostUseCase の新しいインスタンスを生成します。
func NewUpdateFixedCostUseCase(repo updateFixedCostRepository) *UpdateFixedCostUseCase {
	return &UpdateFixedCostUseCase{repo: repo}
}

// Execute は固定費更新ユースケースを実行します。
func (uc *UpdateFixedCostUseCase) Execute(ctx context.Context, userID string, id int, name string, amount int) (domain.FixedCost, error) {
	name = strings.TrimSpace(name)

	if err := validateFixedCostInput(name, amount); err != nil {
		return domain.FixedCost{}, err
	}

	if err := uc.repo.UpdateFixedCost(ctx, int32(id), userID, name, amount); err != nil {
		return domain.FixedCost{}, err
	}

	fixedCosts, err := uc.repo.ListFixedCostsByUser(ctx, userID)
	if err != nil {
		return domain.FixedCost{}, err
	}

	for _, fc := range fixedCosts {
		if fc.ID == id {
			return fc, nil
		}
	}

	return domain.FixedCost{}, &NotFoundError{Message: "固定費が見つかりません"}
}

// validateFixedCostInput は固定費の入力バリデーションを行います
func validateFixedCostInput(name string, amount int) error {
	if name == "" {
		return &ValidationError{Message: "名前を入力してください"}
	}
	if len(name) > FixedCostNameMaxLen {
		return &ValidationError{Message: "名前は100文字以内で入力してください"}
	}

	if amount <= 0 {
		return &ValidationError{Message: "金額は1円以上で入力してください"}
	}
	if amount > BusinessMaxAmount {
		return &ValidationError{Message: "金額は10億円以下で入力してください"}
	}

	return nil
}
