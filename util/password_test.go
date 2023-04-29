package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	password := RandomStringNum(10)

	hashPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NoError(t, CheckPassword(password, hashPassword))

	wrongPassword := RandomStringNum(10)
	err = CheckPassword(wrongPassword, hashPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}