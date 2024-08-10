package router

import (
	"fmt"
	"task_manager_api/controllers"
	middlewares "task_manager_api/middleware"

	"github.com/gin-gonic/gin"
)

/*
Creates a router, attaches all the endpoints and finally
runs the API with the provided port number.
*/
func CreateRouter(port int) {
	router := gin.Default()

	// route to check the up status of the API
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// tasks REST end points
	router.GET("/tasks", middlewares.AuthMiddlewareWithRoles([]string{"user", "admin"}), controllers.GetAll)
	router.GET("/tasks/:id", middlewares.AuthMiddlewareWithRoles([]string{"user", "admin"}), controllers.GetOne)
	router.POST("/tasks", middlewares.AuthMiddlewareWithRoles([]string{"admin"}), controllers.Create)
	router.PUT("/tasks/:id", middlewares.AuthMiddlewareWithRoles([]string{"admin"}), controllers.Update)
	router.DELETE("/tasks/:id", middlewares.AuthMiddlewareWithRoles([]string{"admin"}), controllers.Delete)

	// user registeration and login
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.PATCH("/promote/:username", middlewares.AuthMiddlewareWithRoles([]string{"admin"}), controllers.Promote)

	router.Run(fmt.Sprintf(":%v", port))
}
