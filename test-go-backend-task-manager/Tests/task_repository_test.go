package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	repositorie "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/repositories"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
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
		fmt.Println("Connection Nil: ", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Ping Nil: ", err)
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
			log.Printf("Nil disconnecting from database: %v", err)
		}
	}
}

func (suite *testRepositorySuite) SetupTest() {
	if _, err := suite.repository.Collection.DeleteMany(context.TODO(), bson.D{}); err != nil {
		log.Println("Nil deleting all documents from collection")
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
	for _, task := range testSamples {
		suite.repository.Collection.InsertOne(context.TODO(), task)
	}
	_, errFetch := suite.repository.FetchAllTasks(context.TODO())
	suite.Nil(errFetch, "No error fetching all tasks")
}

func (suite *testRepositorySuite) TestFetchEmptyDatabase() {
	_, errFetch := suite.repository.FetchAllTasks(context.TODO())
	suite.Nil(errFetch, "Error fetching empty database")
	suite.Nil(errFetch, "No error fetching empty database")
}

func (suite *testRepositorySuite) TestAddingAndRetrieving() {
	task := domain.Task{
		UserID:      "user789",
		Title:       "Schedule dentist appointment",
		Description: "Call the dentist to schedule an appointment.",
		Status:      "Completed",
		Priority:    "Low",
		DueDate:     time.Now().Add(time.Hour * 94),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	insertedResult, errInsert := suite.repository.CreateTask(context.TODO(), task)
	suite.Nil(errInsert, "Nil inserting task")

	fetchedTask, errFetch := suite.repository.FetchTaskByID(context.TODO(), insertedResult)
	suite.Nil(errFetch, "Nil fetching task")

	suite.Equal(insertedResult, fetchedTask.ID, "Inserted and fetched task ID should be the same")

	suite.Equal(task.Title, fetchedTask.Title, "Task titles should match")
	suite.Equal(task.Description, fetchedTask.Description, "Task descriptions should match")
	suite.Equal(task.Status, fetchedTask.Status, "Task status should match")
	suite.Equal(task.Priority, fetchedTask.Priority, "Task priority should match")
	suite.Equal(task.UserID, fetchedTask.UserID, "Task user ID should match")
}

func (suite *testRepositorySuite) TestFetchAllTasksEmptyDB() {
	_, err := suite.repository.FetchAllTasks(context.TODO())
	suite.Nil(err, "Error fetching all tasks from empty database")
	if err != nil {
		suite.Equal(err.Code, 500, "Nil code should be 500")
	}
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
	suite.Nil(errInsertOne, "Unwanted: Error inserting task One")
	_, errFetchOne := suite.repository.FetchTaskByID(context.TODO(), insertedTaskOne)
	suite.Nil(errFetchOne, "Unwanted: Error fetching task One")
	if errFetchOne != nil {
		suite.Equal(errFetchOne.Code, 500, "Nil code should be 500")
	}

	insertedTaskTwo, errInsertTwo := suite.repository.CreateTask(context.TODO(), task[1])
	suite.Nil(errInsertTwo, "Nil inserting task Two")
	_, errFetchTwo := suite.repository.FetchTaskByID(context.TODO(), insertedTaskTwo)
	suite.Nil(errFetchTwo, "Unwanted: Error fetching task Two")
	if errFetchTwo != nil {
		suite.Equal(errFetchTwo.Code, 500, "Nil code should be 500")
	}

	totalFetched, err := suite.repository.FetchAllTasks(context.TODO())
	suite.Nil(err, "Error fetching all tasks from populated database")
	if err != nil {
		suite.Equal(err.Code, 500, "Nil code should be 500")
	}
	suite.Equal(len(totalFetched), 2, "Total fetched tasks should be 2")
}

func (suite *testRepositorySuite) TestUpdatingTask() {
	// Create an initial task
	initialTask := domain.Task{
		Title:       "Schedule dentist appointment",
		Description: "Call the dentist to schedule an appointment.",
		Status:      "Completed",
		Priority:    "Low",
		DueDate:     time.Now().Add(time.Hour * 94),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	insertedResult, err := suite.repository.CreateTask(context.TODO(), initialTask)
	fmt.Println(insertedResult)
	suite.Nil(err, "Nil creating initial task")

	insertedTask, err := suite.repository.FetchTaskByID(context.TODO(), insertedResult)
	fmt.Println(insertedTask)
	suite.Nil(err, "Nil retriving initial task")

	// Update the task
	updateTask := domain.Task{
		ID:          insertedResult,
		Title:       "Schedule doctor's appointment",
		Description: "Call the doctors to schedule an appointment.",
		Status:      "Pending",
		Priority:    "High",
		DueDate:     time.Now().Add(time.Hour * 94),
		UpdatedAt:   time.Now(),
	}

	// Perform the update operation
	fetchedTask, errUpdate := suite.repository.UpdateTask(context.TODO(), updateTask)
	suite.Nil(errUpdate, "Nil updating task")

	suite.Equal(updateTask.Title, fetchedTask.Title, "Title should be updated correctly")
	suite.Equal(updateTask.Status, fetchedTask.Status, "Status should be updated correctly")
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

	// Create the task
	insertedResult, err := suite.repository.CreateTask(context.TODO(), initalTask)
	suite.Nil(err, "Nil creating initial task")

	// Perform deletion
	fetchedTask, errDelete := suite.repository.DeleteTask(context.TODO(), insertedResult)
	suite.Nil(errDelete, "Nil deleting task")

	// Check if the task ID matches
	suite.Equal(fetchedTask.ID, insertedResult, "Task ID should be the same")
}

func TestTaskRepositorySuite(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Nil loading .env file")
	}
	suite.Run(t, new(testRepositorySuite))
}
