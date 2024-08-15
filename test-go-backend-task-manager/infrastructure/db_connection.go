package infrastructure

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDatabase() (*mongo.Database, error) {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return &mongo.Database{}, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return &mongo.Database{}, err
	}

	database := client.Database(os.Getenv("DB_NAME"))
	log.Println("Database successfuly connected.")
	return database, nil
}

func EstablisUniqueUsernameIndex(collection *mongo.Collection, index string) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{index: 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	return nil
}
