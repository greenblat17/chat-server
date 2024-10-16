-- +goose Up
-- +goose StatementBegin
ALTER TABLE chats
    RENAME COLUMN "name" TO chat_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chats
    RENAME COLUMN chat_name TO "name";
-- +goose StatementEnd
