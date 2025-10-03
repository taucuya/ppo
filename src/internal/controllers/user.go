package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetUserByPrivatesHandler(ctx *gin.Context) {
	email := ctx.Query("email")
	phone := ctx.Query("phone")
	if phone != "" {
		c.GetUserByPhoneHandler(ctx)
	}
	if email != "" {
		c.GetUserByEmailHandler(ctx)
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone parameter is required"})
}

func (c *Controller) GetUserByEmailHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) GetAllUsersHandler(ctx *gin.Context) {
	good := c.VerifyA(ctx)
	if !good {
		return
	}

	user, err := c.UserService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("[ERROR] Cant get all users: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
