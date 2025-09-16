package handlers

// Business Logic and Handler currently
import (
	"CRMBackendProject/internal/customer"
	"CRMBackendProject/models"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// NOTE: - Handlers file should show handlers first

func ShowHomePage(writer http.ResponseWriter, req *http.Request) {
	path, err := os.Executable()
	if err != nil {
		http.Error(writer, "Unable to get executable path", http.StatusInternalServerError)
		return
	}
	staticPath := filepath.Join(filepath.Dir(path), "static/static.html")
	http.ServeFile(writer, req, staticPath)
}

func GetAllCustomers(writer http.ResponseWriter, _ *http.Request) {
	customers := customer.GetAll()
	writeResponse(writer, customers, http.StatusOK)
}

func GetSingleCustomer(writer http.ResponseWriter, req *http.Request) {
	// Handler logic
	id := extract("id", req)

	customer, err := customer.Get(id)

	if err != nil {
		writeResponse(writer, customer, http.StatusNotFound)
		return
	}
	writeResponse(writer, customer, http.StatusOK)
}

// NOTE: - Keep helpers at the bottom of the page
func writeResponse(writer http.ResponseWriter, data any, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(data)
}

func extract(key string, req *http.Request) string {
	params := mux.Vars(req)
	id := params[key]
	return id
}

//	func CreateNewCustomer(writer http.ResponseWriter, req *http.Request) {
//		// 1. set content-type to JSON
//		writer.Header().Set("Content-Type", "application/json")
//
//		// 2. keep track of new entry so that it can be added to dictionary map
//		var newEntry map[string]models.Customer
//
//		// 3. Read the request
//		reqBody, _ := io.ReadAll(req.Body)
//
//		// 4. Parse JSON Body
//		json.Unmarshal(reqBody, &newEntry)
//
//		// 5. Add new entry to dictionary map if it doesn't already exist
//		for key, value := range newEntry {
//			// - Respond with conflict if entry exists
//			if _, ok := database[key]; ok {
//				writer.WriteHeader(http.StatusConflict)
//			} else {
//				// - Respond with OK if entry does not exist
//				database[key] = value
//				writer.WriteHeader(http.StatusCreated)
//			}
//		}
//
//		// 6. Return updated dictionary
//		json.NewEncoder(writer).Encode(database)
//	}
func DeleteCustomer(writer http.ResponseWriter, req *http.Request) {
	id, database := extract("id", req), customer.GetAll()
	if _, ok := database[id]; ok {
		delete(database, id)
		writeResponse(writer, database, http.StatusNoContent)
		return
	}
	writeResponse(writer, database, http.StatusNotFound)
}

func UpdateCustomer(writer http.ResponseWriter, req *http.Request) {
	id, database := extract("id", req), customer.GetAll()
	// this works fine now, but i assume will have to be pulled from a db later and will need error handling
	var newEntry models.Customer

	// read the request body and handle error
	// unmarshal the request body into newEntry and handle error
	// check if the id exists in the database
	// if it exists, update the entry and respond with status accepted
	// if it doesn't exist, respond with status not found
	if _, ok := database[id]; ok {
		reqBody, _ := io.ReadAll(req.Body)
		err := json.Unmarshal(reqBody, &newEntry)
		if err != nil {
			writeResponse(writer, newEntry, http.StatusUnprocessableEntity)
			return
		}
		database[newEntry.ID] = newEntry
		writeResponse(writer, database, http.StatusAccepted)
		return
	}
	writeResponse(writer, database, http.StatusNotFound)
}

// QUESTIONS:
// 1. Why use unmarshal instead of decode?
// 2. How do I handle the unmarshal error correctly?--nesting seems wrong
