package tests

import (
	"context"
	"fmt"
	"log"
	"os"

	repositorie "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/repositories"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepositorySuite struct {
	suite.Suite
	repository *repositorie.UserRepository
	client     *mongo.Client
}

func (suite *testRepositorySuite) SetupSuite() {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Connection Error: ", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Ping Error: ", err)
	}
	suite.client = client

	database := client.Database(os.Getenv("DB_NAME"))
	log.Println("Database successfuly connected.")
	repo := repositorie.NewTaskRepository(database.Collection(os.Getenv("DB_TASK_COLLECTION_NAME")))
	suite.repository = &repo

}

func (suite *testRepositorySuite) TearDownSuite() {
	if suite.client != nil {
		if err := suite.client.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}
}

func (suite *testRepositorySuite) SetupTest() {
	if _, err := suite.repository.Collection.DeleteMany(context.TODO(), bson.D{}); err != nil {
		log.Println("Error deleting all documents from collection")
	}
}
