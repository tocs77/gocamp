package utility

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userId, username, role string) (string, error) {

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

func VerifyToken(token string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, ErrorHandler(errors.New("JWT_SECRET is not set"), "Internal server error")
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrorHandler(errors.New("invalid signing method"), "invalid signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, ErrorHandler(err, "invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrorHandler(errors.New("invalid token claims"), "invalid token claims")
	}
	fmt.Println("claims: ", claims)
	return claims, nil
}
