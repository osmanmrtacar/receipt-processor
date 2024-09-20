package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	Db *sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &SQLiteDB{Db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS receipts (
			id TEXT PRIMARY KEY,
			retailer TEXT NOT NULL,
			purchase_date TEXT NOT NULL,
			purchase_time TEXT NOT NULL,
			total REAL NOT NULL,
			points INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			receipt_id TEXT,
			short_description TEXT NOT NULL,
			price REAL NOT NULL,
			FOREIGN KEY (receipt_id) REFERENCES receipts(id)
		);
	`)
	return err
}

func (r *SQLiteDB) Close() error {
	return r.Db.Close()
}
