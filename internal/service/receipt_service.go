package service

import (
	"github.com/osmanmrtacar/receipt-processor/internal/mapper"
	"github.com/osmanmrtacar/receipt-processor/internal/repository"
	"github.com/osmanmrtacar/receipt-processor/pkg/dto"
	"github.com/osmanmrtacar/receipt-processor/pkg/logic"
)

type ReceiptService struct {
	repo repository.ReceiptRepository
}

func NewReceiptService(repo repository.ReceiptRepository) *ReceiptService {
	return &ReceiptService{repo: repo}
}

func (s *ReceiptService) ProcessReceipt(receiptDto *dto.ReceiptRequestDto) (string, error) {
	receipt, err := mapper.MapReceiptDtoToModel(*receiptDto)
	if err != nil {
		return "", err
	}

	receipt.Points = logic.CalculatePoints(&receipt)

	if err := s.repo.SaveReceipt(&receipt); err != nil {
		return "", err
	}

	return receipt.ID, nil
}

func (s *ReceiptService) GetPoints(id string) (int, error) {
	receipt, err := s.repo.GetReceiptByID(id)

	if err != nil {
		return 0, err
	}

	return receipt.Points, nil
}
