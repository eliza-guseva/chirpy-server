-- name: CreateUser :one
INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetPassowordHash :one
SELECT hashed_password FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET 
    email = $1,
    hashed_password = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpgradeUser :one
UPDATE users SET
    is_chirpy_red = true
WHERE id = $1
RETURNING *;
