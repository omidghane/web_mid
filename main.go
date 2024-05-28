package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/omidghane/web_mid/database"
	"github.com/omidghane/web_mid/handlers"
	"github.com/omidghane/web_mid/middlewares"
	// "github.com/omidghane/web_mid/utils"
)

func main() {
	handlers.InitValidator()

	database.Init()
	defer database.DB.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/signup", handlers.Signup)
	e.POST("/login", handlers.Login)

	authorized := e.Group("/auth")
	authorized.Use(middlewares.JWTMiddleware)
	{
		authorized.GET("/basket", handlers.GetBaskets)
		authorized.POST("/basket", handlers.CreateBasket)
		authorized.PATCH("/basket/:id", handlers.UpdateBasket)
		authorized.GET("/basket/:id", handlers.GetBasket)
		authorized.DELETE("/basket/:id", handlers.DeleteBasket)
	}

	e.Start(":8080")
}
