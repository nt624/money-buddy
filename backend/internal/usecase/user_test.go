package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pace-wallet-backend/internal/domain"
)

type mockGetUserRepo struct {
	getUserByIDFunc func(ctx context.Context, id string) (domain.User, error)
}

func (m *mockGetUserRepo) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return domain.User{}, errors.New("not implemented")
}

type mockUpdateUserRepo struct {
	updateUserSettingsFunc func(ctx context.Context, id string, income int, savingGoal int) error
}

func (m *mockUpdateUserRepo) UpdateUserSettings(ctx context.Context, id string, income int, savingGoal int) error {
	if m.updateUserSettingsFunc != nil {
		return m.updateUserSettingsFunc(ctx, id, income, savingGoal)
	}
	return errors.New("not implemented")
}

func TestGetUser_Success(t *testing.T) {
	expectedUser := domain.User{
		ID:         "test-user",
		Income:     300000,
		SavingGoal: 50000,
		CreatedAt:  "2024-01-01T00:00:00Z",
		UpdatedAt:  "2024-01-01T00:00:00Z",
	}

	repo := &mockGetUserRepo{
		getUserByIDFunc: func(ctx context.Context, id string) (domain.User, error) {
			assert.Equal(t, "test-user", id)
			return expectedUser, nil
		},
	}

	uc := NewGetUserUseCase(repo)
	user, err := uc.Execute(context.Background(), "test-user")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "test-user", user.ID)
	assert.Equal(t, 300000, user.Income)
	assert.Equal(t, 50000, user.SavingGoal)
}

func TestGetUser_NotFound(t *testing.T) {
	repo := &mockGetUserRepo{
		getUserByIDFunc: func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{}, errors.New("user not found")
		},
	}

	uc := NewGetUserUseCase(repo)
	user, err := uc.Execute(context.Background(), "non-existent-user")

	require.Error(t, err)
	require.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestGetUser_RepositoryError(t *testing.T) {
	repo := &mockGetUserRepo{
		getUserByIDFunc: func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{}, errors.New("database connection error")
		},
	}

	uc := NewGetUserUseCase(repo)
	user, err := uc.Execute(context.Background(), "test-user")

	require.Error(t, err)
	require.Nil(t, user)
	assert.Contains(t, err.Error(), "database connection error")
}

func TestUpdateUserSettings_Success(t *testing.T) {
	called := false
	repo := &mockUpdateUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			assert.Equal(t, "test-user", id)
			assert.Equal(t, 300000, income)
			assert.Equal(t, 50000, savingGoal)
			return nil
		},
	}

	uc := NewUpdateUserSettingsUseCase(repo)
	err := uc.Execute(context.Background(), "test-user", 300000, 50000)

	require.NoError(t, err)
	assert.True(t, called)
}

func TestUpdateUserSettings_InvalidIncome(t *testing.T) {
	testCases := []struct {
		name   string
		income int
	}{
		{"zero income", 0},
		{"negative income", -100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			repo := &mockUpdateUserRepo{
				updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
					called = true
					return nil
				},
			}

			uc := NewUpdateUserSettingsUseCase(repo)
			err := uc.Execute(context.Background(), "test-user", tc.income, 50000)

			require.Error(t, err)
			assert.False(t, called)
			assert.Contains(t, err.Error(), "収入は1円以上で入力してください")
		})
	}
}

func TestUpdateUserSettings_InvalidSavingGoal(t *testing.T) {
	called := false
	repo := &mockUpdateUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return nil
		},
	}

	uc := NewUpdateUserSettingsUseCase(repo)
	err := uc.Execute(context.Background(), "test-user", 300000, -100)

	require.Error(t, err)
	assert.False(t, called)
	assert.Contains(t, err.Error(), "貯金目標は0円以上で入力してください")
}

func TestUpdateUserSettings_IncomeExceedsLimit(t *testing.T) {
	called := false
	repo := &mockUpdateUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return nil
		},
	}

	uc := NewUpdateUserSettingsUseCase(repo)
	err := uc.Execute(context.Background(), "test-user", 1000000001, 50000)

	require.Error(t, err)
	assert.False(t, called)
	assert.Contains(t, err.Error(), "収入は10億円以下で入力してください")
}

func TestUpdateUserSettings_SavingGoalExceedsLimit(t *testing.T) {
	called := false
	repo := &mockUpdateUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return nil
		},
	}

	uc := NewUpdateUserSettingsUseCase(repo)
	err := uc.Execute(context.Background(), "test-user", 300000, 1000000001)

	require.Error(t, err)
	assert.False(t, called)
	assert.Contains(t, err.Error(), "貯金目標は10億円以下で入力してください")
}

func TestUpdateUserSettings_RepositoryError(t *testing.T) {
	called := false
	repo := &mockUpdateUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return errors.New("database connection error")
		},
	}

	uc := NewUpdateUserSettingsUseCase(repo)
	err := uc.Execute(context.Background(), "test-user", 300000, 50000)

	require.Error(t, err)
	assert.True(t, called)
	assert.Contains(t, err.Error(), "database connection error")
}
