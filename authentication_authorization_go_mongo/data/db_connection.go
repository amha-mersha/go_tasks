package data

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var TaskCollection *mongo.Collection
var UserCollection *mongo.Collection

func ConnecDB() error {
	clientOptions := options.Client().ApplyURI(os.Getenv("DB_CONNECTION_STRING"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	Client = client

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	TaskCollection = client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_TASK_COLLECTION_NAME"))
	UserCollection = client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_USER_COLLECTION_NAME"))
	log.Println(TaskCollection.Name(), UserCollection.Name())
	log.Println("Database successfuly connected.")
	return nil

}
