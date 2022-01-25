package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/controllers"
)


func BlogRoutes(routes *gin.Engine) {
	routes.GET("/api/blogs", controllers.GetAllBlogs())
	routes.GET("/api/blogs/:blogid", controllers.GetBlogByID())
	routes.GET("/api/blogs/query", controllers.GetBlogsByParameters())
	routes.POST("/api/blogs/email", controllers.SendEmail())
}