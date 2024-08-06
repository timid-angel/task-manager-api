package router

import (
	"fmt"
	"task_manager_api/controllers"

	"github.com/gin-gonic/gin"
)

/*
Creates a router, attaches all the endpoints and finally
runs the API with the provided port number.
*/
func CreateRouter(port int) {
	router := gin.Default()

	// route to check the
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// tasks REST end points
	router.GET("/tasks", controllers.GetAll)
	router.GET("/tasks/:id", controllers.GetOne)
	router.POST("/tasks", controllers.Create)
	router.PUT("/tasks/:id", controllers.Update)
	router.DELETE("/tasks/:id", controllers.Delete)

	router.Run(fmt.Sprintf(":%v", port))
}
