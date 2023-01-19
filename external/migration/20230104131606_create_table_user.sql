-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (
    id varchar(36) primary key ,
    username varchar(255) ,
    password varchar(255) ,
    user_type smallint ,
    status boolean ,
    created_at timestamp ,
    updated_at timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user;
-- +goose StatementEnd
