package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func main() {
	password := "123456"
	hashedPassword := sha256.Sum256([]byte(password))
	fmt.Println(hashedPassword)
	fmt.Printf("hashedPassword: %x\n", hashedPassword)

	salt, err := generateSalt()
	if err != nil {
		panic(err)
	}
	fmt.Printf("salt: %x\n", salt)
	saltedHash, err := hashPassword(password, salt)
	if err != nil {
		panic(err)
	}
	fmt.Println("saltedHash: ", saltedHash)
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string, salt []byte) (string, error) {
	saltedPassword := append([]byte(password), salt...)
	hashedPassword := sha256.Sum256(saltedPassword)

	return fmt.Sprintf("%x", hashedPassword), nil
}
