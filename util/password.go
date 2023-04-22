package util

import (
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %s", err)
	}
	return string(hashpassword), nil
}

// CheckPassword checks if the provided password is correct or not
// returns nil if the password is correct and an error if it is not
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateRandomString generates a random string of length n
// with characters from a-z, A-Z, and 0-9
func RandomString(n int) string {
	var sb strings.Builder
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	k := len(charset)

	for i := 0; i < n; i++ {
		c := charset[rand.Intn(k - 1)]
		sb.WriteByte(c)
	}

	return sb.String()
}
// GenerateRandomString generates a random string of length n
// with special characters
func GenerateRandomStringSpecial(n int) string {
	var sb strings.Builder
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	k := len(charset)

	for i := 0; i < n; i++ {
		c := charset[rand.Intn(k - 1)]
		sb.WriteByte(c)
	}

	return sb.String()
}
// GenerateRandomString generates a random string of length n
// with characters from a-z and A-Z
func GenerateRandomStringABC(n int) string {
	var sb strings.Builder
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	k := len(charset)

	for i := 0; i < n; i++ {
		c := charset[rand.Intn(k - 1)]
		sb.WriteByte(c)
	}

	return sb.String()
}