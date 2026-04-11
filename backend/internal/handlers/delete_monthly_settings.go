package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
)

type deleteMonthlySettingsUseCase interface {
	Execute(ctx context.Context, userID string, year, month int) error
}

type DeleteMonthlySettingsHandler struct {
	uc deleteMonthlySettingsUseCase
}

func NewDeleteMonthlySettingsHandler(uc deleteMonthlySettingsUseCase) *DeleteMonthlySettingsHandler {
	return &DeleteMonthlySettingsHandler{uc: uc}
}

func (h *DeleteMonthlySettingsHandler) Handle(c *gin.Context) {
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

	if err := h.uc.Execute(c.Request.Context(), userID, year, month); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "月設定の削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "monthly settings deleted successfully"})
}
