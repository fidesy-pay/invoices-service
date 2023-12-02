package main

import (
	"context"
	crypto_service "github.com/fidesy-pay/payment-service/pkg/crypto-service"
	"google.golang.org/grpc"
)

func NewCryptoServiceClient(ctx context.Context) (crypto_service.CryptoServiceClient, error) {
	conn, err := dialConnection(ctx, "crypto-service:24707")
	if err != nil {
		return nil, err
	}

	return crypto_service.NewCryptoServiceClient(conn), nil
}

func dialConnection(ctx context.Context, url string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
