package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/omidghane/web_midterm/web_midterm/database"
	"github.com/omidghane/web_midterm/web_midterm/handlers"
	"github.com/omidghane/web_midterm/web_midterm/middlewares"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)


func main() {
	connectDB()
	
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database.connectDB()

	e.POST("/login", handlers.Login)
	e.POST("/signup", handlers.Signup)

	authorized := e.Group("/auth", JWTMiddleware)
	{
		authorized.GET("/basket", handlers.GetBaskets, middlewares.JWTMiddleware)
		authorized.POST("/basket", handlers.CreateBasket, middlewares.JWTMiddleware)
		authorized.PATCH("/basket/:id", handlers.UpdateBasket, middlewares.JWTMiddleware)
		authorized.GET("/basket/:id", handlers.GetBasket, middlewares.JWTMiddleware)
		authorized.DELETE("/basket/:id", handlers.DeleteBasket, middlewares.JWTMiddleware	)
	}

	
	e.Start(":8080")
}
