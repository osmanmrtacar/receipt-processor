package main

import (
	"log"
	"net/http"

	"github.com/osmanmrtacar/receipt-processor/internal/api"
	"github.com/osmanmrtacar/receipt-processor/internal/config"
	"github.com/osmanmrtacar/receipt-processor/internal/database"
	"github.com/osmanmrtacar/receipt-processor/internal/repository"
	"github.com/osmanmrtacar/receipt-processor/internal/service"
)

func main() {
	config := config.LoadConfig()
	db, err := database.NewSQLiteDB(config.DatabaseURI)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	var repo repository.ReceiptRepository
	repo, err = repository.NewSQLiteRepository(db.Db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	receiptService := service.NewReceiptService(repo)
	router := api.NewRouter(receiptService)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
