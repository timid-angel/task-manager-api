package tests

import (
	domain "task_manager_api/Domain"
	usecase "task_manager_api/Usecase"
	"task_manager_api/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type taskUsecaseSuite struct {
	suite.Suite
	repository *mocks.TaskRepositoryInterface
	usecase    usecase.TaskUsecase
}

func (suite *taskUsecaseSuite) SetupSuite() {
	suite.repository = new(mocks.TaskRepositoryInterface)
	suite.usecase = usecase.TaskUsecase{
		TaskRepository: suite.repository,
		Timeout:        2,
	}
}

func (suite *taskUsecaseSuite) SetupTest() {
	suite.repository = new(mocks.TaskRepositoryInterface)
	suite.usecase.TaskRepository = suite.repository
}

func (suite *taskUsecaseSuite) TestGetAllTasks() {
	suite.repository.On("GetAllTasks", mock.Anything).Return([]domain.Task{}, nil).Twice()
	_, err := suite.usecase.GetAllTasks(context.TODO())

	suite.NoError(err, "no error when function is called")
	suite.repository.AssertCalled(suite.T(), "GetAllTasks", mock.Anything)
}

func (suite *taskUsecaseSuite) TestGetTaskByID() {
	taskID := "sample_id"
	suite.repository.On("GetTaskByID", mock.Anything, taskID).Return(domain.Task{}, nil).Twice()
	_, err := suite.usecase.GetTaskByID(context.TODO(), taskID)

	suite.NoError(err, "no error when function is called")
	suite.repository.AssertCalled(suite.T(), "GetTaskByID", mock.Anything, taskID)
}

// A
func (suite *taskUsecaseSuite) TestAddTask_Negative() {
	newTask := domain.Task{
		ID:          "4",
		Title:       "updated title",
		Description: "updated description",
		Status:      "completed",
		DueDate:     time.Now(),
	}

	suite.repository.On("AddTask", mock.Anything, newTask).Return(nil).Twice()
	suite.repository.On("GetTaskByID", mock.Anything, newTask.ID).Return(domain.Task{}, nil).Twice()
	err := suite.usecase.AddTask(context.TODO(), newTask)

	suite.Error(err, "error when GetTaskByID return no errors")
	suite.repository.AssertNotCalled(suite.T(), "AddTask", mock.Anything, newTask)
	suite.repository.AssertCalled(suite.T(), "GetTaskByID", mock.Anything, newTask.ID)
}

func (suite *taskUsecaseSuite) TestAddTask_Postive() {
	newTask := domain.Task{
		ID:          "4",
		Title:       "updated title",
		Description: "updated description",
		Status:      "completed",
		DueDate:     time.Now(),
	}

	suite.repository.On("AddTask", mock.Anything, newTask).Return(nil).Twice()
	suite.repository.On("GetTaskByID", mock.Anything, newTask.ID).Return(domain.Task{}, domain.UserError{Code: domain.ERR_NOT_FOUND}).Twice()
	err := suite.usecase.AddTask(context.TODO(), newTask)

	suite.NoError(err, "no error when GetTaskByID returns ERR_NOT_FOUND")
	suite.repository.AssertCalled(suite.T(), "GetTaskByID", mock.Anything, newTask.ID)
	suite.repository.AssertCalled(suite.T(), "AddTask", mock.Anything, newTask)
}

func (suite *taskUsecaseSuite) TestUpdateTask() {
	taskUpdates := domain.Task{
		Title:       "updated title",
		Description: "updated description",
		Status:      "completed",
	}

	taskID := "sample_id"
	suite.repository.On("UpdateTask", mock.Anything, taskID, taskUpdates).Return(domain.Task{}, nil).Twice()
	_, err := suite.usecase.UpdateTask(context.TODO(), taskID, taskUpdates)

	suite.NoError(err, "no error when function is called")
	suite.repository.AssertCalled(suite.T(), "UpdateTask", mock.Anything, taskID, taskUpdates)
}

func (suite *taskUsecaseSuite) TestDeleteTask() {
	taskID := "sample_id"
	suite.repository.On("DeleteTask", mock.Anything, taskID).Return(nil).Twice()
	err := suite.usecase.DeleteTask(context.TODO(), taskID)

	suite.NoError(err, "no error when function is called")
	suite.repository.AssertCalled(suite.T(), "DeleteTask", mock.Anything, taskID)
}

func TestTaskUsecase(t *testing.T) {
	suite.Run(t, new(taskUsecaseSuite))
}
