package handlers

import (
	"github.com/omidghane/web_midterm/web_midterm/database"
	"github.com/omidghane/web_midterm/web_midterm/models"
	"github.com/omidghane/web_midterm/web_midterm/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)
func Signup(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := validate.Struct(&user); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := user.HashPassword(user.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
	}

	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, gin.H{"message": "signup successful"})
}

func Login(c echo.Context) error {
	var input User
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var user User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
	}

	return c.JSON(http.StatusOK, gin.H{"token": token})
}