package db

import (
	"context"

	u "github.com/arkarsg/splitapp/utils"
)

func CreateRandomAccount(test_input CreateUserParams) (User, error) {
	user, err := testQueries.CreateUser(context.Background(), test_input)
	return user, err
}

func CreateRandomTransaction() Transaction {
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

func CreateRandomDebt() Debt {
	txn := CreateRandomTransaction()
	debt, _ := testQueries.CreateDebt(
		context.Background(),
		txn.ID,
	)
	return debt
}
