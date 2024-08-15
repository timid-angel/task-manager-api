package usecase

import (
	"context"
	domain "task_manager_api/Domain"
	"time"
)

/* Implements the TaskUsecaseInterface defined in `domain`*/
type TaskUsecase struct {
	TaskRepository domain.TaskRepositoryInterface
	Timeout        time.Duration
}

/* Calls GetAllTasks in the repository after setting the timeout */
func (tU *TaskUsecase) GetAllTasks(c context.Context) ([]domain.Task, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.GetAllTasks(ctx)
}

/* Calls GetTaskById in the repository after setting the timeout */
func (tU *TaskUsecase) GetTaskByID(c context.Context, taskID string) (domain.Task, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.GetTaskByID(ctx, taskID)
}

/*
Checks if a task with a similar ID exists before calling AddTask in the repository
with the provided user data after setting the timeout.
*/
func (tU *TaskUsecase) AddTask(c context.Context, newTask domain.Task) domain.CodedError {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()

	_, err := tU.TaskRepository.GetTaskByID(c, newTask.ID)
	if err == nil {
		return domain.TaskError{Message: "Task with the provided ID already exists", Code: domain.ERR_BAD_REQUEST}
	}

	if err.GetCode() != domain.ERR_NOT_FOUND {
		return err
	}

	return tU.TaskRepository.AddTask(ctx, newTask)
}

/*
Calls UpdateTask in the repository with the provided ID and updated data
after setting the timeout
*/
func (tU *TaskUsecase) UpdateTask(c context.Context, taskID string, updatedTask domain.Task) (domain.Task, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.UpdateTask(ctx, taskID, updatedTask)
}

/* Calls DeleteTask with the provided ID in the repository after setting the timeout */
func (tU *TaskUsecase) DeleteTask(c context.Context, taskID string) domain.CodedError {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.DeleteTask(ctx, taskID)
}
