package repository

import (
	"github.com/osmanmrtacar/receipt-processor/internal/models"
)

type ReceiptRepository interface {
	SaveReceipt(receipt *models.Receipt) error
	GetReceiptByID(string) (*models.Receipt, error)
	Close() error
}
