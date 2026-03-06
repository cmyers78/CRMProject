package handlers

import (
	"bytes"
	"CRMBackendProject/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Create an in-memory database (exists only during test)
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
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
		t.Fatalf("Failed to create customers table: %v", err)
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

	if customers[1].Role != "Manager" {
		t.Errorf("Expected 'Manager', got '%s'", customers[1].Role)
	}
}

func TestRetrieveOneCustomer(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	insertSQL := `INSERT INTO customers (id, name, role, email, phone, contacted)
	              VALUES (?, ?, ?, ?, ?, ?)`

	db.Exec(insertSQL, "1", "John Doe", "Developer", "john@example.com", "555-0001", false)

	// 3. Create your handler with the test database
	h := NewHandlers(db)

	// 4. Create a fake HTTP GET request
	req := httptest.NewRequest("GET", "/customers/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// 5. Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// 6. Call your handler (this is where it runs!)
	h.RetrieveSingleCustomer(rr, req)

	// 7. Check the HTTP status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// 8. Parse the JSON response
	var customer models.Customer
	json.NewDecoder(rr.Body).Decode(&customer)

	// 10. Check the data is correct
	if customer.Name != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", customer.Name)
	}
}

func TestDeleteCustomer(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	insertSQL := `INSERT INTO customers (id, name, role, email, phone, contacted)
	              VALUES (?, ?, ?, ?, ?, ?)`

	db.Exec(insertSQL, "1", "John Doe", "Developer", "john@example.com", "555-0001", false)

	// 3. Create your handler with the test database
	h := NewHandlers(db)

	req := httptest.NewRequest("DELETE", "/customers/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	h.DeleteCustomer(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Parse the JSON response (should return the deleted customer)
	var deletedCustomer models.Customer
	json.NewDecoder(rr.Body).Decode(&deletedCustomer)

	// Verify the response contains the customer that was deleted
	if deletedCustomer.Name != "John Doe" {
		t.Errorf("Expected deleted customer name 'John Doe', got '%s'", deletedCustomer.Name)
	}

	// Verify the customer is actually deleted from the database
	var count int
	db.QueryRow("SELECT COUNT(*) FROM customers WHERE id = ?", "1").Scan(&count)

	if count != 0 {
		t.Errorf("Expected customer to be deleted, but it still exists")
	}
}

func TestCreateNewCustomer(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    h := NewHandlers(db)

    // Create the customer data to send
    newCustomer := models.Customer{
        Name:  "Jane Doe",
        Role:  "Tester",
        Email: "jane@example.com",
        Phone: "555-9999",
    }

    // Convert to JSON
    jsonData, _ := json.Marshal(newCustomer)

    // Create POST request with JSON body
    req := httptest.NewRequest("POST", "/customers", bytes.NewBuffer(jsonData))
    rr := httptest.NewRecorder()

    h.CreateNewCustomer(rr, req)

    // Check status is 201 Created
    if rr.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", rr.Code)
    }

    // Parse response
    var createdCustomer models.Customer
    json.NewDecoder(rr.Body).Decode(&createdCustomer)

    // Verify it has an ID
    if createdCustomer.ID == "" {
        t.Errorf("Expected customer to have an ID, got empty string")
    }

    // Verify name matches
    if createdCustomer.Name != "Jane Doe" {
        t.Errorf("Expected name 'Jane Doe', got '%s'", createdCustomer.Name)
    }

    // Query database to verify it was actually inserted
	var dbCustomer models.Customer

	row := db.QueryRow("SELECT id, name, role, email, phone, contacted FROM customers WHERE id = ?", createdCustomer.ID)
	err := row.Scan(&dbCustomer.ID, &dbCustomer.Name, &dbCustomer.Role, &dbCustomer.Email, &dbCustomer.Phone, &dbCustomer.Contacted)

	if err != nil {
		t.Errorf("Customer was not inserted into database: %v", err)
	}
	
	if dbCustomer.Name != "Jane Doe" {
    	t.Errorf("Database customer name mismatch: expected 'Jane Doe', got '%s'", dbCustomer.Name)
	}
	if dbCustomer.Role != "Tester" {
    	t.Errorf("Database customer role mismatch: expected 'Manager', got '%s'", dbCustomer.Role)
	}
}

func TestUpdateCustomer(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := NewHandlers(db)

	initialCustomer := models.Customer{
    Name:      "Old Name",
    Role:      "Old Role",
    Email:     "old@example.com",
    Phone:     "111-1111",
    Contacted: false,
	}
	customerJSON, _ := json.Marshal(initialCustomer)
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(customerJSON))
	rr := httptest.NewRecorder()
	h.CreateNewCustomer(rr, req)

	var createdCustomer models.Customer
	json.NewDecoder(rr.Body).Decode(&createdCustomer)

	updatedCustomer := models.Customer {
		ID: createdCustomer.ID,
		Name: "New Name",
		Role: "New Role",
		Email: "new@example.com",
		Phone: "222-2222",
		Contacted: true,
	}

	updateJson, _ := json.Marshal(updatedCustomer)
	updateReq := httptest.NewRequest(http.MethodPut, "/customers" + updatedCustomer.ID, bytes.NewBuffer(updateJson))
	updateRR := httptest.NewRecorder()
	h.UpdateCustomer(updateRR, updateReq)

	if updateRR.Code != http.StatusAccepted {
		t.Errorf("Expected status 202, got %d", updateRR.Code)
	}

	var dbCustomer models.Customer
	row := db.QueryRow("SELECT id, name, role, email, phone, contacted FROM customers WHERE id = ?", createdCustomer.ID)
	err := row.Scan(&dbCustomer.ID, &dbCustomer.Name, &dbCustomer.Role, &dbCustomer.Email, &dbCustomer.Phone, &dbCustomer.Contacted)
	if err != nil {
    	t.Errorf("Could not find customer in database: %v", err)
	}

	if dbCustomer.Name != "New Name" {
    	t.Errorf("Expected name 'New Name', got '%s'", dbCustomer.Name)
	}
	if dbCustomer.Role != "New Role" {
    	t.Errorf("Expected role 'New Role', got '%s'", dbCustomer.Role)
	}
	if dbCustomer.Email != "new@example.com" {
    	t.Errorf("Expected email 'new@example.com', got '%s'", dbCustomer.Email)
	}
	if !dbCustomer.Contacted {
    	t.Errorf("Expected Contacted to be true")
	}
}