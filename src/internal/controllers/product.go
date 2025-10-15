package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) CreateProductHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

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
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		ctx.JSON(http.StatusOK, product)
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid query parameter provided"})
}

func (c *Controller) DeleteProductHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (c *Controller) GetProductsByCategoryHandler(ctx *gin.Context) {
	category := ctx.Query("category")
	possible := []string{"уходовая", "декоративная", "парфюмерия", "для волос", "мужская"}
	inarr := false
	for _, i := range possible {
		if category == i {
			inarr = true
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (c *Controller) GetProductsByBrandHandler(ctx *gin.Context) {
	brand := ctx.Param("brand")
	products, err := c.ProductService.GetByBrand(ctx, brand)
	if err != nil {
		log.Printf("[ERROR] Cant get products by brand: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}
