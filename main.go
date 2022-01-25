package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/middleware"
	"github.com/teandresmith/injapan-golang-react-backend/routes"
)

func main() {
	port := os.Getenv("PORT")	
	if port == "" {
		port = "8080"
	}

	router := gin.New()

	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":              "Welcome to InJapan Backend API",
			"public-get-endpoints": "/api/blogs, /api/blogs/:blogid, /api/tag, /api/tag/tag:id",
		})
	})
	
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())

	routes.BlogRoutes(router)
	routes.TagRoutes(router)
	routes.LoginRoutes(router)

	router.Use(middleware.Authentication())
	routes.AdminRoutes(router)

	log.Fatal(router.Run(":" + port))
}