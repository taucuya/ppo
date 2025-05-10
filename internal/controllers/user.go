package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":            user.Id,
		"name":          user.Name,
		"date_of_birth": user.Date_of_birth.Format("2006-01-02"),
	})
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":            user.Id,
		"name":          user.Name,
		"date_of_birth": user.Date_of_birth.Format("2006-01-02"),
	})
}
