package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/controllers"
)


func TagRoutes(routes *gin.Engine) {
	routes.GET("/api/tag", controllers.GetAllTags())
}