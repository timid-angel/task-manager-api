package middlewares

import (
	"fmt"
	"os"
	"strings"
	services "task_manager_api/data"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type MWError struct {
	message string
}

func (error MWError) Error() string {
	return error.message
}

func AuthMiddlewareWithRoles(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		token, err := jwt.Parse(headerSegments[1], func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, MWError{message: fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"])}
			}

			return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
		})

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

		token.Claims.Valid()
		username, ok := token.Claims.(jwt.MapClaims)["username"]
		if !ok {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

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
