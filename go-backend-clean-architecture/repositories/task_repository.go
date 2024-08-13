package repositorie

import (
	"context"
	"net/http"

	domain "github.com/amha-mersha/go_tasks/go-backend-clean-architecture/domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// locally used types
const (
	CREATED_SUCCESSFULLY = "task created successfully"
)

type taskRepository struct {
	database   *mongo.Database
	collection string
}

// creating a task repository instance

func NewTaskRepository(database *mongo.Database, collection string) taskRepository {
	return taskRepository{database: database, collection: collection}
}

func (taskRepo *taskRepository) FetchAllTasks(cxt context.Context) ([]domain.Task, *domain.TaskError) {
	collection := taskRepo.database.Collection(taskRepo.collection)
	filter := bson.D{}
	cursor, err := collection.Find(cxt, filter)
	if err != nil {
		return []domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	var fetchedTasks []domain.Task
	err = cursor.All(cxt, &fetchedTasks)
	if err != nil {
		return []domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return fetchedTasks, nil
}

func (taskRepo *taskRepository) FetchTaskByID(cxt context.Context, ID string) (domain.Task, *domain.TaskError) {
	collection := taskRepo.database.Collection(taskRepo.collection)

	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	filter := bson.D{{"_id", taskID}}
	var fetchedTask domain.Task
	err = collection.FindOne(cxt, filter).Decode(&fetchedTask)
	if err != nil {
		return domain.Task{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return fetchedTask, nil
}

func (taskRepo *taskRepository) CreateTask(cxt context.Context, newTask domain.Task) (domain.TaskSuccess, *domain.TaskError) {
	collection := taskRepo.database.Collection(taskRepo.collection)
	newTask.ID = ""
	insertedTask, err := collection.InsertOne(cxt, newTask)
	if err != nil {
		return domain.TaskSuccess{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return domain.TaskSuccess{Message: CREATED_SUCCESSFULLY, Return: insertedTask}, nil
}

func (taskRepo *taskRepository) UpdateTask(cxt context.Context, updateTask domain.Task) (domain.TaskSuccess, *domain.TaskError) {
	collection := taskRepo.database.Collection(taskRepo.collection)
	filter := bson.D{{"_id", updateTask.ID}}
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
	var returnedtask domain.Task
	err := collection.FindOneAndUpdate(cxt, filter, updateTask, opts).Decode(&returnedtask)
	if err != nil {
		return domain.TaskSuccess{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return domain.TaskSuccess{Message: CREATED_SUCCESSFULLY, Return: returnedtask}, nil

}

func (taskRepo *taskRepository) DeleteTask(cxt context.Context, ID string) (domain.TaskSuccess, *domain.TaskError) {
	collection := taskRepo.database.Collection(taskRepo.collection)
	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.TaskSuccess{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	filter := bson.D{{"_id", taskID}}
	var deletedTask domain.Task
	err = collection.FindOneAndDelete(cxt, filter).Decode(&deletedTask)
	if err != nil {
		return domain.TaskSuccess{}, &domain.TaskError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return domain.TaskSuccess{Message: "task deleted successfully", Return: deletedTask}, nil
}
