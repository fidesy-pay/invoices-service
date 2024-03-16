package dto

import (
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type UpdateInvoiceStatusParams struct {
	InvoiceID uuid.UUID
	Status    desc.InvoiceStatus
}

func UpdateInvoiceStatusParamsFromRequest(req *desc.UpdateInvoiceStatusRequest) (UpdateInvoiceStatusParams, error) {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.InvoiceId, validation.Required, is.UUIDv4),
		validation.Field(&req.Status, validation.Required))
	if err != nil {
		return UpdateInvoiceStatusParams{}, err
	}

	return UpdateInvoiceStatusParams{
		InvoiceID: uuid.MustParse(req.GetInvoiceId()),
		Status:    req.GetStatus(),
	}, nil
}
