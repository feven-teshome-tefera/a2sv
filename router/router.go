package router

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", controllers.GetAll)
	r.GET("/tasks/:id", controllers.GetByID)
	r.POST("/tasks", controllers.Create)
	r.PUT("/tasks/:id", controllers.Update)
	r.DELETE("/tasks/:id", controllers.Delete)

	return r
}
