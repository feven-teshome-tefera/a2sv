package router

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/tasks/register", controllers.Register)
	r.POST("/tasks/login", controllers.Login)
	taskRoutes := r.Group("/tasks")
	taskRoutes.Use(controllers.AuthMiddleware())
	{
		taskRoutes.GET("", controllers.GetAll)
		taskRoutes.GET("/:id", controllers.GetByID)
		taskRoutes.POST("", controllers.Create)
		taskRoutes.PUT("/:id", controllers.Update)
		taskRoutes.DELETE("/:id", controllers.Delete)
	}
	userRoutes := r.Group("/users")
	userRoutes.Use(controllers.AuthMiddleware())
	{
		userRoutes.GET("", controllers.GetAllusers)
		userRoutes.DELETE("/:email", controllers.Deleteuser)
	}

	return r
}
