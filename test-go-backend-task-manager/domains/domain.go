package domain

import (
	"context"
	"time"
)

// task struc
type Task struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string    `json:"userID" bson:"userID"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Status      string    `json:"status,omitempty" bson:"status,omitempty"`
	Priority    string    `json:"priority,omitempty" bson:"priority,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty" bson:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// user structs

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
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
	FetchAllTasks(cxt context.Context) ([]Task, *TaskError)
	FetchTaskByID(cxt context.Context, ID string) (Task, *TaskError)
	CreateTask(cxt context.Context, newTask Task) (string, *TaskError)
	UpdateTask(cxt context.Context, updateTask Task) (Task, *TaskError)
	DeleteTask(cxt context.Context, taskID string) (Task, *TaskError)
}

// task use case interface
type TaskUsecase interface {
	GetAllTasks(cxt context.Context) ([]Task, *TaskError)
	GetTaskByID(cxt context.Context, taskID string) (Task, *TaskError)
	CreateTask(cxt context.Context, newTask Task) (string, *TaskError)
	UpdateTask(cxt context.Context, updateTask Task) (Task, *TaskError)
	DeleteTask(cxt context.Context, taskID string, authority User) (Task, *TaskError)
}

// users use case interface
type UserUsecase interface {
	GetAllUser(cxt context.Context) ([]User, *UserError)
	GetUserByID(cxt context.Context, userID string) (User, *UserError)
	GetUserByUsername(cxt context.Context, username string) (User, *UserError)
	CreateUser(cxt context.Context, newUser User) (string, *UserError)
	UpdateUser(cxt context.Context, userUpdate User) (User, *UserError)
	DeleteUser(cxt context.Context, authority User, deleteID string) (User, *UserError)
	LoginUser(cxt context.Context, loggingUser User) (string, *UserError)
}

// task repository struct
type UserRepository interface {
	FetchAllUsers(cxt context.Context) ([]User, *UserError)
	FetchUserCount(cxt context.Context) (int, *UserError)
	FetchUserByID(cxt context.Context, ID string) (User, *UserError)
	FetchUserByUsername(cxt context.Context, username string) (User, *UserError)
	CreateUser(cxt context.Context, newUser User) (string, *UserError)
	UpdateUser(cxt context.Context, updateUser User) (User, *UserError)
	DeleteUser(cxt context.Context, userID string) (User, *UserError)
}
