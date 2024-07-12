-- name: CreateTransaction :one
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
RETURNING *;

-- name: GetTransactionById :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: GetTransactionsByPayer :many
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
OFFSET $3;

-- name: UpdateTransaction :one
UPDATE transactions
SET amount = coalesce(sqlc.narg('amount'), amount),
    currency = coalesce(sqlc.narg('currency'), currency),
    title = coalesce(sqlc.narg('title'), title)
WHERE id = sqlc.arg('id') AND payer_id = sqlc.arg('payer_id')
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;
