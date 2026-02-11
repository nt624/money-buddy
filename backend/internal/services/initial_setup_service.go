package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"money-buddy-backend/internal/models"
	"money-buddy-backend/internal/repositories"
)

type InitialSetupService interface {
	CompleteInitialSetup(ctx context.Context, userID string, income, savingGoal int, fixedCosts []models.FixedCostInput) error
}

type initialSetupService struct {
	userRepo      repositories.UserRepository
	fixedCostRepo repositories.FixedCostRepository
	txManager     TxManager
}

func NewInitialSetupService(userRepo repositories.UserRepository, fixedCostRepo repositories.FixedCostRepository, txManager TxManager) InitialSetupService {
	return &initialSetupService{
		userRepo:      userRepo,
		fixedCostRepo: fixedCostRepo,
		txManager:     txManager,
	}
}

func (s *initialSetupService) CompleteInitialSetup(ctx context.Context, userID string, income, savingGoal int, fixedCosts []models.FixedCostInput) error {
	if income <= 0 {
		return &ValidationError{Message: "収入は1円以上で入力してください"}
	}
	if savingGoal < 0 {
		return &ValidationError{Message: "貯金目標は0円以上で入力してください"}
	}

	// 固定費の正規化とバリデーション
	normalizedFixedCosts := make([]models.FixedCostInput, len(fixedCosts))
	for i, fc := range fixedCosts {
		// 名前を正規化（前後の空白を除去）
		trimmedName := strings.TrimSpace(fc.Name)

		if fc.Amount <= 0 {
			return &ValidationError{Message: "固定費の金額は1円以上で入力してください"}
		}
		if trimmedName == "" {
			return &ValidationError{Message: "固定費の名前を入力してください"}
		}
		if len(trimmedName) > FixedCostNameMaxLen {
			return &ValidationError{Message: "固定費の名前は100文字以内で入力してください"}
		}

		// 正規化された値を使用
		normalizedFixedCosts[i] = models.FixedCostInput{
			Name:   trimmedName,
			Amount: fc.Amount,
		}
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := tx.Context(ctx)

	user, err := s.userRepo.GetUserByID(txCtx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := s.userRepo.CreateUser(txCtx, userID, income, savingGoal); err != nil {
				_ = tx.Rollback()
				return err
			}
		} else {
			_ = tx.Rollback()
			return err
		}
	} else if user != (models.User{}) {
		if err := s.userRepo.UpdateUserSettings(txCtx, userID, income, savingGoal); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := s.fixedCostRepo.DeleteFixedCostsByUser(txCtx, userID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// 正規化された固定費を使用
	if err := s.fixedCostRepo.BulkCreateFixedCosts(txCtx, userID, normalizedFixedCosts); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
