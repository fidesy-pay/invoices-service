package main

import (
	"context"
	"github.com/fidesy-pay/payment-service/internal/app"
	"github.com/fidesy-pay/payment-service/internal/config"
	"github.com/fidesy-pay/payment-service/internal/pkg/kafka"
	payment_service "github.com/fidesy-pay/payment-service/internal/pkg/payment-service"
	in_memory "github.com/fidesy-pay/payment-service/internal/pkg/storage/in-memory"
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

	err := config.Init()
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

	cryptoServiceClient, err := NewCryptoServiceClient(ctx)
	if err != nil {
		log.Fatalf("NewCryptoServiceClient: %v", err)
	}

	paymentService := payment_service.New(ctx, storage, kafkaConsumer, cryptoServiceClient)

	impl := app.New(paymentService)

	app := scratch.New()

	if err = app.Run(ctx, impl); err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
