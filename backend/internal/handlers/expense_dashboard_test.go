package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pace-wallet-backend/internal/domain"
	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

// setUserID は gin コンテキストにユーザーIDを設定するヘルパーです。
func setUserID(c *gin.Context, userID string) {
	c.Set(string(middleware.UserIDKey), userID)
}

// newRouterWithUserID は認証済みユーザーIDをコンテキストに設定するルーターを返します。
func newRouterWithUserID(userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		setUserID(c, userID)
		c.Next()
	})
	return router
}

// --- Create Expense ---

type mockCreateExpenseUC struct {
	fn func(userID string, input usecase.CreateExpenseInput) (domain.Expense, error)
}

func (m *mockCreateExpenseUC) Execute(userID string, input usecase.CreateExpenseInput) (domain.Expense, error) {
	if m.fn != nil {
		return m.fn(userID, input)
	}
	return domain.Expense{}, nil
}

func TestCreateExpenseHandler_Created(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockCreateExpenseUC{
		fn: func(userID string, input usecase.CreateExpenseInput) (domain.Expense, error) {
			return domain.Expense{
				ID:       1,
				Amount:   *input.Amount,
				Memo:     input.Memo,
				SpentAt:  input.SpentAt,
				Category: domain.Category{ID: *input.CategoryID},
			}, nil
		},
	}
	router.POST("/expenses", NewCreateExpenseHandler(uc).Handle)

	body := `{"amount":1000,"category_id":2,"memo":"lunch","spent_at":"2025-12-30"}`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]domain.Expense
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	expense, ok := resp["expense"]
	require.True(t, ok)
	require.Equal(t, 1, expense.ID)
	require.Equal(t, 1000, expense.Amount)
	require.Equal(t, 2, expense.Category.ID)
	require.Equal(t, "lunch", expense.Memo)
	require.Equal(t, "2025-12-30", expense.SpentAt)
}

func TestCreateExpenseHandler_ValidationError(t *testing.T) {
	t.Run("金額が0の場合", func(t *testing.T) {
		router := newRouterWithUserID("test-user")

		called := false
		amount := 0
		uc := &mockCreateExpenseUC{
			fn: func(userID string, input usecase.CreateExpenseInput) (domain.Expense, error) {
				called = true
				require.NotNil(t, input.Amount)
				require.Equal(t, 0, *input.Amount)
				return domain.Expense{}, &usecase.ValidationError{Message: "amount must be greater than 0"}
			},
		}
		router.POST("/expenses", NewCreateExpenseHandler(uc).Handle)

		body := fmt.Sprintf(`{"amount":%d,"category_id":2,"memo":"lunch","spent_at":"2025-12-30"}`, amount)
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.True(t, called)
		require.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "amount must be greater than 0", resp["error"])
	})
}

func TestCreateExpenseHandler_InternalError(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockCreateExpenseUC{
		fn: func(userID string, input usecase.CreateExpenseInput) (domain.Expense, error) {
			return domain.Expense{}, errors.New("db error")
		},
	}
	router.POST("/expenses", NewCreateExpenseHandler(uc).Handle)

	body := `{"amount":1000,"category_id":2,"memo":"test","spent_at":"2025-12-30"}`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- List Expenses ---

type mockListExpensesUC struct {
	fn func(userID string, filter usecase.MonthFilter) ([]domain.Expense, error)
}

func (m *mockListExpensesUC) Execute(ctx context.Context, userID string, filter usecase.MonthFilter) ([]domain.Expense, error) {
	if m.fn != nil {
		return m.fn(userID, filter)
	}
	return nil, nil
}

func TestListExpensesHandler_Success(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockListExpensesUC{
		fn: func(userID string, filter usecase.MonthFilter) ([]domain.Expense, error) {
			return []domain.Expense{{ID: 1, Amount: 500}, {ID: 2, Amount: 1000}}, nil
		},
	}
	router.GET("/expenses", NewListExpensesHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]domain.Expense
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp["expenses"], 2)
}

func TestListExpensesHandler_WithMonthFilter(t *testing.T) {
	router := newRouterWithUserID("test-user")

	var gotFilter usecase.MonthFilter
	uc := &mockListExpensesUC{
		fn: func(userID string, filter usecase.MonthFilter) ([]domain.Expense, error) {
			gotFilter = filter
			return []domain.Expense{{ID: 1, Amount: 500}}, nil
		},
	}
	router.GET("/expenses", NewListExpensesHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/expenses?year=2025&month=3", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2025, gotFilter.Year)
	assert.Equal(t, 3, gotFilter.Month)
}

func TestListExpensesHandler_InvalidYear(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.GET("/expenses", NewListExpensesHandler(&mockListExpensesUC{}).Handle)

	cases := []struct{ query string }{
		{"/expenses?year=abc&month=4"},
		{"/expenses?year=0&month=4"},
		{"/expenses?year=-1&month=4"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.query, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "query: %s", tc.query)
	}
}

func TestListExpensesHandler_InvalidMonth(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.GET("/expenses", NewListExpensesHandler(&mockListExpensesUC{}).Handle)

	cases := []struct{ query string }{
		{"/expenses?year=2025&month=0"},
		{"/expenses?year=2025&month=13"},
		{"/expenses?year=2025&month=abc"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.query, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "query: %s", tc.query)
	}
}

func TestListExpensesHandler_PartialFilter(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.GET("/expenses", NewListExpensesHandler(&mockListExpensesUC{}).Handle)

	cases := []struct{ query string }{
		{"/expenses?year=2025"},
		{"/expenses?month=4"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.query, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "query: %s", tc.query)
	}
}

// --- Delete Expense ---

type mockDeleteExpenseUC struct {
	fn func(userID string, id int) error
}

func (m *mockDeleteExpenseUC) Execute(userID string, id int) error {
	if m.fn != nil {
		return m.fn(userID, id)
	}
	return nil
}

func TestDeleteExpenseHandler_Success(t *testing.T) {
	router := newRouterWithUserID("test-user")

	called := false
	uc := &mockDeleteExpenseUC{
		fn: func(userID string, id int) error {
			called = true
			require.Equal(t, 42, id)
			return nil
		},
	}
	router.DELETE("/expenses/:id", NewDeleteExpenseHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodDelete, "/expenses/42", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, called)
}

func TestDeleteExpenseHandler_NotFound(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockDeleteExpenseUC{
		fn: func(userID string, id int) error {
			return &usecase.NotFoundError{Message: "支出が見つかりません"}
		},
	}
	router.DELETE("/expenses/:id", NewDeleteExpenseHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodDelete, "/expenses/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteExpenseHandler_InvalidID(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockDeleteExpenseUC{}
	router.DELETE("/expenses/:id", NewDeleteExpenseHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodDelete, "/expenses/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// --- Update Expense ---

type mockUpdateExpenseUC struct {
	fn func(userID string, input usecase.UpdateExpenseInput) (domain.Expense, error)
}

func (m *mockUpdateExpenseUC) Execute(userID string, input usecase.UpdateExpenseInput) (domain.Expense, error) {
	if m.fn != nil {
		return m.fn(userID, input)
	}
	return domain.Expense{}, nil
}

func TestUpdateExpenseHandler_Success(t *testing.T) {
	router := newRouterWithUserID("test-user")

	ret := domain.Expense{ID: 42, Amount: 700, Memo: "updated", SpentAt: "2025-07-01", Status: "confirmed", Category: domain.Category{ID: 5}}
	uc := &mockUpdateExpenseUC{
		fn: func(userID string, input usecase.UpdateExpenseInput) (domain.Expense, error) {
			return ret, nil
		},
	}
	router.PUT("/expenses/:id", NewUpdateExpenseHandler(uc).Handle)

	body := `{"amount":700,"category_id":5,"memo":"updated","spent_at":"2025-07-01","status":"confirmed"}`
	req := httptest.NewRequest(http.MethodPut, "/expenses/42", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]domain.Expense
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	exp := resp["expense"]
	require.Equal(t, 42, exp.ID)
	require.Equal(t, 700, exp.Amount)
}

func TestUpdateExpenseHandler_ValidationError(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockUpdateExpenseUC{
		fn: func(userID string, input usecase.UpdateExpenseInput) (domain.Expense, error) {
			return domain.Expense{}, &usecase.ValidationError{Message: "金額は1円以上で入力してください"}
		},
	}
	router.PUT("/expenses/:id", NewUpdateExpenseHandler(uc).Handle)

	body := `{"amount":0,"category_id":1,"memo":"x","spent_at":"2025-01-01"}`
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateExpenseHandler_StatusTransitionConflict(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockUpdateExpenseUC{
		fn: func(userID string, input usecase.UpdateExpenseInput) (domain.Expense, error) {
			return domain.Expense{}, usecase.ErrInvalidStatusTransition
		},
	}
	router.PUT("/expenses/:id", NewUpdateExpenseHandler(uc).Handle)

	body := `{"amount":1000,"category_id":1,"memo":"x","spent_at":"2025-01-01","status":"planned"}`
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}

// --- Get Dashboard ---

type mockGetDashboardUC struct {
	fn func(ctx context.Context, userID string) (*usecase.Dashboard, error)
}

func (m *mockGetDashboardUC) Execute(ctx context.Context, userID string) (*usecase.Dashboard, error) {
	if m.fn != nil {
		return m.fn(ctx, userID)
	}
	return nil, nil
}

func TestGetDashboardHandler_Success(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockGetDashboardUC{
		fn: func(ctx context.Context, userID string) (*usecase.Dashboard, error) {
			return &usecase.Dashboard{
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
	router.GET("/dashboard", NewGetDashboardHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp DashboardResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, int64(300000), resp.Income)
	require.Equal(t, int64(50000), resp.Remaining)
}

func TestGetDashboardHandler_ServiceError(t *testing.T) {
	router := newRouterWithUserID("test-user")

	uc := &mockGetDashboardUC{
		fn: func(ctx context.Context, userID string) (*usecase.Dashboard, error) {
			return nil, errors.New("db error")
		},
	}
	router.GET("/dashboard", NewGetDashboardHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
