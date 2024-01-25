package invoicesservice

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrInvoiceNotFoundByID = func(invoiceID uuid.UUID) error {
		return fmt.Errorf("invoice not found by id = %q", invoiceID.String())
	}

	ErrInvoiceNotFoundByAddress = func(address string) error {
		return fmt.Errorf("invoice not found by address = %q", address)
	}

	ErrInvoiceAlreadyCompleted = errors.New("invoice is already completed")
)
