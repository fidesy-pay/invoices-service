package app

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/pkg/dto"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"google.golang.org/grpc"
)

type (
	Implementation struct {
		desc.UnimplementedInvoicesServiceServer

		invoicesService InvoicesService
	}
	InvoicesService interface {
		CreateInvoice(ctx context.Context, input *invoicesservice.CreateInvoiceInput) (*models.Invoice, error)
		CheckInvoice(ctx context.Context, invoiceID string) (*models.Invoice, error)
		UpdateInvoice(ctx context.Context, input *invoicesservice.UpdateInvoiceInput) (*models.Invoice, error)
		ListInvoices(ctx context.Context, reqFilter *desc.ListInvoicesRequest_Filter) ([]*models.Invoice, error)
		UpdateInvoiceStatus(ctx context.Context, params dto.UpdateInvoiceStatusParams) error
	}
)

func New(invoicesService InvoicesService) *Implementation {
	return &Implementation{
		invoicesService: invoicesService,
	}
}

func (i *Implementation) GetDescription() *grpc.ServiceDesc {
	return &desc.InvoicesService_ServiceDesc
}
