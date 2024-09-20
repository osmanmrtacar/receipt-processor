package service

import (
	"errors"
	"testing"

	"github.com/osmanmrtacar/receipt-processor/internal/models"
	"github.com/osmanmrtacar/receipt-processor/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReceiptRepository struct {
	mock.Mock
}

func (m *MockReceiptRepository) SaveReceipt(receipt *models.Receipt) error {
	args := m.Called(receipt)
	return args.Error(0)
}

func (m *MockReceiptRepository) GetReceiptByID(id string) (*models.Receipt, error) {
	args := m.Called(id)
	if receipt, ok := args.Get(0).(*models.Receipt); ok {
		return receipt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReceiptRepository) Close() error {
	return nil
}

func createSampleReceiptDTO() *dto.ReceiptRequestDto {
	return &dto.ReceiptRequestDto{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []dto.ReceiptItem{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
}

func TestProcessReceipt_Success(t *testing.T) {

	mockRepo := new(MockReceiptRepository)

	receiptDTO := createSampleReceiptDTO()

	expectedReceipt := models.Receipt{
		ID:           "unique-id-123",
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: 6.49},
			{ShortDescription: "Emils Cheese Pizza", Price: 12.25},
			{ShortDescription: "Knorr Creamy Chicken", Price: 1.26},
			{ShortDescription: "Doritos Nacho Cheese", Price: 3.35},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: 12.00},
		},
		Total:  35.35,
		Points: 28,
	}

	mockRepo.On("SaveReceipt", mock.AnythingOfType("*models.Receipt")).Return(nil).Run(func(args mock.Arguments) {
		receipt := args.Get(0).(*models.Receipt)

		receipt.ID = "unique-id-123"
	})

	service := NewReceiptService(mockRepo)

	receiptID, err := service.ProcessReceipt(receiptDTO)

	assert.NoError(t, err)
	assert.Equal(t, "unique-id-123", receiptID)

	mockRepo.AssertNumberOfCalls(t, "SaveReceipt", 1)

	mockRepo.AssertCalled(t, "SaveReceipt", mock.MatchedBy(func(receipt *models.Receipt) bool {
		return receipt.Retailer == expectedReceipt.Retailer &&
			receipt.PurchaseDate == expectedReceipt.PurchaseDate &&
			receipt.PurchaseTime == expectedReceipt.PurchaseTime &&
			len(receipt.Items) == len(expectedReceipt.Items) &&
			receipt.Total == expectedReceipt.Total &&
			receipt.Points == expectedReceipt.Points
	}))
}

func TestProcessReceipt_MappingError(t *testing.T) {
	mockRepo := new(MockReceiptRepository)

	receiptDTO := &dto.ReceiptRequestDto{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []dto.ReceiptItem{
			{ShortDescription: "Mountain Dew 12PK", Price: "invalid-price"},
		},
		Total: "35.35",
	}

	service := NewReceiptService(mockRepo)

	receiptID, err := service.ProcessReceipt(receiptDTO)

	assert.Error(t, err)
	assert.Equal(t, "", receiptID)

	mockRepo.AssertNotCalled(t, "SaveReceipt", mock.Anything)
}

func TestProcessReceipt_SaveReceiptError(t *testing.T) {

	mockRepo := new(MockReceiptRepository)

	receiptDTO := createSampleReceiptDTO()

	mockRepo.On("SaveReceipt", mock.AnythingOfType("*models.Receipt")).Return(errors.New("database error"))

	service := NewReceiptService(mockRepo)

	receiptID, err := service.ProcessReceipt(receiptDTO)

	assert.Error(t, err)
	assert.Equal(t, "", receiptID)
	assert.EqualError(t, err, "database error")

	mockRepo.AssertNumberOfCalls(t, "SaveReceipt", 1)
}

func TestGetPoints_Success(t *testing.T) {

	mockRepo := new(MockReceiptRepository)

	sampleReceipt := &models.Receipt{
		ID:       "unique-id-123",
		Points:   28,
		Retailer: "Target",
	}

	mockRepo.On("GetReceiptByID", "unique-id-123").Return(sampleReceipt, nil)

	service := NewReceiptService(mockRepo)

	points, err := service.GetPoints("unique-id-123")

	assert.NoError(t, err)
	assert.Equal(t, 28, points)

	mockRepo.AssertNumberOfCalls(t, "GetReceiptByID", 1)
}

func TestGetPoints_ReceiptNotFound(t *testing.T) {

	mockRepo := new(MockReceiptRepository)

	mockRepo.On("GetReceiptByID", "non-existent-id").Return(nil, errors.New("receipt not found"))

	service := NewReceiptService(mockRepo)

	points, err := service.GetPoints("non-existent-id")

	assert.Error(t, err)
	assert.Equal(t, 0, points)
	assert.EqualError(t, err, "receipt not found")

	mockRepo.AssertNumberOfCalls(t, "GetReceiptByID", 1)
}
