-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now(),
    "name" TEXT NOT NUll
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
-- +goose StatementEnd
