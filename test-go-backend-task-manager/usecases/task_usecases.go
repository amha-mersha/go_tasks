package usecases

import (
	"context"
	"time"

	domain "github.com/amha-mersha/go_tasks/go-backend-clean-architecture/domains"
)

type taskUseCase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepo domain.TaskRepository, timeout time.Duration) taskUseCase {
	return taskUseCase{
		taskRepository: taskRepo,
		contextTimeout: timeout,
	}
}

// fetching task

func (taskUC *taskUseCase) GetAllTasks(cxt context.Context) ([]domain.Task, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()
	return taskUC.taskRepository.FetchAllTasks(context)
}

func (taskUC taskUseCase) GetTaskByID(cxt context.Context, taskID string) (domain.Task, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.FetchTaskByID(context, taskID)
}

func (taskUC *taskUseCase) CreateTask(cxt context.Context, newTask domain.Task) (domain.TaskSuccess, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.CreateTask(context, newTask)
}

func (taskUC *taskUseCase) UpdateTask(cxt context.Context, updateTask domain.Task) (domain.TaskSuccess, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.CreateTask(context, updateTask)
}

func (taskUC *taskUseCase) DeleteTask(cxt context.Context, taskID string) (domain.TaskSuccess, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.DeleteTask(context, taskID)
}
