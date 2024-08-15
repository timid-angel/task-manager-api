package tests

import (
	"context"
	"log"
	domain "task_manager_api/Domain"
	repository "task_manager_api/Repository"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepositorySuite struct {
	suite.Suite
	UserRepository *repository.UserRepository
	collection     *mongo.Collection
}

func (suite *userRepositorySuite) SetupSuite() {
	clientOptions := options.Client().ApplyURI(viper.GetString("DB_ADDRESS"))
	client, connectionErr := mongo.Connect(context.TODO(), clientOptions)
	if connectionErr != nil {
		log.Fatalf("Error: %v", connectionErr.Error())
	}

	databse := client.Database(viper.GetString("TEST_DB_NAME"))
	collection := databse.Collection("users")
	suite.collection = collection
	collection.DeleteMany(context.TODO(), bson.D{{}})
	SetupUserCollection(suite.collection)
	suite.UserRepository = &repository.UserRepository{Collection: collection}
}

func (suite *userRepositorySuite) SetupTest() {
	suite.UserRepository.Collection.DeleteMany(context.TODO(), bson.D{{}})
}

// Tests CreateUser function
func (suite *userRepositorySuite) TestSignup() {
	user := domain.User{
		Username: "username12",
		Password: "password12",
		Email:    "mailer@mail.com",
		Role:     "user",
	}

	err := suite.UserRepository.CreateUser(context.TODO(), user)
	suite.NoError(err, "no error when creating account")
}

// Tests CreateUser function with duplicate emails
func (suite *userRepositorySuite) TestSignup_DupEmail() {
	user := domain.User{
		Username: "username12",
		Password: "password3212",
		Email:    "mailer@mail.com",
		Role:     "user",
	}

	err := suite.UserRepository.CreateUser(context.TODO(), user)
	suite.NoError(err, "no error when creating account")
	user.Username = "username34"
	err = suite.UserRepository.CreateUser(context.TODO(), user)
	suite.Error(err, "error when creating account")
}

// Tests CreateUser function with duplicate usernames
func (suite *userRepositorySuite) TestSignup_DupUsername() {
	user := domain.User{
		Username: "username12",
		Password: "password12",
		Email:    "mailer@mail.com",
		Role:     "user",
	}

	err := suite.UserRepository.CreateUser(context.TODO(), user)
	suite.NoError(err, "no error when creating account")
	user.Email = "changedMailer@mail.com"
	err = suite.UserRepository.CreateUser(context.TODO(), user)
	suite.Error(err, "error when creating account")
}

// Tests GetDuplicate function with duplicate usernames
func (suite *userRepositorySuite) TestCheckDuplicates() {
	user := domain.User{
		Username: "username12",
		Password: "password12",
		Email:    "mailer@mail.com",
		Role:     "user",
	}

	err := suite.UserRepository.CreateUser(context.TODO(), user)
	suite.NoError(err, "no error when creating account")

	err = suite.UserRepository.CheckDuplicate(context.TODO(), "email", user.Email, "")
	suite.Error(err, "error when same email is passed")
	err = suite.UserRepository.CheckDuplicate(context.TODO(), "email", "differentEmail@mail.com", "")
	suite.NoError(err, "no error when a different email is passed")

	err = suite.UserRepository.CheckDuplicate(context.TODO(), "username", user.Username, "")
	suite.Error(err, "error when same username is passed")
	err = suite.UserRepository.CheckDuplicate(context.TODO(), "username", "username32", "")
	suite.NoError(err, "no error when a different username is passed")
}

func (suite *userRepositorySuite) TestGetByUsername() {
	user := domain.User{
		Username: "username12",
		Password: "password12",
		Email:    "mailer@mail.com",
		Role:     "user",
	}

	err := suite.UserRepository.CreateUser(context.TODO(), user)
	suite.NoError(err, "no error when creating account")

	foundUser, err := suite.UserRepository.GetByUsername(context.TODO(), user.Username)
	suite.NoError(err, "error when same email is passed")
	suite.Equal(user.Username, foundUser.Username)
}

func (suite *userRepositorySuite) TestPromoteUser_Positive() {
	user := domain.User{
		Username: "username12",
		Password: "password12",
		Email:    "mailer@mail.com",
		Role:     "user",
	}

	err := suite.UserRepository.CreateUser(context.TODO(), user)
	suite.NoError(err, "no error when creating account")

	err = suite.UserRepository.PromoteUser(context.TODO(), user.Username)
	suite.NoError(err, "no error when promoting user")

	user, err = suite.UserRepository.GetByUsername(context.TODO(), user.Username)
	suite.NoError(err, "no error when fetching user")
	suite.Equal("admin", user.Role)
}

func (suite *userRepositorySuite) TestPromoteUser_Negative() {
	username := "solitary_confinement"

	err := suite.UserRepository.PromoteUser(context.TODO(), username)
	suite.Error(err, "no error when promoting user")
	suite.Equal(domain.ERR_NOT_FOUND, err.GetCode())
}

func TestUserRepositorySuite(t *testing.T) {
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	suite.Run(t, new(userRepositorySuite))
}
