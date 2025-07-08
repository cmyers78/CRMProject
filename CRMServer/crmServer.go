package crmServer

import (
	"CRMBackendProject/Models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var database = seedCustomerDatabase()

func showHomePage(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "./static/static.html")
}

func getAllCustomers(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(database)
}

func getSingleCustomer(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id := params["id"]

	if customer, ok := database[id]; ok {
		// delete the entry, return successful response
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode(customer)
	} else { // 4. If not, return an error, but still return the dictionary
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(database)
	}
}

func createNewCustomer(writer http.ResponseWriter, req *http.Request) {
	// 1. set content-type to JSON
	writer.Header().Set("Content-Type", "application/json")

	// 2. keep track of new entry so that it can be added to dictionary map
	var newEntry map[string]Models.Customer

	// 3. Read the request
	reqBody, _ := io.ReadAll(req.Body)

	// 4. Parse JSON Body
	json.Unmarshal(reqBody, &newEntry)

	// 5. Add new entry to dictionary map if it doesn't already exist
	for key, value := range newEntry {
		// - Respond with conflict if entry exists
		if _, ok := database[key]; ok {
			writer.WriteHeader(http.StatusConflict)
		} else {
			// - Respond with OK if entry does not exist
			database[key] = value
			writer.WriteHeader(http.StatusCreated)
		}
	}

	// 6. Return updated dictionary
	json.NewEncoder(writer).Encode(database)
}

func deleteCustomer(writer http.ResponseWriter, req *http.Request) {
	// 1. Set Content Type
	writer.Header().Set("Content-Type", "application/json")
	// 2. Grab the member id from the url params
	params := mux.Vars(req)
	id := params["id"]
	if _, ok := database[id]; ok {
		// delete the entry, return successful response
		delete(database, id)
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode(database)
	} else { // 4. If not, return an error, but still return the dictionary
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(database)
	}
}

// TODO: - update a customer by id
func updateCustomer(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id := params["id"]

	var newEntry map[string]Models.Customer
	if _, ok := database[id]; ok {
		reqBody, _ := io.ReadAll(req.Body)
		json.Unmarshal(reqBody, &newEntry)
		value := newEntry[id]
		database[id] = value
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode(database)
	} else {
		writer.WriteHeader(http.StatusConflict)
		json.NewEncoder(writer).Encode(database)
	}
}

func StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", showHomePage)
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getSingleCustomer).Methods("GET")
	router.HandleFunc("/customers", createNewCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	fmt.Println("Server starting on port 3000")
	http.ListenAndServe(":3000", router)
}

func seedCustomerDatabase() map[string]Models.Customer {
	customers := make(map[string]Models.Customer)

	custID := uuid.New().String()
	customers[custID] = Models.Customer{
		ID:        custID,
		Name:      "Chris Myers",
		Role:      "Engineer",
		Email:     "chris.myers@nosuchco.com",
		Phone:     "765-897-0099",
		Contacted: false,
	}

	custID2 := uuid.New().String()
	customers[custID2] = Models.Customer{
		ID:        custID2,
		Name:      "Neville Myers",
		Role:      "Chief Dog Officer",
		Email:     "give.me.a.bone@nosuchco.com",
		Phone:     "000-000-0000",
		Contacted: false,
	}
	return customers
}
