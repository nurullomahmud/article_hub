-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS paragraphs (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    headline VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    order INT
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS paragraphs;
-- +goose StatementEnd