package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/middleware"
	"money-buddy-backend/internal/usecase"
)

type deleteExpenseUseCase interface {
	Execute(userID string, id int) error
}

type DeleteExpenseHandler struct {
	uc deleteExpenseUseCase
}

func NewDeleteExpenseHandler(uc deleteExpenseUseCase) *DeleteExpenseHandler {
	return &DeleteExpenseHandler{uc: uc}
}

func (h *DeleteExpenseHandler) Handle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "支出IDが正しくありません"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	err = h.uc.Execute(userID, int(id))
	if err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		var ne *usecase.NotFoundError
		if errors.As(err, &ne) {
			c.JSON(http.StatusNotFound, gin.H{"error": ne.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	c.Status(http.StatusNoContent)
}
