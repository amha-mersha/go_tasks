package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// task struc
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"userID" bson:"userID"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Priority    string             `json:"priority"`
	DueDate     time.Time          `json:"due_date"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// user structs

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Role     string             `json:"role"`
}

// success struct

type TaskSuccess struct {
	Message string
	Return  interface{}
}

// error structs

type UserError struct {
	Message string
	Code    int
}

func (usererr *UserError) Error() string {
	return usererr.Message
}

type TaskError struct {
	Message string
	Code    int
}

func (taskerr *TaskError) Error() string {
	return taskerr.Message
}

// task repository struct
type TaskRepository interface {
	FetchAllTask(context.Context) ([]Task, error)
	FetchTaskByID(context.Context, string) (Task, error)
	CreateTask(context.Context, Task) (TaskSuccess, error)
	DeleteTask(context.Context, string) (TaskSuccess, error)
	UpdateTask(context.Context, string) (TaskSuccess, error)
}

// task use case interface
type TaskUsecase interface {
	GetAllTasks(cxt context.Context) ([]Task, error)
	GetTaskByID(cxt context.Context, taskID string) (Task, error)
	CreateTask(cxt context.Context, newTask Task) (TaskSuccess, error)
	UpdateTask(cxt context.Context, updateTask Task) (TaskSuccess, error)
	DeleteTask(cxt context.Context, taskID string) (TaskSuccess, error)
}
