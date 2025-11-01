package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       string `json:"price" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Amount      int    `json:"amount" binding:"required,min=0"`
	IdBrand     string `json:"id_brand" binding:"required"`
	PicLink     string `json:"pic_link"`
	Articule    string `json:"articule"`
}

// CreateProductHandler создает новый продукт
// @Summary Создать продукт
// @Description Создает новый продукт в системе (только для администраторов)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProductRequest true "Данные для создания продукта"
// @Success 201 {object} object "Продукт успешно создан"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера при создании продукта"
// @Router /api/v1/products [post]
func (c *Controller) CreateProductHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to create product")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

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
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := strconv.ParseFloat(input.Price, 64)
	if err != nil {
		log.Printf("[ERROR] Cant parse product price: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id_brnd, err := uuid.Parse(input.IdBrand)
	if err != nil {
		log.Printf("[ERROR] Cant parse brand id: %v", err)
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
		log.Printf("[ERROR] Cant create product: %v", err)
		if errors.Is(err, structs.ErrDuplicateArticule) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product with this articule already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

// GetProductsHandler получает продукты
// @Summary Получить продукты
// @Description Возвращает список продуктов с различными фильтрами: по категории или бренду
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string false "UUID продукта"
// @Param art query string false "Артикул продукта"
// @Param category query string false "Категория продукта" Enums(уход, декоративная, парфюмерия, для волос, мужская)
// @Param brand query string false "Название бренда"
// @Success 200 {object} object "Данные продукта или список продуктов"
// @Failure 400 {object} object "Неверные параметры запроса"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 404 {object} object "Продукт не найден"
// @Failure 500 {object} object "Ошибка сервера при получении продуктов"
// @Router /api/v1/products [get]
func (c *Controller) GetProductsHandler(ctx *gin.Context) {
	category := ctx.Query("category")
	brand := ctx.Query("brand")

	if category == "" && brand == "" {
		c.GetProductHandler(ctx)
	} else if category != "" {
		c.GetProductsByCategoryHandler(ctx)
	} else if brand != "" {
		c.GetProductsByBrandHandler(ctx)
	}
}

func (c *Controller) GetProductHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to get product")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	if id := ctx.Query("id"); id != "" {
		pid, err := uuid.Parse(id)
		if err != nil {
			log.Printf("[ERROR] Cant parse product id: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
			return
		}

		product, err := c.ProductService.GetById(ctx, pid)
		if err != nil {
			log.Printf("[ERROR] Cant get product by id: %v", err)
			if errors.Is(err, structs.ErrProductNotFound) ||
				errors.Is(err, structs.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}

			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		ctx.JSON(http.StatusOK, product)
	}
	// if name := ctx.Query("name"); name != "" {
	// 	product, err := c.ProductService.GetByName(ctx, name)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, product)
	// }
	if art := ctx.Query("art"); art != "" {
		product, err := c.ProductService.GetByArticule(ctx, art)
		if err != nil {
			log.Printf("[ERROR] Cant get product by articule: %v", err)
			if errors.Is(err, structs.ErrProductNotFound) ||
				errors.Is(err, structs.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		ctx.JSON(http.StatusOK, product)
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid query parameter provided"})
}

// DeleteProductHandler удаляет продукт
// @Summary Удалить продукт
// @Description Удаляет продукт по его идентификатору (только для администраторов)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID продукта"
// @Success 200 {object} object "Продукт успешно удален"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Продукт не найден"
// @Failure 500 {object} object "Ошибка сервера при удалении продукта"
// @Router /api/v1/products/{id} [delete]
func (c *Controller) DeleteProductHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to delete product")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse product id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ProductService.Delete(ctx, id); err != nil {
		log.Printf("[ERROR] Cant delete product by id: %v", err)
		if errors.Is(err, structs.ErrProductNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (c *Controller) GetProductsByCategoryHandler(ctx *gin.Context) {
	category := ctx.Query("category")
	possible := []string{"уход", "декоративная", "парфюмерия", "для волос", "мужская"}
	inarr := false
	for _, i := range possible {
		if category == i {
			inarr = true
			break
		}
	}
	if !inarr {
		log.Printf("[ERROR] Cant parse status to get product category")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product category format"})
		return
	}
	products, err := c.ProductService.GetByCategory(ctx, category)
	if err != nil {
		log.Printf("[ERROR] Cant get products by category: %v", err)
		if errors.Is(err, structs.ErrProductNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No products found in this category"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No products found in this category"})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (c *Controller) GetProductsByBrandHandler(ctx *gin.Context) {
	brand := ctx.Query("brand")
	if brand == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Brand parameter is required"})
		return
	}

	products, err := c.ProductService.GetByBrand(ctx, brand)
	if err != nil {
		log.Printf("[ERROR] Cant get products by brand: %v", err)

		if errors.Is(err, structs.ErrProductNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No products found for this brand"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products by brand"})
		return
	}

	if len(products) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No products found for this brand"})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetReviewsForProductHandler получает отзывы для продукта
// @Summary Получить отзывы для продукта
// @Description Возвращает список отзывов для указанного продукта
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "UUID продукта"
// @Success 200 {array} object "Список отзывов для продукта"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 404 {object} object "Отзывы не найден"
// @Failure 500 {object} object "Ошибка сервера при получении отзывов"
// @Router /api/v1/products/{id}/reviews [get]
func (c *Controller) GetReviewsForProductHandler(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Cant parse product id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviews, err := c.ReviewService.ReviewsForProduct(ctx, productID)
	if err != nil {
		log.Printf("[ERROR] Cant get reviews for product: %v", err)

		if errors.Is(err, structs.ErrReviewNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No reviews found for this product"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reviews) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No reviews found for this product"})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}
