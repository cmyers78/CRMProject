package main

import (
	"CRMBackendProject/internal/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	StartServer()
}

func StartServer() {

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ShowHomePage)

	router.HandleFunc("/customers", handlers.GetAllCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", handlers.GetSingleCustomer).Methods("GET")
	//router.HandleFunc("/customers", handlers.CreateNewCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", handlers.DeleteCustomer).Methods("DELETE")
	//router.HandleFunc("/customers/{id}", handlers.UpdateCustomer).Methods("PUT")
	fmt.Println("Server starting on port 3000")
	http.ListenAndServe(":3000", router)
}
