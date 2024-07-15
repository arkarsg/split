package db

import (
	"context"
	"database/sql"
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

func TestGetTransactionById(t *testing.T) {
	expectedTxn := CreateRandomTransaction()
	actualTxn, err := testQueries.GetTransactionById(
		context.Background(),
		expectedTxn.ID,
	)
	require.NoError(t, err)
	assert.Equal(t, expectedTxn, actualTxn)
}

func TestGetTransactionByPayer(t *testing.T) {
	txn := CreateRandomTransaction()
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

func TestUpdateTransaction(t *testing.T) {
	txn := CreateRandomTransaction()
	updateTxnParmas := UpdateTransactionParams{
		Title: sql.NullString{
			String: u.RandomString(10),
			Valid:  true,
		},
		ID: txn.ID,
	}
	updatedTxn, err := testQueries.UpdateTransaction(
		context.Background(),
		updateTxnParmas,
	)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTxn)
	assert.NotEqual(t, txn.Title, updatedTxn.Title)
}

func TestDeleteTransaction(t *testing.T) {
	txn := CreateRandomTransaction()
	var err error
	err = testQueries.DeleteTransaction(
		context.Background(),
		txn.ID,
	)
	require.NoError(t, err)

	deletedTxn, err := testQueries.GetTransactionById(
		context.Background(),
		txn.ID,
	)
	require.Error(t, err)
	require.Empty(t, deletedTxn)
}
