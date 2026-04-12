package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pace-wallet-backend/internal/domain"
)

// --- shared mocks ---

type mockMonthlySettingsRepo struct {
	getMonthlySetting    func(ctx context.Context, userID string, year, month int) (*domain.MonthlySettings, error)
	getUserByID          func(ctx context.Context, id string) (domain.User, error)
	upsertMonthlySetting func(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error)
	deleteMonthlySetting func(ctx context.Context, userID string, year, month int) error
}

func (m *mockMonthlySettingsRepo) GetMonthlySetting(ctx context.Context, userID string, year, month int) (*domain.MonthlySettings, error) {
	if m.getMonthlySetting != nil {
		return m.getMonthlySetting(ctx, userID, year, month)
	}
	return nil, errors.New("not implemented")
}

func (m *mockMonthlySettingsRepo) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	if m.getUserByID != nil {
		return m.getUserByID(ctx, id)
	}
	return domain.User{}, errors.New("not implemented")
}

func (m *mockMonthlySettingsRepo) UpsertMonthlySetting(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error) {
	if m.upsertMonthlySetting != nil {
		return m.upsertMonthlySetting(ctx, userID, year, month, income, savingGoal)
	}
	return nil, errors.New("not implemented")
}

func (m *mockMonthlySettingsRepo) DeleteMonthlySetting(ctx context.Context, userID string, year, month int) error {
	if m.deleteMonthlySetting != nil {
		return m.deleteMonthlySetting(ctx, userID, year, month)
	}
	return errors.New("not implemented")
}

// --- GetMonthlySettings ---

func TestGetMonthlySettings_CustomSettingExists(t *testing.T) {
	repo := &mockMonthlySettingsRepo{
		getMonthlySetting: func(_ context.Context, _ string, year, month int) (*domain.MonthlySettings, error) {
			return &domain.MonthlySettings{Year: year, Month: month, Income: 400000, SavingGoal: 60000}, nil
		},
	}
	uc := NewGetMonthlySettingsUseCase(repo)
	result, err := uc.Execute(context.Background(), "user1", 2026, 4)
	require.NoError(t, err)
	assert.True(t, result.IsCustom)
	assert.Equal(t, 400000, result.Income)
	assert.Equal(t, 60000, result.SavingGoal)
}

func TestGetMonthlySettings_FallbackToDefault(t *testing.T) {
	repo := &mockMonthlySettingsRepo{
		getMonthlySetting: func(_ context.Context, _ string, _, _ int) (*domain.MonthlySettings, error) {
			return nil, nil
		},
		getUserByID: func(_ context.Context, id string) (domain.User, error) {
			return domain.User{Income: 300000, SavingGoal: 50000}, nil
		},
	}
	uc := NewGetMonthlySettingsUseCase(repo)
	result, err := uc.Execute(context.Background(), "user1", 2026, 4)
	require.NoError(t, err)
	assert.False(t, result.IsCustom)
	assert.Equal(t, 300000, result.Income)
	assert.Equal(t, 50000, result.SavingGoal)
}

func TestGetMonthlySettings_InvalidMonth(t *testing.T) {
	repo := &mockMonthlySettingsRepo{}
	uc := NewGetMonthlySettingsUseCase(repo)

	_, err := uc.Execute(context.Background(), "user1", 2026, 0)
	var ve *ValidationError
	require.ErrorAs(t, err, &ve)

	_, err = uc.Execute(context.Background(), "user1", 2026, 13)
	require.ErrorAs(t, err, &ve)
}

func TestGetMonthlySettings_InvalidYear(t *testing.T) {
	repo := &mockMonthlySettingsRepo{}
	uc := NewGetMonthlySettingsUseCase(repo)

	_, err := uc.Execute(context.Background(), "user1", 0, 4)
	var ve *ValidationError
	require.ErrorAs(t, err, &ve)
}

// --- UpsertMonthlySettings ---

func TestUpsertMonthlySettings_Success(t *testing.T) {
	repo := &mockMonthlySettingsRepo{
		upsertMonthlySetting: func(_ context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error) {
			return &domain.MonthlySettings{UserID: userID, Year: year, Month: month, Income: income, SavingGoal: savingGoal}, nil
		},
	}
	uc := NewUpsertMonthlySettingsUseCase(repo)
	ms, err := uc.Execute(context.Background(), "user1", 2026, 4, 300000, 50000)
	require.NoError(t, err)
	assert.Equal(t, 300000, ms.Income)
}

func TestUpsertMonthlySettings_Validation(t *testing.T) {
	repo := &mockMonthlySettingsRepo{}
	uc := NewUpsertMonthlySettingsUseCase(repo)

	cases := []struct {
		year, month, income, savingGoal int
	}{
		{0, 4, 300000, 50000},   // invalid year
		{2026, 0, 300000, 50000},  // invalid month (0)
		{2026, 13, 300000, 50000}, // invalid month (13)
		{2026, 4, 0, 50000},       // income <= 0
		{2026, 4, 300000, -1},     // saving_goal < 0
	}
	for _, c := range cases {
		_, err := uc.Execute(context.Background(), "user1", c.year, c.month, c.income, c.savingGoal)
		var ve *ValidationError
		assert.ErrorAs(t, err, &ve, "expected ValidationError for %+v", c)
	}
}

// --- DeleteMonthlySettings ---

func TestDeleteMonthlySettings_Success(t *testing.T) {
	called := false
	repo := &mockMonthlySettingsRepo{
		deleteMonthlySetting: func(_ context.Context, _ string, _, _ int) error {
			called = true
			return nil
		},
	}
	uc := NewDeleteMonthlySettingsUseCase(repo)
	err := uc.Execute(context.Background(), "user1", 2026, 4)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestDeleteMonthlySettings_InvalidMonth(t *testing.T) {
	repo := &mockMonthlySettingsRepo{}
	uc := NewDeleteMonthlySettingsUseCase(repo)

	err := uc.Execute(context.Background(), "user1", 2026, 0)
	var ve *ValidationError
	require.ErrorAs(t, err, &ve)
}

func TestDeleteMonthlySettings_InvalidYear(t *testing.T) {
	repo := &mockMonthlySettingsRepo{}
	uc := NewDeleteMonthlySettingsUseCase(repo)

	err := uc.Execute(context.Background(), "user1", 0, 4)
	var ve *ValidationError
	require.ErrorAs(t, err, &ve)
}
