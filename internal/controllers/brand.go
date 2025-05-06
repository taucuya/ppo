package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateBrandHandler(ctx *gin.Context) {
	var input struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		PriceCategory string `json:"price_category"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b := structs.Brand{
		Name:          input.Name,
		Description:   input.Description,
		PriceCategory: input.PriceCategory,
	}

	if err := c.BrandService.Create(ctx, b); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Brand created"})
}

func (c *Controller) DeleteBrandHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = c.BrandService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Brand deleted"})
}

func (c *Controller) GetAllBrandsInCategoryHander(ctx *gin.Context) {
	category := ctx.Param("cat")

	res, err := c.BrandService.GetAllBrandsInCategory(ctx, category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, res)
}
