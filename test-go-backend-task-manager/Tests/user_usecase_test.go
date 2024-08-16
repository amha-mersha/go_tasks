package tests

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	mocks "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/Mocks"
	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/infrastructure"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/usecases"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type userUsecaseSuite struct {
	suite.Suite
	repositorie *mocks.UserRepository
	usecase     domain.UserUsecase
	authService *mocks.AuthService
}

func (suite *userUsecaseSuite) SetupTest() {
	repo := new(mocks.UserRepository)
	suite.authService = new(mocks.AuthService)
	userUC := usecases.NewUserUsecase(repo, time.Second*2)
	suite.usecase = userUC
	suite.repositorie = repo
}

func (suite *userUsecaseSuite) TestGetAllUsers() {
	users := []domain.User{
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

	suite.repositorie.On("FetchAllUsers", mock.Anything).Return(users, nil)

	fetchedUsers, err := suite.usecase.GetAllUser(context.TODO())
	suite.Nil(err, "error should be nil")
	suite.Equal(users, fetchedUsers, "users should be equal")
}

func (suite *userUsecaseSuite) TestGetUserByID() {
	user := domain.User{
		ID:       "1",
		Username: "johndoe",
		Password: "password123",
		Role:     "admin",
	}
	suite.repositorie.On("FetchUserByID", mock.Anything, user.ID).Return(user, nil)

	fetchedUser, err := suite.usecase.GetUserByID(context.TODO(), user.ID)
	suite.Nil(err, "error should be nil")
	suite.Equal(user, fetchedUser, "users should be equal")
}

func (suite *userUsecaseSuite) TestGetUserByUsername() {
	user := domain.User{
		ID:       "1",
		Username: "johndoe",
		Password: "password123",
		Role:     "admin",
	}

	suite.repositorie.On("FetchUserByUsername", mock.Anything, user.Username).Return(user, nil)

	fetchedUser, err := suite.usecase.GetUserByUsername(context.TODO(), user.Username)
	suite.Nil(err, "error should be nil")
	suite.Equal(user, fetchedUser, "users should be equal")
}

func (suite *userUsecaseSuite) TestCreateUser_Positive() {
	ID := "1"
	user := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "admin",
	}
	suite.repositorie.On("FetchUserCount", mock.Anything).Return(0, nil)
	suite.repositorie.On("CreateUser", mock.Anything, user).Return(ID, nil)

	fetchID, err := suite.usecase.CreateUser(context.TODO(), user)
	suite.Nil(err, "error should be nil")
	suite.Equal(ID, fetchID, "users should be equal")
}

func (suite *userUsecaseSuite) TestCreateUser_Negative() {
	user := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}

	duplicateKeyErr := mongo.WriteException{
		WriteErrors: []mongo.WriteError{
			{
				Code:    11000,
				Message: "E11000 duplicate key error collection: testdb.testcollection index: _id_ dup key: { _id: \"12345\" }",
			},
		},
	}
	suite.repositorie.On("FetchUserCount", mock.Anything).Return(3, nil)
	suite.repositorie.On("CreateUser", mock.Anything, user).Return("", &domain.UserError{Message: duplicateKeyErr.Error(), Code: 500})
	fetchID, errFetch := suite.usecase.CreateUser(context.TODO(), user)
	fmt.Println(reflect.TypeOf(errFetch))

	suite.NotNil(errFetch, "Error should not be nil because of duplicate key error")
	suite.Empty(fetchID, "fetchID should be empty when an error occurs")
}

func (suite *userUsecaseSuite) TestUpdateUser() {
	user := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}

	suite.repositorie.On("UpdateUser", mock.Anything, user).Return(user, nil)
	userRetrived, err := suite.usecase.UpdateUser(context.TODO(), user)
	suite.Nil(err, "error should be nil")
	suite.Equal(user, userRetrived, "users should be equal")
}

func (suite *userUsecaseSuite) TestDeleteUser_Positive() {
	authorityUser := domain.User{
		ID:       "2",
		Username: "johndoe",
		Password: "password123",
		Role:     "admin",
	}
	deleteUser := domain.User{
		ID:       "1",
		Username: "janedoe",
		Password: "password123",
		Role:     "user",
	}

	suite.repositorie.On("FetchUserByID", mock.Anything, authorityUser.ID).Return(authorityUser, nil)
	suite.repositorie.On("DeleteUser", mock.Anything, deleteUser.ID).Return(deleteUser, nil)
	deletedUser, errDelete := suite.usecase.DeleteUser(context.TODO(), authorityUser, deleteUser.ID)
	suite.Nil(errDelete, "error should be nil")
	suite.Equal(deleteUser, deletedUser, "users should be equal")
}

func (suite *userUsecaseSuite) TestDeleteUserSelfDeletion() {
	authorityUser := domain.User{
		ID:       "1",
		Username: "johndoe",
		Password: "password123",
		Role:     "admin",
	}
	deleteUser := domain.User{
		ID:       "1",
		Username: "janedoe",
		Password: "password123",
		Role:     "user",
	}

	suite.repositorie.On("FetchUserByID", mock.Anything, authorityUser.ID).Return(authorityUser, nil)
	suite.repositorie.On("DeleteUser", mock.Anything, deleteUser.ID).Return(deleteUser, nil)
	_, errDelete := suite.usecase.DeleteUser(context.TODO(), authorityUser, deleteUser.ID)
	suite.NotNil(errDelete, "error should be nil, trying self delete")

}

func (suite *userUsecaseSuite) TestDeleteUserUnauthorizedDeletion() {
	authorityUser := domain.User{
		ID:       "1",
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}
	deleteUser := domain.User{
		ID:       "1",
		Username: "janedoe",
		Password: "password123",
		Role:     "user",
	}

	suite.repositorie.On("FetchUserByID", mock.Anything, authorityUser.ID).Return(authorityUser, nil)
	suite.repositorie.On("DeleteUser", mock.Anything, deleteUser.ID).Return(deleteUser, nil)
	_, errDelete := suite.usecase.DeleteUser(context.TODO(), authorityUser, deleteUser.ID)
	suite.NotNil(errDelete, "error should be nil, trying unauthorized deletion")

}

func (suite *userUsecaseSuite) TestLoginUser() {
	os.Setenv("SIGNITURE_TIME_DURATION", "3600")
	os.Setenv("SIGNITURE_SECRET", "mysecretkey")

	hashedPassword, err := infrastructure.HashPassword("password123")
	suite.Nil(err, "bcrypt hash generation should not fail")
	suite.NotNil(hashedPassword, "hashed password should not be nil")

	user := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}

	storedUser := domain.User{
		Username: "johndoe",
		Password: hashedPassword,
		Role:     "user",
	}

	suite.repositorie.On("FetchUserByUsername", mock.Anything, user.Username).Return(storedUser, nil)
	suite.repositorie.On("ValidatePassword", storedUser.Password, user.Password).Return(nil)

	timeDurationEnv, err := strconv.ParseInt(os.Getenv("SIGNITURE_TIME_DURATION"), 10, 64)
	suite.Nil(err, "parsing SIGNITURE_TIME_DURATION should not fail")
	expectedToken, err := infrastructure.CreateJWTToken(user.Username, user.Role, time.Duration(timeDurationEnv)*time.Second)
	suite.Nil(err, "JWT token creation should not fail")
	suite.NotNil(expectedToken, "expectedToken should not be nil")

	suite.repositorie.On("CreateJWTToken", user.Username, user.Role, time.Duration(timeDurationEnv)*time.Second).Return(expectedToken, nil)
	token, err := suite.usecase.LoginUser(context.TODO(), user)

	suite.Nil(err, "error should be nil")
	suite.Equal(expectedToken, token, "tokens should be equal")
}

func (suite *userUsecaseSuite) TestLoginUserUsernameNotExist() {
	os.Setenv("SIGNITURE_TIME_DURATION", "3600")
	os.Setenv("SIGNITURE_SECRET", "mysecretkey")

	hashedPassword, err := infrastructure.HashPassword("password123")
	suite.Nil(err, "bcrypt hash generation should not fail")
	suite.NotNil(hashedPassword, "hashed password should not be nil")

	user := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}

	storedUser := domain.User{
		Username: "johndoe",
		Password: hashedPassword,
		Role:     "user",
	}

	errFetch := &domain.UserError{Message: "User not found", Code: 404}
	suite.repositorie.On("FetchUserByUsername", mock.Anything, user.Username).Return(domain.User{}, errFetch)
	suite.repositorie.On("ValidatePassword", storedUser.Password, user.Password).Return(nil)

	timeDurationEnv, err := strconv.ParseInt(os.Getenv("SIGNITURE_TIME_DURATION"), 10, 64)
	suite.Nil(err, "parsing SIGNITURE_TIME_DURATION should not fail")
	expectedToken, err := infrastructure.CreateJWTToken(user.Username, user.Role, time.Duration(timeDurationEnv)*time.Second)
	suite.Nil(err, "JWT token creation should not fail")
	suite.NotNil(expectedToken, "expectedToken should not be nil")

	suite.repositorie.On("CreateJWTToken", user.Username, user.Role, time.Duration(timeDurationEnv)*time.Second).Return(expectedToken, nil)
	_, err = suite.usecase.LoginUser(context.TODO(), user)

	suite.NotNil(err, "error should be nil")
}

func (suite *userUsecaseSuite) TestLoginUserValidation_Negative() {
	os.Setenv("SIGNITURE_TIME_DURATION", "3600")
	os.Setenv("SIGNITURE_SECRET", "mysecretkey")

	hashedPassword := "$2a$10$wH9lYx8S.xVx0o5eKf2vve5TTNN/JUVyyX9B3RYOP3KrmRnZRCE1a"
	suite.authService.On("HashPassword", "password123").Return(hashedPassword, nil)

	user := domain.User{
		Username: "johndoe",
		Password: "password123",
		Role:     "user",
	}

	storedUser := domain.User{
		Username: "johndoe",
		Password: hashedPassword,
		Role:     "user",
	}

	suite.repositorie.On("FetchUserByUsername", mock.Anything, user.Username).Return(storedUser, nil)

	errValidation := domain.UserError{Message: "Password validation failed", Code: 500}
	suite.authService.On("ValidatePassword", storedUser.Password, user.Password).Return(&errValidation)

	timeDurationEnv, err := strconv.ParseInt(os.Getenv("SIGNITURE_TIME_DURATION"), 10, 64)
	suite.Nil(err, "parsing SIGNITURE_TIME_DURATION should not fail")
	expectedToken := "some.jwt.token"
	suite.authService.On("CreateJWTToken", user.Username, user.Role, time.Duration(timeDurationEnv)*time.Second).Return(expectedToken, nil)

	_, errLogin := suite.usecase.LoginUser(context.TODO(), user)

	suite.NotNil(errLogin, "error should be not nil")
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(userUsecaseSuite))
}
