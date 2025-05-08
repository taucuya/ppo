package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) CreateWorkerHandler(ctx *gin.Context) {
	var input struct {
		IdUser   uuid.UUID `json:"id_user"`
		JobTitle string    `json:"job_title"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	worker := structs.Worker{
		IdUser:   input.IdUser,
		JobTitle: input.JobTitle,
	}

	if err := c.WorkerService.Create(ctx, worker); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Worker created successfully"})
}

func (c *Controller) DeleteWorkerHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID format"})
		return
	}

	if err := c.WorkerService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Worker deleted successfully"})
}

func (c *Controller) GetWorkerByIdHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID format"})
		return
	}

	worker, err := c.WorkerService.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Worker not found"})
		return
	}

	ctx.JSON(http.StatusOK, worker)
}

func (c *Controller) GetAllWorkersHandler(ctx *gin.Context) {
	workers, err := c.WorkerService.GetAllWorkers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workers)
}

func (c *Controller) AcceptOrderHandler(ctx *gin.Context) {
	var input struct {
		WorkerID string `json:"id_worker"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	workerID, err := uuid.Parse(input.WorkerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID format"})
		return
	}

	if err := c.WorkerService.AcceptOrder(ctx, orderID, workerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order accepted successfully"})
}
