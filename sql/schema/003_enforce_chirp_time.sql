-- +goose Up
ALTER TABLE chirps ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE chirps ALTER COLUMN updated_at SET NOT NULL;

-- +goose Down
ALTER TABLE chirps ALTER COLUMN created_at DROP NOT NULL;
ALTER TABLE chirps ALTER COLUMN updated_at DROP NOT NULL;
