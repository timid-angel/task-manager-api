package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// Mock data for tasks
var tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

// handlers
func getAll(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func getOne(c *gin.Context) {
	id := c.Param("id")
	for _, v := range tasks {
		if v.ID == id {
			c.JSON(http.StatusOK, v)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Task with ID not found"})
}

func create(c *gin.Context) {
	var newTask Task
	if err := c.Bind(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error during object binding"})
		return
	}

	tasks = append(tasks, newTask)
	c.JSON(http.StatusCreated, newTask)
}

func update(c *gin.Context) {
	var updatedTask Task
	id := c.Param("id")
	if err := c.Bind(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error during object binding"})
		return
	}

	idx := -1
	for i, v := range tasks {
		if v.ID == id {
			idx = i
		}
	}

	if idx != -1 {
		if updatedTask.Title != "" {
			tasks[idx].Title = updatedTask.Title
		}
		if updatedTask.Description != "" {
			tasks[idx].Description = updatedTask.Description
		}

		c.JSON(http.StatusOK, tasks[idx])
		return
	}
}

func delete(c *gin.Context) {
	id := c.Param("id")
	for i, v := range tasks {
		if v.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusNoContent, gin.H{"message": "Task removed"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Task with ID not found"})
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/tasks", getAll)
	router.GET("/tasks/:id", getOne)
	router.POST("/tasks", create)
	router.PUT("/tasks/:id", update)
	router.DELETE("/tasks/:id", delete)

	router.Run(":8080")
}
