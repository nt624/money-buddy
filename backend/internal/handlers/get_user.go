package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/domain"
	"money-buddy-backend/internal/middleware"
)

type getUserUseCase interface {
	Execute(ctx context.Context, userID string) (*domain.User, error)
}

type GetUserHandler struct {
	uc getUserUseCase
}

func NewGetUserHandler(uc getUserUseCase) *GetUserHandler {
	return &GetUserHandler{uc: uc}
}

func (h *GetUserHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	user, err := h.uc.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, user)
}
