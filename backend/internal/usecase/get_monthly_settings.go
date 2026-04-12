package usecase

import (
	"context"

	"pace-wallet-backend/internal/domain"
)

type getMonthlySettingsRepository interface {
	GetMonthlySetting(ctx context.Context, userID string, year, month int) (*domain.MonthlySettings, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
}

// MonthlySettingsResult はユーザーへ返す月設定情報を表します。
// IsCustom が true の場合はその月専用の設定が存在し、false の場合はユーザーデフォルト値です。
type MonthlySettingsResult struct {
	Year       int
	Month      int
	Income     int
	SavingGoal int
	IsCustom   bool
}

type GetMonthlySettingsUseCase struct {
	repo getMonthlySettingsRepository
}

func NewGetMonthlySettingsUseCase(repo getMonthlySettingsRepository) *GetMonthlySettingsUseCase {
	return &GetMonthlySettingsUseCase{repo: repo}
}

func (uc *GetMonthlySettingsUseCase) Execute(ctx context.Context, userID string, year, month int) (*MonthlySettingsResult, error) {
	if year < 1 {
		return nil, &ValidationError{Message: "年は1以上で入力してください"}
	}
	if month < 1 || month > 12 {
		return nil, &ValidationError{Message: "月は1〜12の範囲で入力してください"}
	}

	ms, err := uc.repo.GetMonthlySetting(ctx, userID, year, month)
	if err != nil {
		return nil, err
	}

	if ms != nil {
		return &MonthlySettingsResult{
			Year:       ms.Year,
			Month:      ms.Month,
			Income:     ms.Income,
			SavingGoal: ms.SavingGoal,
			IsCustom:   true,
		}, nil
	}

	// 月設定がない場合はユーザーデフォルト値を返す
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &MonthlySettingsResult{
		Year:       year,
		Month:      month,
		Income:     user.Income,
		SavingGoal: user.SavingGoal,
		IsCustom:   false,
	}, nil
}
