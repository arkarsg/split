package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
	debt          Debt
	debtorPayment Payment
}

func (s *Store) SettleDebtPaymentsTx(ctx context.Context, args SettleDebtPaymentTxParams) (SettleDebtPaymentsTxResult, error) {
	var result SettleDebtPaymentsTxResult
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		result.debtorPayment, err = q.CreatePayment(ctx, CreatePaymentParams{
			DebtID:   args.DebtId,
			DebtorID: args.DebtorId,
			Amount:   args.Amount,
			Currency: args.Currency,
		})

		if err != nil {
			return err
		}
	})
	return result, nil
}
