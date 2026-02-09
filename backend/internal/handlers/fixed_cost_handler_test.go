package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"money-buddy-backend/internal/models"
	"money-buddy-backend/internal/services"
)

// fixedCostServiceMock is a mock implementing services.FixedCostService
type fixedCostServiceMock struct {
	ListFixedCostsFunc  func(ctx context.Context, userID string) ([]models.FixedCost, error)
	UpdateFixedCostFunc func(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error)
	DeleteFixedCostFunc func(ctx context.Context, userID string, id int) error
}

func (m *fixedCostServiceMock) ListFixedCosts(ctx context.Context, userID string) ([]models.FixedCost, error) {
	if m.ListFixedCostsFunc != nil {
		return m.ListFixedCostsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *fixedCostServiceMock) UpdateFixedCost(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error) {
	if m.UpdateFixedCostFunc != nil {
		return m.UpdateFixedCostFunc(ctx, userID, id, name, amount)
	}
	return models.FixedCost{}, nil
}

func (m *fixedCostServiceMock) DeleteFixedCost(ctx context.Context, userID string, id int) error {
	if m.DeleteFixedCostFunc != nil {
		return m.DeleteFixedCostFunc(ctx, userID, id)
	}
	return nil
}

// TestListFixedCosts_Success は固定費一覧取得の成功ケースをテストします
func TestListFixedCosts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{
		ListFixedCostsFunc: func(ctx context.Context, userID string) ([]models.FixedCost, error) {
			return []models.FixedCost{
				{ID: 1, UserID: userID, Name: "家賃", Amount: 80000},
				{ID: 2, UserID: userID, Name: "通信費", Amount: 5000},
			}, nil
		},
	}
	NewFixedCostHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/fixed-costs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]models.FixedCost
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	fixedCosts, ok := resp["fixed_costs"]
	require.True(t, ok)
	require.Len(t, fixedCosts, 2)
	require.Equal(t, "家賃", fixedCosts[0].Name)
	require.Equal(t, 80000, fixedCosts[0].Amount)
}

// TestListFixedCosts_ServiceError はサービスエラーをテストします
func TestListFixedCosts_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{
		ListFixedCostsFunc: func(ctx context.Context, userID string) ([]models.FixedCost, error) {
			return nil, &services.InternalError{Message: "database error"}
		},
	}
	NewFixedCostHandler(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/fixed-costs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestUpdateFixedCost_Success は固定費更新の成功ケースをテストします
func TestUpdateFixedCost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{
		UpdateFixedCostFunc: func(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error) {
			return models.FixedCost{
				ID:     id,
				UserID: userID,
				Name:   name,
				Amount: amount,
			}, nil
		},
	}
	NewFixedCostHandler(router, svc)

	body := `{"name":"家賃（更新）","amount":85000}`
	req := httptest.NewRequest(http.MethodPut, "/fixed-costs/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]models.FixedCost
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	fixedCost, ok := resp["fixed_cost"]
	require.True(t, ok)
	require.Equal(t, 1, fixedCost.ID)
	require.Equal(t, "家賃（更新）", fixedCost.Name)
	require.Equal(t, 85000, fixedCost.Amount)
}

// TestUpdateFixedCost_InvalidID はIDが不正な場合をテストします
func TestUpdateFixedCost_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{}
	NewFixedCostHandler(router, svc)

	body := `{"name":"家賃","amount":80000}`
	req := httptest.NewRequest(http.MethodPut, "/fixed-costs/abc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateFixedCost_ValidationError はバリデーションエラーをテストします
func TestUpdateFixedCost_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{
		UpdateFixedCostFunc: func(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error) {
			return models.FixedCost{}, &services.ValidationError{Message: "amount must be greater than 0"}
		},
	}
	NewFixedCostHandler(router, svc)

	body := `{"name":"家賃","amount":0}`
	req := httptest.NewRequest(http.MethodPut, "/fixed-costs/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateFixedCost_NotFound は固定費が見つからない場合をテストします
func TestUpdateFixedCost_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{
		UpdateFixedCostFunc: func(ctx context.Context, userID string, id int, name string, amount int) (models.FixedCost, error) {
			return models.FixedCost{}, &services.NotFoundError{Message: "fixed cost not found"}
		},
	}
	NewFixedCostHandler(router, svc)

	body := `{"name":"家賃","amount":80000}`
	req := httptest.NewRequest(http.MethodPut, "/fixed-costs/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// TestDeleteFixedCost_Success は固定費削除の成功ケースをテストします
func TestDeleteFixedCost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{
		DeleteFixedCostFunc: func(ctx context.Context, userID string, id int) error {
			return nil
		},
	}
	NewFixedCostHandler(router, svc)

	req := httptest.NewRequest(http.MethodDelete, "/fixed-costs/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

// TestDeleteFixedCost_InvalidID はIDが不正な場合をテストします
func TestDeleteFixedCost_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{}
	NewFixedCostHandler(router, svc)

	req := httptest.NewRequest(http.MethodDelete, "/fixed-costs/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeleteFixedCost_NotFound は固定費が見つからない場合をテストします
func TestDeleteFixedCost_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &fixedCostServiceMock{
		DeleteFixedCostFunc: func(ctx context.Context, userID string, id int) error {
			return &services.NotFoundError{Message: "fixed cost not found"}
		},
	}
	NewFixedCostHandler(router, svc)

	req := httptest.NewRequest(http.MethodDelete, "/fixed-costs/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}
