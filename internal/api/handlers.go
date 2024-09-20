package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/osmanmrtacar/receipt-processor/internal/service"
	"github.com/osmanmrtacar/receipt-processor/pkg/dto"
	"github.com/osmanmrtacar/receipt-processor/pkg/validator"

	"github.com/gorilla/mux"
)

type Handlers struct {
	receiptService *service.ReceiptService
}

func NewHandlers(receiptService *service.ReceiptService) *Handlers {
	return &Handlers{receiptService: receiptService}
}

func (h *Handlers) ProcessReceipt(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		fmt.Println("Empty request body")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		fmt.Println("Failed to read request body:", err)
		return
	}

	var receiptRequest dto.ReceiptRequestDto
	if err := json.Unmarshal(bodyBytes, &receiptRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("json decode error", err)
		return
	}

	err = validator.ValidateReceipt(receiptRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("validation error", err)
		return
	}

	id, err := h.receiptService.ProcessReceipt(&receiptRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("process receipt error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := &dto.ReceiptResponseDto{Id: id}

	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) GetPoints(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	points, err := h.receiptService.GetPoints(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]int{"points": points})
}
