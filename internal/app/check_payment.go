package app

import (
	"context"
	desc "github.com/fidesy-pay/payment-service/pkg/payment-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CheckPayment(ctx context.Context, req *desc.CheckPaymentRequest) (*desc.CheckPaymentResponse, error) {
	err := validateCheckPaymentRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	paymentStatus, err := i.paymentService.CheckPayment(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paymentService.CheckPayment: %v", err)
	}

	return &desc.CheckPaymentResponse{
		Status: paymentStatus,
	}, nil
}

func validateCheckPaymentRequest(req *desc.CheckPaymentRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Id, is.UUIDv4))

	return err
}
