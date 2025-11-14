package main

import (
	"CRMBackendProject/internal/database"
	"CRMBackendProject/internal/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	StartServer()
}

func StartServer() {
	// Initialize database
	db, err := initializeDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)
	h := handlers.NewHandlers(db)
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ShowHomePage)

	router.HandleFunc("/customers", h.RetrieveAllCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", h.RetrieveSingleCustomer).Methods("GET")
	router.HandleFunc("/customers", h.CreateNewCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", h.DeleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{id}", h.UpdateCustomer).Methods("PUT")
	fmt.Println("Server starting on port 3000")
	http.ListenAndServe(":3000", router)
}

func initializeDatabase() (*sql.DB, error) {
	db, err := database.InitDB("crm.db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	err = database.CreateTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}
