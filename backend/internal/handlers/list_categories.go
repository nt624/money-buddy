package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/domain"
)

type listCategoriesUseCase interface {
	Execute(ctx context.Context) ([]domain.Category, error)
}

type ListCategoriesHandler struct {
	uc listCategoriesUseCase
}

func NewListCategoriesHandler(uc listCategoriesUseCase) *ListCategoriesHandler {
	return &ListCategoriesHandler{uc: uc}
}

func (h *ListCategoriesHandler) Handle(c *gin.Context) {
	categories, err := h.uc.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "カテゴリの取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
