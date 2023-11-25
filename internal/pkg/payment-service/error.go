package payment_service

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrPaymentNotFoundByID = func(paymentID uuid.UUID) error {
		return fmt.Errorf("payment not found by id = %q", paymentID.String())
	}

	ErrPaymentNotFoundByAddress = func(address string) error {
		return fmt.Errorf("payment not found by address = %q", address)
	}
)
