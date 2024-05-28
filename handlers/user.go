package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/omidghane/web_mid/database"
	"github.com/omidghane/web_mid/models"
	"github.com/omidghane/web_mid/utils"
	// "github.com/go-playground/validator/v10"
)

// var validate *validator.Validate

// func init() {
// 	validate = validator.New()
// }

func Signup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	if err := validate.Struct(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	if err := user.HashPassword(user.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "could not hash password"})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "signup successful"})
}

func Login(c echo.Context) error {
	var input models.User
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "invalid credentials"})
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "invalid credentials"})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "could not create token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"token": token})
}
