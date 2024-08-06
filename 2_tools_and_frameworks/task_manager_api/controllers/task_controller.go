package controllers

import (
	"net/http"
	services "task_manager_api/data"
	"task_manager_api/models"

	"github.com/gin-gonic/gin"
)

// handler for GET /tasks
func GetAll(c *gin.Context) {
	tasks := services.GetAllTasks()
	c.JSON(http.StatusOK, tasks)
}

// handler for GET /tasks/:id
func GetOne(c *gin.Context) {
	id := c.Param("id")
	task, err := services.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

// handler for POST /tasks
func Create(c *gin.Context) {
	var newTask models.Task
	if err := c.Bind(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error during object binding"})
		return
	}

	services.AddTask(newTask)
	c.JSON(http.StatusCreated, newTask)
}

// handler for PUT /tasks/:id
func Update(c *gin.Context) {
	var updatedTask models.Task
	id := c.Param("id")
	if err := c.Bind(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error during object binding"})
		return
	}

	newTask, err := services.UpdateTask(updatedTask, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, newTask)
}

// handler for DELETE /tasks/:id
func Delete(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task with ID not found"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Task removed"})
}
