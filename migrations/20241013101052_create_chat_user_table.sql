-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_users
(
    chat_id  INT REFERENCES chats (id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    PRIMARY KEY (chat_id, username)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chat_users;
-- +goose StatementEnd
