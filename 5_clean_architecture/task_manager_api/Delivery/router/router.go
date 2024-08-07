package router

import (
	"fmt"
	"task_manager_api/Delivery/controllers"
	repository "task_manager_api/Repository"
	usecase "task_manager_api/Usecase"
	middlewares "task_manager_api/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Creates a router, attaches all the endpoints and finally
runs the API with the provided port number.
*/
func CreateRouter(port int, db *mongo.Database) {
	router := gin.Default()

	// task API
	taskRouter := router.Group("/tasks")
	timeout := viper.GetInt("TIMEOUT")
	NewTaskController(time.Duration(timeout)*time.Second, db.Collection(viper.GetString("tasks")), taskRouter)

	// user registeration and login
	authRouter := router.Group("")
	NewAuthController(time.Duration(timeout)*time.Second, db.Collection(viper.GetString("users")), authRouter)

	router.Run(fmt.Sprintf(":%v", port))
}

func NewTaskController(timeout time.Duration, collection *mongo.Collection, group *gin.RouterGroup) {
	taskUsecase := usecase.TaskUsecase{
		TaskRepository: &repository.TaskRepository{
			Collection: collection,
		},
		Timeout: timeout,
	}
	taskController := controllers.TaskController{
		TaskUsecase: &taskUsecase,
	}

	group.GET("", taskController.GetAll)
	group.GET("/:id", taskController.GetOne)
	group.POST("/tasks", middlewares.AuthMiddlewareWithRoles([]string{"user", "admin"}), taskController.Create)
	group.PUT("/tasks/:id", middlewares.AuthMiddlewareWithRoles([]string{"user", "admin"}), taskController.Update)
	group.DELETE("/tasks/:id", middlewares.AuthMiddlewareWithRoles([]string{"admin"}), taskController.Delete)
}

func NewAuthController(timeout time.Duration, collection *mongo.Collection, group *gin.RouterGroup) {
	authUsecase := usecase.UserUsecase{
		UserRespository: &repository.UserRepository{
			Collection: collection,
		},
		Timeout: timeout,
	}
	authController := controllers.UserController{
		UserUsecase: &authUsecase,
	}

	group.POST("/signup", authController.Signup)
	group.POST("/login", authController.Login)
}
