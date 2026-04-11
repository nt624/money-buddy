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

type createFixedCostUseCase interface {
	Execute(ctx context.Context, userID string, name string, amount int) (domain.FixedCost, error)
}

type CreateFixedCostHandler struct {
	uc createFixedCostUseCase
}

func NewCreateFixedCostHandler(uc createFixedCostUseCase) *CreateFixedCostHandler {
	return &CreateFixedCostHandler{uc: uc}
}

type createFixedCostRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func (h *CreateFixedCostHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	var req createFixedCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("create_fixed_cost: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	fixedCost, err := h.uc.Execute(c.Request.Context(), userID, req.Name, req.Amount)
	if err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "固定費の作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"fixed_cost": fixedCost})
}
