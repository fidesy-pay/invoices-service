-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN payer_client_id UUID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices DROP COLUMN payer_client_id;
-- +goose StatementEnd
