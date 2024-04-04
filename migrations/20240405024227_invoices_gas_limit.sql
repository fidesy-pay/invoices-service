-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN gas_limit INT DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices DROP COLUMN gas_limit;
-- +goose StatementEnd