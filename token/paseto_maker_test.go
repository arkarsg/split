package token

import (
	"testing"
	"time"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
)

func TestPasetoMakerCanCreateToken(t *testing.T) {
	_, token, err := createRandomToken(t)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestPasetoMakerCreatesValidToken(t *testing.T) {
	m, token, err := createRandomToken(t)
	payload, err := m.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)
}

func TestExpiredPasetoToken(t *testing.T) {
	m := createPasetoMaker(t)
	token, err := m.CreateToken(u.RandomUser(), -time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	payload, err := m.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrExpiredToken.Error())
	assert.Nil(t, payload)
}

func createPasetoMaker(t *testing.T) TokenMaker {
	m, err := NewPasetoMaker(u.RandomString(32))
	assert.NoError(t, err)
	return m
}

func createRandomToken(t *testing.T) (TokenMaker, string, error) {
	m := createPasetoMaker(t)
	username := u.RandomUser()
	duration := time.Minute
	token, err := m.CreateToken(username, duration)
	return m, token, err
}
