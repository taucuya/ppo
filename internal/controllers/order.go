package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateOrderHandler(ctx *gin.Context) {
	var input struct {
		IdUser   string `json:"id_user"`
		Address  string `json:"address"`
		Price    string `json:"price"`
		IdWorker string `json:"id_worker"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id_usr, err := uuid.Parse(input.IdUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	id_wrkr, err := uuid.Parse(input.IdWorker)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	pr, err := strconv.ParseFloat(input.Price, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	o := structs.Order{
		Date:     time.Now(),
		IdUser:   id_usr,
		Address:  input.Address,
		Status:   "непринятый",
		Price:    pr,
		IdWorker: id_wrkr,
	}

	if err = c.OrderService.Create(ctx, o); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Order created"})
}

func (c *Controller) GetOrderItemsHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := c.OrderService.GetItems(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *Controller) ChangeOrderStatusHandler(ctx *gin.Context) {
	var input struct {
		Status string `json:"status"`
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&input); err != nil || input.Status == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.OrderService.ChangeOrderStatus(ctx, id, input.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func (c *Controller) DeleteOrderHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.OrderService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
