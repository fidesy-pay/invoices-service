package in_memory

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/google/uuid"
)

var (
	ErrInvoiceNotFoundByID = func(invoiceID uuid.UUID) error {
		return fmt.Errorf("invoice not found by id = %q", invoiceID.String())
	}
)

func (s *Storage) CreateInvoice(_ context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	s.mu.Lock()
	s.invoices[invoice.ID] = invoice
	s.mu.Unlock()

	return invoice, nil
}

type ListInvoicesFilter struct {
	// only one at a time
	IDIn      []uuid.UUID
	AddressIn []string
	StatusIn  []desc.InvoiceStatus
}

func (s *Storage) ListInvoices(_ context.Context, filter ListInvoicesFilter) ([]*models.Invoice, error) {
	invoices := make([]*models.Invoice, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(filter.IDIn) > 0 {
		for _, id := range filter.IDIn {
			invoice, ok := s.invoices[id]
			if !ok {
				continue
			}

			invoices = append(invoices, invoice)
		}
	}

	if len(filter.AddressIn) > 0 {
		for _, invoice := range s.invoices {
			for _, address := range filter.AddressIn {
				if invoice.Address == address {
					invoices = append(invoices, invoice)
				}
			}
		}
	}

	return invoices, nil
}

func (s *Storage) UpdateInvoice(_ context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.invoices[invoice.ID]
	if !ok {
		return nil, ErrInvoiceNotFoundByID(invoice.ID)
	}

	s.invoices[invoice.ID] = invoice

	return invoice, nil
}
