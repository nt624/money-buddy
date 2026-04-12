package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "year と month のクエリパラメータは必須です"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year は1以上の整数で指定してください"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month は1〜12の整数で指定してください"})
		return
	}

	if err := h.uc.Execute(c.Request.Context(), userID, year, month); err != nil {
		var validationErr *usecase.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "月設定の削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "月設定を削除しました"})
}
