package app

import (
	"context"
	desc "github.com/fidesy-pay/payment-service/pkg/payment-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreatePayment(ctx context.Context, req *desc.CreatePaymentRequest) (*desc.CreatePaymentResponse, error) {
	payment, err := i.paymentService.CreatePayment(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paymentService.CreatePayment: %v", err)
	}

	return &desc.CreatePaymentResponse{
		Id:      payment.ID.String(),
		Address: payment.Address,
	}, nil
}
