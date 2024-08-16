package tests

import (
	"context"
	"testing"
	"time"

	mocks "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/Mocks"
	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/usecases"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type taskUsecaseSuite struct {
	suite.Suite
	repositorie *mocks.TaskRepository
	usecase     domain.TaskUsecase
}

func (suite *taskUsecaseSuite) SetupTest() {
	repo := new(mocks.TaskRepository)
	taskUC := usecases.NewTaskUsecase(repo, time.Second*2)
	suite.usecase = &taskUC
	suite.repositorie = repo
}

func (suite *taskUsecaseSuite) TestGetAllTasks() {
	tasks := []domain.Task{
		{
			ID:          "task_001",
			UserID:      "user_123",
			Title:       "Complete Go project",
			Description: "Finish writing the Go code and add tests.",
			Status:      "In Progress",
			Priority:    "High",
			DueDate:     time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "task_002",
			UserID:      "user_456",
			Title:       "Review PRs",
			Description: "Review the pull requests for the backend repository.",
			Status:      "Pending",
			Priority:    "Medium",
			DueDate:     time.Date(2024, 8, 18, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "task_003",
			UserID:      "user_789",
			Title:       "Update Documentation",
			Description: "Update the API documentation to reflect recent changes.",
			Status:      "Completed",
			Priority:    "Low",
			DueDate:     time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now().AddDate(0, 0, -2),
			UpdatedAt:   time.Now().AddDate(0, 0, -1),
		},
	}

	suite.repositorie.On("FetchAllTasks", mock.Anything).Return(tasks, nil)

	fetchedTasks, err := suite.usecase.GetAllTasks(context.TODO())
	suite.Nil(err, "error should be nil")
	suite.Equal(tasks, fetchedTasks, "users should be equal")
}

func (suite *taskUsecaseSuite) TestGetTaskByID() {
	tasks := domain.Task{
		ID:          "task_001",
		UserID:      "user_123",
		Title:       "Complete Go project",
		Description: "Finish writing the Go code and add tests.",
		Status:      "In Progress",
		Priority:    "High",
		DueDate:     time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	suite.repositorie.On("FetchTaskByID", mock.Anything, tasks.ID).Return(tasks, nil)

	fetchedTask, err := suite.usecase.GetTaskByID(context.TODO(), tasks.ID)
	suite.Nil(err, "error should be nil")
	suite.Equal(tasks, fetchedTask, "tasks should be equal")
}

func (suite *taskUsecaseSuite) TestCreateTask_Positive() {
	ID := "task_001"
	tasks := domain.Task{
		UserID:      "user_123",
		Title:       "Complete Go project",
		Description: "Finish writing the Go code and add tests.",
		Status:      "In Progress",
		Priority:    "High",
		DueDate:     time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	suite.repositorie.On("CreateTask", mock.Anything, tasks).Return(ID, nil)

	fetchID, err := suite.usecase.CreateTask(context.TODO(), tasks)
	suite.Nil(err, "error should be nil")
	suite.Equal(ID, fetchID, "users should be equal")
}

func (suite *taskUsecaseSuite) TestCreateTask_Negative() {
	ID := "task_001"
	tasks := domain.Task{
		UserID:      "user_123",
		Title:       "",
		Description: "Finish writing the Go code and add tests.",
		Status:      "In Progress",
		Priority:    "High",
		DueDate:     time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	suite.repositorie.On("CreateTask", mock.Anything, tasks).Return(ID, nil)

	_, err := suite.usecase.CreateTask(context.TODO(), tasks)
	suite.NotNil(err, "error should not be nil as the task doesn't have a title")
}

func (suite *taskUsecaseSuite) TestUpdateTask() {
	tasks := domain.Task{
		ID:          "task_001",
		UserID:      "user_123",
		Title:       "",
		Description: "Finish writing the Go code and add tests.",
		Status:      "In Progress",
		Priority:    "High",
		DueDate:     time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	suite.repositorie.On("UpdateTask", mock.Anything, tasks).Return(tasks, nil)

	fetchedTask, err := suite.usecase.UpdateTask(context.TODO(), tasks)
	suite.Nil(err, "error should be nil")
	suite.Equal(tasks, fetchedTask, "tasks should be equal")
}

func (suite *taskUsecaseSuite) TestDeleteTask_Positive() {
	authorityUser := domain.User{
		ID:       "user_123",
		Username: "johndoe",
		Password: "password123",
		Role:     "admin",
	}

	tasks := domain.Task{
		ID:          "task_001",
		UserID:      "user_123",
		Title:       "",
		Description: "Finish writing the Go code and add tests.",
		Status:      "In Progress",
		Priority:    "High",
		DueDate:     time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	suite.repositorie.On("DeleteTask", mock.Anything, tasks.ID).Return(tasks, nil)
	suite.repositorie.On("FetchTaskByID", mock.Anything, tasks.ID).Return(tasks, nil)

	fetchedTask, err := suite.usecase.DeleteTask(context.TODO(), tasks.ID, authorityUser)
	suite.Nil(err, "error should be nil")
	suite.Equal(tasks, fetchedTask, "tasks should be equal")
}

func (suite *taskUsecaseSuite) TestDeleteTask_Negative() {
	authorityUser := domain.User{
		ID:       "user_124",
		Username: "johndoe",
		Password: "password123",
		Role:     "admin",
	}

	tasks := domain.Task{
		ID:          "task_001",
		UserID:      "user_123",
		Title:       "",
		Description: "Finish writing the Go code and add tests.",
		Status:      "In Progress",
		Priority:    "High",
		DueDate:     time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	suite.repositorie.On("DeleteTask", mock.Anything, tasks.ID).Return(tasks, nil)
	suite.repositorie.On("FetchTaskByID", mock.Anything, tasks.ID).Return(tasks, nil)

	_, err := suite.usecase.DeleteTask(context.TODO(), tasks.ID, authorityUser)
	suite.NotNil(err, "error should not be nil as the user is not authorized to delete the task")
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(taskUsecaseSuite))
}
