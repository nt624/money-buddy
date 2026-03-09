package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/middleware"
	"money-buddy-backend/internal/usecase"
)

type deleteFixedCostUseCase interface {
	Execute(ctx context.Context, userID string, id int) error
}

type DeleteFixedCostHandler struct {
	uc deleteFixedCostUseCase
}

func NewDeleteFixedCostHandler(uc deleteFixedCostUseCase) *DeleteFixedCostHandler {
	return &DeleteFixedCostHandler{uc: uc}
}

func (h *DeleteFixedCostHandler) Handle(c *gin.Context) {
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

	err = h.uc.Execute(c.Request.Context(), userID, id)
	if err != nil {
		var ne *usecase.NotFoundError
		if errors.As(err, &ne) {
			c.JSON(http.StatusNotFound, gin.H{"error": ne.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "固定費の削除に失敗しました"})
		return
	}

	c.Status(http.StatusNoContent)
}
