package invoicesservice

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fidesy-pay/invoices-service/internal/config"
	"github.com/fidesy-pay/invoices-service/internal/pkg/common"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	external_api "github.com/fidesy-pay/invoices-service/pkg/external-api"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/fidesy/sdk/common/logger"
	"github.com/fidesy/sdk/common/postgres"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

type (
	Service struct {
		storage             Storage
		cryptoServiceClient CryptoServiceClient
		externalAPI         ExternalAPI
	}

	CryptoServiceClient interface {
		AcceptCrypto(ctx context.Context, in *crypto_service.AcceptCryptoRequest, opts ...grpc.CallOption) (*crypto_service.AcceptCryptoResponse, error)
		Transfer(ctx context.Context, in *crypto_service.TransferRequest, opts ...grpc.CallOption) (*crypto_service.TransferResponse, error)
	}

	ExternalAPI interface {
		GetPrice(ctx context.Context, in *external_api.GetPriceRequest, opts ...grpc.CallOption) (*external_api.GetPriceResponse, error)
	}

	Storage interface {
		CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
		ListInvoices(ctx context.Context, filter storage.ListInvoicesFilter, pagination postgres.Pagination) ([]*models.Invoice, error)
		UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	}
)

func New(
	ctx context.Context,
	storage Storage,
	cryptoServiceClient CryptoServiceClient,
	externalAPI ExternalAPI,
) *Service {
	service := &Service{
		storage:             storage,
		cryptoServiceClient: cryptoServiceClient,
		externalAPI:         externalAPI,
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
	invoices, err := s.storage.ListInvoices(
		ctx,
		storage.ListInvoicesFilter{
			IDIn: []uuid.UUID{input.InvoiceID},
		},
		postgres.NewPagination(1, 1),
	)
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

	tokenPriceResp, err := s.externalAPI.GetPrice(ctx, &external_api.GetPriceRequest{
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
	invoice.PayerClientID = input.PayerClientID

	invoice, err = s.storage.UpdateInvoice(ctx, invoice)
	if err != nil {
		return nil, fmt.Errorf("storage.CreateInvoice: %w", err)
	}

	return invoice, nil
}

func (s *Service) CheckInvoice(ctx context.Context, invoiceIDStr string) (*models.Invoice, error) {
	// we validate ID in handler logic
	invoiceID := uuid.MustParse(invoiceIDStr)

	invoices, err := s.storage.ListInvoices(
		ctx,
		storage.ListInvoicesFilter{
			IDIn: []uuid.UUID{invoiceID},
		},
		postgres.NewPagination(1, 1),
	)
	if err != nil {
		return nil, fmt.Errorf("storage.ListInvoices: %w", err)
	}

	if len(invoices) == 0 {
		return nil, ErrInvoiceNotFoundByID(invoiceID)
	}

	return invoices[0], nil
}

func (s *Service) ListInvoices(ctx context.Context, req *desc.ListInvoicesRequest) ([]*models.Invoice, error) {
	var err error

	reqFilter := req.GetFilter()

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

	invoices, err := s.storage.ListInvoices(ctx, filter, postgres.NewPagination(req.GetPage(), req.GetPerPage()))
	if err != nil {
		return nil, fmt.Errorf("storage.ListInvoices: %w", err)
	}

	return invoices, nil
}

func (s *Service) cleanExpiredInvoicesWorker(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go s.cleanExpiredInvoices(ctx)
		}
	}
}

func (s *Service) cleanExpiredInvoices(ctx context.Context) {
	ctx = context.WithValue(ctx, "skip_span", true)

	invoices, err := s.storage.ListInvoices(
		ctx,
		storage.ListInvoicesFilter{
			StatusIn:    []desc.InvoiceStatus{desc.InvoiceStatus_NEW, desc.InvoiceStatus_PENDING},
			CreatedAtLt: lo.ToPtr(time.Now().Add(-config.Get(config.ExpireInterval).(time.Duration))),
		},
		postgres.NewPagination(1, 100),
	)
	if err != nil {
		logger.Errorf("cleanExpiredInvoices: storage.ListInvoices: %w", err)
		return
	}

	for _, invoice := range invoices {
		invoice.Status = desc.InvoiceStatus_EXPIRED
		_, err = s.storage.UpdateInvoice(ctx, invoice)
		if err != nil {
			logger.Errorf("cleanExpiredInvoices: storage.UpdateInvoice: %w", err)
		}
	}
}

func (s *Service) transferWorker(ctx context.Context) {
	ctx = context.WithValue(ctx, "skip_span", true)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	transferFunc := s.transferCallback()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go transferFunc(ctx)
		}
	}
}

func (s *Service) transferCallback() func(ctx context.Context) {
	var (
		locks = make(map[string]struct{})
		mu    sync.RWMutex
	)

	return func(ctx context.Context) {
		invoices, err := s.storage.ListInvoices(
			ctx,
			storage.ListInvoicesFilter{
				StatusIn: []desc.InvoiceStatus{desc.InvoiceStatus_SENDING_TO_CLIENT},
			},
			postgres.NewPagination(1, 100),
		)
		if err != nil {
			logger.Errorf("transferWorker: storage.ListInvoices: %w", err)
			return
		}

		for _, invoice := range invoices {
			invoice := invoice

			mu.RLock()
			_, ok := locks[invoice.ID.String()]
			mu.RUnlock()
			if ok {
				continue
			}

			mu.Lock()
			locks[invoice.ID.String()] = struct{}{}
			mu.Unlock()

			go func() {
				defer func() {
					mu.Lock()
					delete(locks, invoice.ID.String())
					mu.Unlock()
				}()

				s.completeInvoice(ctx, invoice)
			}()
		}
	}
}

func (s *Service) completeInvoice(ctx context.Context, invoice *models.Invoice) {
	defaultStep := 50000

	for i := 0; i < 10; i++ {
		gasLimit := uint64(50000 + defaultStep*i)
		if invoice.GasLimit != nil {
			gasLimit = uint64(*invoice.GasLimit)
		}

		_, err := s.cryptoServiceClient.Transfer(ctx, &crypto_service.TransferRequest{
			ClientId:  invoice.ClientID.String(),
			InvoiceId: lo.ToPtr(invoice.ID.String()),
			GasLimit:  lo.ToPtr(gasLimit),
		})
		if err != nil {
			continue
		}

		invoice.Status = desc.InvoiceStatus_SUCCESS
		_, err = s.storage.UpdateInvoice(ctx, invoice)
		if err != nil {
			logger.Errorf("storage.UpdateInvoice: %v", err)
			return
		}

		return
	}

	invoice.Status = desc.InvoiceStatus_MANUAL_CONTROL
	_, err := s.storage.UpdateInvoice(ctx, invoice)
	if err != nil {
		logger.Errorf("storage.UpdateInvoice: %v", err)
		return
	}
}
