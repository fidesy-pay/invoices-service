package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fidesy-pay/invoices-service/internal/app"
	"github.com/fidesy-pay/invoices-service/internal/config"
	"github.com/fidesy-pay/invoices-service/internal/pkg/consumers"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	external_api "github.com/fidesy-pay/invoices-service/pkg/external-api"
	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/kafka"
	"github.com/fidesy/sdk/common/logger"
	"github.com/fidesy/sdk/common/outbox_processor"
	"github.com/fidesy/sdk/common/postgres"
)

const (
	balancesTopic = "balances-json"
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

	pool, err := postgres.Connect(ctx, config.Get(config.PgDsn).(string))
	if err != nil {
		logger.Fatalf("postgres.Connect: %v", err)
	}

	storage := storage.New(pool)

	err = kafka.RegisterConsumer(
		ctx,
		consumers.NewWalletBalanceConsumer(storage, cryptoServiceClient),
		config.Get(config.KafkaBrokers).([]string),
		balancesTopic,
	)
	if err != nil {
		logger.Fatalf("consumers.RegisterConsumer: %v", err)
	}

	// Register outbox

	producer, err := kafka.NewProducer(ctx, config.Get(config.KafkaBrokers).([]string))
	if err != nil {
		panic(err)
	}

	outboxProcessor := outbox_processor.New(
		"invoices",
		"invoices-json",
		pool,
		producer,
	)
	go outboxProcessor.Publish(ctx)

	invoicesService := invoicesservice.New(ctx, storage, cryptoServiceClient, externalAPI)

	impl := app.New(invoicesService)

	if err = server.Run(ctx, impl); err != nil {
		logger.Fatalf("app.Run: %v", err)
	}
}
