package main

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/app"
	"github.com/fidesy-pay/invoices-service/internal/config"
	"github.com/fidesy-pay/invoices-service/internal/pkg/consumers"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	admin_service "github.com/fidesy-pay/invoices-service/pkg/admin-service"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	external_api "github.com/fidesy-pay/invoices-service/pkg/external-api"
	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/logger"
	"github.com/fidesy/sdk/common/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	paymentsTopic = "payments-json"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	server, err := grpc.NewServer(
		grpc.WithPort(os.Getenv("GRPC_PORT")),
		grpc.WithMetricsPort(os.Getenv("METRICS_PORT")),
		grpc.WithDomainNameService(ctx, "domain-name-service:10000"),
		grpc.WithGraylog("graylog:5555"),
		grpc.WithTracer("http://jaeger:14268/api/traces"),
	)
	if err != nil {
		log.Fatalf("grpc.NewServer: %v", err)
	}

	err = config.Init()
	if err != nil {
		logger.Fatalf("config.Init: %v", err)
	}

	cryptoServiceClient, err := grpc.NewClient[crypto_service.CryptoServiceClient](
		ctx,
		crypto_service.NewCryptoServiceClient,
		"rpc:///crypto-service",
	)
	if err != nil {
		logger.Fatalf("NewCryptoServiceClient: %v", err)
	}

	externalAPI, err := grpc.NewClient[external_api.ExternalAPIClient](
		ctx,
		external_api.NewExternalAPIClient,
		"rpc:///external-api",
	)
	if err != nil {
		logger.Fatalf("NewCryptoServiceClient: %v", err)
	}

	adminServiceClient, err := grpc.NewClient[admin_service.AdminServiceClient](
		ctx,
		admin_service.NewAdminServiceClient,
		"rpc:///admin-service",
	)
	if err != nil {
		logger.Fatalf("NewAlertsServiceClient: %v", err)
	}

	pool, err := postgres.Connect(ctx, config.Get(config.PgDsn).(string))
	if err != nil {
		logger.Fatalf("postgres.Connect: %v", err)
	}

	storage := storage.New(pool)

	err = consumers.RegisterConsumer(
		ctx,
		consumers.NewWalletBalanceConsumer(storage, cryptoServiceClient),
		paymentsTopic,
	)
	if err != nil {
		logger.Fatalf("consumers.RegisterConsumer: %v", err)
	}

	invoicesService := invoicesservice.New(ctx, storage, cryptoServiceClient, externalAPI, adminServiceClient)

	impl := app.New(invoicesService)

	if err = server.Run(ctx, impl); err != nil {
		logger.Fatalf("app.Run: %v", err)
	}
}
