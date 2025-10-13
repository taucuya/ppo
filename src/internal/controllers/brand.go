package controller

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"

	"github.com/gin-gonic/gin"
)

type CreateBrandRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	PriceCategory string `json:"price_category" binding:"required"`
}

// CreateBrandHandler создает новый бренд
// @Summary Создать бренд
// @Description Создает новый бренд в системе (только для администраторов)
// @Tags brands
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBrandRequest true "Данные для создания бренда"
// @Success 201 {object} object "Бренд успешно создан"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при создании бренда"
// @Router /api/v1/brands [post]
func (c *Controller) CreateBrandHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	var input CreateBrandRequest
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

// GetBrandByIdHandler получает бренд по ID
// @Summary Получить бренд по ID
// @Description Возвращает информацию о бренде по его идентификатору (только для администраторов)
// @Tags brands
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID бренда"
// @Success 200 {object} object "Данные бренда"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Бренд не найден"
// @Router /api/v1/brands/{id} [get]
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

// DeleteBrandHandler удаляет бренд
// @Summary Удалить бренд
// @Description Удаляет бренд по его идентификатору (только для администраторов)
// @Tags brands
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID бренда"
// @Success 200 {object} object "Бренд успешно удален"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при удалении бренда"
// @Router /api/v1/brands/{id} [delete]
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

// GetAllBrandsInCategoryHander получает бренды по категории цен
// @Summary Получить бренды по категории цен
// @Description Возвращает список брендов в указанной ценовой категории
// @Tags brands
// @Accept json
// @Produce json
// @Param category query string true "Ценовая категория" Enums(бюджет, средний сегмент, люкс)
// @Success 200 {array} object "Список брендов в категории"
// @Failure 400 {object} object "Неверная категория"
// @Failure 500 {object} object "Ошибка сервера при получении брендов"
// @Router /api/v1/brands [get]
func (c *Controller) GetAllBrandsInCategoryHander(ctx *gin.Context) {
	category := ctx.Query("category")

	possible := []string{"бюджет", "средний сегмент", "люкс"}
	inarr := false
	for _, i := range possible {
		if category == i {
			inarr = true
		}
	}
	if !inarr {
		log.Printf("[ERROR] Cant parse status to get brand category")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand category format"})
		return
	}

	res, err := c.BrandService.GetAllBrandsInCategory(ctx, category)
	if err != nil {
		log.Printf("[ERROR] Cant get brands by category: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, res)
}
