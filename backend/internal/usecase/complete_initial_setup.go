package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"pace-wallet-backend/internal/domain"
)

const (
	// FixedCostNameMaxLen は固定費名の最大長
	FixedCostNameMaxLen = 100
)

// FixedCostItem は初期設定時の固定費入力DTOです。
type FixedCostItem struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type completeInitialSetupUserRepository interface {
	CreateUser(ctx context.Context, id string, income int, savingGoal int) error
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	UpdateUserSettings(ctx context.Context, id string, income int, savingGoal int) error
}

type completeInitialSetupFixedCostRepository interface {
	DeleteFixedCostsByUser(ctx context.Context, userID string) error
	BulkCreateFixedCosts(ctx context.Context, userID string, fixedCosts []FixedCostItem) error
}

type completeInitialSetupCategoryRepository interface {
	SeedDefaultCategories(ctx context.Context, userID string) error
}

// CompleteInitialSetupUseCase は初期設定完了ユースケースです。
type CompleteInitialSetupUseCase struct {
	userRepo     completeInitialSetupUserRepository
	fixedCostRepo completeInitialSetupFixedCostRepository
	categoryRepo  completeInitialSetupCategoryRepository
	txManager     TxManager
}

// NewCompleteInitialSetupUseCase は CompleteInitialSetupUseCase の新しいインスタンスを生成します。
func NewCompleteInitialSetupUseCase(
	userRepo completeInitialSetupUserRepository,
	fixedCostRepo completeInitialSetupFixedCostRepository,
	categoryRepo completeInitialSetupCategoryRepository,
	txManager TxManager,
) *CompleteInitialSetupUseCase {
	return &CompleteInitialSetupUseCase{
		userRepo:      userRepo,
		fixedCostRepo: fixedCostRepo,
		categoryRepo:  categoryRepo,
		txManager:     txManager,
	}
}

// Execute は初期設定完了ユースケースを実行します。
func (uc *CompleteInitialSetupUseCase) Execute(ctx context.Context, userID string, income, savingGoal int, fixedCosts []FixedCostItem) error {
	if income <= 0 {
		return &ValidationError{Message: "収入は1円以上で入力してください"}
	}
	if income > BusinessMaxAmount {
		return &ValidationError{Message: "収入は10億円以下で入力してください"}
	}
	if savingGoal < 0 {
		return &ValidationError{Message: "貯金目標は0円以上で入力してください"}
	}
	if savingGoal > BusinessMaxAmount {
		return &ValidationError{Message: "貯金目標は10億円以下で入力してください"}
	}

	// 固定費の正規化とバリデーション
	normalizedFixedCosts := make([]FixedCostItem, len(fixedCosts))
	for i, fc := range fixedCosts {
		trimmedName := strings.TrimSpace(fc.Name)

		if fc.Amount <= 0 {
			return &ValidationError{Message: "固定費の金額は1円以上で入力してください"}
		}
		if fc.Amount > BusinessMaxAmount {
			return &ValidationError{Message: "固定費の金額は10億円以下で入力してください"}
		}
		if trimmedName == "" {
			return &ValidationError{Message: "固定費の名前を入力してください"}
		}
		if len(trimmedName) > FixedCostNameMaxLen {
			return &ValidationError{Message: "固定費の名前は100文字以内で入力してください"}
		}

		normalizedFixedCosts[i] = FixedCostItem{
			Name:   trimmedName,
			Amount: fc.Amount,
		}
	}

	tx, err := uc.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := tx.Context(ctx)

	isNewUser := false
	user, err := uc.userRepo.GetUserByID(txCtx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := uc.userRepo.CreateUser(txCtx, userID, income, savingGoal); err != nil {
				_ = tx.Rollback()
				return err
			}
			isNewUser = true
		} else {
			_ = tx.Rollback()
			return err
		}
	} else if user != (domain.User{}) {
		if err := uc.userRepo.UpdateUserSettings(txCtx, userID, income, savingGoal); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if isNewUser {
		if err := uc.categoryRepo.SeedDefaultCategories(txCtx, userID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := uc.fixedCostRepo.DeleteFixedCostsByUser(txCtx, userID); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := uc.fixedCostRepo.BulkCreateFixedCosts(txCtx, userID, normalizedFixedCosts); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
