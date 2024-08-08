package infrastructure

import (
	"fmt"
	domain "task_manager_api/Domain"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

/*
Creates and signs a JWT with the username, role and tokenLifeSpan as the
payloads. Returns the signed token if there aren't any errors.
*/
func SignJWTWithPayload(username string, role string, tokenLifeSpan time.Duration, secret string) (string, domain.CodedError) {
	jwtSecret := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  username,
		"role":      role,
		"expiresAt": time.Now().Add(time.Hour * 2),
	})
	jwtToken, signingErr := token.SignedString(jwtSecret)
	if signingErr != nil {
		return "", domain.UserError{Message: "internal server error: " + signingErr.Error(), Code: domain.ERR_INTERNAL_SERVER}
	}

	return jwtToken, nil
}

/*
Parses the JWT token with the HMAC signing method and returns a pointer
to a jwt.Token struct if the token is valid and not tampered with.
*/
func ValidateAndParseToken(rawToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(viper.GetString("SECRET_TOKEN")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error: " + err.Error())
	}

	if !token.Valid {
		return nil, fmt.Errorf("error: Invalid token,  Potentially malformed")
	}

	return token, nil
}
