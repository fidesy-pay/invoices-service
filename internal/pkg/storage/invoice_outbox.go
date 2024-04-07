package storage

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	"github.com/fidesy/sdk/common/postgres"
)

func (s *Storage) ListInvoicesOutbox(ctx context.Context, limit uint64) ([]*models.InvoiceOutbox, error) {
	query := postgres.Builder().
		Select(invoicesOutboxFields).
		From(invoicesOutboxTable).
		Limit(limit)

	return postgres.Select[models.InvoiceOutbox](ctx, s.pool, query)
}

func (s *Storage) DeleteInvoicesOutbox(ctx context.Context, ids []int) error {
	query := postgres.Builder().
		Delete(invoicesOutboxTable).
		Where(sq.Eq{
			"id": ids,
		})

	_, err := postgres.Exec[models.InvoiceOutbox](ctx, s.pool, query)
	return err
}

func (s *Storage) CreateInvoiceOutbox(ctx context.Context, invoice *models.Invoice) (*models.InvoiceOutbox, error) {
	message, err := json.Marshal(invoice)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	query := postgres.Builder().
		Insert(invoicesOutboxTable).
		SetMap(map[string]interface{}{
			"message": message,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", invoicesOutboxFields))

	return postgres.Exec[models.InvoiceOutbox](ctx, s.pool, query)
}
