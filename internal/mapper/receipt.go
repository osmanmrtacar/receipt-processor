package mapper

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/osmanmrtacar/receipt-processor/internal/models"
	"github.com/osmanmrtacar/receipt-processor/pkg/dto"
)

func MapReceiptDtoToModel(dto dto.ReceiptRequestDto) (models.Receipt, error) {
	total, err := strconv.ParseFloat(dto.Total, 64)
	if err != nil {
		return models.Receipt{}, err
	}

	receiptItems := make([]models.Item, len(dto.Items))

	for i, item := range dto.Items {
		parsedPrice, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return models.Receipt{}, err
		}
		receiptItems[i] = models.Item{
			ShortDescription: item.ShortDescription,
			Price:            parsedPrice,
		}
	}

	return models.Receipt{
		ID:           uuid.New().String(),
		Retailer:     dto.Retailer,
		PurchaseDate: dto.PurchaseDate,
		PurchaseTime: dto.PurchaseTime,
		Total:        total,
		Items:        receiptItems,
	}, nil
}
