package main

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/app"
	"github.com/fidesy-pay/invoices-service/internal/config"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	coingecko_api "github.com/fidesy-pay/invoices-service/pkg/coingecko-api"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/kafka"
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

	kafkaConsumer, err := kafka.NewConsumer(ctx, config.Get(config.KafkaBrokers).([]string), paymentsTopic)
	if err != nil {
		logger.Fatalf("kafka.NewConsumer: %v", err)
	}
	defer func() {
		err = kafkaConsumer.Close()
		if err != nil {
			logger.Fatalf("kafkaConsumer.Close: %v", err)
		}
	}()

	cryptoServiceClient, err := grpc.NewClient[crypto_service.CryptoServiceClient](
		ctx,
		crypto_service.NewCryptoServiceClient,
		"rpc:///crypto-service",
	)
	if err != nil {
		logger.Fatalf("NewCryptoServiceClient: %v", err)
	}

	coinGeckoAPIClient, err := grpc.NewClient[coingecko_api.CoinGeckoAPIClient](
		ctx,
		coingecko_api.NewCoinGeckoAPIClient,
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

	invoicesService := invoicesservice.New(ctx, storage, kafkaConsumer, cryptoServiceClient, coinGeckoAPIClient)

	impl := app.New(invoicesService)

	if err = server.Run(ctx, impl); err != nil {
		logger.Fatalf("app.Run: %v", err)
	}
}
