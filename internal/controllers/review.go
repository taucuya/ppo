package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) CreateReviewHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	fmt.Println(good)
	if !good {
		return
	}

	var input struct {
		Rating int    `json:"rating"`
		Text   string `json:"r_text"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id_prd, err := uuid.Parse(ctx.Param("id_product"))
	if err != nil {
		log.Printf("[ERROR] Cant parse product id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
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

	r := structs.Review{
		IdProduct: id_prd,
		IdUser:    id,
		Rating:    input.Rating,
		Text:      input.Text,
		Date:      time.Now(),
	}

	if err := c.ReviewService.Create(ctx, r); err != nil {
		log.Printf("[ERROR] Cant create review: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Review created"})
}

func (c *Controller) GetReviewByIdHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse review id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID format"})
		return
	}

	review, err := c.ReviewService.GetById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get review by id: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	ctx.JSON(http.StatusOK, review)
}

func (c *Controller) DeleteReviewHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse review id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ReviewService.Delete(ctx, id); err != nil {
		log.Printf("[ERROR] Cant delete review by id: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review deleted"})
}
