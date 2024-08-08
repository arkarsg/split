package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	account := createRandomAccount()
	test_input := CreateUserParams{
		Username: account.Username,
		Email:    account.Email,
	}
	user, err := testQueries.CreateUser(context.Background(), test_input)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	assert.Equal(t, test_input.Username, user.Username)
	assert.Equal(t, test_input.Email, user.Email)
}

func TestGetUserById(t *testing.T) {
	account := createRandomAccount()
	args := CreateUserParams{
		Username: account.Username,
		Email:    account.Email,
	}
	expectedUser, _ := testQueries.CreateUser(context.Background(), args)
	actualUser, err := testQueries.GetUserById(
		context.Background(),
		expectedUser.ID,
	)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)
}

func TestGetUserByUsername(t *testing.T) {
	expectedUser := createRandomUser()
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
	testUser := createRandomUser()
	err := testQueries.DeleteUser(context.Background(), testUser.ID)
	require.NoError(t, err)

	user, err := testQueries.GetUserById(context.Background(), testUser.ID)
	require.Error(t, err)
	assert.Empty(t, user)
}
