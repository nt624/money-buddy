package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	db "pace-wallet-backend/db/generated"
	"pace-wallet-backend/internal/domain"
)




type categoryRepositorySQLC struct {
	q *db.Queries
}

func NewCategoryRepositorySQLC(q *db.Queries) *categoryRepositorySQLC {
	return &categoryRepositorySQLC{q: q}
}

func (r *categoryRepositorySQLC) ListCategories(ctx context.Context, userID string) ([]domain.Category, error) {
	items, err := r.q.ListUserCategories(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]domain.Category, 0, len(items))
	for _, it := range items {
		out = append(out, dbUserCategoryToModel(it))
	}
	return out, nil
}

func (r *categoryRepositorySQLC) CategoryExists(ctx context.Context, userID string, id int32) (bool, error) {
	return r.q.UserCategoryExists(ctx, db.UserCategoryExistsParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *categoryRepositorySQLC) CreateCategory(ctx context.Context, userID string, name string) (domain.Category, error) {
	row, err := r.q.CreateUserCategory(ctx, db.CreateUserCategoryParams{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Category{}, domain.ErrDuplicateCategoryName
		}
		return domain.Category{}, err
	}
	return domain.Category{
		ID:           int(row.ID),
		Name:         row.Name,
		CategoryType: row.CategoryType,
		SortOrder:    int(row.SortOrder),
	}, nil
}

func (r *categoryRepositorySQLC) UpdateCategory(ctx context.Context, userID string, id int32, name string) (domain.Category, error) {
	row, err := r.q.UpdateUserCategory(ctx, db.UpdateUserCategoryParams{
		ID:     id,
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Category{}, domain.ErrDuplicateCategoryName
		}
		return domain.Category{}, err
	}
	return domain.Category{
		ID:           int(row.ID),
		Name:         row.Name,
		CategoryType: row.CategoryType,
		SortOrder:    int(row.SortOrder),
	}, nil
}

func (r *categoryRepositorySQLC) UpdateCategorySortOrder(ctx context.Context, userID string, id int32, sortOrder int32) error {
	return r.q.UpdateUserCategorySortOrder(ctx, db.UpdateUserCategorySortOrderParams{
		ID:        id,
		UserID:    userID,
		SortOrder: sortOrder,
	})
}

func (r *categoryRepositorySQLC) DeleteCategory(ctx context.Context, userID string, id int32) error {
	return r.q.DeleteUserCategory(ctx, db.DeleteUserCategoryParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *categoryRepositorySQLC) CountExpensesByCategory(ctx context.Context, userID string, categoryID int32) (int64, error) {
	return r.q.CountExpensesByUserCategory(ctx, db.CountExpensesByUserCategoryParams{
		UserID:     userID,
		CategoryID: categoryID,
	})
}

func (r *categoryRepositorySQLC) SeedDefaultCategories(ctx context.Context, userID string) error {
	return r.q.SeedDefaultCategoriesForUser(ctx, userID)
}

func dbUserCategoryToModel(c db.ListUserCategoriesRow) domain.Category {
	return domain.Category{
		ID:           int(c.ID),
		Name:         c.Name,
		CategoryType: c.CategoryType,
		SortOrder:    int(c.SortOrder),
	}
}

// isUniqueViolation は PostgreSQL の一意制約違反エラー (23505) かどうかを判定します。
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
