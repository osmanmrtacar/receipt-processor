package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/osmanmrtacar/receipt-processor/internal/database"
	"github.com/osmanmrtacar/receipt-processor/internal/repository"
	"github.com/osmanmrtacar/receipt-processor/internal/service"
	"github.com/osmanmrtacar/receipt-processor/pkg/dto"
	"github.com/stretchr/testify/assert"
)

func SetupRouter(db *sql.DB) *mux.Router {
	repo, _ := repository.NewSQLiteRepository(db)

	receiptService := service.NewReceiptService(repo)

	handlers := NewHandlers(receiptService)

	router := mux.NewRouter()
	router.HandleFunc("/receipts", handlers.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handlers.GetPoints).Methods("GET")
	return router
}

func SampleReceiptRequest() dto.ReceiptRequestDto {
	return dto.ReceiptRequestDto{
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

func TestE2E_ProcessAndGetPoints(t *testing.T) {
	db, err := database.NewSQLiteDB("_test_receipts.db")
	assert.NoError(t, err, "Failed to initialize database")
	defer db.Close()
	defer os.Remove("_test_receipts.db")

	router := SetupRouter(db.Db)

	ts := httptest.NewServer(router)
	defer ts.Close()

	receiptRequest := SampleReceiptRequest()
	requestBody, err := json.Marshal(receiptRequest)
	assert.NoError(t, err, "Failed to marshal receipt request DTO")

	postURL := ts.URL + "/receipts"
	resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "Failed to send POST /receipts request")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 OK for POST /receipts")

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body for POST /receipts")

	var postResponse dto.ReceiptResponseDto
	err = json.Unmarshal(bodyBytes, &postResponse)
	assert.NoError(t, err, "Failed to unmarshal response body for POST /receipts")

	assert.NotEmpty(t, postResponse.Id, "Expected a non-empty receipt ID in response")

	getURL := ts.URL + "/receipts/" + postResponse.Id + "/points"
	getResp, err := http.Get(getURL)
	assert.NoError(t, err, "Failed to send GET /receipts/{id}/points request")
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode, "Expected status code 200 OK for GET /receipts/{id}/points")

	getBodyBytes, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err, "Failed to read response body for GET /receipts/{id}/points")

	var getResponse map[string]int
	err = json.Unmarshal(getBodyBytes, &getResponse)
	assert.NoError(t, err, "Failed to unmarshal response body for GET /receipts/{id}/points")

	points, exists := getResponse["points"]
	assert.True(t, exists, "Response should contain 'points' field")
	assert.Equal(t, 28, points, "Expected points to be 28")
}
