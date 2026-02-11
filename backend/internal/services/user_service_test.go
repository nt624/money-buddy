package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money-buddy-backend/internal/models"
)

type mockUserRepo struct {
	getUserByIDFunc        func(ctx context.Context, id string) (models.User, error)
	updateUserSettingsFunc func(ctx context.Context, id string, income int, savingGoal int) error
}

func (m *mockUserRepo) CreateUser(ctx context.Context, id string, income int, savingGoal int) error {
	return errors.New("not implemented")
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id string) (models.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return models.User{}, errors.New("not implemented")
}

func (m *mockUserRepo) UpdateUserSettings(ctx context.Context, id string, income int, savingGoal int) error {
	if m.updateUserSettingsFunc != nil {
		return m.updateUserSettingsFunc(ctx, id, income, savingGoal)
	}
	return errors.New("not implemented")
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	expectedUser := models.User{
		ID:         "test-user",
		Income:     300000,
		SavingGoal: 50000,
		CreatedAt:  "2024-01-01T00:00:00Z",
		UpdatedAt:  "2024-01-01T00:00:00Z",
	}

	repo := &mockUserRepo{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			assert.Equal(t, "test-user", id)
			return expectedUser, nil
		},
	}

	service := NewUserService(repo)
	user, err := service.GetUserByID(context.Background(), "test-user")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "test-user", user.ID)
	assert.Equal(t, 300000, user.Income)
	assert.Equal(t, 50000, user.SavingGoal)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			return models.User{}, errors.New("user not found")
		},
	}

	service := NewUserService(repo)
	user, err := service.GetUserByID(context.Background(), "non-existent-user")

	require.Error(t, err)
	require.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserService_GetUserByID_RepositoryError(t *testing.T) {
	repo := &mockUserRepo{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			return models.User{}, errors.New("database connection error")
		},
	}

	service := NewUserService(repo)
	user, err := service.GetUserByID(context.Background(), "test-user")

	require.Error(t, err)
	require.Nil(t, user)
	assert.Contains(t, err.Error(), "database connection error")
}

func TestUpdateUserSettings_Success(t *testing.T) {
	called := false
	repo := &mockUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			assert.Equal(t, "test-user", id)
			assert.Equal(t, 300000, income)
			assert.Equal(t, 50000, savingGoal)
			return nil
		},
	}

	service := NewUserService(repo)
	err := service.UpdateUserSettings(context.Background(), "test-user", 300000, 50000)

	require.NoError(t, err)
	assert.True(t, called, "repository method should be called")
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
			repo := &mockUserRepo{
				updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
					called = true
					return nil
				},
			}

			service := NewUserService(repo)
			err := service.UpdateUserSettings(context.Background(), "test-user", tc.income, 50000)

			require.Error(t, err)
			assert.False(t, called, "repository should not be called for invalid input")
			assert.Contains(t, err.Error(), "収入は1円以上で入力してください")
		})
	}
}

func TestUpdateUserSettings_InvalidSavingGoal(t *testing.T) {
	called := false
	repo := &mockUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return nil
		},
	}

	service := NewUserService(repo)
	err := service.UpdateUserSettings(context.Background(), "test-user", 300000, -100)

	require.Error(t, err)
	assert.False(t, called, "repository should not be called for invalid input")
	assert.Contains(t, err.Error(), "貯金目標は0円以上で入力してください")
}

func TestUpdateUserSettings_IncomeExceedsLimit(t *testing.T) {
	called := false
	repo := &mockUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return nil
		},
	}

	service := NewUserService(repo)
	err := service.UpdateUserSettings(context.Background(), "test-user", 1000000001, 50000)

	require.Error(t, err)
	assert.False(t, called, "repository should not be called for invalid input")
	assert.Contains(t, err.Error(), "収入は10億円以下で入力してください")
}

func TestUpdateUserSettings_SavingGoalExceedsLimit(t *testing.T) {
	called := false
	repo := &mockUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return nil
		},
	}

	service := NewUserService(repo)
	err := service.UpdateUserSettings(context.Background(), "test-user", 300000, 1000000001)

	require.Error(t, err)
	assert.False(t, called, "repository should not be called for invalid input")
	assert.Contains(t, err.Error(), "貯金目標は10億円以下で入力してください")
}

func TestUpdateUserSettings_RepositoryError(t *testing.T) {
	called := false
	repo := &mockUserRepo{
		updateUserSettingsFunc: func(ctx context.Context, id string, income int, savingGoal int) error {
			called = true
			return errors.New("database connection error")
		},
	}

	service := NewUserService(repo)
	err := service.UpdateUserSettings(context.Background(), "test-user", 300000, 50000)

	require.Error(t, err)
	assert.True(t, called, "repository method should be called")
	assert.Contains(t, err.Error(), "database connection error")
}
