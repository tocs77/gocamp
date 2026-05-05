package utils

import (
	"golang.org/x/crypto/argon2"

	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrInvalidPasswordResetToken is returned when the reset code from the client is not a 32-byte value encoded as hex (64 hex chars).
var ErrInvalidPasswordResetToken = errors.New("invalid password reset token")

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
		return false, nil
	}
	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}

// PasswordResetTokenHashFromRaw returns the hex-encoded SHA-256 digest of tokenBytes, matching how ForgotPassword stores password_reset_token.
func PasswordResetTokenHashFromRaw(tokenBytes []byte) string {
	sum := sha256.Sum256(tokenBytes)
	return hex.EncodeToString(sum[:])
}

// PasswordResetTokenHashFromPlain decodes resetCode as hex (64 chars = 32 bytes), hashes those bytes with SHA-256, and returns the hex digest for DB lookup.
func PasswordResetTokenHashFromPlain(resetCode string) (string, error) {
	tokenBytes, err := hex.DecodeString(resetCode)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidPasswordResetToken, err)
	}
	if len(tokenBytes) != 32 {
		return "", ErrInvalidPasswordResetToken
	}
	return PasswordResetTokenHashFromRaw(tokenBytes), nil
}
