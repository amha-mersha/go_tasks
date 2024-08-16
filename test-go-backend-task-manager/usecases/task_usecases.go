package usecases

import (
	"context"
	"time"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
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

func (taskUC *taskUseCase) CreateTask(cxt context.Context, newTask domain.Task) (string, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()
	if newTask.Title == "" {
		return "", &domain.TaskError{Message: "Title is required", Code: 400}
	}

	return taskUC.taskRepository.CreateTask(context, newTask)
}

func (taskUC *taskUseCase) UpdateTask(cxt context.Context, updateTask domain.Task) (domain.Task, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	return taskUC.taskRepository.UpdateTask(context, updateTask)
}

func (taskUC *taskUseCase) DeleteTask(cxt context.Context, taskID string, authority domain.User) (domain.Task, *domain.TaskError) {
	context, cancel := context.WithTimeout(cxt, taskUC.contextTimeout)
	defer cancel()

	fetchedTask, errFetch := taskUC.taskRepository.FetchTaskByID(context, taskID)
	if errFetch != nil {
		return domain.Task{}, errFetch
	}

	if authority.ID != fetchedTask.UserID {
		return domain.Task{}, &domain.TaskError{Message: "You are not authorized to delete this task", Code: 403}
	}

	return taskUC.taskRepository.DeleteTask(context, taskID)
}
