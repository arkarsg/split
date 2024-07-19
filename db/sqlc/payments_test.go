package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePayment(t *testing.T) {
	dd := createRandomDebtDebtor()

	p, err := testQueries.CreatePayment(
		context.Background(),
		CreatePaymentParams{
			DebtID:   dd.DebtID,
			DebtorID: dd.DebtorID,
			Amount:   u.RandomAmount(),
			Currency: CurrencySGD,
		},
	)

	require.NoError(t, err)
	require.NotEmpty(t, p)
}

func TestGetPaymentsById(t *testing.T) {
	p := createRandomPayment()
	actual_p, err := testQueries.GetPaymentsById(
		context.Background(),
		p.ID,
	)

	require.NoError(t, err)
	assert.Equal(t, p, actual_p)
}

func TestDeletePayment(t *testing.T) {
	p := createRandomPayment()

	err := testQueries.DeletePayment(
		context.Background(),
		p.ID,
	)

	require.NoError(t, err)

	deletedP, err := testQueries.GetPaymentsById(
		context.Background(),
		p.ID,
	)

	require.Error(t, err)
	require.Empty(t, deletedP)
}
