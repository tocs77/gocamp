package utility

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SingToken(userId, username, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", ErrorHandler(errors.New("JWT_SECRET is not set"), "Internal server error")
	}
	expiresIn := os.Getenv("JWT_EXPIRES_IN")
	if expiresIn == "" {
		return "", ErrorHandler(errors.New("JWT_EXPIRES_IN is not set"), "Internal server error")
	}
	expiresInDuration, err := time.ParseDuration(expiresIn)
	if err != nil {
		return "", ErrorHandler(err, "Internal server error")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  userId,
		"user": username,
		"role": role,
		"exp":  time.Now().Add(expiresInDuration).Unix(),
	})
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", ErrorHandler(err, "Internal server error")
	}
	return signedToken, nil
}
