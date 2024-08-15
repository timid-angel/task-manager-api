package repository

import (
	"context"
	domain "task_manager_api/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* Implements the UserRespositoryInterface defined in `domain`*/
type UserRepository struct {
	Collection *mongo.Collection
}

/*
checks if an object that matches the provided the key-value pair
returns a CodedError if there are any matches or if there is an
error during the DB request
*/
func (uR *UserRepository) CheckDuplicate(c context.Context, key string, value interface{}, errorMessage string) domain.CodedError {
	result := uR.Collection.FindOne(c, bson.D{{Key: key, Value: value}})
	if result.Err() != nil && result.Err().Error() == mongo.ErrNoDocuments.Error() {
		return nil
	}

	if result.Err() == nil {
		return domain.UserError{Message: "Bad request: duplicate " + key, Code: domain.ERR_BAD_REQUEST}
	}

	return domain.UserError{Message: errorMessage, Code: domain.ERR_INTERNAL_SERVER}
}

/* Adds the user to the DB */
func (uR *UserRepository) CreateUser(c context.Context, user domain.User) domain.CodedError {
	_, err := uR.Collection.InsertOne(c, user)
	if err != nil {
		return domain.UserError{Message: "Internal server error: " + err.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return nil
}

/*
queries for a user with the provided username and returns the resulting
user object if it exists
*/
func (uR *UserRepository) GetByUsername(c context.Context, username string) (domain.User, domain.CodedError) {
	var storedUser domain.User
	result := uR.Collection.FindOne(c, bson.D{{Key: "username", Value: username}})
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

/* Promotes an account with role `user` to role `admin` */
func (uR *UserRepository) PromoteUser(c context.Context, username string) domain.CodedError {
	result := uR.Collection.FindOneAndUpdate(context.TODO(), bson.D{{Key: "username", Value: username}}, bson.D{{Key: "$set", Value: bson.D{{Key: "role", Value: "admin"}}}})
	if result.Err() != nil && result.Err().Error() == mongo.ErrNoDocuments.Error() {
		return domain.UserError{Message: "error: user not found", Code: domain.ERR_NOT_FOUND}
	}

	if result.Err() != nil {
		return domain.UserError{Message: "error: " + result.Err().Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return nil
}
