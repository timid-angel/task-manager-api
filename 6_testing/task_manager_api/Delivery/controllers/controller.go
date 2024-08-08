package controllers

import (
	"net/http"
	domain "task_manager_api/Domain"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecaseInterface
}

type UserController struct {
	UserUsecase domain.UserUsecaseInterface
}

/*
Get the HTTP status code of an error based on the incoming error type.
This function checks for the mongoDB errors in particular and returns the
404 status code if the document can not be found
*/
func GetHTTPErrorCode(err domain.CodedError) int {
	switch err.GetCode() {
	case domain.ERR_BAD_REQUEST:
		return http.StatusBadRequest
	case domain.ERR_INTERNAL_SERVER:
		return http.StatusInternalServerError
	case domain.ERR_NOT_FOUND:
		return http.StatusNotFound
	case domain.ERR_UNAUTHORIZED:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// handler for GET /tasks
func (tC *TaskController) GetAll(c *gin.Context) {
	tasks, err := tC.TaskUsecase.GetAllTasks(c)
	if err != nil {
		c.JSON(GetHTTPErrorCode(err), domain.Response{"message": "Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// handler for GET /tasks/:id
func (tC *TaskController) GetOne(c *gin.Context) {
	id := c.Param("id")
	task, err := tC.TaskUsecase.GetTaskByID(c, id)
	if err != nil {
		c.JSON(GetHTTPErrorCode(err), domain.Response{"message": "Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// handler for POST /tasks
func (tC *TaskController) Create(c *gin.Context) {
	var newTask domain.Task
	if err := c.Bind(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{"message": "Error during object binding"})
		return
	}

	err := tC.TaskUsecase.AddTask(c, newTask)
	if err != nil {
		c.JSON(GetHTTPErrorCode(err), domain.Response{"message": "Error; " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

// handler for PUT /tasks/:id
func (tC *TaskController) Update(c *gin.Context) {
	var updatedTask domain.Task
	id := c.Param("id")
	if err := c.Bind(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{"message": "Error during object binding"})
		return
	}

	newTask, err := tC.TaskUsecase.UpdateTask(c, id, updatedTask)
	if err != nil {
		c.JSON(GetHTTPErrorCode(err), domain.Response{"message": "Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTask)
}

// handler for DELETE /tasks/:id
func (tC *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")
	err := tC.TaskUsecase.DeleteTask(c, id)
	if err != nil {
		c.JSON(GetHTTPErrorCode(err), domain.Response{"message": "Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, domain.Response{"message": "Task removed"})
}

// handler for /signup
func (uC *UserController) Signup(c *gin.Context) {
	var user domain.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{"message": "Error during object binding"})
		return
	}

	err := uC.UserUsecase.CreateUser(c, user)
	if err != nil {
		c.JSON(GetHTTPErrorCode(err), domain.Response{"message": "Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.Response{"message": "Signup successful"})
}

// handler for /login
func (uC *UserController) Login(c *gin.Context) {
	var user domain.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{"message": "Error during object binding"})
		return
	}

	token, err := uC.UserUsecase.ValidateAndGetToken(c, user)
	if err != nil {
		c.JSON(GetHTTPErrorCode(err), domain.Response{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response{"message": "User logged in successfully", "token": token})
}
