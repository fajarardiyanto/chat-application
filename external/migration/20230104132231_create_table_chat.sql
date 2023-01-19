-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat (
    id varchar(36) primary key ,
    from_user varchar(36) ,
    to_user varchar(36) ,
    message text ,
    message_type smallint ,
    file varchar(255) ,
    file_suffix varchar(255) ,
    created_at timestamp ,
    updated_at timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat;
-- +goose StatementEnd
