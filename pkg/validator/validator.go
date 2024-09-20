package validator

import (
	"fmt"
	"strconv"
	"time"

	"github.com/osmanmrtacar/receipt-processor/pkg/dto"
)

func ValidateReceipt(receipt dto.ReceiptRequestDto) error {

	_, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return fmt.Errorf("invalid purchase date format: %v", err)
	}

	_, err = time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return fmt.Errorf("invalid purchase time format: %v", err)
	}

	_, err = strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return fmt.Errorf("invalid total format: %v", err)
	}

	for i, item := range receipt.Items {
		_, err = strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return fmt.Errorf("invalid price format for item %d: %v", i, err)
		}
	}

	return nil
}
