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
	"github.com/fidesy/sdk/common/logger"
	"github.com/google/uuid"
	"github.com/samber/lo"
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
		CancelAcceptingCrypto(ctx context.Context, in *crypto_service.CancelAcceptingCryptoRequest, opts ...grpc.CallOption) (*crypto_service.CancelAcceptingCryptoResponse, error)
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
	go service.cleanExpiredInvoicesWorker(ctx)

	return service
}

func (s *Service) CreateInvoice(ctx context.Context, input *CreateInvoiceInput) (*models.Invoice, error) {
	invoice := &models.Invoice{
		ID:             uuid.New(),
		ClientID:       input.ClientID,
		UsdCentsAmount: input.UsdCentsAmount,
		Status:         desc.InvoiceStatus_NEW,
		CreatedAt:      time.Now(),
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

	tokenAmount := float64(invoice.UsdCentsAmount) / (100 * tokenPriceResp.GetPriceUsd())
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

	wallet := new(models.Wallet)
	err := json.Unmarshal(message.Value, &wallet)
	if err != nil {
		logger.Errorf("consumer: json.Unmarshal: %v", err)
		return
	}

	invoices, err := s.storage.ListInvoices(ctx, storage.ListInvoicesFilter{
		AddressIn: []string{strings.ToLower(wallet.Address)},
	})
	if err != nil {
		logger.Errorf("processTopicMessage: storage.ListInvoices: %v", err)
		return
	}

	if len(invoices) == 0 {
		logger.Errorf("%v", ErrInvoiceNotFoundByAddress(wallet.Address))
		return
	}

	invoice := invoices[0]
	if wallet.Chain != invoice.Chain || wallet.Token != invoice.Token {
		return
	}

	if wallet.Balance >= int64(*invoice.TokenAmount*1e18) {
		invoice.Status = desc.InvoiceStatus_SUCCESS
		_, err = s.storage.UpdateInvoice(ctx, invoice)
		if err != nil {
			logger.Errorf("processTopicMessage: storage.UpdateInvoice: %v", err)
			return
		}

		_, err = s.cryptoServiceClient.CancelAcceptingCrypto(ctx, &crypto_service.CancelAcceptingCryptoRequest{
			InvoiceId: invoice.ID.String(),
		})
		if err != nil {
			logger.Errorf("cryptoServiceClient.CancelAcceptingCrypto: %w", err)
		}

		transferResp, err := s.cryptoServiceClient.Transfer(ctx, &crypto_service.TransferRequest{
			ClientId:  invoice.ClientID.String(),
			InvoiceId: invoice.ID.String(),
		})
		if err != nil {
			logger.Errorf("cryptoServiceClient.Transfer: %v", err)
			return
		}

		logger.Info(
			"Transaction hash",
			zap.String("hash", transferResp.TransactionHash),
		)

	}
}

func (s *Service) cleanExpiredInvoicesWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(5 * time.Second):
			go s.cleanExpiredInvoices(ctx)
		}
	}
}

func (s *Service) cleanExpiredInvoices(ctx context.Context) {
	invoices, err := s.storage.ListInvoices(ctx, storage.ListInvoicesFilter{
		StatusIn:    []desc.InvoiceStatus{desc.InvoiceStatus_NEW, desc.InvoiceStatus_PENDING},
		CreatedAtLt: lo.ToPtr(time.Now().Add(-20 * time.Minute)),
	})
	if err != nil {
		logger.Errorf("cleanExpiredInvoices: storage.ListInvoices: %w", err)
		return
	}

	for _, invoice := range invoices {
		invoice.Status = desc.InvoiceStatus_EXPIRED
		_, err = s.storage.UpdateInvoice(ctx, invoice)
		if err != nil {
			logger.Errorf("cleanExpiredInvoices: storage.UpdateInvoice: %w", err)
			continue
		}

		if invoice.Address == "" {
			continue
		}

		_, err = s.cryptoServiceClient.CancelAcceptingCrypto(ctx, &crypto_service.CancelAcceptingCryptoRequest{
			InvoiceId: invoice.ID.String(),
		})
		if err != nil {
			logger.Errorf("cleanExpiredInvoices: cryptoServiceClient.CancelAcceptingCrypto: %w", err)
			continue
		}
	}
}
