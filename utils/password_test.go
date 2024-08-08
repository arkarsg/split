package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(10)
	hashed, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
}

func TestCheckHashPassword(t *testing.T) {
	password := RandomString(10)
	hashed, err := HashPassword(password)
	assert.NoError(t, err)
	assert.True(t, CheckPasswordHash(password, hashed))
}

func TestWrongPassword(t *testing.T) {
	password := RandomString(10)
	hashed, err := HashPassword(password)
	assert.NoError(t, err)
	assert.False(t, CheckPasswordHash(RandomString(10), hashed))
}

func TestHashPasswordTwiceCreatesDifferentHash(t *testing.T) {
	password := RandomString(10)
	hashOne, _ := HashPassword(password)
	hashTwo, _ := HashPassword(password)
	assert.NotEqual(t, hashOne, hashTwo)
}
