package handlers

// Business Logic and Handler currently
import (
	"CRMBackendProject/internal/customer"
	"CRMBackendProject/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// NOTE: - Handlers file should show handlers first
var (
	database = customer.GetDB()
)

func ShowHomePage(writer http.ResponseWriter, req *http.Request) {
	path, err := os.Executable()
	if err != nil {
		http.Error(writer, "Unable to get executable path", http.StatusInternalServerError)
		return
	}
	staticPath := filepath.Join(filepath.Dir(path), "static/static.html")
	http.ServeFile(writer, req, staticPath)
}

func RetrieveAllCustomers(writer http.ResponseWriter, _ *http.Request) {
	customers := customer.GetAll()
	writeResponse(writer, customers, http.StatusOK)
}

func RetrieveSingleCustomer(writer http.ResponseWriter, req *http.Request) {
	// Handler logic
	id := extractOne("id", req)

	cst, err := customer.Get(id)
	if err != nil {
		writeResponse(writer, cst, http.StatusNotFound)
		return
	}
	writeResponse(writer, cst, http.StatusOK)
}

func CreateNewCustomer(writer http.ResponseWriter, req *http.Request) {
	// 1. set content-type to JSON
	writer.Header().Set("Content-Type", "application/json")

	// 2. keep track of new entry so that it can be added to dictionary map
	var newEntry models.Customer

	err := json.NewDecoder(req.Body).Decode(&newEntry)
	if err != nil {
		writeResponse(writer, newEntry, http.StatusUnprocessableEntity)
		return
	}
	// 3. Validate new entry (i.e. name and role cannot be empty)
	validationError := newEntry.Validate()
	if validationError != nil {
		fmt.Printf("Erorr: %s", validationError.Error())
		writeResponse(writer, newEntry, http.StatusBadRequest)
		return
	}
	// 5. Add new entry to dictionary map if it doesn't already exist
	key, err := customer.Insert(newEntry)
	if err != nil {
		fmt.Printf("Error: %s", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	newEntry.ID = key

	// 6. Return updated customer record
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(newEntry)
}

func DeleteCustomer(writer http.ResponseWriter, req *http.Request) {
	id := extractOne("id", req)
	// NOTE: Validation of id
	if id == "" {
		writeResponse(writer, id, http.StatusBadRequest)
		return
	}
	cst, err := customer.Get(id)
	if err != nil {
		writeResponse(writer, cst, http.StatusNotFound)
		return
	}

	_ = customer.Delete(cst.ID)
	writeResponse(writer, database, http.StatusOK)
}

func UpdateCustomer(writer http.ResponseWriter, req *http.Request) {
	// this works fine now, but I assume will have to be pulled from a db later and will need error handling
	var newEntry models.Customer
	err := json.NewDecoder(req.Body).Decode(&newEntry)
	if err != nil {
		writeResponse(writer, newEntry, http.StatusUnprocessableEntity)
		return
	}
	// NOTE: Validation of id
	if newEntry.ID == "" {
		writeResponse(writer, newEntry, http.StatusBadRequest)
		return
	}
	_, err = customer.Get(newEntry.ID)
	if err != nil {
		writeResponse(writer, newEntry.ID, http.StatusNotFound)
		return
	}
	err = customer.Update(newEntry)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writeResponse(writer, database, http.StatusAccepted)
}

// NOTE: - Keep helpers at the bottom of the page
func writeResponse(writer http.ResponseWriter, data any, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(data)
}

func extractOne(key string, req *http.Request) string {
	params := mux.Vars(req) // this only needs to run one time
	value := params[key]
	return value
}

// QUESTIONS:
// 1. Why use unmarshal instead of decode?  answer: For HTML, Decode is preferred, especially for large payloads, as it streams the data directly from the request body.
