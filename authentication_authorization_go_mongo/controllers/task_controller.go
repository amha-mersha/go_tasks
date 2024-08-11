package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/data"
	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetTasks(ctx *gin.Context) {
	tasks, err := data.GetTasks()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

func GetTaskByID(ctx *gin.Context) {
	taskID := ctx.Param("id")
	task, err := data.GetTaskByID(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func UpdateTask(ctx *gin.Context) {
	taskID := ctx.Param("id")
	var updatedTask models.Task
	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MalformedJSON})
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MismatchedFormat})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Error())
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MissingRequireds, "Missing": missingRequireds})
		}
		return
	}
	returnedTask, err := data.UpdateTask(taskID, updatedTask)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, returnedTask)
}

func DeleteTask(ctx *gin.Context) {
	taskID := ctx.Param("id")
	deletedTask, err := data.DeleteTask(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, deletedTask)
}

func PostTask(ctx *gin.Context) {
	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MalformedJSON})
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MismatchedFormat})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Error())
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MissingRequireds, "Missing": missingRequireds})
		}
		return
	}
	result, err := data.PostTask(newTask)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, result)
}

func PostUserRegister(ctx *gin.Context) {
	var newUser models.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MalformedJSON})
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MismatchedFormat})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Error())
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MissingRequireds, "Missing": missingRequireds})
		}
		return
	}
	err := data.PostUserRegister(newUser)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"Server": "User registered successfully"})
}

func PostUserLogin(ctx *gin.Context) {
	var logingUser models.User
	if err := ctx.ShouldBindJSON(&logingUser); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MalformedJSON})
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MismatchedFormat})
		case validator.ValidationErrors:
			missingRequireds := []string{}
			for _, fieldError := range e {
				missingRequireds = append(missingRequireds, fieldError.Error())
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": data.MissingRequireds, "Missing": missingRequireds})
		}
		return
	}
	jwtToken, err := data.PostUserLogin(logingUser)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"Server": "User registered successfully", "token": jwtToken})
}
