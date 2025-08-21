-- name: CreateChirp :one
INSERT INTO chirps (user_id, body) VALUES ($1, $2) RETURNING *;

-- name: GetChirpsASC :many
SELECT 
    *
FROM chirps
ORDER BY chirps.created_at ASC;

-- name: GetChirpsDESC :many
SELECT 
    *
FROM chirps
ORDER BY chirps.created_at DESC;

-- name: GetChirpsByUserIDASC :many
SELECT * FROM chirps
WHERE chirps.user_id = $1
ORDER BY chirps.created_at ASC;

-- name: GetChirpsByUserIDDESC :many
SELECT * FROM chirps
WHERE chirps.user_id = $1                
ORDER BY chirps.created_at DESC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE chirps.id = $1;

-- name: ResetChirps :exec
DELETE FROM chirps;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;
