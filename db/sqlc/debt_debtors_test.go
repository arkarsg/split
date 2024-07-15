package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
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
