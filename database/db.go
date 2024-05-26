package database

import (


	"github.com/omidghane/web_mid/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

var db *gorm.DB
var validate *validator.Validate

func connectDB(){
	var err error
	db, err = gorm.Open("sqlite3", "shoppingcart.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&models.User{}, &models.Basket{})

	validate = validator.New()
	validate.RegisterValidation("customState", func(fl validator.FieldLevel) bool {
		state := fl.Field().String()
		return state == "COMPLETED" || state == "PENDING"
	})

}


