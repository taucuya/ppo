package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type CreateReviewRequest struct {
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
	Text   string `json:"r_text" binding:"required"`
}

// CreateReviewHandler создает новый отзыв
// @Summary Создать отзыв
// @Description Создает новый отзыв для указанного продукта от текущего пользователя
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_product path string true "UUID продукта"
// @Param request body CreateReviewRequest true "Данные для создания отзыва"
// @Success 201 {object} object "Отзыв успешно создан"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 500 {object} object "Ошибка сервера при создании отзыва"
// @Router /api/v1/users/me/products/{id_product}/reviews [post]
func (c *Controller) CreateReviewHandler(ctx *gin.Context) {
	good := c.Verify(ctx)
	if !good {
		return
	}

	var input CreateReviewRequest
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

// GetReviewByIdHandler получает отзыв по ID
// @Summary Получить отзыв по ID
// @Description Возвращает информацию об отзыве по его идентификатору
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID отзыва"
// @Success 200 {object} object "Данные отзыва"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 404 {object} object "Отзыв не найден"
// @Router /api/v1/products/{id}/reviews/{id} [get]
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

// DeleteReviewHandler удаляет отзыв
// @Summary Удалить отзыв
// @Description Удаляет отзыв по его идентификатору (только для администраторов)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID отзыва"
// @Success 200 {object} object "Отзыв успешно удален"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Отзыв не найден"
// @Failure 500 {object} object "Ошибка сервера при удалении отзыва"
// @Router /api/v1/products/{id}/reviews/{id} [delete]
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
