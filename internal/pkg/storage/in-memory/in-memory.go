package in_memory

import (
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	"github.com/google/uuid"
	"sync"
)

type Storage struct {
	mu       sync.RWMutex
	invoices map[uuid.UUID]*models.Invoice
}

func New() *Storage {
	return &Storage{
		invoices: make(map[uuid.UUID]*models.Invoice),
	}
}
