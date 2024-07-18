package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDebtDebtors(t *testing.T) {
	debtor, _ := CreateRandomUser(CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	})
	debt := CreateRandomDebt()

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
	dd := CreateRandomDebtDebtor()
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
	dd := CreateRandomDebtDebtor()
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
