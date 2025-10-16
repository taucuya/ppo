package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type AddFavouriteRequest struct {
	IdProduct string `json:"id_product" binding:"required"`
}

// GetFavouritesHandler получает все товары в избранном пользователя
// @Summary Получить избранные товары
// @Description Возвращает список всех товаров в избранном текущего пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} object "Список товаров в избранном"
// @Failure 400 {object} object "Неверный формат ID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 500 {object} object "Ошибка сервера при получении избранного"
// @Router /api/v1/users/me/favourite/items [get]
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

// AddFavouritesItemHandler добавляет товар в избранное
// @Summary Добавить товар в избранное
// @Description Добавляет товар в избранное текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddFavouriteRequest true "Данные товара для добавления в избранное"
// @Success 201 {object} object "Товар успешно добавлен в избранное"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 500 {object} object "Ошибка сервера при добавлении в избранное"
// @Router /api/v1/users/me/favourite/items [post]
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

	var input AddFavouriteRequest

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

// DeleteFavouritesItemHandler удаляет товар из избранного
// @Summary Удалить товар из избранного
// @Description Удаляет товар из избранного текущего пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_product path string true "UUID товара"
// @Success 200 {object} object "Товар успешно удален из избранного"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 500 {object} object "Ошибка сервера при удалении из избранного"
// @Router /api/v1/users/me/favourite/items/{id_product} [delete]
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

	id_item, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse favourites id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.FavouritesService.DeleteItem(ctx.Request.Context(), id, id_item); err != nil {
		log.Printf("[ERROR] Cant delete item from Favourites: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
