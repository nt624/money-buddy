package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/domain"
	"money-buddy-backend/internal/middleware"
)

type listFixedCostsUseCase interface {
	Execute(ctx context.Context, userID string) ([]domain.FixedCost, error)
}

type ListFixedCostsHandler struct {
	uc listFixedCostsUseCase
}

func NewListFixedCostsHandler(uc listFixedCostsUseCase) *ListFixedCostsHandler {
	return &ListFixedCostsHandler{uc: uc}
}

func (h *ListFixedCostsHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	fixedCosts, err := h.uc.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "固定費の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fixed_costs": fixedCosts})
}
