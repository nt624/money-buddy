package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/domain"
	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type upsertMonthlySettingsUseCase interface {
	Execute(ctx context.Context, userID string, year, month, income, savingGoal int) (*domain.MonthlySettings, error)
}

type UpsertMonthlySettingsHandler struct {
	uc upsertMonthlySettingsUseCase
}

func NewUpsertMonthlySettingsHandler(uc upsertMonthlySettingsUseCase) *UpsertMonthlySettingsHandler {
	return &UpsertMonthlySettingsHandler{uc: uc}
}

type upsertMonthlySettingsRequest struct {
	Year       *int `json:"year"`
	Month      *int `json:"month"`
	Income     *int `json:"income"`
	SavingGoal *int `json:"saving_goal"`
}

func (h *UpsertMonthlySettingsHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	var req upsertMonthlySettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("upsert_monthly_settings: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	if req.Year == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year is required"})
		return
	}
	if req.Month == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month is required"})
		return
	}
	if req.Income == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "income is required"})
		return
	}
	if req.SavingGoal == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "saving_goal is required"})
		return
	}

	ms, err := h.uc.Execute(c.Request.Context(), userID, *req.Year, *req.Month, *req.Income, *req.SavingGoal)
	if err != nil {
		var validationErr *usecase.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "月設定の保存に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"year":        ms.Year,
		"month":       ms.Month,
		"income":      ms.Income,
		"saving_goal": ms.SavingGoal,
		"is_custom":   true,
	})
}
