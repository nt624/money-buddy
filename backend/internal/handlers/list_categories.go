package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/domain"
	"pace-wallet-backend/internal/middleware"
)

type listCategoriesUseCase interface {
	Execute(ctx context.Context, userID string) ([]domain.Category, error)
}

type ListCategoriesHandler struct {
	uc listCategoriesUseCase
}

func NewListCategoriesHandler(uc listCategoriesUseCase) *ListCategoriesHandler {
	return &ListCategoriesHandler{uc: uc}
}

func (h *ListCategoriesHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categories, err := h.uc.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "カテゴリの取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
