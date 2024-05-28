package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/omidghane/web_mid/models"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open("sqlite3", "shoppingcart.db")
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&models.User{}, &models.Basket{})
}
