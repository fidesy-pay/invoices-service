package app

import (
	"context"
	"github.com/fidesy-pay/payment-service/internal/pkg/models"
	desc "github.com/fidesy-pay/payment-service/pkg/payment-service"
	"google.golang.org/grpc"
)

type (
	PaymentService interface {
		CreatePayment(ctx context.Context, req *desc.CreatePaymentRequest) (*models.Payment, error)
		CheckPayment(ctx context.Context, paymentID string) (desc.PaymentStatus, error)
	}
)

type Implementation struct {
	desc.UnimplementedPaymentServiceServer

	paymentService PaymentService
}

func New(paymentService PaymentService) *Implementation {
	return &Implementation{
		paymentService: paymentService,
	}
}

func (i *Implementation) GetDescription() *grpc.ServiceDesc {
	return &desc.PaymentService_ServiceDesc
}
