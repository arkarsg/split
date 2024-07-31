package db

import (
	"context"
	"database/sql"
	"fmt"

	u "github.com/arkarsg/splitapp/utils"
)

type Store interface {
	Querier
	SettleDebtPaymentsTx(ctx context.Context, args SettleDebtPaymentTxParams) (SettleDebtPaymentsTxResult, error)
	CreateTransactionDebtTx(ctx context.Context, args CreateTransactionDebtTxParams) (CreateTransactionDebtTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("Txn Error: %v, Rollback Error: %v", err, rbError)
		}
		return err
	}

	return tx.Commit()
}

type SettleDebtPaymentTxParams struct {
	DebtId   int64
	DebtorId int64
	Amount   string
	Currency Currency
}

type SettleDebtPaymentsTxResult struct {
	Debt          Debt
	DebtorPayment Payment
}

// SettleDebtPaymentsTx creates a payment of SettleDebtPaymentTxParams.Amount,
// and increments Debt.SettledAmount by SettleDebtPaymentTxParams.Amount
func (s *SQLStore) SettleDebtPaymentsTx(ctx context.Context, args SettleDebtPaymentTxParams) (SettleDebtPaymentsTxResult, error) {
	var result SettleDebtPaymentsTxResult
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		result.DebtorPayment, err = q.CreatePayment(ctx, CreatePaymentParams{
			DebtID:   args.DebtId,
			DebtorID: args.DebtorId,
			Amount:   args.Amount,
			Currency: args.Currency,
		})

		if err != nil {
			return err
		}

		originalDebt, err := q.GetDebtByIdForUpdate(ctx, args.DebtId)

		if err != nil {
			return err
		}
		originalAmount := u.StringToMoney(originalDebt.SettledAmount)
		settledAmount := u.StringToMoney(args.Amount)
		newAmount := u.AddMoney(*originalAmount, *settledAmount)

		result.Debt, err = q.UpdateDebt(ctx, UpdateDebtParams{
			ID:               originalDebt.ID,
			NewSettledAmount: newAmount.MoneyToString(),
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

type CreateTransactionDebtTxParams struct {
	Amount   string
	Currency Currency
	Title    string
	PayerID  int64
}

type CreateTransactionDebtTxResult struct {
	Transaction Transaction
	Debt        Debt
}

func (s *SQLStore) CreateTransactionDebtTx(ctx context.Context, args CreateTransactionDebtTxParams) (CreateTransactionDebtTxResult, error) {
	var result CreateTransactionDebtTxResult
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			Amount:   args.Amount,
			Currency: args.Currency,
			Title:    args.Title,
			PayerID:  args.PayerID,
		})
		if err != nil {
			return err
		}

		result.Debt, err = q.CreateDebt(ctx, result.Transaction.ID)
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
