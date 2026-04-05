package repository

import (
	"context"
	"time"

	db "pace-wallet-backend/db/generated"
	"pace-wallet-backend/infra/transaction"
	"pace-wallet-backend/internal/domain"
)

type userRepositorySQLC struct {
	q *db.Queries
}

func NewUserRepositorySQLC(q *db.Queries) *userRepositorySQLC {
	return &userRepositorySQLC{q: q}
}

func (r *userRepositorySQLC) queries(ctx context.Context) *db.Queries {
	if tx, ok := transaction.TxFromContext(ctx); ok {
		return r.q.WithTx(tx)
	}
	return r.q
}

func (r *userRepositorySQLC) CreateUser(ctx context.Context, id string, income int, savingGoal int) error {
	params := db.CreateUserParams{
		ID:         id,
		Income:     int32(income),
		SavingGoal: int32(savingGoal),
	}
	return r.queries(ctx).CreateUser(ctx, params)
}

func (r *userRepositorySQLC) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	row, err := r.queries(ctx).GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return dbUserToModel(row), nil
}

func (r *userRepositorySQLC) UpdateUserSettings(ctx context.Context, id string, income int, savingGoal int) error {
	params := db.UpdateUserSettingsParams{
		ID:         id,
		Income:     int32(income),
		SavingGoal: int32(savingGoal),
	}
	return r.queries(ctx).UpdateUserSettings(ctx, params)
}

func dbUserToModel(u db.User) domain.User {
	createdAt := ""
	if u.CreatedAt.Valid {
		createdAt = u.CreatedAt.Time.Format(time.RFC3339)
	}
	updatedAt := ""
	if u.UpdatedAt.Valid {
		updatedAt = u.UpdatedAt.Time.Format(time.RFC3339)
	}

	return domain.User{
		ID:         u.ID,
		Income:     int(u.Income),
		SavingGoal: int(u.SavingGoal),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}
