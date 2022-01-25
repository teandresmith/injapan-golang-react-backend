package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/controllers"
)

func LoginRoutes(routes *gin.Engine) {
	routes.POST("/api/admin/login", controllers.Login())
}