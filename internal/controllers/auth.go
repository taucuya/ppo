package controller

import (
	"context"
	"log"
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
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		log.Printf("[ERROR] Cant parse date: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := structs.User{
		Name:          input.Name,
		Date_of_birth: t,
		Mail:          input.Mail,
		Password:      input.Password,
		Phone:         input.Phone,
		Address:       input.Address,
		Status:        "новый",
		Role:          "обычный пользователь",
	}

	if err := c.AuthServise.SignIn(ctx.Request.Context(), user); err != nil {
		log.Printf("[ERROR] Cant signup: %v", err)
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
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := c.AuthServise.LogIn(ctx.Request.Context(), input.Email, input.Password)
	if err != nil {
		log.Printf("[ERROR] Cant login: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.SetCookie("access_token", accessToken, 900, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, 604800, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

func (c *Controller) LogoutHandler(ctx *gin.Context) {
	var refreshToken struct {
		Rt string `json:"refresh_token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&refreshToken); err != nil {
		log.Printf("[ERROR] Cant bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if refreshToken.Rt == "" {
		log.Printf("[ERROR] Bad refresh token.")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	if err := c.AuthServise.LogOut(ctx.Request.Context(), refreshToken.Rt); err != nil {
		log.Printf("[ERROR] Cant logout: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (c *Controller) Verify(ctx *gin.Context) bool {
	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		log.Printf("[ERROR] Cant get access token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return false
	}

	rtoken, err := ctx.Cookie("refresh_token")
	if err != nil {
		log.Printf("[ERROR] Cant get refresh token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return false
	}

	new, agood, rgood, err := c.AuthServise.VerifyTokens(context.Background(), atoken, rtoken)
	if agood && rgood && err == nil {
		if new != `` {
			ctx.SetCookie("access_token", new, 900, "/", "localhost", false, true)
		}
		return true
	}
	return false
}

func (c *Controller) VerifyA(ctx *gin.Context) bool {
	if good := c.Verify(ctx); !good {
		return false
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		log.Printf("[ERROR] Cant get access token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return false
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		log.Printf("[ERROR] Cant get id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return false
	}

	good := c.AuthServise.CheckAdmin(context.Background(), id)
	return good
}

func (c *Controller) VerifyW(ctx *gin.Context) bool {
	if good := c.Verify(ctx); !good {
		return false
	}

	atoken, err := ctx.Cookie("access_token")
	if err != nil {
		log.Printf("[ERROR] Cant get access token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return false
	}

	id, err := c.AuthServise.GetId(atoken)
	if err != nil {
		log.Printf("[ERROR] Cant get id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return false
	}

	good := c.AuthServise.CheckWorker(context.Background(), id)
	return good
}
