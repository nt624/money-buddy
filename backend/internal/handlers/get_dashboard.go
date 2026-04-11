package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type getDashboardUseCase interface {
	Execute(ctx context.Context, userID string, year, month int) (*usecase.Dashboard, error)
}

type GetDashboardHandler struct {
	uc getDashboardUseCase
}

func NewGetDashboardHandler(uc getDashboardUseCase) *GetDashboardHandler {
	return &GetDashboardHandler{uc: uc}
}

// DashboardResponse はダッシュボードAPIのレスポンス構造です。
type DashboardResponse struct {
	Income            int64 `json:"income"`
	SavingGoal        int64 `json:"saving_goal"`
	FixedCosts        int64 `json:"fixed_costs"`
	VariableBudget    int64 `json:"variable_budget"`
	ConfirmedExpenses int64 `json:"confirmed_expenses"`
	PlannedExpenses   int64 `json:"planned_expenses"`
	Remaining         int64 `json:"remaining"`
}

func (h *GetDashboardHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if yearStr := c.Query("year"); yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "year must be a valid integer"})
			return
		}
		year = y
	}
	if monthStr := c.Query("month"); monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "month must be a valid integer"})
			return
		}
		month = m
	}

	dashboard, err := h.uc.Execute(c.Request.Context(), userID, year, month)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	response := DashboardResponse{
		Income:            dashboard.Income,
		SavingGoal:        dashboard.SavingGoal,
		FixedCosts:        dashboard.FixedCosts,
		VariableBudget:    dashboard.VariableBudget,
		ConfirmedExpenses: dashboard.ConfirmedExpenses,
		PlannedExpenses:   dashboard.PlannedExpenses,
		Remaining:         dashboard.Remaining,
	}

	c.JSON(http.StatusOK, response)
}
