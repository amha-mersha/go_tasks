package repositorie

import (
	"context"
	"net/http"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// locally used types
const (
	CREATED_SUCCESSFULLY = "task created successfully"
)

type TaskRepository struct {
	Collection *mongo.Collection
}

// creating a task repository instance

func NewTaskRepository(collection *mongo.Collection) TaskRepository {
	return TaskRepository{Collection: collection}
}

func (taskRepo *TaskRepository) FetchAllTasks(cxt context.Context) ([]domain.Task, *domain.TaskError) {
	filter := bson.D{}
	cursor, err := taskRepo.Collection.Find(cxt, filter)
	if err != nil {
		return []domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	var fetchedTasks []domain.Task
	err = cursor.All(cxt, &fetchedTasks)
	defer cursor.Close(cxt)
	if err != nil {
		return []domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return fetchedTasks, nil
}

func (taskRepo *TaskRepository) FetchTaskByID(cxt context.Context, ID string) (domain.Task, *domain.TaskError) {
	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	filter := bson.D{{"_id", taskID}}
	var fetchedTask domain.Task
	err = taskRepo.Collection.FindOne(cxt, filter).Decode(&fetchedTask)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return fetchedTask, nil
}

func (taskRepo *TaskRepository) CreateTask(cxt context.Context, newTask domain.Task) (string, *domain.TaskError) {
	newTask.ID = ""
	insertedTask, err := taskRepo.Collection.InsertOne(cxt, newTask)
	if err != nil {
		return "", &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	result, ok := insertedTask.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return result.Hex(), nil
}

func (taskRepo *TaskRepository) UpdateTask(cxt context.Context, updateTask domain.Task) (domain.Task, *domain.TaskError) {
	objectID, err := primitive.ObjectIDFromHex(updateTask.ID)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: "Invalid ID format", Code: http.StatusBadRequest}
	}
	filter := bson.D{{"_id", objectID}}
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
	inserteTask := domain.Task{
		Title:       updateTask.Title,
		Description: updateTask.Description,
		Status:      updateTask.Status,
		DueDate:     updateTask.DueDate,
		UserID:      updateTask.UserID,
		Priority:    updateTask.Priority,
	}
	update := bson.D{{"$set", inserteTask}}
	var returnedtask domain.Task
	err = taskRepo.Collection.FindOneAndUpdate(cxt, filter, update, opts).Decode(&returnedtask)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return returnedtask, nil
}

func (taskRepo *TaskRepository) DeleteTask(cxt context.Context, ID string) (domain.Task, *domain.TaskError) {
	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	filter := bson.D{{"_id", taskID}}
	var deletedTask domain.Task
	err = taskRepo.Collection.FindOneAndDelete(cxt, filter).Decode(&deletedTask)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return deletedTask, nil
}
