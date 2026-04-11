package usecase

import (
	"context"

	"pace-wallet-backend/internal/domain"
)

type upsertMonthlySettingsRepository interface {
	UpsertMonthlySetting(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error)
}

type UpsertMonthlySettingsUseCase struct {
	repo upsertMonthlySettingsRepository
}

func NewUpsertMonthlySettingsUseCase(repo upsertMonthlySettingsRepository) *UpsertMonthlySettingsUseCase {
	return &UpsertMonthlySettingsUseCase{repo: repo}
}

func (uc *UpsertMonthlySettingsUseCase) Execute(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error) {
	if year < 1 {
		return nil, &ValidationError{Message: "年は1以上で入力してください"}
	}
	if month < 1 || month > 12 {
		return nil, &ValidationError{Message: "月は1〜12の範囲で入力してください"}
	}
	if income <= 0 {
		return nil, &ValidationError{Message: "収入は1円以上で入力してください"}
	}
	if income > BusinessMaxAmount {
		return nil, &ValidationError{Message: "収入は10億円以下で入力してください"}
	}
	if savingGoal < 0 {
		return nil, &ValidationError{Message: "貯金目標は0円以上で入力してください"}
	}
	if savingGoal > BusinessMaxAmount {
		return nil, &ValidationError{Message: "貯金目標は10億円以下で入力してください"}
	}

	return uc.repo.UpsertMonthlySetting(ctx, userID, year, month, income, savingGoal)
}
