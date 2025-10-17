package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type BasketItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Amount    int       `json:"amount" binding:"required,min=1"`
}

type BasketItemDeleteRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
}

// GetBasketItemsHandler получает все товары в корзине пользователя
// @Summary Получить товары корзины
// @Description Возвращает список всех товаров в корзине текущего пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} object "Список товаров в корзине"
// @Failure 400 {object} object "Неверный формат ID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 500 {object} object "Ошибка сервера при получении товаров"
// @Router /api/v1/users/me/basket/items [get]
func (c *Controller) GetBasketItemsHandler(ctx *gin.Context) {
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

	items, err := c.BasketService.GetItems(ctx.Request.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Cant get basket items: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// GetBasketByIdHandler получает корзину по ID пользователя
// @Summary Получить корзину
// @Description Возвращает корзину текущего пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Данные корзины"
// @Failure 400 {object} object "Неверный формат ID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 404 {object} object "Корзина не найдена"
// @Router /api/v1/users/me/basket [get]
func (c *Controller) GetBasketByIdHandler(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid basket ID format"})
		return
	}

	basket, err := c.BasketService.GetById(ctx, id)
	fmt.Println(err)
	if err != nil {
		log.Printf("[ERROR] Cant get basket: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Basket not found"})
		return
	}

	ctx.JSON(http.StatusOK, basket)
}

// AddBasketItemHandler добавляет товар в корзину
// @Summary Добавить товар в корзину
// @Description Добавляет товар с указанным количеством в корзину пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BasketItemRequest true "Данные товара для добавления"
// @Success 201 {object} object "Товар успешно добавлен"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 404 {object} object "Корзина не найдена"
// @Failure 500 {object} object "Ошибка сервера при добавлении товара"
// @Router /api/v1/users/me/basket/items [post]
func (c *Controller) AddBasketItemHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		log.Printf("[ERROR] Cant get access token: %v", err)
		fmt.Println("HERE")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		log.Printf("[ERROR] Cant get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input BasketItemRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := structs.BasketItem{
		IdProduct: input.ProductID,
		Amount:    input.Amount,
	}

	if err := c.BasketService.AddItem(ctx.Request.Context(), item, id); err != nil {
		log.Printf("[ERROR] Cant add item to basket: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Item added"})
}

// DeleteBasketItemHandler удаляет товар из корзины
// @Summary Удалить товар из корзины
// @Description Удаляет указанный товар из корзины пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BasketItemDeleteRequest true "Данные товара для удаления"
// @Success 200 {object} object "Товар успешно удален"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 404 {object} object "Корзина не найдена"
// @Failure 500 {object} object "Ошибка сервера при удалении товара"
// @Router /api/v1/users/me/basket/items [delete]
func (c *Controller) DeleteBasketItemHandler(ctx *gin.Context) {
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

	var input BasketItemDeleteRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BasketService.DeleteItem(ctx.Request.Context(), id, input.ProductID); err != nil {
		log.Printf("[ERROR] Cant delete item from basket: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}

// UpdateBasketItemAmountHandler обновляет количество товара в корзине
// @Summary Обновить количество товара
// @Description Обновляет количество указанного товара в корзине пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BasketItemRequest true "Данные для обновления количества"
// @Success 200 {object} object "Количество товара успешно обновлено"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 404 {object} object "Элемент корзины не найден"
// @Failure 500 {object} object "Ошибка сервера при обновлении количества"
// @Router /api/v1/users/me/basket/items [patch]
func (c *Controller) UpdateBasketItemAmountHandler(ctx *gin.Context) {
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

	var input BasketItemRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BasketService.UpdateItemAmount(ctx.Request.Context(), id, input.ProductID, input.Amount); err != nil {
		log.Printf("[ERROR] Cant update item in basket: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item amount"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item amount updated"})
}
