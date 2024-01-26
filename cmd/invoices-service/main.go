package main

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/app"
	"github.com/fidesy-pay/invoices-service/internal/config"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/kafka"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	coingecko_api "github.com/fidesy-pay/invoices-service/pkg/coingecko-api"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/fidesyx/platform/pkg/scratch"
	"github.com/fidesyx/platform/pkg/scratch/logger"
	postgres "github.com/fidesyx/platform/pkg/scratch/storage"
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

	scratchApp, err := scratch.New(ctx)
	if err != nil {
		log.Fatalf("scratch.New: %v", err)
	}

	err = config.Init()
	if err != nil {
		logger.Fatalf("config.Init: %v", err)
	}

	kafkaConsumer, err := kafka.NewConsumer(ctx, paymentsTopic)
	if err != nil {
		logger.Fatalf("kafka.NewConsumer: %v", err)
	}
	defer func() {
		err = kafkaConsumer.Close()
		if err != nil {
			logger.Fatalf("kafkaConsumer.Close: %v", err)
		}
	}()

	cryptoServiceClient, err := scratch.NewClient[crypto_service.CryptoServiceClient](
		ctx,
		crypto_service.NewCryptoServiceClient,
		"fidesy:///crypto-service",
	)
	if err != nil {
		logger.Fatalf("NewCryptoServiceClient: %v", err)
	}

	coinGeckoAPIClient, err := scratch.NewClient[coingecko_api.CoinGeckoAPIClient](
		ctx,
		coingecko_api.NewCoinGeckoAPIClient,
		"fidesy:///external-api",
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

	// register reverse http proxy
	reverseProxyRouter := scratch.ReverseProxyRouter()
	err = desc.RegisterInvoicesServiceHandlerServer(ctx, reverseProxyRouter, impl)
	if err != nil {
		logger.Fatalf("RegisterInvoicesServiceHandlerServer: %v", err)
	}

	if err = scratchApp.Run(ctx, impl); err != nil {
		logger.Fatalf("app.Run: %v", err)
	}
}
