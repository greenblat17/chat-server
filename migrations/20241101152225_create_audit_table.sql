-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit
(
    id         SERIAL PRIMARY KEY,
    entity     VARCHAR(255) NOT NULL,
    action     VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit;
-- +goose StatementEnd