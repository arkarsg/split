package db

import (
	"context"
	"testing"

	"github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	test_input := CreateAccountParams{
		Username:       utils.RandomString(10),
		HashedPassword: utils.RandomString(10),
		FullName:       utils.RandomUser(),
		Email:          utils.RandomEmail(),
	}

	account, err := testQueries.CreateAccount(context.Background(), test_input)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	assert.Equal(t, test_input.Username, account.Username)
	assert.Equal(t, test_input.HashedPassword, account.HashedPassword)
}
