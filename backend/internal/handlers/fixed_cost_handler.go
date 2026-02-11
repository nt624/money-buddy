package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/services"
)

type FixedCostHandler struct {
	service services.FixedCostService
}

func NewFixedCostHandler(r gin.IRouter, service services.FixedCostService) {
	handler := &FixedCostHandler{service: service}

	r.POST("/fixed-costs", handler.CreateFixedCost)
	r.GET("/fixed-costs", handler.ListFixedCosts)
	r.PUT("/fixed-costs/:id", handler.UpdateFixedCost)
	r.DELETE("/fixed-costs/:id", handler.DeleteFixedCost)
}

// CreateFixedCostRequest は固定費作成のリクエストボディです
type CreateFixedCostRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// CreateFixedCost は固定費を作成します
func (h *FixedCostHandler) CreateFixedCost(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	// リクエストボディ取得
	var req CreateFixedCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 作成実行
	fixedCost, err := h.service.CreateFixedCost(c.Request.Context(), userID, req.Name, req.Amount)
	if err != nil {
		var ve *services.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "固定費の作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"fixed_cost": fixedCost})
}

// ListFixedCosts は固定費一覧を取得します
func (h *FixedCostHandler) ListFixedCosts(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	fixedCosts, err := h.service.ListFixedCosts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "固定費の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fixed_costs": fixedCosts})
}

// UpdateFixedCostRequest は固定費更新のリクエストボディです
type UpdateFixedCostRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// UpdateFixedCost は固定費を更新します
func (h *FixedCostHandler) UpdateFixedCost(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	// IDパラメータ取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが正しくありません"})
		return
	}

	// リクエストボディ取得
	var req UpdateFixedCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新実行
	fixedCost, err := h.service.UpdateFixedCost(c.Request.Context(), userID, id, req.Name, req.Amount)
	if err != nil {
		var ve *services.ValidationError
		var ne *services.NotFoundError
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

// DeleteFixedCost は固定費を削除します
func (h *FixedCostHandler) DeleteFixedCost(c *gin.Context) {
	// TODO: Extract userID from authentication context when auth is implemented
	userID := DummyUserID

	// IDパラメータ取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが正しくありません"})
		return
	}

	// 削除実行
	err = h.service.DeleteFixedCost(c.Request.Context(), userID, id)
	if err != nil {
		var ne *services.NotFoundError
		if errors.As(err, &ne) {
			c.JSON(http.StatusNotFound, gin.H{"error": ne.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "固定費の削除に失敗しました"})
		return
	}

	c.Status(http.StatusNoContent)
}
