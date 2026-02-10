package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"money-buddy-backend/internal/models"
	"money-buddy-backend/internal/services"
)

type userServiceMock struct {
	GetUserByIDFunc        func(ctx context.Context, userID string) (*models.User, error)
	UpdateUserSettingsFunc func(ctx context.Context, userID string, income int, savingGoal int) error
}

func (m *userServiceMock) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *userServiceMock) UpdateUserSettings(ctx context.Context, userID string, income int, savingGoal int) error {
	if m.UpdateUserSettingsFunc != nil {
		return m.UpdateUserSettingsFunc(ctx, userID, income, savingGoal)
	}
	return nil
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

func TestUpdateUserSettingsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	called := false
	svc := &userServiceMock{
		UpdateUserSettingsFunc: func(ctx context.Context, userID string, income int, savingGoal int) error {
			called = true
			require.Equal(t, DummyUserID, userID)
			require.Equal(t, 300000, income)
			require.Equal(t, 50000, savingGoal)
			return nil
		},
	}
	NewUserHandler(router, svc)

	body := `{"income": 300000, "saving_goal": 50000}`
	req := httptest.NewRequest(http.MethodPut, "/user/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, called, "service method should be called")

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "user settings updated successfully", resp["message"])
}

func TestUpdateUserSettingsHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	called := false
	svc := &userServiceMock{
		UpdateUserSettingsFunc: func(ctx context.Context, userID string, income int, savingGoal int) error {
			called = true
			return nil
		},
	}
	NewUserHandler(router, svc)

	body := `{"income": "invalid", "saving_goal": 50000}`
	req := httptest.NewRequest(http.MethodPut, "/user/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, called, "service should not be called for invalid JSON")
}

func TestUpdateUserSettingsHandler_MissingFields(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		errorMessage string
	}{
		{
			name:         "missing income",
			body:         `{"saving_goal": 50000}`,
			errorMessage: "income is required",
		},
		{
			name:         "missing saving_goal",
			body:         `{"income": 300000}`,
			errorMessage: "saving_goal is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			called := false
			svc := &userServiceMock{
				UpdateUserSettingsFunc: func(ctx context.Context, userID string, income int, savingGoal int) error {
					called = true
					return nil
				},
			}
			NewUserHandler(router, svc)

			req := httptest.NewRequest(http.MethodPut, "/user/me", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.False(t, called, "service should not be called")

			var resp map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			require.Equal(t, tc.errorMessage, resp["error"])
		})
	}
}

func TestUpdateUserSettingsHandler_ValidationError(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		errorMessage string
	}{
		{
			name:         "income is 0",
			body:         `{"income": 0, "saving_goal": 50000}`,
			errorMessage: "income must be greater than 0",
		},
		{
			name:         "negative income",
			body:         `{"income": -100, "saving_goal": 50000}`,
			errorMessage: "income must be greater than 0",
		},
		{
			name:         "income exceeds limit",
			body:         `{"income": 1000000001, "saving_goal": 50000}`,
			errorMessage: "income must be 10億 or less",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			called := false
			svc := &userServiceMock{
				UpdateUserSettingsFunc: func(ctx context.Context, userID string, income int, savingGoal int) error {
					called = true
					return &services.ValidationError{Message: tc.errorMessage}
				},
			}
			NewUserHandler(router, svc)

			req := httptest.NewRequest(http.MethodPut, "/user/me", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.True(t, called, "service should be called")

			var resp map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			require.Equal(t, tc.errorMessage, resp["error"])
		})
	}
}

func TestUpdateUserSettingsHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	called := false
	svc := &userServiceMock{
		UpdateUserSettingsFunc: func(ctx context.Context, userID string, income int, savingGoal int) error {
			called = true
			return errors.New("database connection error")
		},
	}
	NewUserHandler(router, svc)

	body := `{"income": 300000, "saving_goal": 50000}`
	req := httptest.NewRequest(http.MethodPut, "/user/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, called, "service method should be called")

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "failed to update user settings", resp["error"])
}
