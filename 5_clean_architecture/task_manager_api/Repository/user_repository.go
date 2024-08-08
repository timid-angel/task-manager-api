package repository

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	domain "task_manager_api/Domain"
	infrastructure "task_manager_api/Infrastructure"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

/*
checks if an object that matches the provided filter object exists

	returns an error if the object exists
*/
func CheckDuplicate(c context.Context, collection *mongo.Collection, filter bson.D, errorMessage string) domain.CodedError {
	result := collection.FindOne(c, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil
	}

	if result.Err() != nil {
		return domain.UserError{Message: "Internal server error", Code: domain.ERR_INTERNAL_SERVER}
	}

	return domain.UserError{Message: errorMessage, Code: domain.ERR_INTERNAL_SERVER}
}

/*
queries for a user with the provided username and returns the
results if the user exists
*/
func GetByUsername(c context.Context, collection *mongo.Collection, username string) (domain.User, error) {
	var storedUser domain.User
	result := collection.FindOne(c, bson.D{{Key: "username", Value: username}})
	if result.Err() != nil && result.Err().Error() == mongo.ErrNoDocuments.Error() {
		return domain.User{}, domain.UserError{Message: "User not found", Code: domain.ERR_BAD_REQUEST}
	}

	if result.Err() != nil {
		return domain.User{}, domain.UserError{Message: "Internal server error: " + result.Err().Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	if err := result.Decode(&storedUser); err != nil {
		return domain.User{}, domain.UserError{Message: "Error during binding", Code: domain.ERR_BAD_REQUEST}
	}

	return storedUser, nil
}

/*
Adds the user to the DB after validating the fields
*/
func (uR *UserRepository) CreateUser(c context.Context, user domain.User) domain.CodedError {
	user.Username = strings.ToLower(strings.TrimSpace(user.Username))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Role = strings.ToLower(strings.TrimSpace(user.Role))

	// validate username
	if len(user.Username) < 3 {
		return domain.UserError{Message: "Username must be atleast 3 characters long", Code: domain.ERR_BAD_REQUEST}
	}

	// validate email
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return domain.UserError{Message: "Invalid email", Code: domain.ERR_BAD_REQUEST}
	}

	// validate passowrd
	if len(user.Password) < 8 {
		return domain.UserError{Message: "Password must be atleast 8 characters long", Code: domain.ERR_BAD_REQUEST}
	}

	// validate role
	if user.Role != "admin" && user.Role != "user" {
		return domain.UserError{Message: "Invalid role: must be either 'user' or 'admin'", Code: domain.ERR_BAD_REQUEST}
	}

	// check for duplicate username
	usernameError := CheckDuplicate(c, uR.Collection, bson.D{{Key: "username", Value: user.Username}}, "An account with the provided username already exists")
	if usernameError != nil {
		return usernameError
	}

	// check for duplicate email
	emailError := CheckDuplicate(c, uR.Collection, bson.D{{Key: "email", Value: user.Email}}, "An account with the provided email already exists")
	if emailError != nil {
		return emailError
	}

	// hash the password before storing
	hashErr := infrastructure.HashUserPassword(&user)
	if hashErr != nil {
		return hashErr
	}

	_, err := uR.Collection.InsertOne(c, user)
	if err != nil {
		return domain.UserError{Message: "Internal server error: " + err.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return nil
}

/*
Checks if the passed user exists in the system before checking with
the hashed password. The function then signs a json-web-token after
signing it with the secret key obtained from the environment variables.

	returns signed token
*/
func (uR *UserRepository) ValidateAndGetToken(c context.Context, user domain.User) (string, domain.CodedError) {
	// query for the user
	var storedUser domain.User
	result := uR.Collection.FindOne(c, bson.D{{Key: "username", Value: user.Username}})
	if result.Err() != nil {
		return "", domain.UserError{Message: "User not found", Code: domain.ERR_NOT_FOUND}
	}

	decodeErr := result.Decode(&storedUser)
	if decodeErr != nil {
		return "", domain.UserError{Message: "Internal server error: " + decodeErr.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	// compare the incoming password and the stored (previously hashed) password
	pwErr := infrastructure.ValidatePassword(&storedUser, &user)
	if pwErr != nil {
		return "", pwErr
	}

	// signs token with the secret token in the env variables
	tkLifespan := time.Minute * time.Duration(viper.GetInt("TOKEN_LIFESPAN_MINUTES"))
	jwtSecret := viper.GetString("SECRET_TOKEN")
	fmt.Println(user.Role)
	return infrastructure.SignJWTWithPayload(storedUser.Username, storedUser.Role, tkLifespan, jwtSecret)
}
