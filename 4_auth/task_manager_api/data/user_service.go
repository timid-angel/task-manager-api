package services

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"task_manager_api/models"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type CodedError interface {
	error
	GetCode() int
}

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

func CheckDuplicate(filter bson.D, errorMessage string, user *models.User) CodedError {
	var foundUser models.User
	result := UserCollection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return UserError{message: "Internal servor error", code: 500}
	}

	result.Decode(&user)
	if foundUser.Username == user.Username {
		return UserError{message: errorMessage, code: 400}
	}

	return nil
}

func CreateUser(user models.User) CodedError {
	user.Username = strings.ToLower(strings.TrimSpace(user.Username))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Password = strings.ToLower(strings.TrimSpace(user.Password))

	// validate username
	if len(user.Username) <= 3 {
		return UserError{message: "Username must be atleast 3 characters long", code: 400}
	}

	// validate email
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return UserError{message: "Invalid email", code: 400}
	}

	// validate passowrd
	if len(user.Password) <= 8 {
		return UserError{message: "Password must be atleast 8 characters long", code: 400}
	}

	// check for duplicate username
	usernameError := CheckDuplicate(bson.D{{Key: "username", Value: user.Username}}, "An account with the provided username already exists", &user)
	if usernameError != nil {
		return usernameError
	}

	// check for duplicate email
	emailError := CheckDuplicate(bson.D{{Key: "email", Value: user.Email}}, "An account with the provided email already exists", &user)
	if emailError != nil {
		return emailError
	}

	var foundUser models.User
	pwResult := UserCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: user.Username}})
	if pwResult.Err() != nil {
		return UserError{message: "Internal servor error", code: 500}
	}

	pwResult.Decode(&user)
	if foundUser.Username == user.Username {
		return UserError{message: "An account with the username exists", code: 400}
	}

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

// TODO: Add more map claims
func ValidateAndGetToken(user models.User) (string, CodedError) {
	var storedUser models.User
	result := UserCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}})
	if result.Err() != nil {
		return "", UserError{message: "User not found", code: 401}
	}

	decodeErr := result.Decode(&storedUser)
	if decodeErr != nil {
		return "", UserError{message: "Internal server error", code: 500}
	}

	compErr := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if compErr != nil {
		return "", UserError{message: "Incorrect password", code: 401}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{})
	jwtToken, signingErr := token.SignedString(os.Getenv("JWT_SECRET_TOKEN"))
	if signingErr != nil {
		return "", UserError{message: "Internal server error", code: 500}
	}

	return jwtToken, nil
}
