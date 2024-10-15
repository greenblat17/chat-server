-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats
(
    id         SERIAL PRIMARY KEY,
    "name"     TEXT NOT NUll,
    created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
-- +goose StatementEnd
