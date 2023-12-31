package app

import (
	"context"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CheckInvoice(ctx context.Context, req *desc.CheckInvoiceRequest) (*desc.CheckInvoiceResponse, error) {
	err := validateCheckInvoiceRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	invoice, err := i.invoicesService.CheckInvoice(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invoicesService.CheckInvoice: %v", err)
	}

	return &desc.CheckInvoiceResponse{
		Invoice: invoice.Proto(),
	}, nil
}

func validateCheckInvoiceRequest(req *desc.CheckInvoiceRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Id, is.UUIDv4))

	return err
}
