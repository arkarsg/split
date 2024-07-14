-- name: CreatePayment :one
INSERT INTO payments (
    transaction_id,
    debtor_id,
    amount
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetPaymentsById :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: GetPaymentsByTransactionId :many
SELECT * FROM payments
WHERE transaction_id = $1;

-- name: GetPaymentsByDebtorId :many
SELECT * FROM payments
WHERE debtor_id = $1;

-- name: UpdatePayment :one
UPDATE payments
SET amount = coalesce(sqlc.narg('amount'), amount)
WHERE transaction_id = sqlc.arg('transactionId') AND debtor_id = sqlc.arg('debtorId')
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;
