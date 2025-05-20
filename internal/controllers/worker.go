package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) CreateWorkerHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	var input struct {
		IdUser   uuid.UUID `json:"id_user"`
		JobTitle string    `json:"job_title"`
	}

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
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
