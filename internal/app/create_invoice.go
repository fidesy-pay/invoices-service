package app

import (
	"context"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateInvoice(ctx context.Context, req *desc.CreateInvoiceRequest) (*desc.CreateInvoiceResponse, error) {
	createInvoiceInput, err := invoicesservice.CreateInvoiceInputFromRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	invoice, err := i.invoicesService.CreateInvoice(ctx, createInvoiceInput)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invoicesService.CreateInvoice: %v", err)
	}

	return &desc.CreateInvoiceResponse{
		Id: invoice.ID.String(),
	}, nil
}
