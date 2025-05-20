package controller

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateBrandHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	var input struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		PriceCategory string `json:"price_category"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b := structs.Brand{
		Name:          input.Name,
		Description:   input.Description,
		PriceCategory: input.PriceCategory,
	}

	if err := c.BrandService.Create(ctx, b); err != nil {
		log.Printf("[ERROR] Cant create brand: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Brand created"})
}

func (c *Controller) GetBrandByIdHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse brand id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID format"})
		return
	}

	brand, err := c.BrandService.GetById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Cant get brand by id: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}

	ctx.JSON(http.StatusOK, brand)
}

func (c *Controller) DeleteBrandHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse brand id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = c.BrandService.Delete(ctx, id); err != nil {
		log.Printf("[ERROR] Cant delete brand: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Brand deleted"})
}

func (c *Controller) GetAllBrandsInCategoryHander(ctx *gin.Context) {
	category := ctx.Param("cat")

	res, err := c.BrandService.GetAllBrandsInCategory(ctx, category)
	if err != nil {
		log.Printf("[ERROR] Cant get brands by category: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, res)
}
