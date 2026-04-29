package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"strings"

	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {	
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}

func ComparePassword(hashedPassword, password string) (bool, error) {
	parts := strings.Split(hashedPassword, ".")
	if len(parts) != 2 {
		return false, HandleError(errors.New("invalid hashed password"), "invalid hashed password")
	}
	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, HandleError(err, "error decoding salt")
	}
	hash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, HandleError(err, "error decoding hash")
	}
	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	if len(hash) != len(computedHash) {
		return false, HandleError(errors.New("invalid password"), "invalid password")
	}
	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}
	return false, HandleError(errors.New("invalid password"), "invalid password")
}