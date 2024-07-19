package db

import (
	"context"

	u "github.com/arkarsg/splitapp/utils"
)

func createRandomUser(test_input CreateUserParams) (User, error) {
	user, err := testQueries.CreateUser(context.Background(), test_input)
	return user, err
}

func createRandomTransaction() Transaction {
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

func createRandomDebt() Debt {
	txn := createRandomTransaction()
	debt, _ := testQueries.CreateDebt(
		context.Background(),
		txn.ID,
	)
	return debt
}

func createRandomDebtDebtor() DebtDebtor {
	debtor, _ := createRandomUser(CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	})
	debt := createRandomDebt()

	createDebtDebtorsParams := CreateDebtDebtorsParams{
		DebtID:   debt.ID,
		DebtorID: debtor.ID,
		Amount:   u.RandomAmount(),
		Currency: CurrencySGD,
	}
	dd, _ := testQueries.CreateDebtDebtors(
		context.Background(),
		createDebtDebtorsParams,
	)

	return dd
}

func createRandomPayment() Payment {
	dd := createRandomDebtDebtor()
	p, _ := testQueries.CreatePayment(
		context.Background(),
		CreatePaymentParams{
			DebtID:   dd.DebtID,
			DebtorID: dd.DebtorID,
			Amount:   u.RandomAmount(),
			Currency: CurrencySGD,
		},
	)

	return p
}
