package storage

import (
	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
	postgres "github.com/fidesyx/platform/pkg/scratch/storage"
)

var (
	invoicesTable = (&models.Invoice{}).TableName()
	invoiceFields = postgres.ModelColumns(&models.Invoice{})
)
