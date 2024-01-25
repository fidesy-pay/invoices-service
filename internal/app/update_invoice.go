package app

import (
	"context"
	"errors"
	
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateInvoice(ctx context.Context, req *desc.UpdateInvoiceRequest) (*desc.UpdateInvoiceResponse, error) {
	updateInvoiceInput, err := invoicesservice.UpdateInvoiceInputFromRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	invoice, err := i.invoicesService.UpdateInvoice(ctx, updateInvoiceInput)
	if err != nil {
		if errors.Is(err, invoicesservice.ErrInvoiceAlreadyCompleted) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "invoicesService.UpdateInvoice: %v", err)
	}

	return &desc.UpdateInvoiceResponse{
		Invoice: invoice.Proto(),
	}, nil
}
