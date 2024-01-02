package app

import (
	"context"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateInvoice(ctx context.Context, req *desc.CreateInvoiceRequest) (*desc.CreateInvoiceResponse, error) {
	if err := validateCreateInvoiceRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	invoice, err := i.invoicesService.CreateInvoice(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invoicesService.CreateInvoice: %v", err)
	}

	return &desc.CreateInvoiceResponse{
		Id: invoice.ID.String(),
	}, nil
}

func validateCreateInvoiceRequest(req *desc.CreateInvoiceRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.ClientId, validation.Required, is.UUIDv4),
		validation.Field(&req.UsdAmount, validation.Required),
	)
	return err
}
