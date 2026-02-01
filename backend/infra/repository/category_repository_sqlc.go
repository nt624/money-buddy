package repository

import (
	"context"

	db "money-buddy-backend/db/generated"
	"money-buddy-backend/internal/models"
	"money-buddy-backend/internal/repositories"
)

type categoryRepositorySQLC struct {
	q *db.Queries
}

func NewCategoryRepositorySQLC(q *db.Queries) repositories.CategoryRepository {
	return &categoryRepositorySQLC{q: q}
}

func (r *categoryRepositorySQLC) ListCategories(ctx context.Context) ([]models.Category, error) {
	items, err := r.q.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	var out []models.Category
	for _, it := range items {
		out = append(out, dbCategoryToModel(it))
	}

	return out, nil
}

func dbCategoryToModel(c db.ListCategoriesRow) models.Category {
	return models.Category{
		ID:   int(c.ID),
		Name: c.Name,
	}
}

func (r *categoryRepositorySQLC) CategoryExists(ctx context.Context, id int32) (bool, error) {
	return r.q.CategoryExists(ctx, id)
}
