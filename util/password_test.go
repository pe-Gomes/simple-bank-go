package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPass, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)

	err = ComparePassword(password, hashedPass)
	require.NoError(t, err)

	wrongPass := RandomString(6)
	err = ComparePassword(wrongPass, hashedPass)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
