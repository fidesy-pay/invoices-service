package app

import (
	"context"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateInvoice(ctx context.Context, req *desc.CreateInvoiceRequest) (*desc.CreateInvoiceResponse, error) {
	invoice, err := i.invoicesService.CreateInvoice(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invoicesService.CreateInvoice: %v", err)
	}

	return &desc.CreateInvoiceResponse{
		Id: invoice.ID.String(),
	}, nil
}
