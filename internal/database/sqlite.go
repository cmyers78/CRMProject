package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("connected to SQLite database at %s", filepath)
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func CreateTables() error {
	createCustomerTable := `CREATE TABLE IF NOT EXISTS customers (
    	id TEXT PRIMARY KEY,
    	name TEXT NOT NULL,
    	role TEXT NOT NULL,
    	email TEXT NOT NULL UNIQUE,
    	phone TEXT NOT NULL,
    	contacted BOOLEAN NOT NULL DEFAULT FALSE
	);`
	_, err := DB.Exec(createCustomerTable)
	if err != nil {
		return fmt.Errorf("failed to create customer table: %w", err)
	}
	log.Println("created customer table")
	return nil
}
