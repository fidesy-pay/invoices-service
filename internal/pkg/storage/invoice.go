package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/fidesy/sdk/common/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type ListInvoicesFilter struct {
	IDIn        []uuid.UUID
	AddressIn   []string
	ClientIDIn  []uuid.UUID
	StatusIn    []desc.InvoiceStatus
	CreatedAtLt *time.Time
}

func (s *Storage) ListInvoices(ctx context.Context, filter ListInvoicesFilter) ([]*models.Invoice, error) {
	query := postgres.Builder().
		Select(invoiceFields).
		From(invoicesTable)

	// does not show expired invoices
	//query = query.Where(sq.NotEq{
	//	"status": desc.InvoiceStatus_EXPIRED,
	//})
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

	if filter.CreatedAtLt != nil {
		query = query.Where(sq.Lt{
			"created_at": filter.CreatedAtLt,
		})
	}

	query = query.OrderBy("CREATED_AT DESC")

	return postgres.Select[models.Invoice](ctx, s.pool, query)
}

func (s *Storage) CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	query := postgres.Builder().
		Insert(invoicesTable).
		SetMap(invoice.ToInsertMap()).
		Suffix(fmt.Sprintf("RETURNING %s", invoiceFields))

	var err error
	err = postgres.WithTransaction(ctx, s.pool, func(tx pgx.Tx) error {
		invoice, err = postgres.Exec[models.Invoice](ctx, tx, query)
		if err != nil {
			return fmt.Errorf("postgres.Exec: %w", err)
		}

		_, err = s.CreateInvoiceOutbox(ctx, invoice)
		if err != nil {
			return fmt.Errorf("s.CreateInvoiceOutbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("postgres.Tx: %w", err)
	}

	return invoice, nil
}

func (s *Storage) UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	query := postgres.Builder().
		Update(invoicesTable).
		SetMap(invoice.ToUpdateMap()).
		Where(sq.Eq{
			"id": invoice.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", invoiceFields))

	var err error
	err = postgres.WithTransaction(ctx, s.pool, func(tx pgx.Tx) error {
		invoice, err = postgres.Exec[models.Invoice](ctx, tx, query)
		if err != nil {
			return fmt.Errorf("postgres.Exec: %w", err)
		}

		_, err = s.CreateInvoiceOutbox(ctx, invoice)
		if err != nil {
			return fmt.Errorf("s.CreateInvoiceOutbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("postgres.Tx: %w", err)
	}

	return invoice, nil
}
