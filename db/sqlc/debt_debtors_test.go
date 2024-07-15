package db

import (
	"context"
	"strconv"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUp() (Transaction, User) {
	debtor, _ := testQueries.CreateUser(
		context.Background(),
		CreateUserParams{
			Username: u.RandomUser(),
			Email:    u.RandomEmail(),
		})

	transaction, _ := testQueries.CreateTransaction(
		context.Background(),
		CreateTransactionParams{
			Amount:   "100.00",
			Currency: CurrencySGD,
			Title:    u.RandomString(5),
			PayerID:  1,
		},
	)

	return transaction, debtor
}

func TestCreateDebtDebtors(t *testing.T) {
	transaction, debtor := setUp()
	testDebtDebtorsParams := CreateDebtDebtorsParams{
		TransactionID: transaction.ID,
		DebtorID:      debtor.ID,
		Amount:        "25.00",
		Currency:      CurrencySGD,
	}
	debtDebtor, err := testQueries.CreateDebtDebtors(
		context.Background(),
		testDebtDebtorsParams,
	)

	require.NoError(t, err)
	require.NotEmpty(t, debtDebtor)
	assert.Equal(t, transaction.ID, debtDebtor.TransactionID)
	assert.Equal(t, debtor.ID, debtDebtor.DebtorID)
	expectedAmount, _ := strconv.Atoi(testDebtDebtorsParams.Amount)
	actualAmount, _ := strconv.Atoi(debtDebtor.Amount)
	assert.Equal(
		t,
		expectedAmount,
		actualAmount,
	)
}
