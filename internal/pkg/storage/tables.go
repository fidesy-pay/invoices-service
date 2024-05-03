package storage

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fidesy-pay/invoices-service/internal/pkg/models"
)

var (
	invoicesTable = (&models.Invoice{}).TableName()
	invoiceFields = modelColumns(&models.Invoice{})
)

type Model interface {
	TableName() string
}

func modelColumns(m Model) string {
	dbTags := make([]string, 0)

	dataType := reflect.TypeOf(m)

	if dataType.Kind() != reflect.Pointer {
		return ""
	}

	data := dataType.Elem()

	for i := 0; i < data.NumField(); i++ {
		field := data.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		dbTags = append(dbTags, fmt.Sprintf("%s.%s", m.TableName(), dbTag))
	}

	return strings.Join(dbTags, ",")
}
