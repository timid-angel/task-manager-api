package router

import (
	"context"
	"fmt"
	"task_manager_api/Delivery/controllers"
	domain "task_manager_api/Domain"
	infrastructure "task_manager_api/Infrastructure"
	repository "task_manager_api/Repository"
	usecase "task_manager_api/Usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	NewTaskController(time.Duration(timeout)*time.Second, db.Collection(domain.CollectionTasks), taskRouter)

	// user registeration and login
	authRouter := router.Group("")
	NewAuthController(time.Duration(timeout)*time.Second, db.Collection(domain.CollectionUsers), authRouter)

	router.Run(fmt.Sprintf(":%v", port))
}

/*
Attaches to the provided router group all the task endpoints with the
appropriate auth middleware configurations and creates all the task controller
that provides the handlers for the endpoints
*/
func NewTaskController(timeout time.Duration, collection *mongo.Collection, group *gin.RouterGroup) {
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		fmt.Println("\n\n Error " + err.Error())
	}

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
	group.POST("", infrastructure.AuthMiddlewareWithRoles([]string{"user", "admin"}), taskController.Create)
	group.PUT(":id", infrastructure.AuthMiddlewareWithRoles([]string{"user", "admin"}), taskController.Update)
	group.DELETE(":id", infrastructure.AuthMiddlewareWithRoles([]string{"admin"}), taskController.Delete)
}

/*
Attaches the `/login` and `/signup` routes along with the controller
that provides the handlers for those endpoints
*/
func NewAuthController(timeout time.Duration, collection *mongo.Collection, group *gin.RouterGroup) {
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		fmt.Println("\n\n Error " + err.Error())
	}

	_, err = collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		fmt.Println("\n\n Error " + err.Error())
	}

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
