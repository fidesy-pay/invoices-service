package main

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/payment-service/internal/app"
	"github.com/fidesy-pay/payment-service/internal/config"
	"github.com/fidesy-pay/payment-service/internal/pkg/kafka"
	payment_service "github.com/fidesy-pay/payment-service/internal/pkg/payment-service"
	in_memory "github.com/fidesy-pay/payment-service/internal/pkg/storage/in-memory"
	desc "github.com/fidesy-pay/payment-service/pkg/payment-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	paymentsTopic = "payments-json"
)

var (
	grpcPort    string
	proxyPort   string
	swaggerPort string
	metricsPort string
)

func main() {
	initEnvs()

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

	errChan := make(chan error)

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

	// run gRPC server
	go func() {
		err = runGrpcServer(ctx, impl)
		if err != nil {
			errChan <- err
		}
	}()

	// run HTTP handler
	go func() {
		log.Printf("http handlers run at :%s", proxyPort)
		err = runHttpReverseProxy(ctx, impl)
		if err != nil {
			errChan <- err
		}
	}()

	// run swagger
	go func() {
		log.Printf("swagger run at :%s", swaggerPort)

		fs := http.FileServer(http.Dir("./swaggerui"))
		http.Handle("/docs/", http.StripPrefix("/docs/", fs))

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", swaggerPort), nil))
	}()

	// run metrics
	//go func() {
	//	log.Printf("metrics run at :%s", metricsPort)
	//	err = metrics.Run(ctx, metricsPort)
	//	if err != nil {
	//		errChan <- err
	//	}
	//}()

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		log.Fatal(err)
	}
}

func runGrpcServer(ctx context.Context, impl *app.Implementation) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterPaymentServiceServer(s, impl)

	log.Printf("grpc run at :%s", grpcPort)

	errChan := make(chan error)
	go func() {
		err = s.Serve(lis)
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}

func runHttpReverseProxy(ctx context.Context, impl *app.Implementation) error {
	router := runtime.NewServeMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
		Debug:          false,
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", proxyPort),
		Handler: corsHandler.Handler(router),
	}

	desc.RegisterPaymentServiceHandlerServer(ctx, router, impl)

	return httpServer.ListenAndServe()
}

func initEnvs() {
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		log.Fatalf("GRPC_PORT ENV variable is required")
	}

	proxyPort = os.Getenv("PROXY_PORT")
	if proxyPort == "" {
		log.Fatalf("PROXY_PORT ENV variable is required")
	}

	swaggerPort = os.Getenv("SWAGGER_PORT")
	if swaggerPort == "" {
		log.Fatalf("SWAGGER_PORT ENV variable is required")
	}

	metricsPort = os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		log.Fatalf("METRICS_PORT ENV variable is required")
	}
}
