package main

import (
	"CRMBackendProject/internal/database"
	"CRMBackendProject/internal/handlers"
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
	if err := initializeDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ShowHomePage)

	router.HandleFunc("/customers", handlers.RetrieveAllCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", handlers.RetrieveSingleCustomer).Methods("GET")
	router.HandleFunc("/customers", handlers.CreateNewCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", handlers.DeleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{id}", handlers.UpdateCustomer).Methods("PUT")
	fmt.Println("Server starting on port 3000")
	http.ListenAndServe(":3000", router)
}

func initializeDatabase() error {
	err := database.InitDB("crm.db")
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	err = database.CreateTables()
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}
