package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateOrderHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		log.Printf("[ERROR] Cant get access token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		log.Printf("[ERROR] Cant get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input struct {
		Address string `json:"address"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o := structs.Order{
		Date:    time.Now(),
		IdUser:  id,
		Address: input.Address,
		Status:  "непринятый",
		Price:   0,
	}

	if err = c.OrderService.Create(ctx, o); err != nil {
		log.Printf("[ERROR] Cant create order: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Order created"})
}

func (c *Controller) GetOrderItemsHandler(ctx *gin.Context) {
	good := c.Verify(ctx)

	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse order id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := c.OrderService.GetItems(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get order items: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *Controller) GetFreeOrdersHandler(ctx *gin.Context) {
	goodW := c.VerifyW(ctx)
	goodA := c.VerifyA(ctx)

	if !goodW && !goodA {
		fmt.Println("HERE", goodW, goodA)
		return
	}

	ords, err := c.OrderService.GetFreeOrders(ctx)
	if err != nil {
		log.Printf("[ERROR] Cant get free orders: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ords)
}

func (c *Controller) GetOrderByIdHandler(ctx *gin.Context) {
	goodW := c.VerifyW(ctx)
	goodA := c.VerifyA(ctx)

	if !goodW && !goodA {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse order id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	order, err := c.OrderService.GetById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get order by id: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (c *Controller) ChangeOrderStatusHandler(ctx *gin.Context) {
	goodW := c.VerifyW(ctx)
	goodA := c.VerifyA(ctx)

	if !goodW && !goodA {
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse order id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&input); err != nil || input.Status == "" {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.OrderService.ChangeOrderStatus(ctx, id, input.Status)
	if err != nil {
		log.Printf("[ERROR] Cant change order status: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func (c *Controller) DeleteOrderHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse order id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.OrderService.Delete(ctx, id); err != nil {
		log.Printf("[ERROR] Cant delete order by id: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}

func (c *Controller) GetOrdersByUser(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		log.Printf("[ERROR] Cant get access token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		log.Printf("[ERROR] Cant get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := c.OrderService.GetOrdersByUser(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get orders by user: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Orders not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
