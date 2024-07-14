package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	test_input := CreateUserParams{
		Username: u.RandomUser(),
		Email: u.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), test_input)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, test_input.Username, user.Username)
	require.Equal(t, test_input.Email, user.Email)
}
