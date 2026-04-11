package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type completeInitialSetupUseCase interface {
	Execute(ctx context.Context, userID string, income, savingGoal int, fixedCosts []usecase.FixedCostItem) error
}

type CompleteInitialSetupHandler struct {
	uc completeInitialSetupUseCase
}

func NewCompleteInitialSetupHandler(uc completeInitialSetupUseCase) *CompleteInitialSetupHandler {
	return &CompleteInitialSetupHandler{uc: uc}
}

type initialSetupRequest struct {
	Income     int                   `json:"income"`
	SavingGoal int                   `json:"savingGoal"`
	FixedCosts []usecase.FixedCostItem `json:"fixedCosts"`
}

func (h *CompleteInitialSetupHandler) Handle(c *gin.Context) {
	var req initialSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("complete_initial_setup: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	err := h.uc.Execute(c.Request.Context(), userID, req.Income, req.SavingGoal, req.FixedCosts)
	if err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		var ne *usecase.NotFoundError
		if errors.As(err, &ne) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": ne.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
