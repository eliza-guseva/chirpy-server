-- +goose Up
ALTER TABLE users ALTER COLUMN is_chirpy_red SET NOT NULL;

-- +goose Down
ALTER TABLE users ALTER COLUMN is_chirpy_red DROP NOT NULL;
