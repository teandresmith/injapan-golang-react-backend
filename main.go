package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/middleware"
	"github.com/teandresmith/injapan-golang-react-backend/routes"
)

func main() {
	port := os.Getenv("PORT")	
	if port == "" {
		port = "8000"
	}

	router := gin.New()

	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())

	routes.BlogRoutes(router)
	routes.TagRoutes(router)
	routes.LoginRoutes(router)

	router.Use(middleware.Authentication())
	routes.AdminRoutes(router)

	log.Fatal(router.Run(":" + port))
}