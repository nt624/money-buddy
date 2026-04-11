package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type reorderCategoriesUseCase interface {
	Execute(ctx context.Context, userID string, items []usecase.ReorderCategoriesInput) error
}

type ReorderCategoriesHandler struct {
	uc reorderCategoriesUseCase
}

func NewReorderCategoriesHandler(uc reorderCategoriesUseCase) *ReorderCategoriesHandler {
	return &ReorderCategoriesHandler{uc: uc}
}

func (h *ReorderCategoriesHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Items []usecase.ReorderCategoriesInput `json:"items"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Execute(c.Request.Context(), userID, body.Items); err != nil {
		var ie *usecase.InternalError
		if errors.As(err, &ie) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "並び替えに失敗しました"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "並び替えに失敗しました"})
		return
	}

	c.Status(http.StatusNoContent)
}
