-- +goose Up
CREATE TABLE chirps (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    user_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    body text NOT NULL
);

-- +goose Down
DROP TABLE chirps;
