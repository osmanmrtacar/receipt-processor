package repository

import (
	"database/sql"
	"fmt"

	"github.com/osmanmrtacar/receipt-processor/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (*SQLiteRepository, error) {

	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepository) SaveReceipt(receipt *models.Receipt) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO receipts (id, retailer, purchase_date, purchase_time, total, points)
		VALUES (?, ?, ?, ?, ?, ?)
	`, receipt.ID, receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total, receipt.Points)
	if err != nil {
		return err
	}

	for _, item := range receipt.Items {
		_, err = tx.Exec(`
			INSERT INTO items (receipt_id, short_description, price)
			VALUES (?, ?, ?)
		`, receipt.ID, item.ShortDescription, item.Price)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepository) GetReceiptByID(id string) (*models.Receipt, error) {
	receipt := &models.Receipt{ID: id}

	err := r.db.QueryRow(`
		SELECT retailer, purchase_date, purchase_time, total, points
		FROM receipts
		WHERE id = ?
	`, id).Scan(&receipt.Retailer, &receipt.PurchaseDate, &receipt.PurchaseTime, &receipt.Total, &receipt.Points)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("receipt not found")
		}
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT short_description, price
		FROM items
		WHERE receipt_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ShortDescription, &item.Price); err != nil {
			return nil, err
		}
		receipt.Items = append(receipt.Items, item)
	}

	return receipt, nil
}
