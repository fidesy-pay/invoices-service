package main

import (
	"context"
	crypto_facade "github.com/fidesy-pay/payment-service/pkg/crypto-facade"
	"google.golang.org/grpc"
)

func NewCryptoServiceClient(ctx context.Context) (crypto_facade.CryptoFacadeClient, error) {
	conn, err := dialConnection(ctx, "crypto-facade:24707")
	if err != nil {
		return nil, err
	}

	return crypto_facade.NewCryptoFacadeClient(conn), nil
}

func dialConnection(ctx context.Context, url string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
