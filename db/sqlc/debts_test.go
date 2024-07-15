package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpRandomTransaction() Transaction {
	var user User
	var err error
	user, err = testQueries.GetUserById(
		context.Background(),
		1,
	)

	if err != nil {
		user, _ = testQueries.CreateUser(
			context.Background(),
			CreateUserParams{
				Username: u.RandomUser(),
				Email:    u.RandomEmail(),
			})
	}

	txnParams := CreateTransactionParams{
		Amount:   u.RandomAmount(),
		Currency: CurrencySGD,
		Title:    u.RandomString(10),
		PayerID:  user.ID,
	}
	txn, _ := testQueries.CreateTransaction(
		context.Background(),
		txnParams,
	)
	return txn
}

func TestCreateDebt(t *testing.T) {
	txn := setUpRandomTransaction()
	debt, err := testQueries.CreateDebt(
		context.Background(),
		txn.ID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, debt)
}

func createRandomDebt() Debt {
	txn := setUpRandomTransaction()
	debt, _ := testQueries.CreateDebt(
		context.Background(),
		txn.ID,
	)
	return debt
}

func TestGetDebtById(t *testing.T) {
	expectedDebt := createRandomDebt()
	actualDebt, err := testQueries.GetDebtById(
		context.Background(),
		expectedDebt.ID,
	)
	require.NoError(t, err)
	assert.Equal(t, expectedDebt, actualDebt)
}

func TestUpdateDebt(t *testing.T) {
	debtToTest := createRandomDebt()
	newAmount := u.RandomAmount()
	updateDebtParams := UpdateDebtParams{
		ID:               debtToTest.ID,
		NewSettledAmount: newAmount,
	}
	newDebt, err := testQueries.UpdateDebt(
		context.Background(),
		updateDebtParams,
	)
	require.NoError(t, err)
	require.NotEmpty(t, newDebt)
	assert.NotEqual(t, debtToTest.SettledAmount, newDebt.SettledAmount)
}
