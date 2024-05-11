package app

import (
	"context"
	"errors"

	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ListInvoices(ctx context.Context, req *desc.ListInvoicesRequest) (*desc.ListInvoicesResponse, error) {
	err := validateListInvoicesRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	invoices, err := i.invoicesService.ListInvoices(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invoicesService.ListInvoices: %v", err)
	}

	return &desc.ListInvoicesResponse{
		Invoices: models.InvoicesToProto(invoices),
	}, nil
}

func validateListInvoicesRequest(req *desc.ListInvoicesRequest) error {
	if req == nil || req.Filter == nil {
		return errors.New("filter is required")
	}

	filter := req.GetFilter()
	err := validation.ValidateStruct(
		filter,
		validation.Field(&filter.ClientIdIn, validation.Each(validation.NotNil, is.UUIDv4)),
	)
	return err
}
