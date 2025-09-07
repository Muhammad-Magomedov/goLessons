-- +goose Up
-- +goose StatementBegin
alter table users add column hashed_password text not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
