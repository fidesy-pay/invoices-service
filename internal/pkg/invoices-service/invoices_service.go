package invoicesservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/fidesy-pay/invoices-service/internal/pkg/common"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	coingecko_api "github.com/fidesy-pay/invoices-service/pkg/coingecko-api"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/fidesyx/platform/pkg/scratch/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strings"
	"time"
)

type (
	Service struct {
		storage             Storage
		kafkaConsumer       KafkaConsumer
		cryptoServiceClient CryptoServiceClient
		coinGeckoAPIClient  CoinGeckoAPIClient
	}

	KafkaConsumer interface {
		Consume() <-chan *sarama.ConsumerMessage
	}

	CryptoServiceClient interface {
		AcceptCrypto(ctx context.Context, in *crypto_service.AcceptCryptoRequest, opts ...grpc.CallOption) (*crypto_service.AcceptCryptoResponse, error)
		Transfer(ctx context.Context, in *crypto_service.TransferRequest, opts ...grpc.CallOption) (*crypto_service.TransferResponse, error)
	}

	CoinGeckoAPIClient interface {
		GetPrice(ctx context.Context, in *coingecko_api.GetPriceRequest, opts ...grpc.CallOption) (*coingecko_api.GetPriceResponse, error)
	}

	Storage interface {
		CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
		ListInvoices(ctx context.Context, filter storage.ListInvoicesFilter) ([]*models.Invoice, error)
		UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	}
)

func New(
	ctx context.Context,
	storage Storage,
	kafkaConsumer KafkaConsumer,
	cryptoServiceClient CryptoServiceClient,
	coinGeckoAPIClient CoinGeckoAPIClient,
) *Service {
	service := &Service{
		storage:             storage,
		kafkaConsumer:       kafkaConsumer,
		cryptoServiceClient: cryptoServiceClient,
		coinGeckoAPIClient:  coinGeckoAPIClient,
	}

	go service.consumeTransactions(ctx)

	return service
}

func (s *Service) CreateInvoice(ctx context.Context, input *CreateInvoiceInput) (*models.Invoice, error) {
	invoice := &models.Invoice{
		ID:        uuid.New(),
		ClientID:  input.ClientID,
		USDAmount: input.USDAmount,
		Status:    desc.InvoiceStatus_NEW,
		CreatedAt: time.Now(),
	}

	var err error
	invoice, err = s.storage.CreateInvoice(ctx, invoice)
	if err != nil {
		return nil, fmt.Errorf("storage.CreateInvoice: %w", err)
	}

	return invoice, nil
}

func (s *Service) UpdateInvoice(ctx context.Context, input *UpdateInvoiceInput) (*models.Invoice, error) {
	invoices, err := s.storage.ListInvoices(ctx, storage.ListInvoicesFilter{
		IDIn: []uuid.UUID{input.InvoiceID},
	})
	if err != nil {
		return nil, fmt.Errorf("storage.ListInvoices: %w", err)
	}
	if len(invoices) == 0 {
		return nil, ErrInvoiceNotFoundByID(input.InvoiceID)
	}

	invoice := invoices[0]

	if invoice.Status == desc.InvoiceStatus_SUCCESS {
		return nil, ErrInvoiceAlreadyCompleted
	}

	acceptCryptoResp, err := s.cryptoServiceClient.AcceptCrypto(ctx, &crypto_service.AcceptCryptoRequest{
		InvoiceId: input.InvoiceID.String(),
		Chain:     input.Chain,
		Token:     input.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("cryptoServiceClient.AcceptCrypto: %w", err)
	}

	tokenPriceResp, err := s.coinGeckoAPIClient.GetPrice(ctx, &coingecko_api.GetPriceRequest{
		Symbol: input.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("coinGeckoAPIClient.GetPrice: %w", err)
	}

	tokenAmount := invoice.USDAmount / tokenPriceResp.GetPriceUsd()
	invoice.TokenAmount = &tokenAmount
	invoice.Chain = input.Chain
	invoice.Token = input.Token
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

	invoices, err := s.storage.ListInvoices(ctx, storage.ListInvoicesFilter{
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

func (s *Service) ListInvoices(ctx context.Context, reqFilter *desc.ListInvoicesRequest_Filter) ([]*models.Invoice, error) {
	var err error

	filter := storage.ListInvoicesFilter{}
	if len(reqFilter.ClientIdIn) > 0 {
		filter.ClientIDIn, err = common.ConvertToUUIDs(reqFilter.GetClientIdIn())
		if err != nil {
			return nil, fmt.Errorf("common.ConvertToUUIDs: %w", err)
		}
	}

	if len(reqFilter.IdIn) > 0 {
		filter.IDIn, err = common.ConvertToUUIDs(reqFilter.GetIdIn())
		if err != nil {
			return nil, fmt.Errorf("common.ConvertToUUIDs: %w", err)
		}
	}

	invoices, err := s.storage.ListInvoices(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("storage.ListInvoices: %w", err)
	}

	return invoices, nil
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
	if message == nil {
		return
	}

	transaction := new(models.Transaction)
	err := json.Unmarshal(message.Value, &transaction)
	if err != nil {
		logger.Errorf(zap.Error(
			fmt.Errorf("consumer: json.Unmarshal: %v", err),
		))
		return
	}

	logger.Info("Transaction", zap.ByteString("message", message.Value))

	invoices, err := s.storage.ListInvoices(ctx, storage.ListInvoicesFilter{
		AddressIn: []string{strings.ToLower(transaction.Receiver)},
	})
	if err != nil {
		logger.Errorf(zap.Error(
			fmt.Errorf("processTopicMessage: storage.ListInvoices: %v", err),
		))
		return
	}

	if len(invoices) == 0 {
		logger.Errorf(zap.Error(
			fmt.Errorf("%v", ErrInvoiceNotFoundByAddress(transaction.Receiver)),
		))
		return
	}

	invoice := invoices[0]

	if transaction.Chain != invoice.Chain || transaction.Token != invoice.Token {
		return
	}

	if transaction.Amount >= *invoice.TokenAmount {
		invoice.Status = desc.InvoiceStatus_SUCCESS
		_, err = s.storage.UpdateInvoice(ctx, invoice)
		if err != nil {
			logger.Errorf(zap.Error(
				fmt.Errorf("processTopicMessage: storage.UpdateInvoice: %v", err),
			))
			return
		}

		transferResp, err := s.cryptoServiceClient.Transfer(ctx, &crypto_service.TransferRequest{
			ClientId:  invoice.ClientID.String(),
			InvoiceId: invoice.ID.String(),
			Chain:     invoice.Chain,
			Token:     invoice.Token,
		})
		if err != nil {
			logger.Errorf(zap.Error(
				fmt.Errorf("cryptoServiceClient.Transfer: %v", err),
			))
			return
		}

		logger.Info(
			"Transaction hash",
			zap.String("hash", transferResp.TransactionHash),
		)
	}
}
