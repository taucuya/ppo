package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) CreateReviewHandler(ctx *gin.Context) {
	var input struct {
		IdProduct string `json:"id_product"`
		IdUser    string `json:"id_user"`
		Rating    int    `json:"rating"`
		Text      string `json:"r_text"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id_prd, err := uuid.Parse(input.IdProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	id_usr, err := uuid.Parse(input.IdUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	r := structs.Review{
		IdProduct: id_prd,
		IdUser:    id_usr,
		Rating:    input.Rating,
		Text:      input.Text,
		Date:      time.Now(),
	}

	if err := c.ReviewService.Create(ctx, r); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Review created"})
}

func (c *Controller) DeleteReviewHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ReviewService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review deleted"})
}
