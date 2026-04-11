package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type getMonthlySettingsUseCase interface {
	Execute(ctx context.Context, userID string, year, month int) (*usecase.MonthlySettingsResult, error)
}

type GetMonthlySettingsHandler struct {
	uc getMonthlySettingsUseCase
}

func NewGetMonthlySettingsHandler(uc getMonthlySettingsUseCase) *GetMonthlySettingsHandler {
	return &GetMonthlySettingsHandler{uc: uc}
}

func (h *GetMonthlySettingsHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month query parameters are required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year must be a valid integer"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month must be a valid integer"})
		return
	}

	result, err := h.uc.Execute(c.Request.Context(), userID, year, month)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"year":        result.Year,
		"month":       result.Month,
		"income":      result.Income,
		"saving_goal": result.SavingGoal,
		"is_custom":   result.IsCustom,
	})
}
