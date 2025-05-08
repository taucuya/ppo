package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetUserByEmailHandler(ctx *gin.Context) {
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

func (c *Controller) GetUserByPhoneHandler(ctx *gin.Context) {
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
