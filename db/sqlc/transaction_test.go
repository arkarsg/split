package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	test_input := CreateTransactionParams{
		Amount:   "100.00",
		Currency: CurrencySGD,
		Title:    u.RandomString(5),
		PayerID:  u.RandomInt(0, 1),
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
		PayerID:  u.RandomInt(0, 3),
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
	assert.Equal(t, txn.Title, txnRows[0].TransactionTitle)
}
