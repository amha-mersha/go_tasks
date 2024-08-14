package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	repositorie "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/repositories"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testRepositorySuite struct {
	suite.Suite
	repository *repositorie.TaskRepository
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

func (suite *testRepositorySuite) TestFetchAllTasks() {
	testSamples := []domain.Task{
		{
			UserID:      "user789",
			Title:       "Schedule dentist appointment",
			Description: "Call the dentist to schedule an appointment.",
			Status:      "Completed",
			Priority:    "Low",
			DueDate:     time.Now().Add(time.Hour * 94),
			CreatedAt:   time.Now().Add(time.Hour * 94),
			UpdatedAt:   time.Now().Add(time.Hour * 94),
		},
		{
			UserID:      "user012",
			Title:       "Prepare for presentation",
			Description: "Work on the slides and practice the speech.",
			Status:      "In Progress",
			Priority:    "High",
			DueDate:     time.Now().Add(time.Hour * 92),
			CreatedAt:   time.Now().Add(time.Hour * 92),
			UpdatedAt:   time.Now().Add(time.Hour * 92),
		},
		{
			UserID:      "user345",
			Title:       "Read a book",
			Description: "Finish reading 'The Great Gatsby'.",
			Status:      "Not Started",
			Priority:    "Low",
			DueDate:     time.Now().Add(time.Hour * 90),
			CreatedAt:   time.Now().Add(time.Hour * 90),
			UpdatedAt:   time.Now().Add(time.Hour * 90),
		},
	}
	suite.repository.Collection.InsertMany(context.TODO(), testSamples)
	_, errFetch := suite.repository.FetchAllTasks(context.TODO())
	suite.Error(errFetch, "Error fetching all tasks")
	suite.NoError(errFetch, "No error fetching all tasks")
}

func (suite *testRepositorySuite) TestFetchEmptyDatabase() {
	_, errFetch := suite.repository.FetchAllTasks(context.TODO())
	suite.Error(errFetch, "Error fetching empty database")
	suite.NoError(errFetch, "No error fetching empty database")
}

func (suite *testRepositorySuite) TestAddingAndRetriving() {
	task := domain.Task{
		UserID:      "user789",
		Title:       "Schedule dentist appointment",
		Description: "Call the dentist to schedule an appointment.",
		Status:      "Completed",
		Priority:    "Low",
		DueDate:     time.Now().Add(time.Hour * 94),
		CreatedAt:   time.Now().Add(time.Hour * 94),
		UpdatedAt:   time.Now().Add(time.Hour * 94),
	}
	insertedTask, errInsert := suite.repository.CreateTask(context.TODO(), task)
	suite.Error(errInsert, "Unwanted: Error inserting task")
	suite.NoError(errInsert, insertedTask.Return)
	fetchedTask, errFetch := suite.repository.FetchTaskByID(context.TODO(), insertedTask.Return.(primitive.ObjectID).Hex())
	suite.Error(errFetch, "Unwanted: Error fetching task")
	suite.NoError(errFetch, fetchedTask)
	suite.Equal(insertedTask.Return.(primitive.ObjectID).Hex(), fetchedTask.ID, "Inserted and fetched task ID should be the same")
}

func (suite *testRepositorySuite) TestFetchAllTasksEmptyDB() {
	_, err := suite.repository.FetchAllTasks(context.TODO())
	suite.Error(err, "Error fetching all tasks from empty database")
	suite.Equal(err.Code, 500, "Error code should be 500")
}

func (suite *testRepositorySuite) TestFetchAllTasksFilledDB() {
	task := []domain.Task{
		{
			Title:       "Schedule dentist appointment",
			Description: "Call the dentist to schedule an appointment.",
			Status:      "Completed",
			Priority:    "Low",
			DueDate:     time.Now().Add(time.Hour * 94),
			CreatedAt:   time.Now().Add(time.Hour * 94),
			UpdatedAt:   time.Now().Add(time.Hour * 94),
		},
		{
			UserID:      "user234",
			Title:       "Fix the leaking faucet",
			Description: "Repair the kitchen faucet to stop the leak.",
			Status:      "Not Started",
			Priority:    "Medium",
			DueDate:     time.Now().Add(time.Hour * 83),
			CreatedAt:   time.Now().Add(time.Hour * 83),
			UpdatedAt:   time.Now().Add(time.Hour * 83),
		},
	}
	insertedTaskOne, errInsertOne := suite.repository.CreateTask(context.TODO(), task[0])
	suite.Error(errInsertOne, "Unwanted: Error inserting task One")
	_, errFetchOne := suite.repository.FetchTaskByID(context.TODO(), insertedTaskOne.Return.(primitive.ObjectID).Hex())
	suite.Error(errFetchOne, "Unwanted: Error fetching task One")
	suite.Equal(errFetchOne.Code, 500, "Error code should be 500")

	insertedTaskTwo, errInsertTwo := suite.repository.CreateTask(context.TODO(), task[1])
	suite.Error(errInsertTwo, "Unwanted: Error inserting task Two")
	_, errFetchTwo := suite.repository.FetchTaskByID(context.TODO(), insertedTaskTwo.Return.(primitive.ObjectID).Hex())
	suite.Error(errFetchTwo, "Unwanted: Error fetching task Two")
	suite.Equal(errFetchTwo.Code, 500, "Error code should be 500")

	totalFetched, err := suite.repository.FetchAllTasks(context.TODO())
	suite.Error(err, "Error fetching all tasks from populated database")
	suite.Equal(err.Code, 500, "Error code should be 500")
	suite.Equal(len(totalFetched), 2, "Total fetched tasks should be 2")
}

func (suite *testRepositorySuite) TestUpdatingTask() {
	initalTask := domain.Task{
		Title:       "Schedule dentist appointment",
		Description: "Call the dentist to schedule an appointment.",
		Status:      "Completed",
		Priority:    "Low",
		DueDate:     time.Now().Add(time.Hour * 94),
		CreatedAt:   time.Now().Add(time.Hour * 94),
		UpdatedAt:   time.Now().Add(time.Hour * 94),
	}
	insertedResult, err := suite.repository.CreateTask(context.TODO(), initalTask)
	suite.Error(err, "Error creating initial task")
	updateTask := domain.Task{
		ID:          insertedResult.Return.(primitive.ObjectID).Hex(),
		Title:       "Schedule doctors appointment",
		Description: "Call the doctors to schedule an appointment.",
		Status:      "Pending",
		Priority:    "High",
		DueDate:     time.Now().Add(time.Hour * 94),
		CreatedAt:   time.Now().Add(time.Hour * 94),
		UpdatedAt:   time.Now().Add(time.Hour * 94),
	}

	fetchedTask, errUpdate := suite.repository.UpdateTask(context.TODO(), updateTask)
	suite.Error(errUpdate, "Error updating task")
	suite.Equal(errUpdate.Code, 500, "Error code should be 500")
	suite.Equal(fetchedTask.Message, "task created successfully", "Task updated successfully")
	suite.Equal(fetchedTask.Return.(domain.Task).ID, insertedResult.Return.(primitive.ObjectID).Hex(), "Task ID should be the same")
}

func (suite *testRepositorySuite) TestDeletingTask() {
	initalTask := domain.Task{
		Title:       "Schedule dentist appointment",
		Description: "Call the dentist to schedule an appointment.",
		Status:      "Completed",
		Priority:    "Low",
		DueDate:     time.Now().Add(time.Hour * 94),
		CreatedAt:   time.Now().Add(time.Hour * 94),
		UpdatedAt:   time.Now().Add(time.Hour * 94),
	}
	insertedResult, err := suite.repository.CreateTask(context.TODO(), initalTask)
	suite.Error(err, "Error creating initial task")

	fetchedTask, errUpdate := suite.repository.DeleteTask(context.TODO(), insertedResult.Return.(primitive.ObjectID).Hex())
	suite.Error(errUpdate, "Error deleting task")
	suite.Equal(errUpdate.Code, 500, "Error code should be 500")
	suite.Equal(fetchedTask.Return.(domain.Task).ID, insertedResult.Return.(primitive.ObjectID).Hex(), "Task ID should be the same")
}
