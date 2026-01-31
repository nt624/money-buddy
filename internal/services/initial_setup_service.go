package services

import (
	"context"
	"database/sql"

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
		return &ValidationError{Message: "income must be greater than 0"}
	}
	if savingGoal < 0 {
		return &ValidationError{Message: "saving_goal must be greater than or equal to 0"}
	}
	for _, fc := range fixedCosts {
		if fc.Amount <= 0 {
			return &ValidationError{Message: "fixed_cost.amount must be greater than 0"}
		}
		if fc.Name == "" {
			return &ValidationError{Message: "fixed_cost.name must be provided"}
		}
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := s.userRepo.CreateUser(ctx, userID, income, savingGoal); err != nil {
				_ = tx.Rollback()
				return err
			}
		} else {
			_ = tx.Rollback()
			return err
		}
	} else if user != (models.User{}) {
		if err := s.userRepo.UpdateUserSettings(ctx, userID, income, savingGoal); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := s.fixedCostRepo.DeleteFixedCostsByUser(ctx, userID); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := s.fixedCostRepo.BulkCreateFixedCosts(ctx, userID, fixedCosts); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
