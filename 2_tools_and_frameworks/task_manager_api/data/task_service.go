package services

import (
	"task_manager_api/models"
	"time"
)

type ServiceError struct {
	message string
}

func (err ServiceError) Error() string {
	return err.message
}

var tasks = []models.Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func GetAllTasks() []models.Task {
	return tasks
}

func GetTaskByID(id string) (models.Task, error) {
	for _, v := range tasks {
		if v.ID == id {
			return v, nil
		}
	}

	return models.Task{}, ServiceError{message: "Task with ID not found"}
}

func AddTask(newTask models.Task) {
	tasks = append(tasks, newTask)
}

func UpdateTask(updatedTask models.Task, id string) (models.Task, error) {
	idx := -1
	for i, v := range tasks {
		if v.ID == id {
			idx = i
		}
	}

	if idx == -1 {
		return models.Task{}, ServiceError{message: "Task with ID not found"}
	}

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

func DeleteTask(id string) error {
	for i, v := range tasks {
		if v.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}

	return ServiceError{message: "Task with ID not found"}
}
