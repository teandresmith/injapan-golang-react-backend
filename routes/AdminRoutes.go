package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/controllers"
)


func AdminRoutes(routes *gin.Engine) {
	routes.POST("/api/admin/create", controllers.CreateAdmin())
	routes.POST("/api/admin/blogs/create", controllers.CreateBlog())
	routes.PATCH("/api/admin/blogs/update/:blogid", controllers.UpdateBlog())
	routes.DELETE("/api/admin/blogs/delete/:blogid", controllers.DeleteBlog())
	routes.POST("/api/admin/tags/create", controllers.CreateTag())
	routes.PATCH("/api/admin/tags/update/:tagid", controllers.UpdateTag())
	routes.DELETE("/api/admin/tags/delete/:tagid", controllers.DeleteTag())
}