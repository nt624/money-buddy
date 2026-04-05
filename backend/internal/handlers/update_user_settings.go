package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type updateUserSettingsUseCase interface {
	Execute(ctx context.Context, userID string, income int, savingGoal int) error
}

type UpdateUserSettingsHandler struct {
	uc updateUserSettingsUseCase
}

func NewUpdateUserSettingsHandler(uc updateUserSettingsUseCase) *UpdateUserSettingsHandler {
	return &UpdateUserSettingsHandler{uc: uc}
}

type updateUserSettingsRequest struct {
	Income     *int `json:"income"`
	SavingGoal *int `json:"saving_goal"`
}

func (h *UpdateUserSettingsHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	var req updateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
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

	err := h.uc.Execute(c.Request.Context(), userID, *req.Income, *req.SavingGoal)
	if err != nil {
		var validationErr *usecase.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー設定の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user settings updated successfully"})
}
