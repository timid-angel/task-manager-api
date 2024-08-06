package router

import (
	"task_manager_api/controllers"

	"github.com/gin-gonic/gin"
)

// Creates a router, attaches all the endpoints and starts running the app
func CreateRouter() {
	router := gin.Default()

	// route to check the up status of the API
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// tasks REST end points
	router.GET("/tasks", controllers.GetAll)
	router.GET("/tasks/:id", controllers.GetOne)
	router.POST("/tasks", controllers.Create)
	router.PUT("/tasks/:id", controllers.Update)
	router.DELETE("/tasks/:id", controllers.Delete)

	router.Run(":8080")
}
