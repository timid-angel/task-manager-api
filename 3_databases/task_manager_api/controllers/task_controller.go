package controllers

import (
	"fmt"
	"log"
	"net/http"
	services "task_manager_api/data"
	"task_manager_api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetErrorCode(err error) int {
	log.Printf("%v, %T", err, err)
	switch err {
	case mongo.ErrNoDocuments, mongo.ErrNilDocument, mongo.ErrNilCursor:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func GetAll(c *gin.Context) {
	tasks, err := services.GetAllTasks()
	if err != nil {
		c.JSON(GetErrorCode(err), gin.H{"message": fmt.Sprintf("Error: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func GetOne(c *gin.Context) {
	id := c.Param("id")
	task, err := services.GetTaskByID(id)
	if err != nil {
		c.JSON(GetErrorCode(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func Create(c *gin.Context) {
	var newTask models.Task
	if err := c.Bind(&newTask); err != nil {
		c.JSON(GetErrorCode(err), gin.H{"message": "Error during object binding"})
		return
	}

	services.AddTask(newTask)
	c.JSON(http.StatusCreated, newTask)
}

func Update(c *gin.Context) {
	var updatedTask models.Task
	id := c.Param("id")
	if err := c.Bind(&updatedTask); err != nil {
		c.JSON(GetErrorCode(err), gin.H{"message": "Error during object binding"})
		return
	}

	newTask, err := services.UpdateTask(updatedTask, id)
	if err != nil {
		c.JSON(GetErrorCode(err), gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTask)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteTask(id)
	if err != nil {
		c.JSON(GetErrorCode(err), gin.H{"message": "Task with ID not found"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Task removed"})
}
