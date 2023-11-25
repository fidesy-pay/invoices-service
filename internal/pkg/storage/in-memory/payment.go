package in_memory

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/payment-service/internal/pkg/models"
	"github.com/google/uuid"
	"time"
)

var (
	ErrPaymentNotFoundByID = func(paymentID uuid.UUID) error {
		return fmt.Errorf("payment not found by id = %q", paymentID.String())
	}
)

func (s *Storage) CreatePayment(_ context.Context, payment *models.Payment) (*models.Payment, error) {
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}

	payment.CreatedAt = time.Now().UTC()

	s.mu.Lock()
	s.payments[payment.ID] = payment
	s.mu.Unlock()

	return payment, nil
}

type ListPaymentsFilter struct {
	// only one at a time
	IDIn      []uuid.UUID
	AddressIn []string
}

func (s *Storage) ListPayments(_ context.Context, filter ListPaymentsFilter) ([]*models.Payment, error) {
	payments := make([]*models.Payment, 0)

	s.mu.RLock()
	if len(filter.IDIn) > 0 {
		for _, id := range filter.IDIn {
			payment, ok := s.payments[id]
			if !ok {
				continue
			}

			payments = append(payments, payment)
		}
	}

	if len(filter.AddressIn) > 0 {
		for _, payment := range s.payments {
			for _, address := range filter.AddressIn {
				if payment.Address == address {
					payments = append(payments, payment)
				}
			}
		}
	}

	s.mu.RUnlock()

	return payments, nil
}

func (s *Storage) UpdatePayment(_ context.Context, payment *models.Payment) (*models.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.payments[payment.ID]
	if !ok {
		return nil, ErrPaymentNotFoundByID(payment.ID)
	}

	s.payments[payment.ID] = payment

	return payment, nil
}
