package controllers

import (
	domain "github.com/amha-mersha/go_tasks/go-backend-clean-architecture/domains"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

func (taskController *TaskController) GetTasks(cxt *gin.Context) {
  tasks, err := taskController.TaskUsecase.GetAllTasks(cxt)
  if err != nil{
    cxt.JSON(http.)
  }

}
