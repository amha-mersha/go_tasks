package repositorie

import (
	"context"
	"net/http"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return UserRepository{
		Collection: collection,
	}
}

func (userRepo *UserRepository) FetchAllUsers(cxt context.Context) ([]domain.User, *domain.UserError) {
	filter := bson.D{{}}
	cursor, err := userRepo.Collection.Find(cxt, filter)
	if err != nil {
		return []domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	var users []domain.User
	err = cursor.All(cxt, &users)
	if err != nil {
		return []domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return users, nil
}

func (userRepo *UserRepository) FetchUserByID(cxt context.Context, ID string) (domain.User, *domain.UserError) {
	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	filter := bson.D{{"_id", taskID}}
	var retrivedUser domain.User
	err = userRepo.Collection.FindOne(cxt, filter).Decode(&retrivedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	return retrivedUser, nil
}

func (userRepo *UserRepository) FetchUserByUsername(cxt context.Context, username string) (domain.User, *domain.UserError) {
	filter := bson.D{{"username", username}}
	var retrivedUser domain.User
	err := userRepo.Collection.FindOne(cxt, filter).Decode(&retrivedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	return retrivedUser, nil
}

func (userRepo *UserRepository) FetchUserCount(cxt context.Context) (int, *domain.UserError) {
	usersCount, err := userRepo.Collection.EstimatedDocumentCount(cxt)
	if err != nil {
		return 0, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return int(usersCount), nil
}

func (userRepo *UserRepository) CreateUser(cxt context.Context, newUser domain.User) (string, *domain.UserError) {
	createdUser, err := userRepo.Collection.InsertOne(cxt, newUser)
	if err != nil {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	insertedID, ok := createdUser.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return insertedID.Hex(), nil
}

func (userRepo *UserRepository) UpdateUser(cxt context.Context, updateUser domain.User) (domain.User, *domain.UserError) {
	objectID, err := primitive.ObjectIDFromHex(updateUser.ID)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: "Invalid ID format", Code: http.StatusBadRequest}
	}
	filter := bson.D{{"_id", objectID}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	inserteUser := domain.User{
		Username: updateUser.Username,
		Role:     updateUser.Role,
		Password: updateUser.Password,
	}
	var returnedUser domain.User
	err = userRepo.Collection.FindOneAndUpdate(cxt, filter, bson.D{{"$set", inserteUser}}, opts).Decode(&returnedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return returnedUser, nil
}

func (userRepo *UserRepository) DeleteUser(cxt context.Context, ID string) (domain.User, *domain.UserError) {
	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	var returnedUser domain.User
	err = userRepo.Collection.FindOneAndDelete(cxt, bson.D{{"_id", taskID}}).Decode(&returnedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return returnedUser, nil
}
