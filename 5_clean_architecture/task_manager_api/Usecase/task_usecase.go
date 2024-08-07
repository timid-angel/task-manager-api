package usecase

import (
	"context"
	domain "task_manager_api/Domain"
	"time"
)

type TaskUsecase struct {
	TaskRepository domain.TaskRepositoryInterface
	Timeout        time.Duration
}

func (tU *TaskUsecase) GetAllTasks(c context.Context) ([]domain.Task, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.GetAllTasks(ctx)
}

func (tU *TaskUsecase) GetTaskByID(c context.Context, taskID string) (domain.Task, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.GetTaskByID(ctx, taskID)
}

func (tU *TaskUsecase) AddTask(c context.Context, newTask domain.Task) domain.CodedError {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.AddTask(ctx, newTask)
}

func (tU *TaskUsecase) UpdateTask(c context.Context, taskID string, updatedTask domain.Task) (domain.Task, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.UpdateTask(ctx, taskID, updatedTask)
}

func (tU *TaskUsecase) DeleteTask(c context.Context, taskID string) domain.CodedError {
	ctx, cancel := context.WithTimeout(c, tU.Timeout)
	defer cancel()
	return tU.TaskRepository.DeleteTask(ctx, taskID)
}
