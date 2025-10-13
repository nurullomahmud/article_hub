-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image TEXT,
    title TEXT NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd