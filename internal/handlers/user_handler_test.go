package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"money-buddy-backend/internal/models"
)

type userServiceMock struct {
	GetUserByIDFunc func(ctx context.Context, userID string) (*models.User, error)
}

func (m *userServiceMock) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, userID)
	}
	return nil, nil
}

func TestUserHandler_GetCurrentUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &userServiceMock{
		GetUserByIDFunc: func(ctx context.Context, userID string) (*models.User, error) {
			require.Equal(t, DummyUserID, userID)
			return &models.User{
				ID:         "test-user",
				Income:     300000,
				SavingGoal: 50000,
				CreatedAt:  "2024-01-01T00:00:00Z",
				UpdatedAt:  "2024-01-01T00:00:00Z",
			}, nil
		},
	}
	NewUserHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var user models.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(t, err)
	require.Equal(t, "test-user", user.ID)
	require.Equal(t, 300000, user.Income)
	require.Equal(t, 50000, user.SavingGoal)
}

func TestUserHandler_GetCurrentUser_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &userServiceMock{
		GetUserByIDFunc: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}
	NewUserHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "user not found", resp["error"])
}

func TestUserHandler_GetCurrentUser_RepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &userServiceMock{
		GetUserByIDFunc: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, errors.New("database connection error")
		},
	}
	NewUserHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "user not found", resp["error"])
}
