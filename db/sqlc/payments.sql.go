// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: payments.sql

package db

import (
	"context"
	"database/sql"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (
    debt_id,
    debtor_id,
    amount
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, debt_id, debtor_id, amount, created_at
`

type CreatePaymentParams struct {
	DebtID   int64
	DebtorID int64
	Amount   string
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, createPayment, arg.DebtID, arg.DebtorID, arg.Amount)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.DebtID,
		&i.DebtorID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1
`

func (q *Queries) DeletePayment(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePayment, id)
	return err
}

const getPaymentsByDebtId = `-- name: GetPaymentsByDebtId :many
SELECT id, debt_id, debtor_id, amount, created_at FROM payments
WHERE debt_id = $1
`

func (q *Queries) GetPaymentsByDebtId(ctx context.Context, debtID int64) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, getPaymentsByDebtId, debtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Payment
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.DebtID,
			&i.DebtorID,
			&i.Amount,
			&i.CreatedAt,
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

const getPaymentsByDebtorId = `-- name: GetPaymentsByDebtorId :many
SELECT id, debt_id, debtor_id, amount, created_at FROM payments
WHERE debtor_id = $1
`

func (q *Queries) GetPaymentsByDebtorId(ctx context.Context, debtorID int64) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, getPaymentsByDebtorId, debtorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Payment
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.DebtID,
			&i.DebtorID,
			&i.Amount,
			&i.CreatedAt,
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

const getPaymentsById = `-- name: GetPaymentsById :one
SELECT id, debt_id, debtor_id, amount, created_at FROM payments
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPaymentsById(ctx context.Context, id int64) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPaymentsById, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.DebtID,
		&i.DebtorID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const updatePayment = `-- name: UpdatePayment :one
UPDATE payments
SET amount = coalesce($1, amount)
WHERE debt_id = $2 AND debtor_id = $3
RETURNING id, debt_id, debtor_id, amount, created_at
`

type UpdatePaymentParams struct {
	Amount   sql.NullString
	DebtId   int64
	DebtorId int64
}

func (q *Queries) UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, updatePayment, arg.Amount, arg.DebtId, arg.DebtorId)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.DebtID,
		&i.DebtorID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
