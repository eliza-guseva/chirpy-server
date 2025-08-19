-- name: CreateUser :one
INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetPassowordHash :one
SELECT hashed_password FROM users WHERE email = $1;
