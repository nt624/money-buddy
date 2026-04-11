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
	return uc.repo.DeleteMonthlySetting(ctx, userID, year, month)
}
