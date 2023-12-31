package main

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/app"
	"github.com/fidesy-pay/invoices-service/internal/config"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/kafka"
	in_memory "github.com/fidesy-pay/invoices-service/internal/pkg/storage/in-memory"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	"github.com/fidesyx/platform/pkg/scratch"
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
		log.Fatalf("config.Init: %v", err)
	}

	storage := in_memory.New()

	kafkaConsumer, err := kafka.NewConsumer(ctx, paymentsTopic)
	if err != nil {
		log.Fatalf("kafka.NewConsumer: %v", err)
	}
	defer func() {
		err = kafkaConsumer.Close()
		if err != nil {
			log.Fatalf("kafkaConsumer.Close: %v", err)
		}
	}()

	cryptoServiceClient, err := scratch.NewClient[crypto_service.CryptoServiceClient](
		ctx,
		crypto_service.NewCryptoServiceClient,
		"crypto-service",
	)
	if err != nil {
		log.Fatalf("NewCryptoServiceClient: %v", err)
	}

	invoicesService := invoicesservice.New(ctx, storage, kafkaConsumer, cryptoServiceClient)

	impl := app.New(invoicesService)

	if err = scratchApp.Run(ctx, impl); err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
