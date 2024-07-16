package db

import (
	"context"
	"testing"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettleDebtPaymentsTx(t *testing.T) {
	store := NewStore(testDb)

	debt1 := CreateRandomDebt()
	debtor1, _ := CreateRandomUser(CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	})

	// run $n$ concurrent payments
	n := 5
	amount := "10.00"

	errs := make(chan error)
	results := make(chan SettleDebtPaymentsTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.SettleDebtPaymentsTx(
				context.Background(),
				SettleDebtPaymentTxParams{
					DebtId:   debt1.ID,
					DebtorId: debtor1.ID,
					Amount:   amount,
					Currency: CurrencySGD,
				})

			errs <- err
			results <- result
		}()
	}

	totalSettledAmount := u.ZeroMoneyAmount()

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		// Check Payment
		createdPayment := res.DebtorPayment
		actualPayment, err := store.GetPaymentsById(
			context.Background(),
			createdPayment.ID,
		)
		require.NoError(t, err)
		assert.True(
			t,
			u.StringToMoney(amount).Amount.Equal(u.StringToMoney(actualPayment.Amount).Amount),
		)

		totalSettledAmount = u.AddMoney(
			totalSettledAmount,
			u.StringToMoney(createdPayment.Amount),
		)
	}

	updatedDebt, err := store.GetDebtById(context.Background(), debt1.ID)
	require.NoError(t, err)
	assert.Equal(
		t,
		totalSettledAmount,
		u.StringToMoney(updatedDebt.SettledAmount),
	)
}
