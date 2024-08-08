package db

import (
	"context"
	"database/sql"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDebtDebtors(t *testing.T) {
	debtor := createRandomUser()
	debt := createRandomDebt()

	createDebtDebtorsParams := CreateDebtDebtorsParams{
		DebtID:   debt.ID,
		DebtorID: debtor.ID,
		Amount:   u.RandomAmount(),
		Currency: CurrencySGD,
	}
	dd, err := testQueries.CreateDebtDebtors(
		context.Background(),
		createDebtDebtorsParams,
	)

	require.NoError(t, err)
	require.NotEmpty(t, dd)
}

func TestGetDebtDebtorsByDebtAndDebtor(t *testing.T) {
	dd := createRandomDebtDebtor()
	dd_actual, err := testQueries.GetDebtDebtorsByDebtAndDebtor(
		context.Background(),
		GetDebtDebtorsByDebtAndDebtorParams{
			DebtID:   dd.DebtID,
			DebtorID: dd.DebtorID,
		},
	)

	require.NoError(t, err)
	assert.Equal(t, dd.Amount, dd_actual.Amount)
}

func TestDeleteDebtDebtor(t *testing.T) {
	dd := createRandomDebtDebtor()
	err := testQueries.DeleteDebtDebtor(
		context.Background(),
		DeleteDebtDebtorParams{
			DebtID:   dd.DebtID,
			DebtorID: dd.DebtorID,
		},
	)

	require.NoError(t, err)

	debtDebtors, err := testQueries.GetDebtDebtorsByDebtAndDebtor(
		context.Background(),
		GetDebtDebtorsByDebtAndDebtorParams{
			DebtID:   dd.DebtID,
			DebtorID: dd.DebtorID,
		},
	)
	require.Error(t, err)
	require.Empty(t, debtDebtors)
}

func TestGetDebtsOfDebtorId(t *testing.T) {
	debtor := createRandomUser()

	// associate 10 debts with debtor
	n := 10
	for i := 0; i < n; i++ {
		debt := createRandomDebt()
		testQueries.CreateDebtDebtors(
			context.Background(),
			CreateDebtDebtorsParams{
				DebtID:   debt.ID,
				DebtorID: debtor.ID,
				Amount:   u.RandomAmount(),
				Currency: CurrencySGD,
			},
		)
	}

	rows, err := testQueries.GetDebtsOfDebtorId(
		context.Background(),
		debtor.ID,
	)

	require.NoError(t, err)
	require.Len(t, rows, n)
}

func TestUpdateAllFieldsDebtDebtor(t *testing.T) {
	dd := createRandomDebtDebtor()
	newAmount := u.RandomAmount()
	updateDebtDebtorParams := UpdateDebtDebtorParams{
		Amount: sql.NullString{
			String: newAmount,
			Valid:  true,
		},
		Currency: NullCurrency{
			Currency: CurrencyUSD,
			Valid:    true,
		},
		DebtId:   dd.DebtID,
		DebtorId: dd.DebtorID,
	}

	newDebtDebtor, err := testQueries.UpdateDebtDebtor(
		context.Background(),
		updateDebtDebtorParams,
	)
	require.NoError(t, err)
	assert.Equal(t, newAmount, newDebtDebtor.Amount)
	assert.Equal(t, CurrencyUSD, newDebtDebtor.Currency)
}
