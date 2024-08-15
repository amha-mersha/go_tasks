package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/infrastructure"
	repositorie "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/repositories"
	"github.com/joho/godotenv"
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

func (suite *userRepositorySuite) SetupSuite() {
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
	repo := repositorie.NewUserRepository(database.Collection(os.Getenv("DB_USER_COLLECTION_NAME")))
	err = infrastructure.EstablisUniqueUsernameIndex(repo.Collection, "username")
	if err != nil {
		panic(err)
	}
	suite.repository = &repo

}

func (suite *userRepositorySuite) TearDownSuite() {
	if suite.client != nil {
		if err := suite.client.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}
}

func (suite *userRepositorySuite) SetupTest() {
	if _, err := suite.repository.Collection.DeleteMany(context.TODO(), bson.D{}); err != nil {
		log.Println("Error deleting all documents from collection")
	}
}

func (suite *userRepositorySuite) TestFetchAllUsers() {
	testUsers := []domain.User{
		{
			Username: "johndoe",
			Password: "password123",
			Role:     "admin",
		},
		{
			Username: "janedoe",
			Password: "securepass",
			Role:     "user",
		},
		{
			Username: "guestuser",
			Password: "guestpass",
			Role:     "guest",
		},
	}
	for _, task := range testUsers {
		suite.repository.Collection.InsertOne(context.TODO(), task)
	}
	fetchedUsers, errFetch := suite.repository.FetchAllUsers(context.TODO())
	suite.Nil(errFetch, "No error fetching all tasks")
	suite.Equal(3, len(fetchedUsers), "Fetched users should be 3")
	count, err := suite.repository.FetchUserCount(context.TODO())
	suite.Nil(err, "No error fetching user count")
	suite.Equal(count, 3, "User count should be 3")
}

func (suite *userRepositorySuite) TestFetchEmptyDatabase() {
	_, errFetch := suite.repository.FetchAllUsers(context.TODO())
	suite.Nil(errFetch, "Error fetching empty database")
}

func (suite *userRepositorySuite) TestWithUserCreation() {
	user1 := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}
	insertedUser, errInsert := suite.repository.CreateUser(context.TODO(), user1)
	suite.Nil(errInsert, "Nil inserting user for the first time")

	retrivedUserOne, errRetrive := suite.repository.FetchUserByID(context.TODO(), insertedUser)
	suite.Nil(errRetrive, "Nil fetching user for the first time")
	suite.Equal(user1.Username, retrivedUserOne.Username, "Usernames should match")
	suite.Equal(user1.Role, retrivedUserOne.Role, "Roles should match")
	suite.Equal(user1.Password, retrivedUserOne.Password, "Passwords should match")

	user2 := domain.User{
		Username: "janedoe",
		Password: "securepass",
		Role:     "admin",
	}

	_, errInsert = suite.repository.CreateUser(context.TODO(), user2)
	suite.Nil(errInsert, "Nil inserting user for the second time")
}

func (suite *userRepositorySuite) TestUserCreation() {
	user1 := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}
	_, errInsert := suite.repository.CreateUser(context.TODO(), user1)
	suite.Nil(errInsert, "Nil inserting user for the first time")

	user2 := domain.User{
		Username: "johndoe",
		Password: "password456",
		Role:     "admin",
	}
	_, errInsert = suite.repository.CreateUser(context.TODO(), user2)
	suite.NotNil(errInsert, errInsert)
}

func (suite *userRepositorySuite) TestUpdateUser() {
	user1 := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}

	insertedUser, errInsert := suite.repository.CreateUser(context.TODO(), user1)
	suite.Nil(errInsert, "Nil inserting user for the first time")

	user2 := domain.User{
		ID:       insertedUser,
		Username: "johndoe",
		Password: "password456",
		Role:     "admin",
	}

	_, errInser := suite.repository.UpdateUser(context.TODO(), user2)
	suite.Nil(errInser, "Nil updating user")
	retrivedUserOne, errRetrive := suite.repository.FetchUserByID(context.TODO(), insertedUser)
	suite.Nil(errRetrive, "Nil fetching user for the first time")
	suite.Equal(user2.Username, retrivedUserOne.Username, "Usernames should match")
	suite.Equal(user2.Role, retrivedUserOne.Role, "Roles should match")
	suite.Equal(user2.Password, retrivedUserOne.Password, "Passwords should match")
}

func (suite *userRepositorySuite) TestDeleteUser() {
	user1 := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}

	insertedUser, errInsert := suite.repository.CreateUser(context.TODO(), user1)
	suite.Nil(errInsert, "Nil inserting user for the first time")

	deletedUser, errDelete := suite.repository.DeleteUser(context.TODO(), insertedUser)
	suite.Nil(errDelete, "Nil deleting user")
	user1.ID = insertedUser
	suite.Equal(user1, deletedUser, "Deleted user should be the same as inserted user")

	count, errRetrive := suite.repository.FetchUserCount(context.TODO())
	suite.Nil(errRetrive, "Nil fetching user count")
	suite.Equal(0, count, "User count should be 0")
}

func TestUserRepositorySuite(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Nil loading .env file")
	}
	suite.Run(t, new(userRepositorySuite))
}
