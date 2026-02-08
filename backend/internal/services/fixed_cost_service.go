package services

import (
	"context"
	"strings"

	"money-buddy-backend/internal/models"
	"money-buddy-backend/internal/repositories"
)

const (
	// FixedCostNameMaxLen は固定費名の最大長
	FixedCostNameMaxLen = 100
)

type FixedCostService interface {
	ListFixedCosts(ctx context.Context, userID string) ([]models.FixedCost, error)
	UpdateFixedCost(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error)
	DeleteFixedCost(ctx context.Context, userID string, id int) error
}

type fixedCostService struct {
	repo repositories.FixedCostRepository
}

func NewFixedCostService(repo repositories.FixedCostRepository) FixedCostService {
	return &fixedCostService{repo: repo}
}

func (s *fixedCostService) ListFixedCosts(ctx context.Context, userID string) ([]models.FixedCost, error) {
	return s.repo.ListFixedCostsByUser(ctx, userID)
}

func (s *fixedCostService) UpdateFixedCost(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error) {
	// バリデーション
	if err := validateFixedCostInput(name, amount); err != nil {
		return models.FixedCost{}, err
	}

	// 更新実行
	if err := s.repo.UpdateFixedCost(ctx, int32(id), userID, name, amount); err != nil {
		return models.FixedCost{}, err
	}

	// 更新後のデータを取得して返す
	fixedCosts, err := s.repo.ListFixedCostsByUser(ctx, userID)
	if err != nil {
		return models.FixedCost{}, err
	}

	// 更新した固定費を探す
	for _, fc := range fixedCosts {
		if fc.ID == id {
			return fc, nil
		}
	}

	return models.FixedCost{}, &NotFoundError{Message: "fixed cost not found after update"}
}

func (s *fixedCostService) DeleteFixedCost(ctx context.Context, userID string, id int) error {
	return s.repo.DeleteFixedCost(ctx, int32(id), userID)
}

// validateFixedCostInput は固定費の入力バリデーションを行います
func validateFixedCostInput(name string, amount int) error {
	// 名前チェック
	name = strings.TrimSpace(name)
	if name == "" {
		return &ValidationError{Message: "name is required"}
	}
	if len(name) > FixedCostNameMaxLen {
		return &ValidationError{Message: "name is too long"}
	}

	// 金額チェック
	if amount <= 0 {
		return &ValidationError{Message: "amount must be greater than 0"}
	}
	if amount > BusinessMaxAmount {
		return &ValidationError{Message: "amount exceeds maximum allowed"}
	}

	return nil
}
