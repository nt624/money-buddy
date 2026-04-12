package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pace-wallet-backend/internal/domain"
	"pace-wallet-backend/internal/usecase"
)

// --- GetMonthlySettings handler ---

type mockGetMonthlySettingsUC struct {
	fn func(ctx context.Context, userID string, year, month int) (*usecase.MonthlySettingsResult, error)
}

func (m *mockGetMonthlySettingsUC) Execute(ctx context.Context, userID string, year, month int) (*usecase.MonthlySettingsResult, error) {
	if m.fn != nil {
		return m.fn(ctx, userID, year, month)
	}
	return &usecase.MonthlySettingsResult{Year: year, Month: month, Income: 300000, SavingGoal: 50000, IsCustom: false}, nil
}

func TestGetMonthlySettingsHandler_Success(t *testing.T) {
	router := newRouterWithUserID("test-user")
	uc := &mockGetMonthlySettingsUC{
		fn: func(_ context.Context, _ string, year, month int) (*usecase.MonthlySettingsResult, error) {
			return &usecase.MonthlySettingsResult{Year: year, Month: month, Income: 400000, SavingGoal: 60000, IsCustom: true}, nil
		},
	}
	router.GET("/monthly-settings", NewGetMonthlySettingsHandler(uc).Handle)

	req := httptest.NewRequest(http.MethodGet, "/monthly-settings?year=2026&month=4", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(400000), resp["income"])
	assert.Equal(t, true, resp["is_custom"])
}

func TestGetMonthlySettingsHandler_MissingParams(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.GET("/monthly-settings", NewGetMonthlySettingsHandler(&mockGetMonthlySettingsUC{}).Handle)

	for _, query := range []string{"/monthly-settings", "/monthly-settings?year=2026", "/monthly-settings?month=4"} {
		req := httptest.NewRequest(http.MethodGet, query, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "query: %s", query)
	}
}

func TestGetMonthlySettingsHandler_InvalidParams(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.GET("/monthly-settings", NewGetMonthlySettingsHandler(&mockGetMonthlySettingsUC{}).Handle)

	cases := []string{
		"/monthly-settings?year=abc&month=4",
		"/monthly-settings?year=2026&month=abc",
		"/monthly-settings?year=0&month=4",
		"/monthly-settings?year=2026&month=0",
		"/monthly-settings?year=2026&month=13",
	}
	for _, query := range cases {
		req := httptest.NewRequest(http.MethodGet, query, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "query: %s", query)
	}
}

// --- UpsertMonthlySettings handler ---

type mockUpsertMonthlySettingsUC struct {
	fn func(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error)
}

func (m *mockUpsertMonthlySettingsUC) Execute(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error) {
	if m.fn != nil {
		return m.fn(ctx, userID, year, month, income, savingGoal)
	}
	return &domain.MonthlySettings{UserID: userID, Year: year, Month: month, Income: income, SavingGoal: savingGoal}, nil
}

func TestUpsertMonthlySettingsHandler_Success(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.PUT("/monthly-settings", NewUpsertMonthlySettingsHandler(&mockUpsertMonthlySettingsUC{}).Handle)

	body, _ := json.Marshal(map[string]int{"year": 2026, "month": 4, "income": 300000, "saving_goal": 50000})
	req := httptest.NewRequest(http.MethodPut, "/monthly-settings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestUpsertMonthlySettingsHandler_MissingFields(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.PUT("/monthly-settings", NewUpsertMonthlySettingsHandler(&mockUpsertMonthlySettingsUC{}).Handle)

	bodies := []string{
		`{"month":4,"income":300000,"saving_goal":50000}`,
		`{"year":2026,"income":300000,"saving_goal":50000}`,
		`{"year":2026,"month":4,"saving_goal":50000}`,
		`{"year":2026,"month":4,"income":300000}`,
	}
	for _, b := range bodies {
		req := httptest.NewRequest(http.MethodPut, "/monthly-settings", bytes.NewBufferString(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "body: %s", b)
	}
}

// --- DeleteMonthlySettings handler ---

type mockDeleteMonthlySettingsUC struct {
	fn func(ctx context.Context, userID string, year, month int) error
}

func (m *mockDeleteMonthlySettingsUC) Execute(ctx context.Context, userID string, year, month int) error {
	if m.fn != nil {
		return m.fn(ctx, userID, year, month)
	}
	return nil
}

func TestDeleteMonthlySettingsHandler_Success(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.DELETE("/monthly-settings", NewDeleteMonthlySettingsHandler(&mockDeleteMonthlySettingsUC{}).Handle)

	req := httptest.NewRequest(http.MethodDelete, "/monthly-settings?year=2026&month=4", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["message"])
}

func TestDeleteMonthlySettingsHandler_InvalidParams(t *testing.T) {
	router := newRouterWithUserID("test-user")
	router.DELETE("/monthly-settings", NewDeleteMonthlySettingsHandler(&mockDeleteMonthlySettingsUC{}).Handle)

	cases := []string{
		"/monthly-settings",
		"/monthly-settings?year=2026",
		"/monthly-settings?month=4",
		"/monthly-settings?year=0&month=4",
		"/monthly-settings?year=2026&month=13",
		"/monthly-settings?year=abc&month=4",
	}
	for _, query := range cases {
		req := httptest.NewRequest(http.MethodDelete, query, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "query: %s", query)
	}
}
