// Code generated by sqlc. DO NOT EDIT.
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
	Amount   string   `json:"amount"`
	Currency Currency `json:"currency"`
	Title    string   `json:"title"`
	PayerID  int64    `json:"payer_id"`
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
    t.id = $1
AND
    t.payer_id = u.id
`

type GetTransactionByIdRow struct {
	TransactionID        int64     `json:"transaction_id"`
	TransactionAmount    string    `json:"transaction_amount"`
	TransactionCurrency  Currency  `json:"transaction_currency"`
	TransactionTitle     string    `json:"transaction_title"`
	TransactionCreatedAt time.Time `json:"transaction_created_at"`
	PayerID              int64     `json:"payer_id"`
	PayerUsername        string    `json:"payer_username"`
}

func (q *Queries) GetTransactionById(ctx context.Context, id int64) (GetTransactionByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getTransactionById, id)
	var i GetTransactionByIdRow
	err := row.Scan(
		&i.TransactionID,
		&i.TransactionAmount,
		&i.TransactionCurrency,
		&i.TransactionTitle,
		&i.TransactionCreatedAt,
		&i.PayerID,
		&i.PayerUsername,
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
	PayerID int64 `json:"payer_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

type GetTransactionsByPayerRow struct {
	TransactionID        int64     `json:"transaction_id"`
	TransactionAmount    string    `json:"transaction_amount"`
	TransactionCurrency  Currency  `json:"transaction_currency"`
	TransactionTitle     string    `json:"transaction_title"`
	TransactionCreatedAt time.Time `json:"transaction_created_at"`
	PayerID              int64     `json:"payer_id"`
	PayerUsername        string    `json:"payer_username"`
}

func (q *Queries) GetTransactionsByPayer(ctx context.Context, arg GetTransactionsByPayerParams) ([]GetTransactionsByPayerRow, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionsByPayer, arg.PayerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTransactionsByPayerRow{}
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
WHERE id = $4
RETURNING id, amount, currency, title, created_at, payer_id
`

type UpdateTransactionParams struct {
	Amount   sql.NullString `json:"amount"`
	Currency NullCurrency   `json:"currency"`
	Title    sql.NullString `json:"title"`
	ID       int64          `json:"id"`
}

func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, updateTransaction,
		arg.Amount,
		arg.Currency,
		arg.Title,
		arg.ID,
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
