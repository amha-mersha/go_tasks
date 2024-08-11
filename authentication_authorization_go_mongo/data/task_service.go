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
	IDNotFound           = "no item found with the specified id"
	TaskAlreadyExists    = "the task already exists in the database"
	MalformedJSON        = "sent a malfomed json"
	MismatchedFormat     = "the task have a mismatched stucture"
	MissingRequireds     = "there are some missing required feilds"
	MalformedID          = "the id sent is malformed"
	MalformedUsername    = "the username sent is malfomed"
	UserAlreadyExist     = "user by that username exists"
	InternalServerError  = "internal server error occured"
	UserNotFound         = "no user found with the specified username"
	IncorrectCredentials = "invalid password or username"
)

func GetTasks() (interface{}, error) {
	var result []models.TaskMongo
	curr, err := TaskCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return []models.Task{}, err
	}
	err = curr.All(context.TODO(), &result)
	if err != nil {
		return []models.Task{}, err
	}
	return result, nil
}

func GetTaskByID(taskID string) (interface{}, error) {
	var result models.TaskMongo
	ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, fmt.Errorf(MalformedID)
	}
	err = TaskCollection.FindOne(context.TODO(), bson.D{{"_id", ID}}).Decode(&result)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return models.Task{}, fmt.Errorf(IDNotFound)
	} else if err != nil {
		return models.Task{}, err
	}
	return result, nil
}

func UpdateTask(taskID string, updatedTask models.Task) (interface{}, error) {
	var result models.TaskMongo
	ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, fmt.Errorf(MalformedID)
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{{"_id", ID}}

	err = TaskCollection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": updatedTask}, opts).Decode(&result)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return models.Task{}, fmt.Errorf(IDNotFound)
	} else if err != nil {
		return models.Task{}, err
	}
	return result, nil
}

func DeleteTask(taskID string) (models.TaskMongo, error) {
	var result models.TaskMongo
	ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.TaskMongo{}, fmt.Errorf(MalformedID)
	}
	filter := bson.D{{"_id", ID}}
	err = TaskCollection.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return models.TaskMongo{}, fmt.Errorf(IDNotFound)
	} else if err != nil {
		return models.TaskMongo{}, err
	}
	return result, nil
}

func PostTask(newTask models.Task) (*mongo.InsertOneResult, error) {
	result, err := TaskCollection.InsertOne(context.TODO(), newTask)
	if err != nil {
		return &mongo.InsertOneResult{}, fmt.Errorf(MalformedID)
	}
	return result, err
}
