package repositories

import (
	"context"

	"money-buddy-backend/internal/models"
)

type FixedCostRepository interface {
	CreateFixedCost(ctx context.Context, userID string, name string, amount int) (models.FixedCost, error)
	ListFixedCostsByUser(ctx context.Context, userID string) ([]models.FixedCost, error)
	DeleteFixedCostsByUser(ctx context.Context, userID string) error
	BulkCreateFixedCosts(ctx context.Context, userID string, fixedCosts []models.FixedCostInput) error
	UpdateFixedCost(ctx context.Context, id int32, userID string, name string, amount int) error
	DeleteFixedCost(ctx context.Context, id int32, userID string) error
}
