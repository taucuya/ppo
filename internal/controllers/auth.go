package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taucuya/ppo/internal/core/structs"
)

func (c *Controller) SignupHandler(ctx *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		DateOfBirth string `json:"date_of_birth" binding:"required"`
		Mail        string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		Phone       string `json:"phone" binding:"required"`
		Address     string `json:"address" binding:"required"`
		Status      string `json:"status"`
		Role        string `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		return
	}

	user := structs.User{
		Name:          input.Name,
		Date_of_birth: t,
		Mail:          input.Mail,
		Password:      input.Password,
		Phone:         input.Phone,
		Address:       input.Address,
		Status:        input.Status,
		Role:          input.Role,
	}

	if err := c.AuthServise.SignIn(ctx.Request.Context(), user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func (c *Controller) LoginHandler(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := c.AuthServise.LogIn(ctx.Request.Context(), input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *Controller) LogoutHandler(ctx *gin.Context) {
	var refreshToken struct {
		Rt string `json:"refresh_token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&refreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	fmt.Println(refreshToken.Rt)
	if refreshToken.Rt == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	if err := c.AuthServise.LogOut(ctx.Request.Context(), refreshToken.Rt); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
