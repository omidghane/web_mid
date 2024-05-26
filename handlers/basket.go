package handlers

import (
	"github.com/omidghane/web_mid/database"
	"github.com/omidghane/web_mid/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func validState(fl validator.FieldLevel) bool {
	state := fl.Field().String()
	return state == "PENDING" || state == "COMPLETED"
}

func GetBaskets(c echo.Context) error {
	userID := c.Get("userID").(uint)
	var baskets []models.Basket
	database.db.Where("user_id = ?", userID).Find(&baskets)
	return c.JSON(http.StatusOK, baskets)
}

func CreateBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	var input models.BasketInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := validate.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	basket := models.Basket{
		UserID: userID,
		Data:   input.Data,
		State:  input.State,
	}

	if err := database.db.Create(&basket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, basket)
}

func UpdateBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	id := c.Param("id")
	var basket models.Basket
	if err := database.db.Where("id = ? AND user_id = ?", id, userID).First(&basket).Error; err != nil {
		return c.JSON(http.StatusNotFound, gin.H{"error": "basket not found"})
	}
	if basket.State == "COMPLETED" {
		return c.JSON(http.StatusForbidden, gin.H{"error": "cannot modify a completed basket"})
	}

	var input models.BasketInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := validate.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	basket.Data = input.Data
	basket.State = input.State
	basket.UpdatedAt = time.Now()

	if err := database.db.Save(&basket).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, basket)
}

func GetBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	id := c.Param("id")
	var basket models.Basket
	if err := database.db.Where("id = ? AND user_id = ?", id, userID).First(&basket).Error; err != nil {
		return c.JSON(http.StatusNotFound, gin.H{"error": "basket not found"})
	}
	return c.JSON(http.StatusOK, basket)
}

func DeleteBasket(c echo.Context) error {
	userID := c.Get("userID").(uint)
	id := c.Param("id")
	if err := database.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Basket{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, gin.H{"message": "basket deleted"})
}