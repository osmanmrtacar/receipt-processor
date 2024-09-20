package api

import (
	"github.com/gorilla/mux"
	"github.com/osmanmrtacar/receipt-processor/internal/service"
)

func NewRouter(receiptService *service.ReceiptService) *mux.Router {
	router := mux.NewRouter()
	handlers := NewHandlers(receiptService)

	router.HandleFunc("/receipts/process", handlers.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handlers.GetPoints).Methods("GET")

	return router
}
