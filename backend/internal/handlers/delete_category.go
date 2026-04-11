package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type deleteCategoryUseCase interface {
	Execute(ctx context.Context, userID string, id int) error
}

type DeleteCategoryHandler struct {
	uc deleteCategoryUseCase
}

func NewDeleteCategoryHandler(uc deleteCategoryUseCase) *DeleteCategoryHandler {
	return &DeleteCategoryHandler{uc: uc}
}

func (h *DeleteCategoryHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なカテゴリIDです"})
		return
	}

	if err := h.uc.Execute(c.Request.Context(), userID, id); err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "カテゴリの削除に失敗しました"})
		return
	}

	c.Status(http.StatusNoContent)
}
