package logic

import (
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/osmanmrtacar/receipt-processor/internal/models"
)

func CalculatePoints(receipt *models.Receipt) int {
	points := 0

	points += calculateRetailerNamePoints(receipt.Retailer)

	points += calculateTotalPoints(receipt.Total)

	points += calculateItemsPoints(receipt.Items)

	points += calculateDatePoints(receipt.PurchaseDate)

	points += calculateTimePoints(receipt.PurchaseTime)

	return points
}

func calculateRetailerNamePoints(retailer string) int {
	return len(strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, retailer))
}

func calculateTotalPoints(total float64) int {
	points := 0
	if total == float64(int(total)) {
		points += 50
	}
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}
	return points
}

func calculateItemsPoints(items []models.Item) int {
	points := 0
	points += (len(items) / 2) * 5

	for _, item := range items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	return points
}

func calculateDatePoints(purchaseDate string) int {
	date, err := time.Parse("2006-01-02", purchaseDate)
	if err != nil {
		return 0
	}
	if date.Day()%2 != 0 {
		return 6
	}
	return 0
}

func calculateTimePoints(purchaseTime string) int {
	t, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		return 0
	}
	if t.Hour() >= 14 && t.Hour() < 16 {
		return 10
	}
	return 0
}
