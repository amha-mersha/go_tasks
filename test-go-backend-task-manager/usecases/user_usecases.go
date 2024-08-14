package usecases

import (
	"context"
	"net/http"
	"time"

	domain "github.com/amha-mersha/go_tasks/go-backend-clean-architecture/domains"
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

func (userUC userUsercase) CreateUser(cxt context.Context, newUser domain.User) (domain.UserSuccess, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	documentCount, err := userUC.userRepository.FetchUserCount(context)
	if err != nil {
		return domain.UserSuccess{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	newUser.ID = ""
	_, err = userUC.userRepository.FetchUserByUsername(context, newUser.Username)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			if documentCount == 0 {
				newUser.Role = "admin"
				return userUC.userRepository.CreateUser(context, newUser)
			}
			if newUser.Role == "admin" {
				return domain.UserSuccess{}, &domain.UserError{Message: "Admin already exists", Code: http.StatusConflict}
			}
			return userUC.userRepository.CreateUser(context, newUser)
		}
		return domain.UserSuccess{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return domain.UserSuccess{}, &domain.UserError{Message: "Username already exists", Code: http.StatusConflict}
}

func (userUC userUsercase) UpdateUser(cxt context.Context, updateUser domain.User) (domain.UserSuccess, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	return userUC.userRepository.UpdateUser(context, updateUser)
}

func (userUC userUsercase) DeleteUser(cxt context.Context, authority domain.User, deleteID string) (domain.UserSuccess, *domain.UserError) {
	context, cancel := context.WithTimeout(cxt, userUC.timeout)
	defer cancel()
	fetchedAuthority, err := userUC.userRepository.FetchUserByID(context, authority.ID)
	if err != nil {
		return domain.UserSuccess{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	if fetchedAuthority.Role != "admin" {
		return domain.UserSuccess{}, &domain.UserError{Message: "Unauthorized to make this update", Code: http.StatusUnauthorized}
	}
	return userUC.userRepository.DeleteUser(context, deleteID)
}
