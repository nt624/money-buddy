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

type createCategoryUseCase interface {
	Execute(ctx context.Context, userID string, name string) (domain.Category, error)
}

type CreateCategoryHandler struct {
	uc createCategoryUseCase
}

func NewCreateCategoryHandler(uc createCategoryUseCase) *CreateCategoryHandler {
	return &CreateCategoryHandler{uc: uc}
}

func (h *CreateCategoryHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("create_category: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	cat, err := h.uc.Execute(c.Request.Context(), userID, body.Name)
	if err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "カテゴリの作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"category": cat})
}
