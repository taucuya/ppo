package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	Address string `json:"address" binding:"required"`
}

type OrderResponse struct {
	ID      uuid.UUID `json:"id"`
	Date    time.Time `json:"date"`
	IdUser  uuid.UUID `json:"id_user"`
	Address string    `json:"address"`
	Status  string    `json:"status "`
	Price   float64   `json:"price"`
}

type OrderItemResponse struct {
	ID        uuid.UUID `json:"id"`
	IdProduct uuid.UUID `json:"id_product"`
	IdOrder   uuid.UUID `json:"id_order"`
	Amount    int       `json:"amount"`
	Price     float64   `json:"price"`
	Name      string    `json:"name"`
}

// CreateOrderHandler создает новый заказ
// @Summary Создать заказ
// @Description Создает новый заказ для текущего пользователя
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrderRequest true "Данные для создания заказа"
// @Success 201 {object} object "Заказ успешно создан"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 500 {object} object "Ошибка сервера при создании заказа"
// @Router /api/v1/orders [post]
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

	var input CreateOrderRequest

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

// GetOrderItemsHandler получает товары в заказе
// @Summary Получить товары заказа
// @Description Возвращает список товаров в указанном заказе
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID заказа"
// @Success 200 {array} object "Список товаров в заказе"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 404 {object} object "Отзывы не найдены"
// @Failure 500 {object} object "Ошибка сервера при получении товаров"
// @Router /api/v1/users/me/orders/{id}/items [get]
func (c *Controller) GetOrderItemsHandler(ctx *gin.Context) {
	good := c.Verify(ctx)

	if !good {
		log.Printf("[ERROR] Cant autorize to get order items")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
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
		if errors.Is(err, structs.ErrOrderNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order items not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(items) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No items in order"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// GetOrdersHandler получает заказы
// @Summary Получить заказы
// @Description Возвращает список заказов. Без параметров - заказы текущего пользователя, с status=непринятый - свободные заказы для работников
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Статус заказа для фильтрации" Enums(непринятый)
// @Success 200 {array} object "Список заказов"
// @Failure 400 {object} object "Неверный параметр статуса"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Заказы не найдены"
// @Failure 500 {object} object "Ошибка сервера при получении заказов"
// @Router /api/v1/users/me/orders [get]
func (c *Controller) GetOrdersHandler(ctx *gin.Context) {
	status := ctx.Query("status")
	if status != "" {
		c.GetFreeOrdersHandler(ctx)
	} else {
		c.GetOrdersByUserHandler(ctx)
	}
}

func (c *Controller) GetFreeOrdersHandler(ctx *gin.Context) {
	goodW := c.VerifyW(ctx)
	goodA := c.VerifyA(ctx)

	if !goodW && !goodA {
		log.Printf("[ERROR] Cant autorize to get free orders")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	status := ctx.Query("status")
	if status != "непринятый" {
		log.Printf("[ERROR] Cant parse status to get free orders")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status format"})
		return
	}
	ords, err := c.OrderService.GetFreeOrders(ctx)
	if err != nil {
		log.Printf("[ERROR] Cant get free orders: %v", err)
		if errors.Is(err, structs.ErrOrderNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No free orders found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ords)
}

// GetOrderByIdHandler получает заказ по ID
// @Summary Получить заказ по ID
// @Description Возвращает информацию о заказе по его идентификатору (для работников и администраторов)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID заказа"
// @Success 200 {object} object "Данные заказа"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Заказ не найден"
// @Router /api/v1/users/me/orders/{id} [get]
func (c *Controller) GetOrderByIdHandler(ctx *gin.Context) {
	goodW := c.VerifyW(ctx)
	goodA := c.VerifyA(ctx)

	if !goodW && !goodA {
		log.Printf("[ERROR] Cant autorize to get order")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
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

		if errors.Is(err, structs.ErrOrderNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// ChangeOrderStatusHandler изменяет статус заказа
// @Summary Изменить статус заказа
// @Description Обновляет статус заказа (для работников и администраторов)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID заказа"
// @Param status query string true "Новый статус заказа" Enums(некорректный, непринятый, принятый, собранный, отданный)
// @Success 200 {object} object "Статус заказа успешно обновлен"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Заказ не найден"
// @Failure 500 {object} object "Ошибка сервера при обновлении статуса"
// @Router /api/v1/users/me/orders/{id} [patch]
func (c *Controller) ChangeOrderStatusHandler(ctx *gin.Context) {
	goodW := c.VerifyW(ctx)
	goodA := c.VerifyA(ctx)

	if !goodW && !goodA {
		log.Printf("[ERROR] Cant autorize to change order status")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	status := ctx.Query("status")

	possible := []string{"некорректный", "непринятый", "принятый", "собранный", "отданный"}

	inarr := false
	for _, i := range possible {
		if status == i {
			inarr = true
			break
		}
	}
	if !inarr {
		log.Printf("[ERROR] Cant parse status to get order status")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status format"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse order id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.OrderService.ChangeOrderStatus(ctx, id, status)
	if err != nil {
		log.Printf("[ERROR] Cant change order status: %v", err)

		if errors.Is(err, structs.ErrOrderNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

// DeleteOrderHandler удаляет заказ
// @Summary Удалить заказ
// @Description Удаляет заказ по его идентификатору (только для администраторов)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID заказа"
// @Success 200 {object} object "Заказ успешно удален"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при удалении заказа"
// @Router /api/v1/users/me/orders/{id} [delete]
func (c *Controller) DeleteOrderHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to delete order")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
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
		if errors.Is(err, structs.ErrOrderNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}

func (c *Controller) GetOrdersByUserHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to get orders")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
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

	orders, err := c.OrderService.GetOrdersByUser(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get orders by user: %v", err)
		if errors.Is(err, structs.ErrOrderNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No orders found for user"})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Orders not found"})
		return
	}

	if len(orders) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No orders found for user"})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
