package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type CreateWorkerRequest struct {
	IdUser   uuid.UUID `json:"id_user" binding:"required"`
	JobTitle string    `json:"job_title" binding:"required"`
}

// CreateWorkerHandler создает нового работника
// @Summary Создать работника
// @Description Создает нового работника в системе (только для администраторов)
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateWorkerRequest true "Данные для создания работника"
// @Success 201 {object} object "Работник успешно создан"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при создании работника"
// @Router /api/v1/workers [post]
func (c *Controller) CreateWorkerHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	var input CreateWorkerRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	worker := structs.Worker{
		IdUser:   input.IdUser,
		JobTitle: input.JobTitle,
	}

	if err := c.WorkerService.Create(ctx, worker); err != nil {
		log.Printf("[ERROR] Cant create worker: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Worker created successfully"})
}

// DeleteWorkerHandler удаляет работника
// @Summary Удалить работника
// @Description Удаляет работника по его идентификатору (только для администраторов)
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID работника"
// @Success 200 {object} object "Работник успешно удален"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при удалении работника"
// @Router /api/v1/workers/{id} [delete]
func (c *Controller) DeleteWorkerHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse worker id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID format"})
		return
	}

	if err := c.WorkerService.Delete(ctx, id); err != nil {
		log.Printf("[ERROR] Cant delete worker by id: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Worker deleted successfully"})
}

// GetWorkerByIdHandler получает работника по ID
// @Summary Получить работника по ID
// @Description Возвращает информацию о работнике по его идентификатору (только для администраторов)
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID работника" example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 {object} object "Данные работника"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Работник не найден"
// @Router /api/v1/workers/{id} [get]
func (c *Controller) GetWorkerByIdHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse worker id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID format"})
		return
	}

	worker, err := c.WorkerService.GetById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get worker by id: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Worker not found"})
		return
	}

	ctx.JSON(http.StatusOK, worker)
}

// GetWorkerOrders получает заказы работника
// @Summary Получить заказы работника
// @Description Возвращает список заказов текущего работника (только для работников)
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} object "Список заказов работника"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Заказы не найдены"
// @Router /api/v1/workers/me/orders [get]
func (c *Controller) GetWorkerOrders(ctx *gin.Context) {
	good := c.VerifyW(ctx)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	worker, err := c.WorkerService.GetOrders(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get worker orders: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Orders not found"})
		return
	}

	ctx.JSON(http.StatusOK, worker)
}

// GetAllWorkersHandler получает всех работников
// @Summary Получить всех работников
// @Description Возвращает список всех работников системы (только для администраторов)
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} object "Список всех работников"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при получении работников"
// @Router /api/v1/workers [get]
func (c *Controller) GetAllWorkersHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	workers, err := c.WorkerService.GetAllWorkers(ctx)
	if err != nil {
		log.Printf("[ERROR] Cant get all workers: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workers)
}

// AcceptOrderHandler принимает заказ работником
// @Summary Принять заказ
// @Description Принимает заказ для выполнения текущим работником (только для работников)
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id query string true "UUID заказа"
// @Success 200 {object} object "Заказ успешно принят"
// @Failure 400 {object} object "Неверный формат UUID заказа"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при принятии заказа"
// @Router /api/v1/workers/me/orders [post]
func (c *Controller) AcceptOrderHandler(ctx *gin.Context) {
	good := c.VerifyW(ctx)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	orderID := ctx.Query("order_id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order id parameter is required"})
		return
	}

	id_ord, err := uuid.Parse(orderID)
	if err != nil {
		log.Printf("[ERROR] Cant parse order id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.WorkerService.AcceptOrder(ctx, id_ord, id); err != nil {
		log.Printf("[ERROR] Cant accept order: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order accepted successfully"})
}
