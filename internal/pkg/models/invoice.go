package models

import (
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Invoice struct {
	ID        uuid.UUID
	ClientID  uuid.UUID
	Amount    float64
	Chain     string
	Token     string
	Status    desc.InvoiceStatus
	Address   string
	CreatedAt time.Time
}

func (i *Invoice) Proto() *desc.Invoice {
	if i == nil {
		return nil
	}

	return &desc.Invoice{
		Id:        i.ID.String(),
		Amount:    i.Amount,
		Chain:     i.Chain,
		Token:     i.Token,
		Status:    i.Status,
		Address:   i.Address,
		CreatedAt: timestamppb.New(i.CreatedAt),
	}
}

func InvoicesToProto(invoices []*Invoice) []*desc.Invoice {
	if invoices == nil {
		return []*desc.Invoice{}
	}

	result := make([]*desc.Invoice, len(invoices))
	for i := 0; i < len(invoices); i++ {
		result[i] = invoices[i].Proto()
	}

	return result
}
