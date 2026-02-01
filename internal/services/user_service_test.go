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
	getUserByIDFunc func(ctx context.Context, id string) (models.User, error)
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
