package utils

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userId, username, role string) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", HandleError(errors.New("JWT_SECRET is not set"), "Internal server error")
	}
	expiresIn := os.Getenv("JWT_EXPIRES_IN")
	if expiresIn == "" {
		return "", HandleError(errors.New("JWT_EXPIRES_IN is not set"), "Internal server error")
	}
	expiresInDuration, err := time.ParseDuration(expiresIn)
	if err != nil {
		return "", HandleError(err, "Internal server error")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  userId,
		"user": username,
		"role": role,
		"exp":  time.Now().Add(expiresInDuration).Unix(),
	})
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", HandleError(err, "Internal server error")
	}
	return signedToken, nil
}

func VerifyToken(token string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, HandleError(errors.New("JWT_SECRET is not set"), "Internal server error")
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, HandleError(errors.New("invalid signing method"), "invalid signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, HandleError(err, "invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, HandleError(errors.New("invalid token claims"), "invalid token claims")
	}
	fmt.Println("claims: ", claims)
	return claims, nil
}

type JWTStore struct {
	tokens map[string]time.Time
	mu     sync.Mutex
}

func (js *JWTStore) AddToken(token string) {
	js.mu.Lock()
	defer js.mu.Unlock()
	claims, err := VerifyToken(token)
	if err != nil {
		return
	}
	expiresAt := claims["exp"].(float64)
	js.tokens[token] = time.Unix(int64(expiresAt), 0)
}

func (js *JWTStore) HasToken(token string) bool {
	js.mu.Lock()
	defer js.mu.Unlock()
	_, ok := js.tokens[token]
	return ok
}

func (js *JWTStore) DeleteToken(token string) {
	js.mu.Lock()
	defer js.mu.Unlock()
	delete(js.tokens, token)
}

func (js *JWTStore) CleanUpExpiredTokens() {

	for {
		time.Sleep(1 * time.Minute)
		js.mu.Lock()
		for token, expiresAt := range js.tokens {
			if time.Now().After(expiresAt) {
				delete(js.tokens, token)
			}
		}
		js.mu.Unlock()
	}
}

func NewJWTStore() *JWTStore {
	return &JWTStore{
		tokens: make(map[string]time.Time),
		mu:     sync.Mutex{},
	}
}

var JWTStorage = NewJWTStore()

func init() {
	fmt.Println("Initializing JWTStorage")
	go JWTStorage.CleanUpExpiredTokens()
}
