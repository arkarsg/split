package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	test_input := CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}
	user, err := CreateRandomUser(test_input)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	assert.Equal(t, test_input.Username, user.Username)
	assert.Equal(t, test_input.Email, user.Email)
}

func TestGetUserById(t *testing.T) {
	test_input := CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}

	expectedUser, _ := CreateRandomUser(test_input)
	actualUser, err := testQueries.GetUserById(
		context.Background(),
		expectedUser.ID,
	)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)
}

func TestGetUserByUsername(t *testing.T) {
	test_input := CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}

	expectedUser, _ := CreateRandomUser(test_input)
	actualUser, err := testQueries.GetUserByUsername(
		context.Background(),
		expectedUser.Username,
	)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)
}

func TestListUsers(t *testing.T) {
	test_input := ListUsersParams{
		Limit:  2,
		Offset: 0,
	}
	users, err := testQueries.ListUsers(context.Background(), test_input)
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestDeleteUser(t *testing.T) {
	testUserParams := CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}
	testUser, _ := CreateRandomUser(testUserParams)
	err := testQueries.DeleteUser(context.Background(), testUser.ID)
	require.NoError(t, err)

	user, err := testQueries.GetUserById(context.Background(), testUser.ID)
	require.Error(t, err)
	assert.Empty(t, user)
}
