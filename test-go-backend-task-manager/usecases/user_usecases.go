package usecases

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
)

type userUsercase struct {
	userRepository domain.UserRepository
	timeout        time.Duration
}

func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration) userUsercase {
	return userUsercase{
		userRepository: userRepo,
		timeout:        timeout,
	}
}

func (userUC userUsercase) GetAllUser(cxt context.Context) ([]domain.User, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	return userUC.userRepository.FetchAllUsers(context)

}

func (userUC userUsercase) GetUserByID(cxt context.Context, userID string) (domain.User, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	return userUC.userRepository.FetchUserByID(context, userID)

}

func (userUC userUsercase) GetUserByUsername(cxt context.Context, username string) (domain.User, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	return userUC.userRepository.FetchUserByUsername(context, username)
}

func (userUC userUsercase) CreateUser(cxt context.Context, newUser domain.User) (string, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	documentCount, err := userUC.userRepository.FetchUserCount(context)
	if err != nil {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	newUser.ID = ""
	if documentCount == 0 {
		newUser.Role = "admin"
	} else {
		newUser.Role = "user"
	}
	inserted, err := userUC.userRepository.CreateUser(context, newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", &domain.UserError{Message: "Username already exists", Code: http.StatusConflict}
		} else {
			return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
		}
	}
	return inserted, nil

}

func (userUC userUsercase) UpdateUser(cxt context.Context, updateUser domain.User) (domain.User, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	return userUC.userRepository.UpdateUser(context, updateUser)
}

func (userUC userUsercase) DeleteUser(cxt context.Context, authority domain.User, deleteID string) (domain.User, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	fetchedAuthority, err := userUC.userRepository.FetchUserByID(context, authority.ID)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	if fetchedAuthority.Role != "admin" {
		return domain.User{}, &domain.UserError{Message: "Unauthorized to make this update", Code: http.StatusUnauthorized}
	}
	return userUC.userRepository.DeleteUser(context, deleteID)
}

func (userUC userUsercase) LoginUser(cxt context.Context, loggingUser domain.User) (string, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	result, err := userUC.userRepository.FetchUserByUsername(context, loggingUser.Username)
	if err != nil {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	if err := infrastructure.ValidatePassword(result.Password, loggingUser.Password); result.Role != loggingUser.Role || err != nil {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	timeDurationEnv, errDuration := strconv.ParseInt(os.Getenv("SIGNITURE_TIME_DURATION"), 10, 64)
	if errDuration != nil {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	token, errToken := infrastructure.CreateJWTToken(result.Username, result.Role, time.Duration(timeDurationEnv)*time.Second)
	if errToken != nil {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return token, nil
}
