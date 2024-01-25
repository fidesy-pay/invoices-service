package main

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/app"
	"github.com/fidesy-pay/invoices-service/internal/config"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/kafka"
	in_memory "github.com/fidesy-pay/invoices-service/internal/pkg/storage/in-memory"
	coingecko_api "github.com/fidesy-pay/invoices-service/pkg/coingecko-api"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	"github.com/fidesyx/platform/pkg/scratch"
	"github.com/fidesyx/platform/pkg/scratch/logger"
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

	storage := in_memory.New()

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

	invoicesService := invoicesservice.New(ctx, storage, kafkaConsumer, cryptoServiceClient, coinGeckoAPIClient)

	impl := app.New(invoicesService)

	if err = scratchApp.Run(ctx, impl); err != nil {
		logger.Fatalf("app.Run: %v", err)
	}
}
