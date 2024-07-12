-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  username, email
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
