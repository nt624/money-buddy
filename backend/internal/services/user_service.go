package services

import (
	"context"

	"money-buddy-backend/internal/models"
	"money-buddy-backend/internal/repositories"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	UpdateUserSettings(ctx context.Context, userID string, income int, savingGoal int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) UpdateUserSettings(ctx context.Context, userID string, income int, savingGoal int) error {
	// Validate income
	if income <= 0 {
		return &ValidationError{Message: "収入は1円以上で入力してください"}
	}
	if income > BusinessMaxAmount {
		return &ValidationError{Message: "収入は10億円以下で入力してください"}
	}

	// Validate saving goal
	if savingGoal < 0 {
		return &ValidationError{Message: "貯金目標は0円以上で入力してください"}
	}
	if savingGoal > BusinessMaxAmount {
		return &ValidationError{Message: "貯金目標は10億円以下で入力してください"}
	}

	// Update user settings
	return s.userRepo.UpdateUserSettings(ctx, userID, income, savingGoal)
}
