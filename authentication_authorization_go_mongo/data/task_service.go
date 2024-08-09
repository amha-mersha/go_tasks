package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskError struct {
	message string
}

func (err TaskError) Error() string {
	return err.message
}

const (
	IDNotFound        = "no item found with the specified id"
	TaskAlreadyExists = "the task already exists in the database"
	MalformedJSON     = "sent a malfomed json"
	MismatchedFormat  = "the task have a mismatched stucture"
	MissingRequireds  = "there are some missing required feilds"
	MalformedID       = "the id sent is malformed"
)

func GetTasks() ([]models.Task, error) {
	var result []models.Task
	curr, err := Collection.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		return []models.Task{}, err
	}
	err = curr.All(context.TODO(), &result)
	if err != nil {
		return []models.Task{}, err
	}
	return result, nil
}

func GetTaskByID(taskID string) (models.Task, error) {
	var result models.Task
	ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, fmt.Errorf(MalformedID)
	}
	err = Collection.FindOne(context.TODO(), bson.D{{"_id", ID}}).Decode(&result)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return models.Task{}, fmt.Errorf(IDNotFound)
	} else if err != nil {
		return models.Task{}, err
	}
	return result, nil
}

func UpdateTask(taskID string, updatedTask models.Task) (models.Task, error) {
	var result models.Task
	ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, fmt.Errorf(MalformedID)
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{{"_id", ID}}

	err = Collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": updatedTask}, opts).Decode(&result)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return models.Task{}, fmt.Errorf(IDNotFound)
	} else if err != nil {
		return models.Task{}, err
	}
	return result, nil
}

func DeleteTask(taskID string) (models.Task, error) {
	var result models.Task
	ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, fmt.Errorf(MalformedID)
	}
	filter := bson.D{{"_id", ID}}
	err = Collection.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return models.Task{}, fmt.Errorf(IDNotFound)
	} else if err != nil {
		return models.Task{}, err
	}
	return result, nil
}

func PostTask(newTask models.Task) (*mongo.InsertOneResult, error) {
	result, err := Collection.InsertOne(context.TODO(), newTask)
	if err != nil {
		return &mongo.InsertOneResult{}, fmt.Errorf(MalformedID)
	}
	return result, err
}
