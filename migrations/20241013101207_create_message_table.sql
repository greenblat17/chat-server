-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages
(
    id            SERIAL PRIMARY KEY,
    from_username TEXT      NOT NULL,
    text          TEXT      NOT NULL,
    sent_at       TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
