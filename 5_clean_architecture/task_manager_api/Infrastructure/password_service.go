package infrastructure

import (
	domain "task_manager_api/Domain"

	"golang.org/x/crypto/bcrypt"
)

/*
Accepts a reference to a user object. The function hashed the password
using bcrypt and replaced the `Password` field inplace.
*/
func HashUserPassword(user *domain.User) domain.CodedError {
	hashedPwd, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return domain.UserError{Message: "Internal server error: " + hashErr.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	user.Password = string(hashedPwd)
	return nil
}

/*
Accepts references to a stored user object (with a previously hashed
password) and checks it with the user that is currently being authenticated.
*/
func ValidatePassword(storedUser *domain.User, currUser *domain.User) domain.CodedError {
	compErr := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(currUser.Password))
	if compErr != nil {
		return domain.UserError{Message: "Incorrect password", Code: domain.ERR_UNAUTHORIZED}
	}

	return nil
}
