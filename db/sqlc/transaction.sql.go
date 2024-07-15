// CoPayde generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: transaction.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    amount,
    currency,
    title,
    payer_id
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, amount, currency, title, created_at, payer_id
`

type CreateTransactionParams struct {
	Amount   string
	Currency Currency
	Title    string
	PayerID  int64
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.Amount,
		arg.Currency,
		arg.Title,
		arg.PayerID,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Currency,
		&i.Title,
		&i.CreatedAt,
		&i.PayerID,
	)
	return i, err
}

const deleteTransaction = `-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1
`

func (q *Queries) DeleteTransaction(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransaction, id)
	return err
}

const getTransactionById = `-- name: GetTransactionById :one
SELECT id, amount, currency, title, created_at, payer_id FROM transactions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransactionById(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransactionById, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Currency,
		&i.Title,
		&i.CreatedAt,
		&i.PayerID,
	)
	return i, err
}

const getTransactionsByPayer = `-- name: GetTransactionsByPayer :many
SELECT
    t.id as transaction_id,
    t.amount as transaction_amount,
    t.currency as transaction_currency,
    t.title as transaction_title,
    t.created_at as transaction_created_at,
    u.id as payer_id,
    u.username as payer_username
FROM
    transactions t, users u
WHERE
    t.payer_id = $1
AND
    t.payer_id = u.id
LIMIT $2
OFFSET $3
`

type GetTransactionsByPayerParams struct {
	PayerID int64
	Limit   int32
	Offset  int32
}

type GetTransactionsByPayerRow struct {
	TransactionID        int64
	TransactionAmount    string
	TransactionCurrency  Currency
	TransactionTitle     string
	TransactionCreatedAt time.Time
	PayerID              int64
	PayerUsername        string
}

func (q *Queries) GetTransactionsByPayer(ctx context.Context, arg GetTransactionsByPayerParams) ([]GetTransactionsByPayerRow, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionsByPayer, arg.PayerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransactionsByPayerRow
	for rows.Next() {
		var i GetTransactionsByPayerRow
		if err := rows.Scan(
			&i.TransactionID,
			&i.TransactionAmount,
			&i.TransactionCurrency,
			&i.TransactionTitle,
			&i.TransactionCreatedAt,
			&i.PayerID,
			&i.PayerUsername,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTransaction = `-- name: UpdateTransaction :one
UPDATE transactions
SET amount = coalesce($1, amount),
    currency = coalesce($2, currency),
    title = coalesce($3, title)
WHERE id = $4 AND payer_id = $5
RETURNING id, amount, currency, title, created_at, payer_id
`

type UpdateTransactionParams struct {
	Amount   sql.NullString
	Currency NullCurrency
	Title    sql.NullString
	ID       int64
	PayerID  int64
}

func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, updateTransaction,
		arg.Amount,
		arg.Currency,
		arg.Title,
		arg.ID,
		arg.PayerID,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Currency,
		&i.Title,
		&i.CreatedAt,
		&i.PayerID,
	)
	return i, err
}
