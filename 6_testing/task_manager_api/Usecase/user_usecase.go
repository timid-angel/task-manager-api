package usecase

import (
	"context"
	"net/mail"
	"strings"
	domain "task_manager_api/Domain"
	"time"

	"github.com/spf13/viper"
)

/* Implements the UserUsecaseInterface defined in `domain`*/
type UserUsecase struct {
	UserRespository    domain.UserRepositoryInterface
	Timeout            time.Duration
	HashUserPassword   func(user *domain.User) domain.CodedError
	SignJWTWithPayload func(username string, role string, tokenLifeSpan time.Duration, secret string) (string, domain.CodedError)
	ValidatePassword   func(storedUser *domain.User, currUser *domain.User) domain.CodedError
}

/* Validates the user data with business rules and calls the create function in the repository */
func (uC *UserUsecase) CreateUser(c context.Context, user domain.User) domain.CodedError {
	ctx, cancel := context.WithTimeout(c, uC.Timeout)
	defer cancel()

	user.Username = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(user.Username)), " ", "")
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
	usernameError := uC.UserRespository.CheckDuplicate(c, "username", user.Username, "An account with the provided username already exists")
	if usernameError != nil {
		return usernameError
	}

	// check for duplicate email
	emailError := uC.UserRespository.CheckDuplicate(c, "email", user.Email, "An account with the provided email already exists")
	if emailError != nil {
		return emailError
	}

	// hash the password before storing
	hashErr := uC.HashUserPassword(&user)
	if hashErr != nil {
		return hashErr
	}

	return uC.UserRespository.CreateUser(ctx, user)
}

/*
Validates the user credentials with the ones in the DB and returns
a signed JWT if the credentials match
*/
func (uC *UserUsecase) ValidateAndGetToken(c context.Context, user domain.User) (string, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, uC.Timeout)
	defer cancel()

	storedUser, queryErr := uC.UserRespository.GetByUsername(ctx, user.Username)
	if queryErr != nil {
		return "", queryErr
	}

	// compare the incoming password and the stored (previously hashed) password
	pwErr := uC.ValidatePassword(&storedUser, &user)
	if pwErr != nil {
		return "", pwErr
	}

	// signs token with the secret token in the env variables
	tkLifespan := time.Minute * time.Duration(viper.GetInt("TOKEN_LIFESPAN_MINUTES"))
	jwtSecret := viper.GetString("SECRET_TOKEN")
	return uC.SignJWTWithPayload(storedUser.Username, storedUser.Role, tkLifespan, jwtSecret)
}
