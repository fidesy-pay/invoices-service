package payment_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/fidesy-pay/payment-service/internal/pkg/models"
	in_memory "github.com/fidesy-pay/payment-service/internal/pkg/storage/in-memory"
	crypto_service "github.com/fidesy-pay/payment-service/pkg/crypto-service"
	desc "github.com/fidesy-pay/payment-service/pkg/payment-service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"strings"
)

type (
	KafkaConsumer interface {
		Consume() <-chan *sarama.ConsumerMessage
	}
	CryptoServiceClient interface {
		AcceptCrypto(ctx context.Context, req *crypto_service.AcceptCryptoRequest, opts ...grpc.CallOption) (*crypto_service.AcceptCryptoResponse, error)
	}

	Storage interface {
		CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
		ListPayments(ctx context.Context, filter in_memory.ListPaymentsFilter) ([]*models.Payment, error)
		UpdatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	}
)

type Service struct {
	storage             Storage
	kafkaConsumer       KafkaConsumer
	cryptoServiceClient CryptoServiceClient
}

func New(
	ctx context.Context,
	storage Storage,
	kafkaConsumer KafkaConsumer,
	cryptoServiceClient CryptoServiceClient,
) *Service {
	service := &Service{
		storage:             storage,
		kafkaConsumer:       kafkaConsumer,
		cryptoServiceClient: cryptoServiceClient,
	}

	go service.consumePayments(ctx)

	return service
}

func (s *Service) CreatePayment(ctx context.Context, req *desc.CreatePaymentRequest) (*models.Payment, error) {
	acceptCryptoResp, err := s.cryptoServiceClient.AcceptCrypto(ctx, &crypto_service.AcceptCryptoRequest{
		Chain: crypto_service.Chain(req.GetChain()),
		Token: crypto_service.Token(req.GetToken()),
	})
	if err != nil {
		return nil, fmt.Errorf("cryptoServiceClient.AcceptCrypto: %w", err)
	}

	payment := &models.Payment{
		Amount:  req.GetAmount(),
		Chain:   req.GetChain(),
		Token:   req.GetToken(),
		Status:  desc.PaymentStatus_PENDING,
		Address: strings.ToLower(acceptCryptoResp.GetAddress()),
	}

	payment, err = s.storage.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("storage.CreatePayment: %w", err)
	}

	return payment, nil
}

func (s *Service) CheckPayment(ctx context.Context, paymentIDStr string) (desc.PaymentStatus, error) {
	// we validate ID in handler logic
	paymentID := uuid.MustParse(paymentIDStr)

	payments, err := s.storage.ListPayments(ctx, in_memory.ListPaymentsFilter{
		IDIn: []uuid.UUID{paymentID},
	})
	if err != nil {
		return desc.PaymentStatus_UNKNOWN_STATUS, fmt.Errorf("storage.ListPayments: %w", err)
	}

	if len(payments) == 0 {
		return desc.PaymentStatus_UNKNOWN_STATUS, ErrPaymentNotFoundByID(paymentID)
	}

	return payments[0].Status, nil
}

func (s *Service) consumePayments(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-s.kafkaConsumer.Consume():
			go s.processTopicMessage(ctx, message)
		}
	}
}

func (s *Service) processTopicMessage(ctx context.Context, message *sarama.ConsumerMessage) {
	transaction := new(models.Transaction)

	err := json.Unmarshal(message.Value, &transaction)
	if err != nil {
		log.Printf("consumer: json.Unmarshal: %v", err)
	}

	log.Println("New transaction:", transaction)

	payments, err := s.storage.ListPayments(ctx, in_memory.ListPaymentsFilter{
		AddressIn: []string{strings.ToLower(transaction.Receiver)},
	})
	if err != nil {
		log.Printf("processTopicMessage: storage.ListPayments: %v", err)
		return
	}

	if len(payments) == 0 {
		log.Printf("%v", ErrPaymentNotFoundByAddress(transaction.Receiver))
		return
	}

	payment := payments[0]

	if transaction.Chain != payment.Chain.String() || transaction.Token != payment.Token.String() {
		return
	}

	if transaction.Amount >= payment.Amount {
		payment.Status = desc.PaymentStatus_SUCCESS
		_, err = s.storage.UpdatePayment(ctx, payment)
		if err != nil {
			log.Printf("processTopicMessage: storage.UpdatePayment: %v", err)
			return
		}
	}
}
