package in_memory

import (
	"github.com/fidesy-pay/payment-service/internal/pkg/models"
	"github.com/google/uuid"
	"sync"
)

type Storage struct {
	mu       sync.RWMutex
	payments map[uuid.UUID]*models.Payment
}

func New() *Storage {
	return &Storage{
		payments: make(map[uuid.UUID]*models.Payment),
	}
}
