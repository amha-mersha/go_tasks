package usecases

import (
	"context"
	"time"

	domain "github.com/amha-mersha/go_tasks/go-backend-clean-architecture/delivery/domains"
)

type taskUseCase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewUseCase() {

}

// fetching task

func (taskUC *taskUseCase) GetAllTasks(cxt context.Context) ([]domain.Task, error) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()
	return taskUC.taskRepository.FetchAllTask(context)
}

func (taskUC taskUseCase) GetTaskByID(cxt context.Context, taskID string) (domain.Task, error) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.FetchTaskByID(context, taskID)
}

func (taskUC *taskUseCase) CreateTask(cxt context.Context, newTask domain.Task) (domain.TaskSuccess, error) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.CreateTask(context, newTask)
}

func (taskUC *taskUseCase) UpdateTask(cxt context.Context, updateTask domain.Task) (domain.TaskSuccess, error) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.CreateTask(context, updateTask)
}

func (taskUC *taskUseCase) DeleteTask(cxt context.Context, taskID string) (domain.TaskSuccess, error) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.DeleteTask(context, taskID)
}
