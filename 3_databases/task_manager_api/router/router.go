package router

import (
	"task_manager_api/controllers"

	"github.com/gin-gonic/gin"
)

func CreateRouter() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/tasks", controllers.GetAll)
	router.GET("/tasks/:id", controllers.GetOne)
	router.POST("/tasks", controllers.Create)
	router.PUT("/tasks/:id", controllers.Update)
	router.DELETE("/tasks/:id", controllers.Delete)

	router.Run(":8080")
}
