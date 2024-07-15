package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	var testUser User
	user, err := testQueries.GetUserById(
		context.Background(),
		1,
	)
	testUser = user
	if err != nil {
		dummyUserParams := CreateUserParams{
			Username: u.RandomUser(),
			Email:    u.RandomEmail(),
		}
		user, _ := testQueries.CreateUser(
			context.Background(),
			dummyUserParams,
		)
		testUser = user
	}

	test_input := CreateTransactionParams{
		Amount:   "100.00",
		Currency: CurrencySGD,
		Title:    u.RandomString(5),
		PayerID:  testUser.ID,
	}
	transaction, err := testQueries.CreateTransaction(
		context.Background(),
		test_input,
	)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)
}

func createRandomTransaction() Transaction {
	createTransactionParmas := CreateTransactionParams{
		Amount:   "100.00",
		Currency: CurrencyUSD,
		Title:    u.RandomString(10),
		PayerID:  1,
	}

	transaction, _ := testQueries.CreateTransaction(
		context.Background(),
		createTransactionParmas,
	)
	return transaction
}

func TestGetTransactionByPayer(t *testing.T) {
	txn := createRandomTransaction()
	getQueryParams := GetTransactionsByPayerParams{
		PayerID: txn.PayerID,
		Limit:   1,
		Offset:  0,
	}

	txnRows, err := testQueries.GetTransactionsByPayer(
		context.Background(),
		getQueryParams,
	)
	payer, _ := testQueries.GetUserById(
		context.Background(),
		txn.PayerID,
	)

	require.NoError(t, err)
	require.NotEmpty(t, txnRows)
	require.Len(t, txnRows, 1)
	assert.Equal(t, payer.Username, txnRows[0].PayerUsername)
}
