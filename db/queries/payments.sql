-- name: CreatePayment :one
INSERT INTO payments (
    debt_id,
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

-- name: GetPaymentsByDebtId :many
SELECT * FROM payments
WHERE debt_id = $1;

-- name: GetPaymentsByDebtorId :many
SELECT * FROM payments
WHERE debtor_id = $1;

-- name: UpdatePayment :one
UPDATE payments
SET amount = coalesce(sqlc.narg('amount'), amount)
WHERE debt_id = sqlc.arg('debtId') AND debtor_id = sqlc.arg('debtorId')
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;
