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
	CreateFixedCost(ctx context.Context, userID string, name string, amount int) (models.FixedCost, error)
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

func (s *fixedCostService) CreateFixedCost(ctx context.Context, userID string, name string, amount int) (models.FixedCost, error) {
	// 名前を正規化（前後の空白を除去）
	name = strings.TrimSpace(name)

	// バリデーション
	if err := validateFixedCostInput(name, amount); err != nil {
		return models.FixedCost{}, err
	}

	// 作成実行
	return s.repo.CreateFixedCost(ctx, userID, name, amount)
}

func (s *fixedCostService) ListFixedCosts(ctx context.Context, userID string) ([]models.FixedCost, error) {
	return s.repo.ListFixedCostsByUser(ctx, userID)
}

func (s *fixedCostService) UpdateFixedCost(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error) {
	// 名前を正規化（前後の空白を除去）
	name = strings.TrimSpace(name)

	// バリデーション
	if err := validateFixedCostInput(name, amount); err != nil {
		return models.FixedCost{}, err
	}

	// 更新実行（トリム済みのnameを使用）
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
	// 削除前に対象が存在するか確認
	fixedCosts, err := s.repo.ListFixedCostsByUser(ctx, userID)
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
		return &NotFoundError{Message: "fixed cost not found"}
	}

	// 削除実行
	return s.repo.DeleteFixedCost(ctx, int32(id), userID)
}

// validateFixedCostInput は固定費の入力バリデーションを行います
func validateFixedCostInput(name string, amount int) error {
	// 名前チェック（呼び出し側で既にTrimSpaceされていることを前提）
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
