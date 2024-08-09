package tests

import (
	"strings"
	infrastructure "task_manager_api/Infrastructure"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/suite"
)

type jwtServiceSuite struct {
	suite.Suite
}

type passwordServiceSuite struct {
	suite.Suite
}

func (suite *jwtServiceSuite) TestSignWithJWTPayload_Positive() {
	username := "suser"
	role := "admin"
	secret := "secret"

	token, err := infrastructure.SignJWTWithPayload(username, role, time.Minute, secret)

	suite.NoError(err, "no error when given valid inputs")
	suite.True(strings.HasPrefix(token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."))
}

func (suite *jwtServiceSuite) TestSignWithJWTPayload_Negative() {
	username := "suser"
	role := "admin"
	secret := ""

	token, err := infrastructure.SignJWTWithPayload(username, role, time.Minute, secret)

	suite.Error(err, "error when given valid inputs")
	suite.Equal(token, "")
}

func (suite *jwtServiceSuite) TestValidateAndParseToken_Positive() {
	username := "suser"
	role := "admin"
	secret := "secret"

	token, _ := infrastructure.SignJWTWithPayload(username, role, time.Minute, secret)
	jwtToken, err := infrastructure.ValidateAndParseToken(token, secret)

	suite.NoError(err, "no error when given the same secret key")
	suite.Equal(username, jwtToken.Claims.(jwt.MapClaims)["username"])
	suite.Equal(role, jwtToken.Claims.(jwt.MapClaims)["role"])
}

func (suite *jwtServiceSuite) TestValidateAndParseToken_Negative() {
	username := "suser"
	role := "admin"
	secret := "secret"

	token, _ := infrastructure.SignJWTWithPayload(username, role, time.Minute, secret)
	_, err := infrastructure.ValidateAndParseToken(token, "modified secret")

	suite.Error(err, "error when given the same secret key")
}

func (suite *passwordServiceSuite) TestHashPassword() {
	password := "plain_text password"

	hashedPwd, err := infrastructure.HashPassword(password)

	suite.NoError(err, "no errors when given password")
	suite.NotEqual(hashedPwd, password)
}

func (suite *passwordServiceSuite) TestValidatePassword_Positive() {
	password := "plain_text password"

	hashedPwd, _ := infrastructure.HashPassword(password)
	err := infrastructure.ValidatePassword(hashedPwd, password)

	suite.NoError(err, "no errors when given correct password")
}

func (suite *passwordServiceSuite) TestValidatePassword_Negative() {
	password := "plain_text password"

	hashedPwd, _ := infrastructure.HashPassword(password)
	err := infrastructure.ValidatePassword(hashedPwd, "wrong plain_text password")

	suite.Error(err, "error when given incorrect password")
}

func TestInfrastructureSuite(t *testing.T) {
	suite.Run(t, new(jwtServiceSuite))
	suite.Run(t, new(passwordServiceSuite))
}
