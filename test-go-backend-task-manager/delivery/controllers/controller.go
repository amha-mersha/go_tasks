package controllers

import (
	"encoding/json"
	"net/http"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Controller struct {
	TaskUsecase domain.TaskUsecase
	UserUsecase domain.UserUsecase
}

func NewController(taskUC domain.TaskUsecase, userUC domain.UserUsecase) Controller {
	return Controller{
		TaskUsecase: taskUC,
		UserUsecase: userUC,
	}

}

func (controller *Controller) GetTasks(cxt *gin.Context) {
	tasks, err := controller.TaskUsecase.GetAllTasks(cxt)
	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Message})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (controller *Controller) GetTaskByID(cxt *gin.Context) {
	taskID := cxt.Param("id")
	task, err := controller.TaskUsecase.GetTaskByID(cxt, taskID)
	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Error()})
		return
	}
	cxt.JSON(http.StatusOK, task)
}

func (controller *Controller) UpdateTask(cxt *gin.Context) {
	var updatedTask domain.Task
	if err := cxt.ShouldBindJSON(&updatedTask); err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Malformed JSON"})
		case *json.UnmarshalTypeError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Mismatched format"})
		default:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request payload"})
		}
		return
	}

	validate := validator.New()
	if err := validate.Struct(updatedTask); err != nil {
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			missingFields := []string{}
			for _, fieldError := range fieldErrors {
				missingFields = append(missingFields, fieldError.Field())
			}
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Missing required fields", "Missing": missingFields})
		} else {
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Validation error"})
		}
		return
	}

	returnedTask, err := controller.TaskUsecase.UpdateTask(cxt, updatedTask)
	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Error()})
		return
	}

	cxt.JSON(http.StatusOK, returnedTask)
}

func (controller *Controller) DeleteTask(cxt *gin.Context) {
	taskID := cxt.Param("id")
	authorityID := cxt.Param("userid")
	deletedTask, err := controller.TaskUsecase.DeleteTask(cxt, taskID, authorityID)
	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Error()})
		return
	}
	cxt.JSON(http.StatusOK, deletedTask)
}

func (controller *Controller) PostTask(cxt *gin.Context) {
	var newTask domain.Task
	if err := cxt.ShouldBindJSON(&newTask); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Malformed JSON"})
		case *json.UnmarshalTypeError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Mismatched format"})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Field())
			}
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Missing required fields", "Missing": missingRequireds})
		}
		return
	}
	result, err := controller.TaskUsecase.CreateTask(cxt, newTask)

	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Error()})
		return
	}
	cxt.JSON(http.StatusAccepted, result)
}
func (controller *Controller) PostUserAssign(cxt *gin.Context) {
	var updateUser domain.User
	if err := cxt.ShouldBind(&updateUser); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Malformed JSON"})
		case *json.UnmarshalTypeError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Mismatched format"})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Field())
			}
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Missing required fields", "Missing": missingRequireds})
		}
		return
	}
	result, err := controller.UserUsecase.UpdateUser(cxt, updateUser)
	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Error()})
		return
	}
	cxt.JSON(http.StatusAccepted, result)
}

func (controller *Controller) PostUserRegister(cxt *gin.Context) {
	var registeringUser domain.User
	if err := cxt.ShouldBind(&registeringUser); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Malformed JSON"})
		case *json.UnmarshalTypeError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Mismatched format"})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Field())
			}
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Missing required fields", "Missing": missingRequireds})
		}
		return
	}
	result, err := controller.UserUsecase.CreateUser(cxt, registeringUser)
	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Error()})
		return
	}
	cxt.JSON(http.StatusAccepted, result)
}

func (controller *Controller) PostUserLogin(cxt *gin.Context) {
	var loggingUser domain.User
	if err := cxt.ShouldBind(&loggingUser); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Malformed JSON"})
		case *json.UnmarshalTypeError:
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Mismatched format"})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Field())
			}
			cxt.JSON(http.StatusBadRequest, gin.H{"Error": "Missing required fields", "Missing": missingRequireds})
		}
		return
	}
	token, err := controller.UserUsecase.LoginUser(cxt, loggingUser)
	if err != nil {
		cxt.JSON(err.Code, gin.H{"Error": err.Error()})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{"token": token})
}
