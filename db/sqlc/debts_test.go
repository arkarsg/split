package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDebt(t *testing.T) {
	txn := createRandomTransaction()
	debt, err := testQueries.CreateDebt(
		context.Background(),
		txn.ID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, debt)
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

func TestGetDebtByTransactionID(t *testing.T) {
	expectedDebt := createRandomDebt()
	actualDebt, err := testQueries.GetDebtByTransactionId(
		context.Background(),
		expectedDebt.TransactionID,
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
