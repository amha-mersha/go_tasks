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
	userDatabase   *mongo.Database
	userCollection string
}

func NewUserRepository(database *mongo.Database, collection string) UserRepository {
	return UserRepository{
		userDatabase:   database,
		userCollection: collection,
	}
}

func (userRepo *UserRepository) FetchAllUsers(cxt context.Context) ([]domain.User, *domain.UserError) {
	collection := userRepo.userDatabase.Collection(userRepo.userCollection)
	filter := bson.D{{}}
	cursor, err := collection.Find(cxt, filter)
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
	collection := userRepo.userDatabase.Collection(userRepo.userCollection)
	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	filter := bson.D{{"_id", taskID}}
	var retrivedUser domain.User
	err = collection.FindOne(cxt, filter).Decode(&retrivedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	return retrivedUser, nil
}

func (userRepo *UserRepository) FetchUserByUsername(cxt context.Context, username string) (domain.User, *domain.UserError) {
	collection := userRepo.userDatabase.Collection(userRepo.userCollection)
	filter := bson.D{{"username", username}}
	var retrivedUser domain.User
	err := collection.FindOne(cxt, filter).Decode(&retrivedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	return retrivedUser, nil
}

func (userRepo *UserRepository) FetchUserCount(cxt context.Context) (int, *domain.UserError) {
	collection := userRepo.userDatabase.Collection(userRepo.userCollection)
	usersCount, err := collection.EstimatedDocumentCount(cxt)
	if err != nil {
		return 0, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return int(usersCount), nil
}

func (userRepo *UserRepository) CreateUser(cxt context.Context, newUser domain.User) (string, *domain.UserError) {
	collection := userRepo.userDatabase.Collection(userRepo.userCollection)
	createdUser, err := collection.InsertOne(cxt, newUser)
	if err != nil {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	insertedID, ok := createdUser.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return insertedID.Hex(), nil
}

func (userRepo *UserRepository) UpdateUser(cxt context.Context, updateTask domain.User) (domain.User, *domain.UserError) {
	collection := userRepo.userDatabase.Collection(userRepo.userCollection)
	filter := bson.D{{"_id", updateTask.ID}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var returnedUser domain.User
	err := collection.FindOneAndUpdate(cxt, filter, bson.D{{"$set", updateTask}}, opts).Decode(&returnedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return returnedUser, nil
}

func (userRepo *UserRepository) DeleteUser(cxt context.Context, ID string) (domain.User, *domain.UserError) {
	collection := userRepo.userDatabase.Collection(userRepo.userCollection)
	taskID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	var returnedUser domain.User
	err = collection.FindOneAndDelete(cxt, bson.D{{"_id", taskID}}).Decode(&returnedUser)
	if err != nil {
		return domain.User{}, &domain.UserError{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return returnedUser, nil
}
