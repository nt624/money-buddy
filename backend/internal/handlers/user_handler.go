package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(r gin.IRouter, service services.UserService) {
	h := &UserHandler{service: service}
	r.GET("/user/me", h.GetCurrentUser)
	r.PUT("/user/me", h.UpdateUserSettings)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateUserSettingsRequest struct {
	Income     *int `json:"income"`
	SavingGoal *int `json:"saving_goal"`
}

func (h *UserHandler) UpdateUserSettings(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	var req UpdateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	// 必須フィールドのチェック
	if req.Income == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "income is required"})
		return
	}
	if req.SavingGoal == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "saving_goal is required"})
		return
	}

	err := h.service.UpdateUserSettings(c.Request.Context(), userID, *req.Income, *req.SavingGoal)
	if err != nil {
		var validationErr *services.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー設定の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user settings updated successfully"})
}
