package http

import (
	"async-api-task-manager/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", taskHandler.Health)

	tasks := r.Group("/tasks")
	{
		tasks.POST("/", taskHandler.CreateTask)
		tasks.GET("/:id", taskHandler.GetTaskByID)
		tasks.GET("/", taskHandler.ListTasks)
		tasks.PATCH("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}

	return r
}
