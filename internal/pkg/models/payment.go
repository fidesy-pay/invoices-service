package models

import (
	desc "github.com/fidesy-pay/payment-service/pkg/payment-service"
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID        uuid.UUID
	Amount    float64
	Chain     desc.Chain
	Token     desc.Token
	Status    desc.PaymentStatus
	Address   string
	CreatedAt time.Time
}
