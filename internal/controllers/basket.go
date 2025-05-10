package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) GetBasketItemsHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	items, err := c.BasketService.GetItems(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *Controller) GetBasketByIdHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid basket ID format"})
		return
	}

	basket, err := c.BasketService.GetById(ctx, id)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Basket not found"})
		return
	}

	ctx.JSON(http.StatusOK, basket)
}

func (c *Controller) AddBasketItemHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input struct {
		IdProduct string `json:"id_product"`
		Amount    int    `json:"amount"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idProduct, err := uuid.Parse(input.IdProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_product"})
		return
	}

	item := structs.BasketItem{
		IdProduct: idProduct,
		Amount:    input.Amount,
	}

	if err := c.BasketService.AddItem(ctx.Request.Context(), item, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Item added"})
}

func (c *Controller) DeleteBasketItemHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input struct {
		ProductID uuid.UUID `json:"product_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BasketService.DeleteItem(ctx.Request.Context(), id, input.ProductID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}

func (c *Controller) UpdateBasketItemAmountHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input struct {
		ProductID uuid.UUID `json:"product_id" binding:"required"`
		Amount    int       `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BasketService.UpdateItemAmount(ctx.Request.Context(), id, input.ProductID, input.Amount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item amount"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item amount updated"})
}
