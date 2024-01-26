-- +goose Up
-- +goose StatementBegin
create table invoices
(
    id         uuid      default uuid_generate_v4() not null
        primary key,
    client_id    uuid not null ,
    usd_amount real not null ,
    token_amount      real,
    chain text,
    token text,
    status     int,
    address text,
    created_at timestamp default now()              not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table invoices;
-- +goose StatementEnd
