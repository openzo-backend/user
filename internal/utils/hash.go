package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashNumberWithSecret(number int, secret string) string {
	// Convert the number to a byte slice
	numberBytes := []byte(fmt.Sprintf("%d", number))

	// Create a new HMAC by defining the hash type and the key (secret)
	hmac := hmac.New(sha256.New, []byte(secret))

	// Write the number bytes to the HMAC object
	hmac.Write(numberBytes)

	// Compute the HMAC hash
	hash := hmac.Sum(nil)

	// Encode the hash to a hexadecimal string and return
	return hex.EncodeToString(hash)
}
