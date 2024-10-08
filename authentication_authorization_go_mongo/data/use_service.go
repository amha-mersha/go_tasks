package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserCustomClaim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserRole string `json:"userrole"`
	jwt.RegisteredClaims
}

func PostUserRegister(newUser models.User) error {
	newUser.Username = strings.TrimSpace(newUser.Username)
	newUser.Password = strings.TrimSpace(newUser.Password)
	// checking if the collection is empty or not
	documentCount, err := UserCollection.EstimatedDocumentCount(context.TODO(), options.EstimatedDocumentCount().SetMaxTime(2*time.Second))
	if err != nil {
		return fmt.Errorf(InternalServerError)
	}
	// searching for the user
	filter := bson.D{{"username", newUser.Username}}
	var retrivedUser models.User
	err = UserCollection.FindOne(context.TODO(), filter).Decode(&retrivedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			if documentCount > 0 && newUser.Role == "admin" {
				return fmt.Errorf("an admin needs to promot you to admin")
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf(InternalServerError)
			}
			newUser.Password = string(hashedPassword)
			if documentCount == 0 {
				newUser.Role = "admin"
			}
			_, err = UserCollection.InsertOne(context.TODO(), newUser)
			if err != nil {
				return fmt.Errorf(InternalServerError)
			}
			return nil
		} else {
			return fmt.Errorf(MalformedUsername)
		}
	}
	return fmt.Errorf(UserAlreadyExist)
}

func PostUserLogin(newUser models.User) (interface{}, error) {
	jwtSecret := []byte(os.Getenv("SIGNITURE_SECRET"))

	// filter out the user with the username
	filter := bson.D{{"username", newUser.Username}}
	var retrivedUser models.User

	err := UserCollection.FindOne(context.TODO(), filter).Decode(&retrivedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return nil, fmt.Errorf(UserNotFound)
		}
		return nil, fmt.Errorf(InternalServerError)
	}

	// checking if the username and password are correct
	if bcrypt.CompareHashAndPassword([]byte(retrivedUser.Password), []byte(newUser.Password)) != nil {
		return nil, fmt.Errorf(IncorrectCredentials)
	}

	// generating the jwt token
	claim := UserCustomClaim{
		Username: retrivedUser.Username,
		Password: retrivedUser.Password,
		UserRole: newUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(720 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf(InternalServerError)
	}
	return jwtToken, nil

}

func PostUserAssign(issuedUser models.IssuedUser) (interface{}, error) {

	// filter out the user with the username
	var retrivedUser models.IssuedUser

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{{"username", issuedUser.Username}}
	update := bson.M{"$set": bson.M{"role": issuedUser.Role}}

	err := UserCollection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&retrivedUser)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return models.Task{}, fmt.Errorf(IDNotFound)
	} else if err != nil {
		return models.Task{}, err
	}
	return retrivedUser, nil

}

func GetUserByUsername(username string) (models.User, error) {
	filter := bson.D{{"username", username}}
	var retrivedUser models.User
	err := UserCollection.FindOne(context.TODO(), filter).Decode(&retrivedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, fmt.Errorf(UserNotFound)
		} else {
			return models.User{}, fmt.Errorf(MalformedUsername)
		}
	}
	return retrivedUser, nil
}
