package outbox_processor

import (
	"context"
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	"github.com/fidesy/sdk/common/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"time"
)

const invoicesTopic = "invoices-json"

type (
	Service struct {
		storage       Storage
		kafkaProducer KafkaProducer
	}

	Storage interface {
		ListInvoicesOutbox(ctx context.Context, limit uint64) ([]*models.InvoiceOutbox, error)
		DeleteInvoicesOutbox(ctx context.Context, userIDs []int) error
	}

	KafkaProducer interface {
		ProduceMessage(topic string, messageBytes []byte)
	}
)

func New(storage Storage, producer KafkaProducer) *Service {
	return &Service{
		storage:       storage,
		kafkaProducer: producer,
	}
}

func (s *Service) Publish(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(500 * time.Millisecond):
			s.publish(ctx)
		}
	}
}

func (s *Service) publish(ctx context.Context) {
	invoicesOutbox, err := s.storage.ListInvoicesOutbox(ctx, 100)
	if err != nil {
		logger.Errorf("storage.ListInvoicesOutbox: %v", err)
		return
	}

	if len(invoicesOutbox) == 0 {
		return
	}

	for _, invoice := range invoicesOutbox {
		s.kafkaProducer.ProduceMessage(invoicesTopic, []byte(invoice.Message))
	}

	err = s.storage.DeleteInvoicesOutbox(ctx, lo.Map(invoicesOutbox, func(invoice *models.InvoiceOutbox, _ int) int {
		return invoice.ID
	}))
	if err != nil {
		logger.Errorf("storage.DeleteInvoicesOutbox: %v", err,
			zap.Int("invoices_length:", len(invoicesOutbox)),
		)
	}
}