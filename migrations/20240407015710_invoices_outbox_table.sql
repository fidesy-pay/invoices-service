-- +goose Up
-- +goose StatementBegin
CREATE TABLE invoices_outbox (
    id SERIAL PRIMARY KEY,
    message TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE invoices_outbox;
-- +goose StatementEnd
