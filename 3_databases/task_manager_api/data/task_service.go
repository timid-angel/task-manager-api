package services

import (
	"context"
	"task_manager_api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Interface used to define structs that compose the standard error interface
with an function used to obtain an error code.
*/
type CodedError interface {
	error
	GetCode() int
}

/*
A struct that implements the `error` interface.
Created to enable the exchange of error messages
and signals between the different sections of the
services sub-package.
*/
type ServiceError struct {
	message string
	code    int
}

func (err ServiceError) Error() string {
	return err.message
}

func (err ServiceError) GetCode() int {
	return err.code
}

// retrieves all the tasks in the db
func GetAllTasks() ([]models.Task, error) {
	cursor, err := TaskCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return []models.Task{}, err
	}

	tasks := []models.Task{}

	// loop through the cursor and add items to the slice
	for cursor.Next(context.TODO()) {
		var task models.Task
		err := cursor.Decode(&task)
		if err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}

	cursor.Close(context.TODO())
	return tasks, nil
}

// retrieves the task associated with the provided id if it exists
func GetTaskByID(id string) (models.Task, error) {
	var task models.Task
	result := TaskCollection.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	err := result.Decode(&task)
	if err != nil {
		return task, err
	}

	return task, nil
}

// adds the provided task to the database
func AddTask(newTask models.Task) CodedError {
	_, queryErr := GetTaskByID(newTask.ID)
	if queryErr == nil {
		return ServiceError{message: "A task with the provided ID already exists", code: 400}
	}

	if queryErr.Error() != mongo.ErrNoDocuments.Error() {
		return ServiceError{message: "Internal server error: " + queryErr.Error(), code: 500}
	}

	_, err := TaskCollection.InsertOne(context.TODO(), newTask)
	if err != nil {
		return ServiceError{message: "Internal server error: " + queryErr.Error(), code: 500}
	}

	return nil
}

// updates the task associated with the provided id with the parameters provided in the provided task struct
func UpdateTask(updatedTask models.Task, id string) (models.Task, error) {
	var setAttributes bson.D
	var task models.Task

	// check if the fields of the incoming struct are not the default values
	// and append the results to setAttributes accordingly
	if updatedTask.Title != "" {
		setAttributes = append(setAttributes, bson.E{Key: "title", Value: updatedTask.Title})
	}
	if updatedTask.Description != "" {
		setAttributes = append(setAttributes, bson.E{Key: "description", Value: updatedTask.Description})
	}
	if updatedTask.Status != "" {
		setAttributes = append(setAttributes, bson.E{Key: "status", Value: updatedTask.Status})
	}
	if !updatedTask.DueDate.IsZero() {
		setAttributes = append(setAttributes, bson.E{Key: "due_date", Value: updatedTask.DueDate})
	}

	result := TaskCollection.FindOneAndUpdate(context.TODO(), bson.D{{Key: "id", Value: id}}, bson.D{
		{Key: "$set", Value: setAttributes},
	})

	if result.Err() != nil {
		return task, result.Err()
	}

	// fetch the updated task to get the latest version of the updated task
	newTask, err := GetTaskByID(id)
	if result.Err() != nil {
		return newTask, err
	}

	return newTask, nil
}

// deletes the task associated with the provided id if it exists
func DeleteTask(id string) error {
	result := TaskCollection.FindOneAndDelete(context.TODO(), bson.D{{Key: "id", Value: id}})
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
