package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taucuya/ppo/internal/core/structs"
)

// GetUserByPrivatesHandler получает пользователя по email или телефону
// @Summary Получить пользователя по email или телефону
// @Description Возвращает информацию о пользователе по email или номеру телефона (только для администраторов) или информацию о всех userах
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param email query string false "Email пользователя"
// @Param phone query string false "Телефон пользователя"
// @Success 200 {object} object "Данные пользователя"
// @Failure 400 {object} object "Не указаны параметры email или phone"
// @Failure 401 {object} object "Неавторизованный доступ"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 404 {object} object "Пользователь не найден"
// @Router /api/v1/users [get]
func (c *Controller) GetUserByPrivatesHandler(ctx *gin.Context) {
	email := ctx.Query("email")
	phone := ctx.Query("phone")
	if phone != "" {
		c.GetUserByPhoneHandler(ctx)
	} else if email != "" {
		c.GetUserByEmailHandler(ctx)
	} else {
		c.GetAllUsersHandler(ctx)
	}
}

func (c *Controller) GetUserByEmailHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to get user")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}

	user, err := c.UserService.GetByMail(ctx, email)
	if err != nil {
		log.Printf("[ERROR] Cant get user by mail: %v", err)

		if errors.Is(err, structs.ErrUserNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) GetAllUsersHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to get users")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	users, err := c.UserService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("[ERROR] Cant get all users: %v", err)
		if errors.Is(err, structs.ErrUserNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// func (c *Controller) GetUserByIdHandler(ctx *gin.Context) {
// 	good := c.VerifyA(ctx)
// 	if !good {
// 		return
// 	}

// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
// 		return
// 	}

// 	user, err := c.UserService.GetById(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, user)
// }

func (c *Controller) GetUserByPhoneHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		log.Printf("[ERROR] Cant autorize to get user")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	phone := ctx.Query("phone")
	if phone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone parameter is required"})
		return
	}

	user, err := c.UserService.GetByPhone(ctx, phone)
	if err != nil {
		log.Printf("[ERROR] Cant get user by phone: %v", err)

		if errors.Is(err, structs.ErrUserNotFound) ||
			errors.Is(err, structs.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
