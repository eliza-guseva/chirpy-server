
-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1 AND revoked_at is NULL;

-- name: ExpireRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW() WHERE token = $1;

-- name: GetUserByRefreshToken :one
SELECT * FROM users WHERE id = (SELECT user_id FROM refresh_tokens WHERE token = $1);
