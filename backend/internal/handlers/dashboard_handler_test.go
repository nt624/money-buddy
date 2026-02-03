package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money-buddy-backend/internal/services"
)

// dashboardServiceMock は DashboardService のモック実装です
type dashboardServiceMock struct {
	GetDashboardFunc func(ctx context.Context, userID string) (*services.Dashboard, error)
}

func (m *dashboardServiceMock) GetDashboard(ctx context.Context, userID string) (*services.Dashboard, error) {
	if m.GetDashboardFunc != nil {
		return m.GetDashboardFunc(ctx, userID)
	}
	return nil, nil
}

// TestDashboardHandler_GetDashboard_Success は正常系のテストです
func TestDashboardHandler_GetDashboard_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &dashboardServiceMock{
		GetDashboardFunc: func(ctx context.Context, userID string) (*services.Dashboard, error) {
			require.Equal(t, DummyUserID, userID)
			return &services.Dashboard{
				Income:            300000,
				SavingGoal:        50000,
				FixedCosts:        100000,
				VariableBudget:    150000,
				ConfirmedExpenses: 80000,
				PlannedExpenses:   20000,
				Remaining:         50000,
			}, nil
		},
	}
	NewDashboardHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp DashboardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, int64(300000), resp.Income)
	assert.Equal(t, int64(50000), resp.SavingGoal)
	assert.Equal(t, int64(100000), resp.FixedCosts)
	assert.Equal(t, int64(150000), resp.VariableBudget)
	assert.Equal(t, int64(80000), resp.ConfirmedExpenses)
	assert.Equal(t, int64(20000), resp.PlannedExpenses)
	assert.Equal(t, int64(50000), resp.Remaining)
}

// TestDashboardHandler_GetDashboard_UserNotFound はユーザーが存在しない場合のテストです
func TestDashboardHandler_GetDashboard_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &dashboardServiceMock{
		GetDashboardFunc: func(ctx context.Context, userID string) (*services.Dashboard, error) {
			return nil, sql.ErrNoRows
		},
	}
	NewDashboardHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "user not found", resp["error"])
}

// TestDashboardHandler_GetDashboard_ServiceError はサービス層でエラーが発生した場合のテストです
func TestDashboardHandler_GetDashboard_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &dashboardServiceMock{
		GetDashboardFunc: func(ctx context.Context, userID string) (*services.Dashboard, error) {
			return nil, errors.New("database connection error")
		},
	}
	NewDashboardHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "internal server error", resp["error"])
}

// TestDashboardHandler_GetDashboard_NegativeRemaining は残額がマイナスの場合のテストです
func TestDashboardHandler_GetDashboard_NegativeRemaining(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &dashboardServiceMock{
		GetDashboardFunc: func(ctx context.Context, userID string) (*services.Dashboard, error) {
			return &services.Dashboard{
				Income:            300000,
				SavingGoal:        50000,
				FixedCosts:        100000,
				VariableBudget:    150000,
				ConfirmedExpenses: 120000,
				PlannedExpenses:   80000,
				Remaining:         -50000, // マイナス
			}, nil
		},
	}
	NewDashboardHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp DashboardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// マイナスの残額も正しくレスポンスに含まれることを確認
	assert.Equal(t, int64(-50000), resp.Remaining)
}
