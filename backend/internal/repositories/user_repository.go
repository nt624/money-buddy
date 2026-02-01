package repositories

import (
	"context"

	"money-buddy-backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, id string, income int, savingGoal int) error
	GetUserByID(ctx context.Context, id string) (models.User, error)
	UpdateUserSettings(ctx context.Context, id string, income int, savingGoal int) error
}
