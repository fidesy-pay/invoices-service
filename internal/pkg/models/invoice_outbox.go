package models

type InvoiceOutbox struct {
	ID      int    `db:"id"`
	Message string `db:"message"`
}

func (io *InvoiceOutbox) TableName() string {
	return "invoices_outbox"
}
