package services

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"task_manager_api/models"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

/*
Implements the CodedError interface and serves to communicate the nature of
certain errors directly to the controller.
*/
type UserError struct {
	message string
	code    int
}

func (err UserError) Error() string {
	return err.message
}

func (err UserError) GetCode() int {
	return err.code
}

/*
checks if an object that matches the provided filter object exists

	returns an error if the object exists
*/
func CheckDuplicate(filter bson.D, errorMessage string) CodedError {
	result := UserCollection.FindOne(context.TODO(), filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil
	}

	if result.Err() != nil {
		return UserError{message: "Internal server error", code: 500}
	}

	return UserError{message: errorMessage, code: 400}
}

/*
queries for a user with the provided username and returns the
results if the user exists
*/
func GetByUsername(username string) (models.User, error) {
	var storedUser models.User
	result := UserCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: username}})
	if result.Err() != nil {
		return models.User{}, UserError{message: "User not found", code: 404}
	}

	if err := result.Decode(&storedUser); err != nil {
		return models.User{}, UserError{message: "Error during binding", code: 404}
	}

	return storedUser, nil
}

/*
Adds the user to the DB after validating the fields
*/
func CreateUser(user models.User) CodedError {
	user.Username = strings.ToLower(strings.TrimSpace(user.Username))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Role = strings.ToLower(strings.TrimSpace(user.Role))

	// validate username
	if len(user.Username) < 3 {
		return UserError{message: "Username must be atleast 3 characters long", code: 400}
	}

	// validate email
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return UserError{message: "Invalid email", code: 400}
	}

	// validate passowrd
	if len(user.Password) < 8 {
		return UserError{message: "Password must be atleast 8 characters long", code: 400}
	}

	// validate role
	if user.Role != "admin" && user.Role != "user" {
		return UserError{message: "Invalid role: must be either 'user' or 'admin'", code: 400}
	}

	// check for duplicate username
	usernameError := CheckDuplicate(bson.D{{Key: "username", Value: user.Username}}, "An account with the provided username already exists")
	if usernameError != nil {
		return usernameError
	}

	// check for duplicate email
	emailError := CheckDuplicate(bson.D{{Key: "email", Value: user.Email}}, "An account with the provided email already exists")
	if emailError != nil {
		return emailError
	}

	// hash the password before storing
	hashedPwd, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return UserError{message: "Internal server error", code: 500}
	}

	user.Password = string(hashedPwd)
	_, err := UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return UserError{message: fmt.Sprintf("Internal server error: %v", err.Error()), code: 500}
	}

	return nil
}

/*
Checks if the passed user exists in the system before checking with
the hashed password. The function then signs a json-web-token after
signing it with the secret key obtained from the environment variables.

	returns signed token
*/
func ValidateAndGetToken(user models.User) (string, CodedError) {
	// query for the user
	var storedUser models.User
	result := UserCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: user.Username}})
	if result.Err() != nil {
		return "", UserError{message: "User not found", code: 404}
	}

	decodeErr := result.Decode(&storedUser)
	if decodeErr != nil {
		return "", UserError{message: "Internal server error", code: 500}
	}

	// compare the incoming password and the stored (previously hashed) password
	compErr := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if compErr != nil {
		return "", UserError{message: "Incorrect password", code: 401}
	}

	// signs token with the secret token in the env variables
	jwtSecret := []byte(os.Getenv("JWT_SECRET_TOKEN"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  user.Username,
		"expiresAt": time.Now().Add(time.Hour * 2),
	})
	jwtToken, signingErr := token.SignedString(jwtSecret)
	if signingErr != nil {
		return "", UserError{message: signingErr.Error(), code: 500}
	}

	return jwtToken, nil
}

/* Promotes an account with role `user` to role `admin` */
func PromoteUser(userID string) CodedError {
	result := UserCollection.FindOneAndUpdate(context.TODO(), bson.D{{Key: "id", Value: userID}}, bson.D{{Key: "role", Value: "admin"}})
	if result.Err() != nil && result.Err().Error() == mongo.ErrNoDocuments.Error() {
		return UserError{message: "error: user not found", code: 404}
	}

	if result.Err() != nil {
		return UserError{message: "error: " + result.Err().Error(), code: 500}
	}

	return nil
}
