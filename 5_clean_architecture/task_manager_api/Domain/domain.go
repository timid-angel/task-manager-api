package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CollectionTasks     = "tasks"
	CollectionUsers     = "users"
	ERR_NOT_FOUND       = "not_found"
	ERR_INTERNAL_SERVER = "internal_server_error"
	ERR_BAD_REQUEST     = "bad_request"
	ERR_UNAUTHORIZED    = "unauthorized"
)

/*
Interface used to define structs that compose the standard error interface
with an function used to obtain an error code.
*/
type CodedError interface {
	error
	GetCode() string
}

/*
A struct that implements the `CodedError` interface. Created to enable the
exchange of error messages and signals between the different sections of
the task API.
*/
type TaskError struct {
	Message string
	Code    string
}

func (err TaskError) Error() string {
	return err.Message
}

func (err TaskError) GetCode() string {
	return err.Code
}

/*
A struct that implements the `CodedError` interface. Created to enable the
exchange of error messages and signals between the different sections of
the functionalities of user auth.
*/
type UserError struct {
	Message string
	Code    string
}

func (err UserError) Error() string {
	return err.Message
}

func (err UserError) GetCode() string {
	return err.Code
}

/*
This is the definition of the task struct that will be used
throughout the application. Along with the field names, the
json labels are provided to facilitate the binding process
between the model itself and the JSON format.
*/
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

/*
This is the definition of the user struct used for the authentication
and authorization aspects of the project. The email and user name will
be unique through entries. Along with the field names, the json labels
are provided to facilitate the binding process between the model itself
and the JSON format.
*/
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type TaskUsecaseInterface interface {
	GetAllTasks(c context.Context) ([]Task, CodedError)
	GetTaskByID(c context.Context, taskID string) (Task, CodedError)
	AddTask(c context.Context, newTask Task) CodedError
	UpdateTask(c context.Context, taskID string, updatedTask Task) (Task, CodedError)
	DeleteTask(c context.Context, taskID string) CodedError
}

type TaskRepositoryInterface interface {
	GetAllTasks(c context.Context) ([]Task, CodedError)
	GetTaskByID(c context.Context, taskID string) (Task, CodedError)
	AddTask(c context.Context, newTask Task) CodedError
	UpdateTask(c context.Context, taskID string, updatedTask Task) (Task, CodedError)
	DeleteTask(c context.Context, taskID string) CodedError
}

type UserUsecaseInterface interface {
	CreateUser(c context.Context, user User) CodedError
	ValidateAndGetToken(c context.Context, user User) (string, CodedError)
}

type UserRepositoryInterface interface {
	CreateUser(c context.Context, user User) CodedError
	ValidateAndGetToken(c context.Context, user User) (string, CodedError)
}

type Response gin.H
