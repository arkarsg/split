-- name: CreateDebt :one
INSERT INTO debts (
    transaction_id
) VALUES (
    $1
)
RETURNING *;

-- name: GetDebtById :one
SELECT * FROM debts
WHERE id = $1 LIMIT 1;

-- name: GetDebtByTransactionId :one
SELECT * FROM debts
WHERE transaction_id = $1 LIMIT 1;

-- name: UpdateDebt :one
UPDATE debts
SET settled_amount = sqlc.arg('newSettledAmount')
WHERE id = $1
RETURNING *;

-- name: DeleteDebt :exec
DELETE FROM debts
WHERE id = $1;
