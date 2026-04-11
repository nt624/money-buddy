package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "pace-wallet-backend/db/generated"
	"pace-wallet-backend/internal/domain"
)

type monthlySettingsRepositorySQLC struct {
	q *db.Queries
}

func NewMonthlySettingsRepositorySQLC(q *db.Queries) *monthlySettingsRepositorySQLC {
	return &monthlySettingsRepositorySQLC{q: q}
}

func (r *monthlySettingsRepositorySQLC) UpsertMonthlySetting(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error) {
	row, err := r.q.UpsertMonthlySetting(ctx, db.UpsertMonthlySettingParams{
		UserID:     userID,
		Year:       int32(year),
		Month:      int32(month),
		Income:     int32(income),
		SavingGoal: int32(savingGoal),
	})
	if err != nil {
		return nil, err
	}

	return dbMonthlySettingToModel(row), nil
}

func (r *monthlySettingsRepositorySQLC) GetMonthlySetting(ctx context.Context, userID string, year, month int) (*domain.MonthlySettings, error) {
	row, err := r.q.GetMonthlySetting(ctx, db.GetMonthlySettingParams{
		UserID: userID,
		Year:   int32(year),
		Month:  int32(month),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return dbMonthlySettingToModel(row), nil
}

func (r *monthlySettingsRepositorySQLC) DeleteMonthlySetting(ctx context.Context, userID string, year, month int) error {
	return r.q.DeleteMonthlySetting(ctx, db.DeleteMonthlySettingParams{
		UserID: userID,
		Year:   int32(year),
		Month:  int32(month),
	})
}

// monthlySettingsFullRepo は GetMonthlySettingsUseCase が要求する
// GetMonthlySetting と GetUserByID の両方を提供する複合リポジトリです。
type monthlySettingsFullRepo struct {
	ms   *monthlySettingsRepositorySQLC
	user *userRepositorySQLC
}

func NewMonthlySettingsFullRepo(ms *monthlySettingsRepositorySQLC, user *userRepositorySQLC) *monthlySettingsFullRepo {
	return &monthlySettingsFullRepo{ms: ms, user: user}
}

func (r *monthlySettingsFullRepo) GetMonthlySetting(ctx context.Context, userID string, year, month int) (*domain.MonthlySettings, error) {
	return r.ms.GetMonthlySetting(ctx, userID, year, month)
}

func (r *monthlySettingsFullRepo) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return r.user.GetUserByID(ctx, id)
}

func dbMonthlySettingToModel(m db.MonthlySetting) *domain.MonthlySettings {
	return &domain.MonthlySettings{
		ID:         int(m.ID),
		UserID:     m.UserID,
		Year:       int(m.Year),
		Month:      int(m.Month),
		Income:     int(m.Income),
		SavingGoal: int(m.SavingGoal),
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  m.UpdatedAt.Format(time.RFC3339),
	}
}
