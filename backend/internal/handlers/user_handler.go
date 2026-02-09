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

func NewUserHandler(r *gin.Engine, service services.UserService) {
	h := &UserHandler{service: service}
	r.GET("/user/me", h.GetCurrentUser)
	r.PUT("/user/me", h.UpdateUserSettings)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateUserSettingsRequest struct {
	Income     int `json:"income" binding:"required"`
	SavingGoal int `json:"saving_goal" binding:"required"`
}

func (h *UserHandler) UpdateUserSettings(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	var req UpdateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.service.UpdateUserSettings(c.Request.Context(), userID, req.Income, req.SavingGoal)
	if err != nil {
		var validationErr *services.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user settings updated successfully"})
}
