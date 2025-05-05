package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) GetBasketItemsHandler(ctx *gin.Context) {

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	items, err := c.BasketService.GetItems(ctx.Request.Context(), id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *Controller) AddBasketItemHandler(ctx *gin.Context) {
	var input struct {
		IdProduct string `json:"id_product"`
		IdBasket  string `json:"id_basket"`
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

	idBasket, err := uuid.Parse(input.IdBasket)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_basket"})
		return
	}

	item := structs.BasketItem{
		IdProduct: idProduct,
		IdBasket:  idBasket,
		Amount:    input.Amount,
	}

	if err := c.BasketService.AddItem(ctx.Request.Context(), item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Item added"})
}

func (c *Controller) DeleteBasketItemHandler(ctx *gin.Context) {
	idStr := ctx.Param("item_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := c.BasketService.DeleteItem(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}

func (c *Controller) UpdateBasketItemAmountHandler(ctx *gin.Context) {
	var input struct {
		ItemID uuid.UUID `json:"item_id" binding:"required"`
		Amount int       `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BasketService.UpdateItemAmount(ctx.Request.Context(), input.ItemID, input.Amount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item amount"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item amount updated"})
}
