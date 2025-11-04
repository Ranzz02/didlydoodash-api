-- name: CreateUser :one
INSERT INTO users (id, email, password, username)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUsers :many
SELECT
    id,
    username,
    avatar
FROM users
WHERE
    (sqlc.arg('search') = '' OR username ILIKE '%' || sqlc.arg('search') || '%')
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset'); 

-- name: GetByEmail :one
SELECT
    *
FROM users
WHERE email = $1;