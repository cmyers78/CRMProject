package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes and returns a connection to the SQLite database instead of a global variable
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("connected to SQLite database at %s", filepath)
	return db, nil
}

func CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func CreateTables(db *sql.DB) error {
	createCustomerTable := `CREATE TABLE IF NOT EXISTS customers (
    	id TEXT PRIMARY KEY,
    	name TEXT NOT NULL,
    	role TEXT NOT NULL,
    	email TEXT NOT NULL UNIQUE,
    	phone TEXT NOT NULL,
    	contacted BOOLEAN NOT NULL DEFAULT FALSE
	);`
	_, err := db.Exec(createCustomerTable)
	if err != nil {
		return fmt.Errorf("failed to create customer table: %w", err)
	}
	log.Println("created customer table")
	return nil
}
