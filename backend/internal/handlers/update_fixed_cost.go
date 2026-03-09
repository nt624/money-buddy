package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/domain"
	"money-buddy-backend/internal/middleware"
	"money-buddy-backend/internal/usecase"
)

type updateFixedCostUseCase interface {
	Execute(ctx context.Context, userID string, id int, name string, amount int) (domain.FixedCost, error)
}

type UpdateFixedCostHandler struct {
	uc updateFixedCostUseCase
}

func NewUpdateFixedCostHandler(uc updateFixedCostUseCase) *UpdateFixedCostHandler {
	return &UpdateFixedCostHandler{uc: uc}
}

type updateFixedCostRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func (h *UpdateFixedCostHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが正しくありません"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが正しくありません"})
		return
	}

	var req updateFixedCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fixedCost, err := h.uc.Execute(c.Request.Context(), userID, id, req.Name, req.Amount)
	if err != nil {
		var ve *usecase.ValidationError
		var ne *usecase.NotFoundError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		if errors.As(err, &ne) {
			c.JSON(http.StatusNotFound, gin.H{"error": ne.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "固定費の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fixed_cost": fixedCost})
}
