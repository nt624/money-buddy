package usecase

import "context"

type updateUserRepository interface {
	UpdateUserSettings(ctx context.Context, id string, income int, savingGoal int) error
}

// UpdateUserSettingsUseCase はユーザー設定更新ユースケースです。
type UpdateUserSettingsUseCase struct {
	repo updateUserRepository
}

// NewUpdateUserSettingsUseCase は UpdateUserSettingsUseCase の新しいインスタンスを生成します。
func NewUpdateUserSettingsUseCase(repo updateUserRepository) *UpdateUserSettingsUseCase {
	return &UpdateUserSettingsUseCase{repo: repo}
}

// Execute はユーザー設定更新ユースケースを実行します。
func (uc *UpdateUserSettingsUseCase) Execute(ctx context.Context, userID string, income int, savingGoal int) error {
	// Validate income
	if income <= 0 {
		return &ValidationError{Message: "収入は1円以上で入力してください"}
	}
	if income > BusinessMaxAmount {
		return &ValidationError{Message: "収入は10億円以下で入力してください"}
	}

	// Validate saving goal
	if savingGoal < 0 {
		return &ValidationError{Message: "貯金目標は0円以上で入力してください"}
	}
	if savingGoal > BusinessMaxAmount {
		return &ValidationError{Message: "貯金目標は10億円以下で入力してください"}
	}

	return uc.repo.UpdateUserSettings(ctx, userID, income, savingGoal)
}
