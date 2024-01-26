package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	postgres "github.com/fidesyx/platform/pkg/scratch/storage"
	"github.com/google/uuid"
)

type ListInvoicesFilter struct {
	IDIn       []uuid.UUID
	AddressIn  []string
	ClientIDIn []uuid.UUID
	StatusIn   []desc.InvoiceStatus
}

func (s *Storage) ListInvoices(ctx context.Context, filter ListInvoicesFilter) ([]*models.Invoice, error) {
	query := postgres.Builder().
		Select(invoiceFields).
		From(invoicesTable)

	if len(filter.IDIn) > 0 {
		query = query.Where(sq.Eq{
			"id": filter.IDIn,
		})
	}

	if len(filter.AddressIn) > 0 {
		query = query.Where(sq.Eq{
			"address": filter.AddressIn,
		})
	}

	if len(filter.ClientIDIn) > 0 {
		query = query.Where(sq.Eq{
			"client_id": filter.ClientIDIn,
		})
	}

	if len(filter.StatusIn) > 0 {
		query = query.Where(sq.Eq{
			"status": filter.StatusIn,
		})
	}

	return postgres.Selectx[models.Invoice](ctx, s.pool, query)
}

func (s *Storage) CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	query := postgres.Builder().
		Insert(invoicesTable).
		SetMap(invoice.ToInsertMap()).
		Suffix(fmt.Sprintf("RETURNING %s", invoiceFields))

	return postgres.Getx[models.Invoice](ctx, s.pool, query)
}

func (s *Storage) UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	query := postgres.Builder().
		Update(invoicesTable).
		SetMap(invoice.ToUpdateMap()).
		Where(sq.Eq{
			"id": invoice.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", invoiceFields))

	return postgres.Getx[models.Invoice](ctx, s.pool, query)
}
