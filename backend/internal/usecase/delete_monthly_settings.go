package usecase

import "context"

type deleteMonthlySettingsRepository interface {
	DeleteMonthlySetting(ctx context.Context, userID string, year, month int) error
}

type DeleteMonthlySettingsUseCase struct {
	repo deleteMonthlySettingsRepository
}

func NewDeleteMonthlySettingsUseCase(repo deleteMonthlySettingsRepository) *DeleteMonthlySettingsUseCase {
	return &DeleteMonthlySettingsUseCase{repo: repo}
}

func (uc *DeleteMonthlySettingsUseCase) Execute(ctx context.Context, userID string, year, month int) error {
	if year < 1 {
		return &ValidationError{Message: "年は1以上で入力してください"}
	}
	if month < 1 || month > 12 {
		return &ValidationError{Message: "月は1〜12の範囲で入力してください"}
	}

	return uc.repo.DeleteMonthlySetting(ctx, userID, year, month)
}
