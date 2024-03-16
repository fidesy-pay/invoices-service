package invoicesservice

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/invoices-service/internal/pkg/common"
	"github.com/fidesy-pay/invoices-service/internal/pkg/dto"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	admin_service "github.com/fidesy-pay/invoices-service/pkg/admin-service"
	coingecko_api "github.com/fidesy-pay/invoices-service/pkg/coingecko-api"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/fidesy/sdk/common/logger"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"strings"
	"sync"
	"time"
)

type (
	Service struct {
		storage             Storage
		cryptoServiceClient CryptoServiceClient
		coinGeckoAPIClient  CoinGeckoAPIClient
		adminServiceClient  AdminServiceClient
	}

	CryptoServiceClient interface {
		AcceptCrypto(ctx context.Context, in *crypto_service.AcceptCryptoRequest, opts ...grpc.CallOption) (*crypto_service.AcceptCryptoResponse, error)
		CancelAcceptingCrypto(ctx context.Context, in *crypto_service.CancelAcceptingCryptoRequest, opts ...grpc.CallOption) (*crypto_service.CancelAcceptingCryptoResponse, error)
		Transfer(ctx context.Context, in *crypto_service.TransferRequest, opts ...grpc.CallOption) (*crypto_service.TransferResponse, error)
	}

	CoinGeckoAPIClient interface {
		GetPrice(ctx context.Context, in *coingecko_api.GetPriceRequest, opts ...grpc.CallOption) (*coingecko_api.GetPriceResponse, error)
	}

	AdminServiceClient interface {
		SendAlert(ctx context.Context, in *admin_service.SendAlertRequest, opts ...grpc.CallOption) (*admin_service.SendAlertResponse, error)
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
	cryptoServiceClient CryptoServiceClient,
	coinGeckoAPIClient CoinGeckoAPIClient,
	adminServiceClient AdminServiceClient,
) *Service {
	service := &Service{
		storage:             storage,
		cryptoServiceClient: cryptoServiceClient,
		coinGeckoAPIClient:  coinGeckoAPIClient,
		adminServiceClient:  adminServiceClient,
	}

	go service.cleanExpiredInvoicesWorker(ctx)
	go service.transferWorker(ctx)

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

	if len(reqFilter.InvoiceStatusIn) > 0 {
		filter.StatusIn = reqFilter.InvoiceStatusIn
	}

	invoices, err := s.storage.ListInvoices(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("storage.ListInvoices: %w", err)
	}

	return invoices, nil
}

func (s *Service) UpdateInvoiceStatus(ctx context.Context, params dto.UpdateInvoiceStatusParams) error {
	filter := storage.ListInvoicesFilter{
		IDIn: []uuid.UUID{params.InvoiceID},
	}

	invoices, err := s.storage.ListInvoices(ctx, filter)
	if err != nil {
		return fmt.Errorf("storage.ListInvoices: %w", err)
	}

	if len(invoices) == 0 {
		return ErrInvoiceNotFoundByID(params.InvoiceID)
	}

	invoice := invoices[0]
	invoice.Status = params.Status

	_, err = s.storage.UpdateInvoice(ctx, invoice)
	if err != nil {
		return fmt.Errorf("storage.UpdateInvoice: %w", err)
	}

	return nil
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

func (s *Service) transferWorker(ctx context.Context) {
	transferFunc := s.transferCallback()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(5 * time.Second):
			go transferFunc(ctx)
		}
	}
}

func (s *Service) transferCallback() func(ctx context.Context) {
	invoiceRetries := map[uuid.UUID]int{}
	var mu sync.RWMutex
	defaultStep := 2500
	maxRetries := 10

	return func(ctx context.Context) {
		invoices, err := s.storage.ListInvoices(ctx, storage.ListInvoicesFilter{
			StatusIn: []desc.InvoiceStatus{desc.InvoiceStatus_SENDING_TO_CLIENT},
		})
		if err != nil {
			logger.Errorf("transferWorker: storage.ListInvoices: %w", err)
			return
		}

		for _, invoice := range invoices {
			mu.RLock()
			retriesAmount := invoiceRetries[invoice.ID]
			mu.RUnlock()

			_, err = s.cryptoServiceClient.Transfer(ctx, &crypto_service.TransferRequest{
				ClientId:  invoice.ClientID.String(),
				InvoiceId: invoice.ID.String(),
				GasLimit:  lo.ToPtr(uint64(50000 + defaultStep*retriesAmount)),
			})
			if err != nil {
				logger.Errorf("cryptoServiceClient.Transfer: %v", err)

				if retriesAmount < maxRetries {
					mu.Lock()
					invoiceRetries[invoice.ID]++
					mu.Unlock()
					return
				}

				_, err = s.adminServiceClient.SendAlert(ctx, &admin_service.SendAlertRequest{
					Message: alertMessage(invoice),
				})
				if err != nil {
					logger.Errorf("alertsServiceClient.SendAlert: %w", err)
					return
				}

				invoice.Status = desc.InvoiceStatus_MANUAL_CONTROL
				_, err = s.storage.UpdateInvoice(ctx, invoice)
				if err != nil {
					logger.Errorf("storage.UpdateInvoice: %v", err)
					return
				}

				return
			}

			invoice.Status = desc.InvoiceStatus_SUCCESS
			_, err = s.storage.UpdateInvoice(ctx, invoice)
			if err != nil {
				logger.Errorf("storage.UpdateInvoice: %v", err)
				return
			}
		}
	}
}

func alertMessage(invoice *models.Invoice) string {
	return fmt.Sprintf("MANUAL_CONTROL\nInvoiceID: %s \nClientID: %s \nTotal: $%.2f \nChain: %s",
		invoice.ID,
		invoice.ClientID,
		float64(invoice.UsdCentsAmount)/100,
		invoice.Chain,
	)
}
