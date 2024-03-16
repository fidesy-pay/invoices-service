package app

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/pkg/dto"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateInvoiceStatus(ctx context.Context, req *desc.UpdateInvoiceStatusRequest) (*desc.UpdateInvoiceStatusResponse, error) {
	updateInvoiceStatusParams, err := dto.UpdateInvoiceStatusParamsFromRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = i.invoicesService.UpdateInvoiceStatus(ctx, updateInvoiceStatusParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invoicesService.UpdateInvoiceStatus: %v", err)
	}

	return &desc.UpdateInvoiceStatusResponse{}, nil
}
