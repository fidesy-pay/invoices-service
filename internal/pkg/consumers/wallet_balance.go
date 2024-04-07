package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	invoicesservice "github.com/fidesy-pay/invoices-service/internal/pkg/invoices-service"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	"github.com/fidesy-pay/invoices-service/internal/pkg/storage"
	crypto_service "github.com/fidesy-pay/invoices-service/pkg/crypto-service"
	desc "github.com/fidesy-pay/invoices-service/pkg/invoices-service"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"strings"
)

type (
	WalletBalanceConsumer struct {
		storage             Storage
		cryptoServiceClient CryptoServiceClient
	}

	Storage interface {
		ListInvoices(ctx context.Context, filter storage.ListInvoicesFilter) ([]*models.Invoice, error)
		UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	}

	CryptoServiceClient interface {
		Transfer(ctx context.Context, in *crypto_service.TransferRequest, opts ...grpc.CallOption) (*crypto_service.TransferResponse, error)
	}
)

func NewWalletBalanceConsumer(
	storage Storage,
	cryptoServiceClient CryptoServiceClient,
) *WalletBalanceConsumer {
	return &WalletBalanceConsumer{
		storage:             storage,
		cryptoServiceClient: cryptoServiceClient,
	}
}

func (c *WalletBalanceConsumer) Consume(ctx context.Context, msg []byte) error {
	wallet := new(models.WalletMessage)
	err := json.Unmarshal(msg, &wallet)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	invoices, err := c.storage.ListInvoices(ctx, storage.ListInvoicesFilter{
		AddressIn: []string{strings.ToLower(wallet.Address)},
	})
	if err != nil {
		return fmt.Errorf("storage.ListInvoices: %v", err)
	}

	if len(invoices) == 0 {
		return invoicesservice.ErrInvoiceNotFoundByAddress(wallet.Address)
	}

	invoice := invoices[0]
	if wallet.Chain != invoice.Chain || wallet.Token != invoice.Token {
		return nil
	}

	if wallet.Balance < int64(*invoice.TokenAmount*1e18) {
		return nil
	}

	invoice.Status = desc.InvoiceStatus_SENDING_TO_CLIENT
	_, err = c.storage.UpdateInvoice(ctx, invoice)
	if err != nil {
		return fmt.Errorf("storage.UpdateInvoice: %v", err)
	}

	_, err = c.cryptoServiceClient.Transfer(ctx, &crypto_service.TransferRequest{
		ClientId:  invoice.ClientID.String(),
		InvoiceId: lo.ToPtr(invoice.ID.String()),
	})
	if err != nil {
		return fmt.Errorf("cryptoServiceClient.Transfer: %v", err)
	}

	invoice.Status = desc.InvoiceStatus_SUCCESS
	_, err = c.storage.UpdateInvoice(ctx, invoice)
	if err != nil {
		return fmt.Errorf("storage.UpdateInvoice: %v", err)
	}

	return nil
}
