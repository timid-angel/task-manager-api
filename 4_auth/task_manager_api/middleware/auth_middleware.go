package middlewares

import (
	"fmt"
	"os"
	"strings"
	services "task_manager_api/data"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// a struct that implements the error interface
type MWError struct {
	message string
}

func (error MWError) Error() string {
	return error.message
}

/*
This is the authorization middleware used for the endpoints. It accepts a set of
roles for which the endpoint is open.

WORKFLOW:
  - Obtains the JWT from the authorization header
  - Parses the JWT and verifies the signature
  - Checks the role of the user associated with the token
  - Calls `c.Next()` if the querying user has permission to access the endpoint
*/
func AuthMiddlewareWithRoles(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// obtain token from the request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"message": "Authorization header not found"})
			c.Abort()
			return
		}

		headerSegments := strings.Split(authHeader, " ")
		if len(headerSegments) != 2 || strings.ToLower(headerSegments[0]) != "bearer" {
			c.JSON(401, gin.H{"message": "Authorization header is invalid"})
			c.Abort()
			return
		}

		// parses token with the correct signing method
		token, err := jwt.Parse(headerSegments[1], func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, MWError{message: fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"])}
			}

			return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
		})

		// checks for errors and token validity
		if err != nil {
			c.JSON(401, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// get the username from the claims of the JWT
		username, ok := token.Claims.(jwt.MapClaims)["username"]
		if !ok {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// check the expiry date of the token
		expiresAt, ok := token.Claims.(jwt.MapClaims)["expiresAt"]
		if !ok {
			c.JSON(401, gin.H{"message": "Expiry date not found"})
			c.Abort()
			return
		}

		expiresAtTime, convErr := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiresAt))
		if convErr != nil {
			c.JSON(401, gin.H{"message": "Error during parsing expiry date: " + convErr.Error()})
			c.Abort()
			return
		}

		if expiresAtTime.Compare(time.Now()) == -1 {
			c.JSON(401, gin.H{"message": "Token expired"})
			c.Abort()
			return
		}

		// query the Db for the user associated with the claim and check perms
		storedUser, err := services.GetByUsername(fmt.Sprintf("%v", username))
		if err != nil {
			c.JSON(401, gin.H{"message": "User with the provided name does not exist"})
			c.Abort()
			return
		}

		valid := false
		for _, role := range validRoles {
			if storedUser.Role == role {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(403, gin.H{"message": fmt.Sprintf("%v roles are not allowed to access this endpoint", storedUser.Role)})
			c.Abort()
			return
		}

		c.Next()
	}
}
