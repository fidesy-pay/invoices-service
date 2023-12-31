package app

import (
	"context"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateInvoice(ctx context.Context, req *desc.UpdateInvoiceRequest) (*desc.UpdateInvoiceResponse, error) {
	err := validateUpdateInvoiceRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	invoice, err := i.invoicesService.UpdateInvoice(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invoicesService.UpdateInvoice: %v", err)
	}

	return &desc.UpdateInvoiceResponse{
		Invoice: invoice.Proto(),
	}, nil
}

func validateUpdateInvoiceRequest(req *desc.UpdateInvoiceRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Id, validation.Required, is.UUIDv4),
		validation.Field(&req.Amount, validation.Required),
		validation.Field(&req.Chain, validation.Required),
		validation.Field(&req.Token, validation.Required),
	)
	return err
}
