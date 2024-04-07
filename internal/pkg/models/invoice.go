package models

import (
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Invoice struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	ClientID       uuid.UUID          `db:"client_id" json:"client_id"`
	UsdCentsAmount int64              `db:"usd_cents_amount" json:"usd_cents_amount"`
	TokenAmount    *float64           `db:"token_amount" json:"token_amount"`
	Chain          string             `db:"chain" json:"chain"`
	Token          string             `db:"token" json:"token"`
	Status         desc.InvoiceStatus `db:"status" json:"status"`
	Address        string             `db:"address" json:"address"`
	CreatedAt      time.Time          `db:"created_at" json:"created_at"`
	PayerClientID  *string            `db:"payer_client_id" json:"payer_client_id"`
	GasLimit       *int               `db:"gas_limit" json:"gas_limit"`
}

func (i *Invoice) TableName() string {
	return "invoices"
}

func (i *Invoice) ToInsertMap() map[string]interface{} {
	return map[string]interface{}{
		"client_id":        i.ClientID.String(),
		"usd_cents_amount": i.UsdCentsAmount,
		"token_amount":     i.TokenAmount,
		"chain":            i.Chain,
		"token":            i.Token,
		"status":           i.Status,
		"address":          i.Address,
	}
}

func (i *Invoice) ToUpdateMap() map[string]interface{} {
	updateData := map[string]interface{}{}

	if i.Chain != "" {
		updateData["chain"] = i.Chain
	}

	if i.Token != "" {
		updateData["token"] = i.Token
	}

	if i.TokenAmount != nil {
		updateData["token_amount"] = *i.TokenAmount
	}

	if i.Address != "" {
		updateData["address"] = i.Address
	}

	if i.Status != desc.InvoiceStatus_UNKNOWN_STATUS {
		updateData["status"] = i.Status
	}

	if i.PayerClientID != nil {
		updateData["payer_client_id"] = i.PayerClientID
	}

	return updateData
}

func (i *Invoice) Proto() *desc.Invoice {
	if i == nil {
		return nil
	}

	invoice := &desc.Invoice{
		Id:        i.ID.String(),
		ClientId:  i.ClientID.String(),
		UsdAmount: float64(i.UsdCentsAmount) / 100,
		Chain:     i.Chain,
		Token:     i.Token,
		Status:    i.Status,
		Address:   i.Address,
		CreatedAt: timestamppb.New(i.CreatedAt),
	}

	if i.TokenAmount != nil {
		invoice.TokenAmount = *i.TokenAmount
	}

	if i.PayerClientID != nil {
		invoice.PayerClientId = *i.PayerClientID
	}

	return invoice
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
