package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/Mocks"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/delivery/controllers"
	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type controllerTestSuite struct {
	suite.Suite
	taskUsecase *mocks.TaskUsecase
	userUsecase *mocks.UserUsecase
	controller  controllers.Controller
	router      *gin.Engine
	authUsecase *mocks.AuthService
}

type TaskResponse struct {
	Tasks []domain.Task `json:"tasks"`
}

func (suite *controllerTestSuite) SetupTest() {
	taskUC := new(mocks.TaskUsecase)
	userUC := new(mocks.UserUsecase)
	authUC := new(mocks.AuthService)
	suite.authUsecase = authUC
	suite.userUsecase = userUC
	suite.taskUsecase = taskUC
	suite.controller = controllers.NewController(taskUC, userUC)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default() // Make sure router is assigned to suite.router

	suite.router.POST("/task", suite.controller.PostTask)
	suite.router.PUT("/task", suite.controller.UpdateTask)
	suite.router.DELETE("/task/:id/:userid", suite.controller.DeleteTask)
	suite.router.POST("/user/assign", suite.controller.PostUserAssign)
	suite.router.POST("/user/register", suite.controller.PostUserRegister)
	suite.router.POST("/user/login", suite.controller.PostUserLogin)
	suite.router.GET("/task", suite.controller.GetTasks)
	suite.router.GET("/task/:id", suite.controller.GetTaskByID)
}

func (suite *controllerTestSuite) TestGetAllTasks_Positive() {
	tasks := []domain.Task{
		{
			ID:          "1",
			UserID:      "user_123",
			Title:       "Task 1",
			Description: "Description for Task 1",
			Status:      "Pending",
			Priority:    "High",
			DueDate:     time.Now().Truncate(time.Minute),
			CreatedAt:   time.Now().Truncate(time.Minute),
			UpdatedAt:   time.Now().Truncate(time.Minute),
		},
		{
			ID:          "2",
			UserID:      "user_456",
			Title:       "Task 2",
			Description: "Description for Task 2",
			Status:      "Completed",
			Priority:    "Low",
			DueDate:     time.Now().Truncate(time.Minute),
			CreatedAt:   time.Now().Truncate(time.Minute),
			UpdatedAt:   time.Now().Truncate(time.Minute),
		},
	}

	suite.taskUsecase.On("GetAllTasks", mock.Anything).Return(tasks, nil)
	req, _ := http.NewRequest(http.MethodGet, "/task", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	var tasksResponse TaskResponse
	fmt.Println(resp.Body.String())
	err := json.Unmarshal(resp.Body.Bytes(), &tasksResponse)
	suite.Nil(err)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(2, len(tasksResponse.Tasks), "Expected 2 tasks, but got %d", len(tasksResponse.Tasks))

}

func (suite *controllerTestSuite) TestGetTaskByID_Positive() {
	task := domain.Task{
		ID:          "1",
		UserID:      "user_123",
		Title:       "Task 1",
		Description: "Description for Task 1",
		Status:      "Pending",
		Priority:    "High",
		DueDate:     time.Now().Truncate(time.Minute),
		CreatedAt:   time.Now().Truncate(time.Minute),
		UpdatedAt:   time.Now().Truncate(time.Minute),
	}

	suite.taskUsecase.On("GetTaskByID", mock.Anything, "1").Return(task, nil)

	req, _ := http.NewRequest(http.MethodGet, "/task/1", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	var taskResponse domain.Task
	err := json.Unmarshal(resp.Body.Bytes(), &taskResponse)
	suite.Nil(err)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(task.ID, taskResponse.ID)
	suite.Equal(task.Title, taskResponse.Title)
	suite.Equal(task.Description, taskResponse.Description)
}

func (suite *controllerTestSuite) TestGetTaskByID_Negative() {
	expectedError := domain.TaskError{Message: "Task not found", Code: http.StatusInternalServerError}

	suite.taskUsecase.On("GetTaskByID", mock.Anything, "1").Return(domain.Task{}, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/task/1", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusInternalServerError, resp.Code, resp.Code)
}

func (suite *controllerTestSuite) TestUpdateTask_Positive() {
	updatedTask := domain.Task{
		ID:          "1",
		UserID:      "user_123",
		Title:       "Updated Task",
		Description: "Updated description",
		Status:      "In Progress",
		Priority:    "Medium",
		DueDate:     time.Now().Truncate(time.Minute),
		CreatedAt:   time.Now().Truncate(time.Minute),
		UpdatedAt:   time.Now().Truncate(time.Minute),
	}

	suite.taskUsecase.On("UpdateTask", mock.Anything, updatedTask).Return(updatedTask, nil)

	taskJSON, _ := json.Marshal(updatedTask)
	req, _ := http.NewRequest(http.MethodPut, "/task", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	var returnedTask domain.Task
	err := json.Unmarshal(resp.Body.Bytes(), &returnedTask)
	suite.Nil(err)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(updatedTask.ID, returnedTask.ID)
	suite.Equal(updatedTask.Title, returnedTask.Title)
	suite.Equal(updatedTask.Description, returnedTask.Description)
}

func (suite *controllerTestSuite) TestUpdateTask_Negative_Invalid_Payload() {
	malformedJSON := `{"id": "1", "title": "Updated Task", "description": "Updated description",`

	req, _ := http.NewRequest(http.MethodPut, "/task", bytes.NewBufferString(malformedJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	var errorResponse map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	suite.Nil(err)

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal("Invalid request payload", errorResponse["Error"])
}

func (suite *controllerTestSuite) TestUpdateTask_Negative_ValidationError() {
	invalidTask := domain.Task{
		ID: "1",
	}

	taskJSON, _ := json.Marshal(invalidTask)
	req, _ := http.NewRequest(http.MethodPut, "/task", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	suite.Nil(err)

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal("Missing required fields", errorResponse["Error"], "Expected error message: Missing required fields, but got %s", errorResponse["Error"])
	suite.Contains(errorResponse["Missing"], "Title")
}

func (suite *controllerTestSuite) TestDeleteTask_Positive() {
	taskID := "1"
	authorityID := "user_123"

	deletedTask := domain.Task{
		ID:          taskID,
		UserID:      authorityID,
		Title:       "Task to delete",
		Description: "Description of the task to delete",
		Status:      "Pending",
		Priority:    "High",
		DueDate:     time.Now().Truncate(time.Minute),
		CreatedAt:   time.Now().Truncate(time.Minute),
		UpdatedAt:   time.Now().Truncate(time.Minute),
	}

	suite.taskUsecase.On("DeleteTask", mock.Anything, taskID, authorityID).Return(deletedTask, nil)

	req, _ := http.NewRequest(http.MethodDelete, "/task/"+taskID+"/"+authorityID, nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	var returnedTask domain.Task
	err := json.Unmarshal(resp.Body.Bytes(), &returnedTask)
	suite.Nil(err)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(deletedTask.ID, returnedTask.ID)
	suite.Equal(deletedTask.Title, returnedTask.Title)
	suite.Equal(deletedTask.Description, returnedTask.Description)
}

func (suite *controllerTestSuite) TestDeleteTask_Negative_Error() {
	taskID := "1"
	authorityID := "user_123"

	suite.taskUsecase.On("DeleteTask", mock.Anything, taskID, authorityID).Return(domain.Task{}, &domain.TaskError{Code: http.StatusNotFound, Message: "Task not found"})

	req, _ := http.NewRequest(http.MethodDelete, "/task/"+taskID+"/"+authorityID, nil)

	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	var errorResponse map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	suite.Nil(err)

	suite.Equal(http.StatusNotFound, resp.Code)
	suite.Equal("Task not found", errorResponse["Error"])
}

func (suite *controllerTestSuite) TestPostTask_Success() {
	newTask := domain.Task{
		ID:          "3",
		UserID:      "user_789",
		Title:       "New Task",
		Description: "Description for the new task",
		Status:      "Pending",
		Priority:    "Medium",
		DueDate:     time.Now().Truncate(time.Minute),
		CreatedAt:   time.Now().Truncate(time.Minute),
		UpdatedAt:   time.Now().Truncate(time.Minute),
	}

	suite.taskUsecase.On("CreateTask", mock.Anything, newTask).Return(newTask.ID, nil)

	taskJSON, _ := json.Marshal(newTask)
	req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	var feedback string
	err := json.Unmarshal(resp.Body.Bytes(), &feedback)
	suite.Nil(err, "Error unmarshalling response")
	suite.Equal(http.StatusAccepted, resp.Code)
	suite.Equal(feedback, newTask.ID)
}

func (suite *controllerTestSuite) TestPostTask_MalformedJSON() {
	malformedJSON := `{"ID"  "1"  "UserID": "user_123", "Title": "Task", "Description": "Description", "Status": "Pending", "Priority": "High", "DueDate": "2024-08-16T00:00:00Z", "CreatedAt": "2024-08-16T00:00:00Z", "UpdatedAt": "2024-08-16T00:00:00Z"`

	req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBuffer([]byte(malformedJSON)))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	suite.Nil(err, "Error unmarshalling response body")

	suite.Equal("Malformed JSON", errorResponse["Error"], "Expected error message 'Malformed JSON', but got '%s'", errorResponse["Error"])
}

func (suite *controllerTestSuite) TestPostUser_Success() {
	newUser := domain.User{
		ID:       "3",
		Username: "user_789",
		Password: "password",
	}

	suite.userUsecase.On("UpdateUser", mock.Anything, newUser).Return(newUser, nil)

	userJSON, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, "/user/assign", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	var fetchedUser domain.User
	err := json.Unmarshal(resp.Body.Bytes(), &fetchedUser)
	suite.Nil(err, "Error unmarshalling response")
	suite.Equal(http.StatusAccepted, resp.Code)
	suite.Equal(fetchedUser.ID, newUser.ID)
}

func (suite *controllerTestSuite) TestPostUser_Malform() {
	malformedJSON := `{"ID"  "1"  "UserID": "user_123", "Title": "Task", "Description": "Description", "Status": "Pending", "Priority": "High", "DueDate": "2024-08-16T00:00:00Z", "CreatedAt": "2024-08-16T00:00:00Z", "UpdatedAt": "2024-08-16T00:00:00Z"`

	req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBuffer([]byte(malformedJSON)))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	suite.Nil(err, "Error unmarshalling response body")

	suite.Equal("Malformed JSON", errorResponse["Error"], "Expected error message 'Malformed JSON', but got '%s'", errorResponse["Error"])
}

func (suite *controllerTestSuite) TestPostUserRegister_Success() {
	newUser := domain.User{
		ID:       "3",
		Username: "new_user",
		Password: "$2a$10$hashedpassword",
	}

	suite.userUsecase.On("CreateUser", mock.Anything, newUser).Return(newUser.ID, nil)

	userJSON, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	var registeredUserID string
	err := json.Unmarshal(resp.Body.Bytes(), &registeredUserID)
	suite.Nil(err, "Error unmarshalling response")
	suite.Equal(http.StatusAccepted, resp.Code)
	suite.Equal(newUser.ID, registeredUserID)
}

func (suite *controllerTestSuite) TestPostUserRegister_Fail() {
	malformedJSON := `{"ID"  "1"  "UserID": "user_123", "Title": "Task", "Description": "Description", "Status": "Pending", "Priority": "High", "DueDate": "2024-08-16T00:00:00Z", "CreatedAt": "2024-08-16T00:00:00Z", "UpdatedAt": "2024-08-16T00:00:00Z"`

	req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBuffer([]byte(malformedJSON)))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	suite.Nil(err, "Error unmarshalling response body")

	suite.Equal("Malformed JSON", errorResponse["Error"], "Expected error message 'Malformed JSON', but got '%s'", errorResponse["Error"])
}

func (suite *controllerTestSuite) TestPostUserLogin_Success() {
	user := domain.User{
		ID:       "3",
		Username: "new_user",
		Password: "$2a$10$hashedpassword",
	}
	token := "token"
	suite.userUsecase.On("LoginUser", mock.Anything, user).Return(token, nil)

	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)
	fmt.Println(resp.Body.String())

	var regisetToken map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &regisetToken)
	suite.Nil(err, "Error unmarshalling response")
	suite.Equal(http.StatusOK, resp.Code, "Expected status code 200, but got %d", resp.Code)
	suite.Equal(token, regisetToken["token"], "Expected token to be returned, but got %s", token)

}

func (suite *controllerTestSuite) TestPostUserLogin_Fail() {
	malformedJSON := `{"ID"  "1"  "UserID": "user_123", "Title": "Task", "Description": "Description", "Status": "Pending", "Priority": "High", "DueDate": "2024-08-16T00:00:00Z", "CreatedAt": "2024-08-16T00:00:00Z", "UpdatedAt": "2024-08-16T00:00:00Z"`

	req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBuffer([]byte(malformedJSON)))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	suite.Nil(err, "Error unmarshalling response body")

	suite.Equal("Malformed JSON", errorResponse["Error"], "Expected error message 'Malformed JSON', but got '%s'", errorResponse["Error"])
}
func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(controllerTestSuite))
}
