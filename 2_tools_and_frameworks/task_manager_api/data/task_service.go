package services

import (
	"task_manager_api/models"
	"time"
)

/*
A struct that implements the `error` interface.
Created to enable the exchange of error messages
and signals between the different sections of the
services sub-package.
*/
type ServiceError struct {
	message string
}

func (err ServiceError) Error() string {
	return err.message
}

// Sample data for the in-memory data store
var tasks = []models.Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

// retrieves all the tasks in the db
func GetAllTasks() []models.Task {
	return tasks
}

// retrieves the task associated with the provided id if it exists
func GetTaskByID(id string) (models.Task, error) {
	for _, v := range tasks {
		if v.ID == id {
			return v, nil
		}
	}

	return models.Task{}, ServiceError{message: "Task with ID not found"}
}

// adds the provided task to the database
func AddTask(newTask models.Task) {
	tasks = append(tasks, newTask)
}

// updates the task associated with the provided id with the parameters provided in the provided task struct
func UpdateTask(updatedTask models.Task, id string) (models.Task, error) {
	idx := -1
	for i, v := range tasks {
		if v.ID == id {
			idx = i
		}
	}

	// task id not found
	if idx == -1 {
		return models.Task{}, ServiceError{message: "Task with ID not found"}
	}

	// check if the fields of the incoming struct are not the default values
	// and make the changes accordingly
	if updatedTask.Title != "" {
		tasks[idx].Title = updatedTask.Title
	}
	if updatedTask.Description != "" {
		tasks[idx].Description = updatedTask.Description
	}
	if updatedTask.Status != "" {
		tasks[idx].Status = updatedTask.Status
	}
	if !updatedTask.DueDate.IsZero() {
		tasks[idx].DueDate = updatedTask.DueDate
	}

	return tasks[idx], nil
}

// deletes the task associated with the provided id if it exists
func DeleteTask(id string) error {
	for i, v := range tasks {
		if v.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}

	return ServiceError{message: "Task with ID not found"}
}
