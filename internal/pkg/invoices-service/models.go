package invoicesservice

import (
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type CreateInvoiceInput struct {
	ClientID       uuid.UUID
	UsdCentsAmount int64
}

func CreateInvoiceInputFromRequest(req *desc.CreateInvoiceRequest) (*CreateInvoiceInput, error) {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.ClientId, validation.Required, is.UUIDv4),
		validation.Field(&req.UsdAmount, validation.Required),
	)
	if err != nil {
		return nil, err
	}

	return &CreateInvoiceInput{
		ClientID:       uuid.MustParse(req.GetClientId()),
		UsdCentsAmount: int64(req.GetUsdAmount() * 100),
	}, nil
}

type UpdateInvoiceInput struct {
	InvoiceID     uuid.UUID
	Chain         string
	Token         string
	PayerClientID *string
}

func UpdateInvoiceInputFromRequest(req *desc.UpdateInvoiceRequest) (*UpdateInvoiceInput, error) {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Id, validation.Required, is.UUIDv4),
		validation.Field(&req.Chain, validation.Required),
		validation.Field(&req.Token, validation.Required),
	)
	if err != nil {
		return nil, err
	}

	return &UpdateInvoiceInput{
		InvoiceID:     uuid.MustParse(req.GetId()),
		Chain:         req.GetChain(),
		Token:         req.GetToken(),
		PayerClientID: req.PayerClientId,
	}, nil
}
