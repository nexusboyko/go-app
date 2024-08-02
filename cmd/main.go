package main

import (
	"app/internal/db"
	"app/internal/middleware"
	"app/internal/routes"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	// router config + middleware
	router := gin.Default()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true
	router.Use(cors.New(cors.Config{
		// AllowOrigins:              []string{"http://localhost:8080", "https://accounts.google.com/"},
		AllowAllOrigins:           true,
		AllowMethods:              []string{"PUT", "GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:              []string{"Origin"},
		ExposeHeaders:             []string{"Content-Length"},
		OptionsResponseStatusCode: 204,
		AllowCredentials:          true,
	}))
	router.Static("/static", "./static")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("", routes.Index)
	main := router.Group("/")

	//  authenticated routes
	main.Use(middleware.AuthMiddleware())
	{
		main.GET("/home", routes.Home)
	}

	router.GET("/login", routes.Login)
	router.GET("/logout", func(c *gin.Context) {
		middleware.ClearSession(c)
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	router.GET("/callback", routes.Callback)

	log.Fatal(router.Run(":8080"))
}
