package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type MWError struct {
	message string
}

func (error MWError) Error() string {
	return error.message
}

func AuthMiddleware() gin.HandlerFunc {
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

			return os.Getenv("JWT_SECRET_TOKEN"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
