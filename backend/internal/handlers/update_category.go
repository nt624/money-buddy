package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/domain"
	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type updateCategoryUseCase interface {
	Execute(ctx context.Context, userID string, id int, name string) (domain.Category, error)
}

type UpdateCategoryHandler struct {
	uc updateCategoryUseCase
}

func NewUpdateCategoryHandler(uc updateCategoryUseCase) *UpdateCategoryHandler {
	return &UpdateCategoryHandler{uc: uc}
}

func (h *UpdateCategoryHandler) Handle(c *gin.Context) {
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

	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("update_category: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	cat, err := h.uc.Execute(c.Request.Context(), userID, id, body.Name)
	if err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		var nfe *usecase.NotFoundError
		if errors.As(err, &nfe) {
			c.JSON(http.StatusNotFound, gin.H{"error": nfe.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "カテゴリの更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": cat})
}
