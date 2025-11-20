package handlers

import (
	"CRMBackendProject/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Create an in-memory database (exists only during test)
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %w", err)
	}

	// Create the customers table
	createTableSQL := `
	CREATE TABLE customers (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		role TEXT,
		email TEXT,
		phone TEXT,
		contacted BOOLEAN
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create customers table: %w", err)
	}
	return db
}

func TestRetrieveAllCustomers(t *testing.T) {
	// 1. Set up: Create a test database
	db := setupTestDB(t)
	defer db.Close() // This runs at the end of the function

	// 2. Insert test data (2 customers)
	insertSQL := `INSERT INTO customers (id, name, role, email, phone, contacted)
	              VALUES (?, ?, ?, ?, ?, ?)`

	db.Exec(insertSQL, "1", "John Doe", "Developer", "john@example.com", "555-0001", false)
	db.Exec(insertSQL, "2", "Jane Smith", "Manager", "jane@example.com", "555-0002", true)

	// 3. Create your handler with the test database
	h := NewHandlers(db)

	// 4. Create a fake HTTP GET request
	req := httptest.NewRequest("GET", "/customers", nil)

	// 5. Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// 6. Call your handler (this is where it runs!)
	h.RetrieveAllCustomers(rr, req)

	// 7. Check the HTTP status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	// 8. Parse the JSON response
	var customers []models.Customer
	json.NewDecoder(rr.Body).Decode(&customers)

	// 9. Verify we got 2 customers back
	if len(customers) != 2 {
		t.Errorf("Expected 2 customers, got %d", len(customers))
	}

	// 10. Check the data is correct
	if customers[0].Name != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", customers[0].Name)
	}

	if customers[1].Role != "Manaager" {
		t.Errorf("Expected 'Manager', got '%s'", customers[1].Role)
	}
}
