package infrastructure

import (
	domain "task_manager_api/Domain"

	"golang.org/x/crypto/bcrypt"
)

/*
Accepts a reference to a user object. The function hashed the password
using bcrypt and replaced the `Password` field inplace.
*/
func HashPassword(password string) (string, domain.CodedError) {
	hashedPwd, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", domain.UserError{Message: "Internal server error: " + hashErr.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return string(hashedPwd), nil
}

/*
Accepts references to a stored user object (with a previously hashed
password) and checks it with the user that is currently being authenticated.
*/
func ValidatePassword(storedPassword string, currentPassword string) domain.CodedError {
	compErr := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(currentPassword))
	if compErr != nil {
		return domain.UserError{Message: "Incorrect password", Code: domain.ERR_UNAUTHORIZED}
	}

	return nil
}
