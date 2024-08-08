package repository

import (
	"context"
	domain "task_manager_api/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	Collection *mongo.Collection
}

// retrieves all the tasks in the db
func (tR *TaskRepository) GetAllTasks(c context.Context) ([]domain.Task, domain.CodedError) {
	cursor, queryErr := tR.Collection.Find(c, bson.D{{}})
	if queryErr != nil {
		return []domain.Task{}, domain.TaskError{Message: "Internal server error: " + queryErr.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	tasks := []domain.Task{}
	bindErr := cursor.All(c, &tasks)
	if bindErr != nil {
		return []domain.Task{}, domain.TaskError{Message: "Internal server error: " + bindErr.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	cursor.Close(c)
	return tasks, nil
}

// retrieves the task associated with the provided id if it exists
func (tR *TaskRepository) GetTaskByID(c context.Context, taskID string) (domain.Task, domain.CodedError) {
	var task domain.Task
	result := tR.Collection.FindOne(c, bson.D{{Key: "id", Value: taskID}})
	if result.Err() != nil && result.Err().Error() == mongo.ErrNoDocuments.Error() {
		return task, domain.TaskError{Message: "Task not found", Code: domain.ERR_NOT_FOUND}
	}

	bindErr := result.Decode(&task)
	if bindErr != nil {
		return task, domain.TaskError{Message: "Internal server error: " + bindErr.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return task, nil
}

// adds the provided task to the database
func (tR *TaskRepository) AddTask(c context.Context, newTask domain.Task) domain.CodedError {
	result := tR.Collection.FindOne(c, bson.D{{Key: "id", Value: newTask.ID}})
	if result.Err() == nil {
		return domain.TaskError{Message: "Task with the provided ID already exists", Code: domain.ERR_BAD_REQUEST}
	}

	if result.Err().Error() != mongo.ErrNoDocuments.Error() {
		return domain.TaskError{Message: "Internal server error: " + result.Err().Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	_, err := tR.Collection.InsertOne(c, newTask)
	if err != nil {
		return domain.TaskError{Message: "Internal server error: " + err.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return nil
}

// updates the task associated with the provided id with the parameters provided in the provided task struct
func (tR *TaskRepository) UpdateTask(c context.Context, taskID string, updatedTask domain.Task) (domain.Task, domain.CodedError) {
	var setAttributes bson.D
	var task domain.Task

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

	result := tR.Collection.FindOneAndUpdate(c, bson.D{{Key: "id", Value: taskID}}, bson.D{
		{Key: "$set", Value: setAttributes},
	})

	if result.Err() != nil && result.Err().Error() == mongo.ErrNoDocuments.Error() {
		return task, domain.TaskError{Message: "Task not found", Code: domain.ERR_NOT_FOUND}
	}

	// fetch the updated task to get the latest version of the updated task
	newTask, _ := tR.GetTaskByID(c, taskID)
	return newTask, nil
}

// deletes the task associated with the provided id if it exists
func (tR *TaskRepository) DeleteTask(c context.Context, taskID string) domain.CodedError {
	result := tR.Collection.FindOneAndDelete(c, bson.D{{Key: "id", Value: taskID}})

	if result.Err() != nil && result.Err().Error() == mongo.ErrNoDocuments.Error() {
		return domain.TaskError{Message: "Task not found", Code: domain.ERR_NOT_FOUND}
	}

	if result.Err() != nil {
		return domain.TaskError{Message: "Internal server error: " + result.Err().Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return nil
}
