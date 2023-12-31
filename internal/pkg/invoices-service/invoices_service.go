package invoices_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	inmemory "github.com/fidesy-pay/invoices-service/internal/pkg/storage/in-memory"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)

type (
	KafkaConsumer interface {
		Consume() <-chan *sarama.ConsumerMessage
	}
	CryptoServiceClient interface {
		AcceptCrypto(ctx context.Context, req *crypto_service.AcceptCryptoRequest, opts ...grpc.CallOption) (*crypto_service.AcceptCryptoResponse, error)
	}

	Storage interface {
		CreateInvoice(ctx context.Context, payment *models.Invoice) (*models.Invoice, error)
		ListInvoices(ctx context.Context, filter inmemory.ListInvoicesFilter) ([]*models.Invoice, error)
		UpdateInvoice(ctx context.Context, payment *models.Invoice) (*models.Invoice, error)
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

	go service.consumeTransactions(ctx)

	return service
}

func (s *Service) CreateInvoice(ctx context.Context, req *desc.CreateInvoiceRequest) (*models.Invoice, error) {
	invoice := &models.Invoice{
		ID:        uuid.New(),
		Status:    desc.InvoiceStatus_NEW,
		CreatedAt: time.Now(),
	}

	invoice, err := s.storage.CreateInvoice(ctx, invoice)
	if err != nil {
		return nil, fmt.Errorf("storage.CreateInvoice: %w", err)
	}

	return invoice, nil
}

func (s *Service) UpdateInvoice(ctx context.Context, req *desc.UpdateInvoiceRequest) (*models.Invoice, error) {
	invoiceID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("uuid.Parse invoiceID: %w", err)
	}

	invoices, err := s.storage.ListInvoices(ctx, inmemory.ListInvoicesFilter{
		IDIn: []uuid.UUID{invoiceID},
	})
	if err != nil {
		return nil, fmt.Errorf("storage.ListInvoices: %w", err)
	}
	invoice := invoices[0]

	acceptCryptoResp, err := s.cryptoServiceClient.AcceptCrypto(ctx, &crypto_service.AcceptCryptoRequest{
		Chain: req.GetChain(),
		Token: req.GetToken(),
	})
	if err != nil {
		return nil, fmt.Errorf("cryptoServiceClient.AcceptCrypto: %w", err)
	}

	invoice.Amount = req.GetAmount()
	invoice.Chain = req.GetChain()
	invoice.Token = req.GetToken()
	invoice.Status = desc.InvoiceStatus_PENDING
	invoice.Address = strings.ToLower(acceptCryptoResp.GetAddress())

	invoice, err = s.storage.UpdateInvoice(ctx, invoice)
	if err != nil {
		return nil, fmt.Errorf("storage.CreateInvoice: %w", err)
	}

	return invoice, nil
}

func (s *Service) CheckInvoice(ctx context.Context, invoiceIDStr string) (*models.Invoice, error) {
	// we validate ID in handler logic
	invoiceID := uuid.MustParse(invoiceIDStr)

	invoices, err := s.storage.ListInvoices(ctx, inmemory.ListInvoicesFilter{
		IDIn: []uuid.UUID{invoiceID},
	})
	if err != nil {
		return nil, fmt.Errorf("storage.ListInvoices: %w", err)
	}

	if len(invoices) == 0 {
		return nil, ErrInvoiceNotFoundByID(invoiceID)
	}

	return invoices[0], nil
}

func (s *Service) consumeTransactions(ctx context.Context) {
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

	invoices, err := s.storage.ListInvoices(ctx, inmemory.ListInvoicesFilter{
		AddressIn: []string{strings.ToLower(transaction.Receiver)},
	})
	if err != nil {
		log.Printf("processTopicMessage: storage.ListInvoices: %v", err)
		return
	}

	if len(invoices) == 0 {
		log.Printf("%v", ErrInvoiceNotFoundByAddress(transaction.Receiver))
		return
	}

	invoice := invoices[0]

	if transaction.Chain != invoice.Chain || transaction.Token != invoice.Token {
		return
	}

	if transaction.Amount >= invoice.Amount {
		invoice.Status = desc.InvoiceStatus_SUCCESS
		_, err = s.storage.UpdateInvoice(ctx, invoice)
		if err != nil {
			log.Printf("processTopicMessage: storage.UpdateInvoice: %v", err)
			return
		}
	}
}
