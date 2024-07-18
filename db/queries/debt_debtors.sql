-- name: CreateDebtDebtors :one
INSERT INTO debt_debtors (
    debt_id,
    debtor_id,
    amount,
    currency
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetDebtDebtorsByDebtId :many
SELECT * FROM debt_debtors
WHERE debt_id = $1;

-- name: GetDebtsOfDebtorId :many
SELECT * FROM debt_debtors
WHERE debtor_id = $1;

-- name: GetDebtDebtorsByDebtAndDebtor :one
SELECT * FROM debt_debtors
WHERE debt_id = $1 AND debtor_id = $2;

-- name: UpdateDebtDebtor :one
UPDATE debt_debtors
SET amount = coalesce(sqlc.narg('amount'), amount),
    currency = coalesce(sqlc.narg('currency'), currency)
WHERE debt_id = sqlc.arg('debtId') AND debtor_id = sqlc.arg('debtorId')
RETURNING *;

-- name: DeleteDebtDebtor :exec
DELETE FROM debt_debtors
WHERE debt_id = $1 AND debtor_id = $2;
