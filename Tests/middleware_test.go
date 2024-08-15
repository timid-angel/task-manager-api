package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	infrastructure "task_manager_api/Infrastructure"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/suite"
)

type authMiddlewareSuite struct {
	suite.Suite
}

func (suite *authMiddlewareSuite) TestAuthMiddleware_Positive() {
	validateToken := func(rawToken string, secret string) (*jwt.Token, error) {
		mockClaim := jwt.MapClaims{
			"expiresAt": time.Now().Add(time.Hour).Round(0).Format(time.RFC3339Nano),
			"role":      "user",
		}

		if secret == "secret" {
			return &jwt.Token{Raw: rawToken, Claims: mockClaim}, nil
		}
		return nil, fmt.Errorf("Error")
	}

	router := gin.Default()
	router.GET("/", infrastructure.AuthMiddlewareWithRoles([]string{"user", "admin"}, "secret", validateToken), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})

	testingServer := httptest.NewServer(router)
	defer testingServer.Close()

	rawToken := "lksadfjlsadfj.elrqlkrjqwaklsdfm.lsiasdflkxzjcvlk"
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, testingServer.URL+"/", nil)
	suite.NoError(err, "no error during request creation")
	request.Header.Add("Authorization", "bearer "+rawToken)

	response, err := client.Do(request)
	suite.NoError(err, "no error during request")
	if response != nil {
		defer response.Body.Close()
	}

	suite.Equal(http.StatusOK, response.StatusCode)
}

func (suite *authMiddlewareSuite) TestAuthMiddleware_Negative_Expired() {
	validateToken := func(rawToken string, secret string) (*jwt.Token, error) {
		mockClaim := jwt.MapClaims{
			"expiresAt": time.Now().Add(time.Hour * -2).Round(0).Format(time.RFC3339Nano),
			"role":      "user",
		}

		if secret == "secret" {
			return &jwt.Token{Raw: rawToken, Claims: mockClaim}, nil
		}
		return nil, fmt.Errorf("Error")
	}

	router := gin.Default()
	router.GET("/", infrastructure.AuthMiddlewareWithRoles([]string{"user", "admin"}, "secret", validateToken), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})

	testingServer := httptest.NewServer(router)
	defer testingServer.Close()

	rawToken := "lksadfjlsadfj.elrqlkrjqwaklsdfm.lsiasdflkxzjcvlk"
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, testingServer.URL+"/", nil)
	suite.NoError(err, "no error during request creation")
	request.Header.Add("Authorization", "bearer "+rawToken)

	response, err := client.Do(request)
	suite.NoError(err, "no error during request")
	if response != nil {
		defer response.Body.Close()
	}

	suite.Equal(http.StatusUnauthorized, response.StatusCode)
}

func (suite *authMiddlewareSuite) TestAuthMiddleware_Negative_Forbidden() {
	validateToken := func(rawToken string, secret string) (*jwt.Token, error) {
		mockClaim := jwt.MapClaims{
			"expiresAt": time.Now().Add(time.Hour).Round(0).Format(time.RFC3339Nano),
			"role":      "user",
		}

		if secret == "secret" {
			return &jwt.Token{Raw: rawToken, Claims: mockClaim}, nil
		}
		return nil, fmt.Errorf("Error")
	}

	router := gin.Default()
	router.GET("/", infrastructure.AuthMiddlewareWithRoles([]string{"admin"}, "secret", validateToken), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})

	testingServer := httptest.NewServer(router)
	defer testingServer.Close()

	rawToken := "lksadfjlsadfj.elrqlkrjqwaklsdfm.lsiasdflkxzjcvlk"
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, testingServer.URL+"/", nil)
	suite.NoError(err, "no error during request creation")
	request.Header.Add("Authorization", "bearer "+rawToken)

	response, err := client.Do(request)
	suite.NoError(err, "no error during request")
	if response != nil {
		defer response.Body.Close()
	}

	suite.Equal(http.StatusForbidden, response.StatusCode)
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(authMiddlewareSuite))
}
