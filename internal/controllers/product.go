package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) CreateProductHandler(ctx *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Category    string `json:"category"`
		Amount      int    `json:"amount"`
		IdBrand     string `json:"id_brand"`
		PicLink     string `json:"pic_link"`
		Articule    string `json:"articule"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := strconv.ParseFloat(input.Price, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id_brnd, err := uuid.Parse(input.IdBrand)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	p := structs.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       pr,
		Category:    input.Category,
		Amount:      input.Amount,
		IdBrand:     id_brnd,
		PicLink:     input.PicLink,
		Articule:    input.Articule,
	}

	if err := c.ProductService.Create(ctx, p); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

func (c *Controller) DeleteProductHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ProductService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (c *Controller) GetProductsByCategoryHandler(ctx *gin.Context) {
	category := ctx.Param("category")
	products, err := c.ProductService.GetByCategory(ctx, category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (c *Controller) GetReviewsForProductHandler(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviews, err := c.ReviewService.ReviewsForProduct(ctx, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}
