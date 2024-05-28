package handlers

import (
	"net/http"
	// "strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/omidghane/web_mid/database"
	"github.com/omidghane/web_mid/models"
)

// var validate *validator.Validate

// func init() {
// 	validate = validator.New()
// 	validate.RegisterValidation("customState", validState)
// }

func validState(fl validator.FieldLevel) bool {
	state := fl.Field().String()
	return state == "PENDING" || state == "COMPLETED"
}

func GetBaskets(c echo.Context) error {
	userID := c.Get("userID").(uint)
	var baskets []models.Basket
	database.DB.Where("user_id = ?", userID).Find(&baskets)
	return c.JSON(http.StatusOK, baskets)
}

func CreateBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	var input models.BasketInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	if err := validate.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	basket := models.Basket{
		UserID: userID,
		Data:   input.Data,
		State:  input.State,
	}

	if err := database.DB.Create(&basket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, basket)
}

func UpdateBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	id := c.Param("id")
	var basket models.Basket
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&basket).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "basket not found"})
	}
	if basket.State == "COMPLETED" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": "cannot modify a completed basket"})
	}

	var input models.BasketInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	if err := validate.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	basket.Data = input.Data
	basket.State = input.State
	basket.UpdatedAt = time.Now()

	if err := database.DB.Save(&basket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, basket)
}

func GetBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	id := c.Param("id")
	var basket models.Basket
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&basket).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "basket not found"})
	}
	return c.JSON(http.StatusOK, basket)
}

func DeleteBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	id := c.Param("id")
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Basket{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "basket deleted"})
}
