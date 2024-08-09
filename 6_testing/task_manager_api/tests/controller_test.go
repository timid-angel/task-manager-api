package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task_manager_api/Delivery/controllers"
	domain "task_manager_api/Domain"
	"task_manager_api/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type controllerSuite struct {
	suite.Suite
	taskUsecase    *mocks.TaskUsecaseInterface
	userUsecase    *mocks.UserUsecaseInterface
	taskController controllers.TaskController
	userController controllers.UserController
	testingServer  *httptest.Server
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (suite *controllerSuite) SetupSuite() {
	suite.taskUsecase = new(mocks.TaskUsecaseInterface)
	suite.userUsecase = new(mocks.UserUsecaseInterface)
	suite.taskController = controllers.TaskController{
		TaskUsecase: suite.taskUsecase,
	}

	suite.userController = controllers.UserController{
		UserUsecase: suite.userUsecase,
	}

	router := gin.Default()

	router.GET("/tasks", suite.taskController.GetAll)
	router.GET("/tasks/:id", suite.taskController.GetOne)
	router.POST("/tasks", suite.taskController.Create)
	router.PUT("/tasks/:id", suite.taskController.Update)
	router.DELETE("/tasks/:id", suite.taskController.Delete)

	router.POST("/signup", suite.userController.Signup)
	router.POST("/login", suite.userController.Login)

	suite.testingServer = httptest.NewServer(router)
}

func (suite *controllerSuite) SetupTest() {
	suite.taskUsecase = new(mocks.TaskUsecaseInterface)
	suite.userUsecase = new(mocks.UserUsecaseInterface)
	suite.taskController.TaskUsecase = suite.taskUsecase
	suite.userController.UserUsecase = suite.userUsecase
}

func (suite *controllerSuite) TearDownSuite() {
	defer suite.testingServer.Close()
}

func (suite *controllerSuite) TestGetHTTPErrorCodes() {
	testParams := map[domain.CodedError]int{
		domain.TaskError{Code: domain.ERR_BAD_REQUEST}:     400,
		domain.TaskError{Code: domain.ERR_INTERNAL_SERVER}: 500,
		domain.TaskError{Code: domain.ERR_NOT_FOUND}:       404,
		domain.TaskError{Code: domain.ERR_UNAUTHORIZED}:    401,
	}

	for domainErr, statusCode := range testParams {
		suite.Equal(statusCode, controllers.GetHTTPErrorCode(domainErr))
	}
}

func (suite *controllerSuite) TestGetAllTasks_Positive() {
	task := domain.Task{
		ID:          "1",
		Title:       "title",
		Description: "description",
		DueDate:     time.Now().Round(0),
		Status:      "pending",
	}

	suite.taskUsecase.On("GetAllTasks", mock.Anything).Return([]domain.Task{task}, nil)
	response, err := http.Get(suite.testingServer.URL + "/tasks")
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())

	var tasks []domain.Task
	err = json.NewDecoder(response.Body).Decode(&tasks)
	suite.NoError(err, "no error during body decoding")
	suite.Equal(1, len(tasks), "sends data correctly")
}

func (suite *controllerSuite) TestGetAllTasks_Negative() {
	task := domain.Task{
		ID:          "1",
		Title:       "title",
		Description: "description",
		DueDate:     time.Now().Round(0),
		Status:      "pending",
	}

	sampleErr := domain.TaskError{Message: "msg123", Code: domain.ERR_INTERNAL_SERVER}
	suite.taskUsecase.On("GetAllTasks", mock.Anything).Return([]domain.Task{task}, sampleErr)
	response, err := http.Get(suite.testingServer.URL + "/tasks")
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(controllers.GetHTTPErrorCode(sampleErr), response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestGetTaskByID_Positive() {
	task := domain.Task{
		ID:          "1",
		Title:       "title",
		Description: "description",
		DueDate:     time.Now().Round(0),
		Status:      "pending",
	}

	suite.taskUsecase.On("GetTaskByID", mock.Anything, task.ID).Return(task, nil)
	response, err := http.Get(suite.testingServer.URL + "/tasks/" + task.ID)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())

	var fetchedTask domain.Task
	err = json.NewDecoder(response.Body).Decode(&fetchedTask)
	suite.NoError(err, "no error during body decoding")
	suite.Equal(task.ID, fetchedTask.ID, "sends task with correct ID")
	suite.Equal(task.Title, fetchedTask.Title, "sends task with correct Title")
	suite.Equal(task.Description, fetchedTask.Description, "sends task with correct Description")
	suite.Equal(task.DueDate.Unix(), fetchedTask.DueDate.Unix(), "sends task with correct DueDate")
	suite.Equal(task.Status, fetchedTask.Status, "sends task with correct Status")
}

func (suite *controllerSuite) TestGetTaskByID_Negative() {
	wrongID := "2"
	sampleErr := domain.TaskError{Message: "msg123", Code: domain.ERR_BAD_REQUEST}
	suite.taskUsecase.On("GetTaskByID", mock.Anything, wrongID).Return(domain.Task{}, sampleErr)
	response, err := http.Get(suite.testingServer.URL + "/tasks/" + wrongID)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(controllers.GetHTTPErrorCode(sampleErr), response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestAdd_Positive() {
	newTask := domain.Task{}
	client := http.Client{}
	suite.taskUsecase.On("AddTask", mock.Anything, newTask).Return(nil)

	requestBody, err := json.Marshal(&newTask)
	suite.NoError(err, "can not marshal struct to json")

	request, _ := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/tasks", bytes.NewBuffer(requestBody))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(http.StatusCreated, response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestAdd_Negative() {
	newTask := domain.Task{}
	client := http.Client{}
	sampleErr := domain.TaskError{Message: "msg123", Code: domain.ERR_BAD_REQUEST}
	suite.taskUsecase.On("AddTask", mock.Anything, newTask).Return(sampleErr)
	requestBody, err := json.Marshal(&newTask)
	suite.NoError(err, "can not marshal struct to json")

	request, _ := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/tasks", bytes.NewBuffer(requestBody))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(controllers.GetHTTPErrorCode(sampleErr), response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestUpdate_Positive() {
	taskID := "1"
	taskUpdates := domain.Task{
		Title:       "title",
		Description: "description",
		Status:      "pending",
	}

	client := http.Client{}
	suite.taskUsecase.On("UpdateTask", mock.Anything, taskID, taskUpdates).Return(taskUpdates, nil)

	requestBody, err := json.Marshal(&taskUpdates)
	suite.NoError(err, "can not marshal struct to json")

	request, _ := http.NewRequest(http.MethodPut, suite.testingServer.URL+"/tasks/"+taskID, bytes.NewBuffer(requestBody))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())

	var fetchedTask domain.Task
	err = json.NewDecoder(response.Body).Decode(&fetchedTask)
	suite.NoError(err, "no error during body decoding")
	suite.Equal(taskUpdates.Title, fetchedTask.Title, "sends task with correct Title")
	suite.Equal(taskUpdates.Description, fetchedTask.Description, "sends task with correct Description")
	suite.Equal(taskUpdates.Status, fetchedTask.Status, "sends task with correct Status")
}

func (suite *controllerSuite) TestUpdate_Negative() {
	taskID := "1"
	taskUpdates := domain.Task{}
	client := http.Client{}
	sampleErr := domain.TaskError{Message: "msg123", Code: domain.ERR_BAD_REQUEST}
	suite.taskUsecase.On("UpdateTask", mock.Anything, taskID, mock.AnythingOfType("Task")).Return(taskUpdates, sampleErr)

	requestBody, err := json.Marshal(&taskUpdates)
	suite.NoError(err, "can not marshal struct to json")

	request, _ := http.NewRequest(http.MethodPut, suite.testingServer.URL+"/tasks/"+taskID, bytes.NewBuffer(requestBody))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(controllers.GetHTTPErrorCode(sampleErr), response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestDelete_Positive() {
	taskID := "1"
	client := http.Client{}
	suite.taskUsecase.On("DeleteTask", mock.Anything, taskID).Return(nil)
	request, _ := http.NewRequest(http.MethodDelete, suite.testingServer.URL+"/tasks/"+taskID, nil)
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(http.StatusNoContent, response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestDelete_Negative() {
	wrongID := "1"
	sampleErr := domain.TaskError{Message: "msg123", Code: domain.ERR_BAD_REQUEST}
	client := http.Client{}
	suite.taskUsecase.On("DeleteTask", mock.Anything, wrongID).Return(sampleErr)
	request, _ := http.NewRequest(http.MethodDelete, suite.testingServer.URL+"/tasks/"+wrongID, nil)
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(controllers.GetHTTPErrorCode(sampleErr), response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestSignup_Positive() {
	user := domain.User{
		Username: "lksdajf",
		Email:    "valid@mail.com",
		Password: "dorwssap",
		Role:     "admin",
	}

	client := http.Client{}
	suite.userUsecase.On("CreateUser", mock.Anything, user).Return(nil)

	requestBody, err := json.Marshal(&user)
	suite.NoError(err, "can not marshal struct to json")

	request, _ := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/signup", bytes.NewBuffer(requestBody))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(http.StatusCreated, response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestSignup_Negative() {
	user := domain.User{}
	client := http.Client{}
	sampleErr := domain.TaskError{Message: "msg123", Code: domain.ERR_INTERNAL_SERVER}
	suite.userUsecase.On("CreateUser", mock.Anything, user).Return(sampleErr)

	requestBody, err := json.Marshal(&user)
	suite.NoError(err, "can not marshal struct to json")

	request, _ := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/signup", bytes.NewBuffer(requestBody))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	suite.Equal(controllers.GetHTTPErrorCode(sampleErr), response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func (suite *controllerSuite) TestLogin_Positive() {
	user := domain.User{}
	sToken := "lskad123i12.3123123sadf"
	client := http.Client{}
	suite.userUsecase.On("ValidateAndGetToken", mock.Anything, user).Return(sToken, nil)

	requestBody, err := json.Marshal(&user)
	suite.NoError(err, "can not marshal struct to json")

	request, _ := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/login", bytes.NewBuffer(requestBody))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	suite.NoError(err, "no errors in request")
	responseBody := TokenResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	suite.NoError(err, "no error during body decoding")
	suite.Equal(sToken, responseBody.Token)
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.taskUsecase.AssertExpectations(suite.T())
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(controllerSuite))
}
