package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/middleware"
	"money-buddy-backend/internal/services"
)

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

type DashboardHandler struct {
	service services.DashboardService
}

func NewDashboardHandler(r gin.IRouter, service services.DashboardService) {
	h := &DashboardHandler{service: service}
	r.GET("/dashboard", h.GetDashboard)
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID := middleware.GetUserID(c)

	dashboard, err := h.service.GetDashboard(c.Request.Context(), userID)
	if err != nil {
		// ユーザーが存在しない場合
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
			return
		}
		// その他のエラー
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	// レスポンスを構築
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
