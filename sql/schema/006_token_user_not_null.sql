-- +goose Up
ALTER TABLE refresh_tokens ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
ALTER TABLE refresh_tokens ALTER COLUMN user_id DROP NOT NULL;
