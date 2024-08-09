package tests

import (
	"context"
	"strings"
	domain "task_manager_api/Domain"
	mocks "task_manager_api/Mocks"
	usecase "task_manager_api/Usecase"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userUsecaseSuite struct {
	suite.Suite
	repository *mocks.UserRepositoryInterface
	usecase    usecase.UserUsecase
}

func SanitizeUser(user *domain.User) {
	user.Username = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(user.Username)), " ", "")
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Role = strings.TrimSpace(strings.ToLower(user.Role))
}

func (suite *userUsecaseSuite) SetupSuite() {
	suite.repository = new(mocks.UserRepositoryInterface)
	suite.usecase = usecase.UserUsecase{
		UserRespository: suite.repository,
		Timeout:         2,
	}

	// mocks HashUserPassword to return no errors and make no changes to the user struct
	suite.usecase.HashUserPassword = func(password string) (string, domain.CodedError) {
		return password, nil
	}

	// mocks ValidatePassword to return a nil error only when the password is "valid_password"
	suite.usecase.ValidatePassword = func(storedPassword string, currentPassword string) domain.CodedError {
		if storedPassword == "valid_password" {
			return nil
		}

		return domain.UserError{Message: "Incorrect password", Code: domain.ERR_UNAUTHORIZED}
	}

	// mocks SignJWTWithPayload with return types dependent on the role
	suite.usecase.SignJWTWithPayload = func(username, role string, tokenLifeSpan time.Duration, secret string) (string, domain.CodedError) {
		if role == "admin" || role == "user" {
			return "valid_token", nil
		}

		return "", domain.UserError{Message: "internal server error", Code: domain.ERR_INTERNAL_SERVER}
	}
}

func (suite *userUsecaseSuite) SetupTest() {
	suite.repository = new(mocks.UserRepositoryInterface)
	suite.usecase.UserRespository = suite.repository
}

func (suite *userUsecaseSuite) TestCreateUser_Positive() {
	user := domain.User{
		Username: " useSDSALJDFname  ",
		Email:    "  valid@mail.com   ",
		Password: "password123",
		Role:     " aDmIn ",
	}

	sanitizedUser := user
	SanitizeUser(&sanitizedUser)

	suite.repository.On("CheckDuplicate", mock.Anything, "email", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	suite.repository.On("CheckDuplicate", mock.Anything, "username", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	suite.repository.On("CreateUser", mock.Anything, mock.AnythingOfType("User")).Return(nil)
	err := suite.usecase.CreateUser(context.TODO(), sanitizedUser)

	suite.NoError(err, "no error when given valid data")
	suite.repository.AssertCalled(suite.T(), "CreateUser", mock.Anything, sanitizedUser)
}

func (suite *userUsecaseSuite) TestCreateUser_UsernameValidation() {
	user := domain.User{
		Username: " ut     ",
		Email:    "  valid@mail.com   ",
		Password: "password123",
		Role:     " aDmIn ",
	}

	sanitizedUser := user
	SanitizeUser(&sanitizedUser)

	suite.repository.On("CheckDuplicate", mock.Anything, "email", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	suite.repository.On("CheckDuplicate", mock.Anything, "username", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := suite.usecase.CreateUser(context.TODO(), user)
	suite.Error(err, "error when given invalid username")
	suite.Equal(err.GetCode(), domain.ERR_BAD_REQUEST)
	suite.repository.AssertNotCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))

	user.Username = "valid_username"
	sanitizedUser.Username = user.Username
	SanitizeUser(&sanitizedUser)
	suite.repository.On("CreateUser", mock.Anything, sanitizedUser).Return(nil)
	err = suite.usecase.CreateUser(context.TODO(), user)

	suite.NoError(err, "no error when given a valid username")
	suite.repository.AssertCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))
}

func (suite *userUsecaseSuite) TestCreateUser_EmailValidation() {
	user := domain.User{
		Username: " valid_usernmae     ",
		Email:    "  invalid email   ",
		Password: "password123",
		Role:     " aDmIn ",
	}

	sanitizedUser := user
	SanitizeUser(&sanitizedUser)

	suite.repository.On("CheckDuplicate", mock.Anything, "email", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	suite.repository.On("CheckDuplicate", mock.Anything, "username", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := suite.usecase.CreateUser(context.TODO(), user)
	suite.Error(err, "error when given invalid email")
	suite.Equal(err.GetCode(), domain.ERR_BAD_REQUEST)
	suite.repository.AssertNotCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))

	user.Email = "valid_email@gmail.com"
	sanitizedUser.Email = user.Email
	SanitizeUser(&sanitizedUser)
	suite.repository.On("CreateUser", mock.Anything, sanitizedUser).Return(nil)
	err = suite.usecase.CreateUser(context.TODO(), user)

	suite.NoError(err, "no error when given a valid email")
	suite.repository.AssertCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))
}

func (suite *userUsecaseSuite) TestCreateUser_PasswordValidation() {
	user := domain.User{
		Username: " valid_usernmae     ",
		Email:    "  valid@email.com   ",
		Password: "invalid",
		Role:     " aDmIn ",
	}

	sanitizedUser := user
	SanitizeUser(&sanitizedUser)

	suite.repository.On("CheckDuplicate", mock.Anything, "email", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	suite.repository.On("CheckDuplicate", mock.Anything, "username", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := suite.usecase.CreateUser(context.TODO(), user)
	suite.Error(err, "error when given invalid password")
	suite.Equal(err.GetCode(), domain.ERR_BAD_REQUEST)
	suite.repository.AssertNotCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))

	user.Password = "valid_password123"
	sanitizedUser.Password = user.Password
	SanitizeUser(&sanitizedUser)
	suite.repository.On("CreateUser", mock.Anything, sanitizedUser).Return(nil)
	err = suite.usecase.CreateUser(context.TODO(), user)

	suite.NoError(err, "no error when given a valid password")
	suite.repository.AssertCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))
}

func (suite *userUsecaseSuite) TestCreateUser_RoleValidation() {
	user := domain.User{
		Username: " valid_usernmae     ",
		Email:    "  valid@email.com   ",
		Password: "password123",
		Role:     " invalid role ",
	}

	sanitizedUser := user
	SanitizeUser(&sanitizedUser)

	suite.repository.On("CheckDuplicate", mock.Anything, "email", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	suite.repository.On("CheckDuplicate", mock.Anything, "username", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := suite.usecase.CreateUser(context.TODO(), user)
	suite.Error(err, "error when given invalid role")
	suite.Equal(err.GetCode(), domain.ERR_BAD_REQUEST)
	suite.repository.AssertNotCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))

	user.Role = "user"
	sanitizedUser.Role = user.Role
	SanitizeUser(&sanitizedUser)
	suite.repository.On("CreateUser", mock.Anything, sanitizedUser).Return(nil)
	err = suite.usecase.CreateUser(context.TODO(), user)

	suite.NoError(err, "no error when given a valid role")
	suite.repository.AssertCalled(suite.T(), "CreateUser", mock.Anything, mock.AnythingOfType("User"))
}

func (suite *userUsecaseSuite) TestValidateAndGetToken_Positive() {
	storedUser := domain.User{
		Username: "valid_username",
		Email:    "valid@email.com",
		Password: "valid_password",
		Role:     "user",
	}

	suite.repository.On("GetByUsername", mock.Anything, storedUser.Username).Return(storedUser, nil)
	str, err := suite.usecase.ValidateAndGetToken(context.TODO(), storedUser)

	suite.NoError(err, "no error when the correct password and a valid role is provided")
	suite.Greater(len(str), 0)
}

func (suite *userUsecaseSuite) TestValidateAndGetToken_RoleValidation() {
	storedUser := domain.User{
		Username: "valid_username",
		Email:    "valid@email.com",
		Password: "valid_password",
		Role:     "invalid_role",
	}

	suite.repository.On("GetByUsername", mock.Anything, storedUser.Username).Return(storedUser, nil)
	str, err := suite.usecase.ValidateAndGetToken(context.TODO(), storedUser)

	suite.Error(err, "error when an invalid role is provided")
	suite.Equal(len(str), 0)
}

func (suite *userUsecaseSuite) TestValidateAndGetToken_PasswordValidation() {
	storedUser := domain.User{
		Username: "valid_username",
		Email:    "valid@email.com",
		Password: "valid_password",
		Role:     "invalid_role",
	}

	suite.repository.On("GetByUsername", mock.Anything, storedUser.Username).Return(storedUser, nil)
	storedUser.Password = "incorrect_password"
	str, err := suite.usecase.ValidateAndGetToken(context.TODO(), storedUser)

	suite.Error(err, "error when an invalid password is provided")
	suite.Equal(len(str), 0)
}

func TestUserUsecase(t *testing.T) {
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	suite.Run(t, new(userUsecaseSuite))
}
