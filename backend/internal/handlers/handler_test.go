package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"money-buddy-backend/internal/domain"
	"money-buddy-backend/internal/usecase"
)

// --- Create Fixed Cost ---

type mockCreateFixedCostUC struct {
	fn func(ctx context.Context, userID string, name string, amount int) (domain.FixedCost, error)
}

func (m *mockCreateFixedCostUC) Execute(ctx context.Context, userID string, name string, amount int) (domain.FixedCost, error) {
	if m.fn != nil {
		return m.fn(ctx, userID, name, amount)
	}
	return domain.FixedCost{}, nil
}

func TestCreateFixedCostHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockCreateFixedCostUC{
		fn: func(ctx context.Context, userID string, name string, amount int) (domain.FixedCost, error) {
			return domain.FixedCost{ID: 1, UserID: userID, Name: name, Amount: amount}, nil
		},
	}
	router.POST("/fixed-costs", NewCreateFixedCostHandler(uc).Handle)

	body := `{"name":"家賃","amount":80000}`
	req := httptest.NewRequest(http.MethodPost, "/fixed-costs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]domain.FixedCost
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	fc := resp["fixed_cost"]
	require.Equal(t, 1, fc.ID)
	require.Equal(t, "家賃", fc.Name)
	require.Equal(t, 80000, fc.Amount)
}

func TestCreateFixedCostHandler_ValidationError(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockCreateFixedCostUC{
		fn: func(ctx context.Context, userID string, name string, amount int) (domain.FixedCost, error) {
			return domain.FixedCost{}, &usecase.ValidationError{Message: "金額は1円以上で入力してください"}
		},
	}
	router.POST("/fixed-costs", NewCreateFixedCostHandler(uc).Handle)

	body := `{"name":"家賃","amount":0}`
	req := httptest.NewRequest(http.MethodPost, "/fixed-costs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// --- List Fixed Costs ---

type mockListFixedCostsUC struct {
	fn func(ctx context.Context, userID string) ([]domain.FixedCost, error)
}

func (m *mockListFixedCostsUC) Execute(ctx context.Context, userID string) ([]domain.FixedCost, error) {
	if m.fn != nil {
		return m.fn(ctx, userID)
	}
	return nil, nil
}

func TestListFixedCostsHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockListFixedCostsUC{
		fn: func(ctx context.Context, userID string) ([]domain.FixedCost, error) {
			return []domain.FixedCost{
				{ID: 1, Name: "家賃", Amount: 80000},
				{ID: 2, Name: "通信費", Amount: 5000},
			}, nil
		},
	}
	router.GET("/fixed-costs", NewListFixedCostsHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/fixed-costs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string][]domain.FixedCost
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp["fixed_costs"], 2)
}

// --- Update Fixed Cost ---

type mockUpdateFixedCostUC struct {
	fn func(ctx context.Context, userID string, id int, name string, amount int) (domain.FixedCost, error)
}

func (m *mockUpdateFixedCostUC) Execute(ctx context.Context, userID string, id int, name string, amount int) (domain.FixedCost, error) {
	if m.fn != nil {
		return m.fn(ctx, userID, id, name, amount)
	}
	return domain.FixedCost{}, nil
}

func TestUpdateFixedCostHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockUpdateFixedCostUC{
		fn: func(ctx context.Context, userID string, id int, name string, amount int) (domain.FixedCost, error) {
			return domain.FixedCost{ID: id, UserID: userID, Name: name, Amount: amount}, nil
		},
	}
	router.PUT("/fixed-costs/:id", NewUpdateFixedCostHandler(uc).Handle)

	body := `{"name":"家賃（更新）","amount":85000}`
	req := httptest.NewRequest(http.MethodPut, "/fixed-costs/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]domain.FixedCost
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	fc := resp["fixed_cost"]
	require.Equal(t, "家賃（更新）", fc.Name)
}

func TestUpdateFixedCostHandler_NotFound(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockUpdateFixedCostUC{
		fn: func(ctx context.Context, userID string, id int, name string, amount int) (domain.FixedCost, error) {
			return domain.FixedCost{}, &usecase.NotFoundError{Message: "固定費が見つかりません"}
		},
	}
	router.PUT("/fixed-costs/:id", NewUpdateFixedCostHandler(uc).Handle)

	body := `{"name":"家賃","amount":80000}`
	req := httptest.NewRequest(http.MethodPut, "/fixed-costs/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// --- Delete Fixed Cost ---

type mockDeleteFixedCostUC struct {
	fn func(ctx context.Context, userID string, id int) error
}

func (m *mockDeleteFixedCostUC) Execute(ctx context.Context, userID string, id int) error {
	if m.fn != nil {
		return m.fn(ctx, userID, id)
	}
	return nil
}

func TestDeleteFixedCostHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	called := false
	uc := &mockDeleteFixedCostUC{
		fn: func(ctx context.Context, userID string, id int) error {
			called = true
			require.Equal(t, 1, id)
			return nil
		},
	}
	router.DELETE("/fixed-costs/:id", NewDeleteFixedCostHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodDelete, "/fixed-costs/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, called)
}

func TestDeleteFixedCostHandler_NotFound(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockDeleteFixedCostUC{
		fn: func(ctx context.Context, userID string, id int) error {
			return &usecase.NotFoundError{Message: "固定費が見つかりません"}
		},
	}
	router.DELETE("/fixed-costs/:id", NewDeleteFixedCostHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodDelete, "/fixed-costs/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// --- Initial Setup ---

type mockCompleteInitialSetupUC struct {
	fn func(ctx context.Context, userID string, income, savingGoal int, fixedCosts []usecase.FixedCostItem) error
}

func (m *mockCompleteInitialSetupUC) Execute(ctx context.Context, userID string, income, savingGoal int, fixedCosts []usecase.FixedCostItem) error {
	if m.fn != nil {
		return m.fn(ctx, userID, income, savingGoal, fixedCosts)
	}
	return nil
}

func TestCompleteInitialSetupHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	called := false
	uc := &mockCompleteInitialSetupUC{
		fn: func(ctx context.Context, userID string, income, savingGoal int, fixedCosts []usecase.FixedCostItem) error {
			called = true
			require.Equal(t, 300000, income)
			return nil
		},
	}
	router.POST("/setup", NewCompleteInitialSetupHandler(uc).Handle)

	body := `{"income":300000,"savingGoal":50000,"fixedCosts":[{"name":"家賃","amount":80000}]}`
	req := httptest.NewRequest(http.MethodPost, "/setup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, called)
}

func TestCompleteInitialSetupHandler_ValidationError(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockCompleteInitialSetupUC{
		fn: func(ctx context.Context, userID string, income, savingGoal int, fixedCosts []usecase.FixedCostItem) error {
			return &usecase.ValidationError{Message: "収入は1円以上で入力してください"}
		},
	}
	router.POST("/setup", NewCompleteInitialSetupHandler(uc).Handle)

	body := `{"income":0,"savingGoal":0,"fixedCosts":[]}`
	req := httptest.NewRequest(http.MethodPost, "/setup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// --- Get User ---

type mockGetUserUC struct {
	fn func(ctx context.Context, userID string) (*domain.User, error)
}

func (m *mockGetUserUC) Execute(ctx context.Context, userID string) (*domain.User, error) {
	if m.fn != nil {
		return m.fn(ctx, userID)
	}
	return nil, nil
}

func TestGetUserHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockGetUserUC{
		fn: func(ctx context.Context, userID string) (*domain.User, error) {
			return &domain.User{ID: userID, Income: 300000, SavingGoal: 50000}, nil
		},
	}
	router.GET("/user/me", NewGetUserHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var user domain.User
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &user))
	require.Equal(t, "user1", user.ID)
}

func TestGetUserHandler_NotFound(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockGetUserUC{
		fn: func(ctx context.Context, userID string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}
	router.GET("/user/me", NewGetUserHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// --- Update User Settings ---

type mockUpdateUserSettingsUC struct {
	fn func(ctx context.Context, userID string, income int, savingGoal int) error
}

func (m *mockUpdateUserSettingsUC) Execute(ctx context.Context, userID string, income int, savingGoal int) error {
	if m.fn != nil {
		return m.fn(ctx, userID, income, savingGoal)
	}
	return nil
}

func TestUpdateUserSettingsHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	called := false
	uc := &mockUpdateUserSettingsUC{
		fn: func(ctx context.Context, userID string, income int, savingGoal int) error {
			called = true
			require.Equal(t, 300000, income)
			require.Equal(t, 50000, savingGoal)
			return nil
		},
	}
	router.PUT("/user/me", NewUpdateUserSettingsHandler(uc).Handle)

	body := `{"income":300000,"saving_goal":50000}`
	req := httptest.NewRequest(http.MethodPut, "/user/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, called)
}

func TestUpdateUserSettingsHandler_MissingIncome(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockUpdateUserSettingsUC{}
	router.PUT("/user/me", NewUpdateUserSettingsHandler(uc).Handle)

	body := `{"saving_goal":50000}`
	req := httptest.NewRequest(http.MethodPut, "/user/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "income is required", resp["error"])
}

// --- List Categories ---

type mockListCategoriesUC struct {
	fn func(ctx context.Context) ([]domain.Category, error)
}

func (m *mockListCategoriesUC) Execute(ctx context.Context) ([]domain.Category, error) {
	if m.fn != nil {
		return m.fn(ctx)
	}
	return nil, nil
}

func TestListCategoriesHandler_Success(t *testing.T) {
	router := newRouterWithUserID("user1")

	uc := &mockListCategoriesUC{
		fn: func(ctx context.Context) ([]domain.Category, error) {
			return []domain.Category{{ID: 1, Name: "食費"}, {ID: 2, Name: "交通費"}}, nil
		},
	}
	router.GET("/categories", NewListCategoriesHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string][]domain.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp["categories"], 2)
}
