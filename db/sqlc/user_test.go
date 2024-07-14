package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(test_input CreateUserParams) (User, error) {
	user, err := testQueries.CreateUser(context.Background(), test_input)
	return user, err
}

func TestCreateUser(t *testing.T) {
	test_input := CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}
	user, err := createRandomAccount(test_input)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, test_input.Username, user.Username)
	require.Equal(t, test_input.Email, user.Email)
}

func TestGetUserById(t *testing.T) {
	test_input := CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}

	expectedUser, _ := createRandomAccount(test_input)
	actualUser, err := testQueries.GetUserById(
		context.Background(),
		expectedUser.ID,
	)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)
}
