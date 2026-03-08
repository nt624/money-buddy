package repository

import (
	"context"

	db "money-buddy-backend/db/generated"
	"money-buddy-backend/internal/domain"
)

type categoryRepositorySQLC struct {
	q *db.Queries
}

func NewCategoryRepositorySQLC(q *db.Queries) *categoryRepositorySQLC {
	return &categoryRepositorySQLC{q: q}
}

func (r *categoryRepositorySQLC) ListCategories(ctx context.Context) ([]domain.Category, error) {
	items, err := r.q.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	var out []domain.Category
	for _, it := range items {
		out = append(out, dbCategoryToModel(it))
	}

	return out, nil
}

func dbCategoryToModel(c db.ListCategoriesRow) domain.Category {
	return domain.Category{
		ID:   int(c.ID),
		Name: c.Name,
	}
}

func (r *categoryRepositorySQLC) CategoryExists(ctx context.Context, id int32) (bool, error) {
	return r.q.CategoryExists(ctx, id)
}
