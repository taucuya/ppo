package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) GetFavouritesHandler(ctx *gin.Context) {
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

	items, err := c.FavouritesService.GetItems(ctx.Request.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Cant get favourites items: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *Controller) AddFavouritesItemHandler(ctx *gin.Context) {
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
		IdProduct string `json:"id_product"`
		Amount    int    `json:"amount"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idProduct, err := uuid.Parse(input.IdProduct)
	if err != nil {
		log.Printf("[ERROR] Cant get product id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_product"})
		return
	}

	item := structs.FavouritesItem{
		IdProduct: idProduct,
	}

	if err := c.FavouritesService.AddItem(ctx.Request.Context(), item, id); err != nil {
		log.Printf("[ERROR] Cant add item to favourites: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Item added"})
}

func (c *Controller) DeleteFavouritesItemHandler(ctx *gin.Context) {
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
		ProductID uuid.UUID `json:"product_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.FavouritesService.DeleteItem(ctx.Request.Context(), id, input.ProductID); err != nil {
		log.Printf("[ERROR] Cant delete item from Favourites: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
