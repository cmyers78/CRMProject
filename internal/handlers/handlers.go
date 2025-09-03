package handlers

// Business Logic and Handler currently
import (
	"CRMBackendProject/internal/customer"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// NOTE: - Handlers file should show handlers first

func ShowHomePage(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "./static/static.html")
}

func GetAllCustomers(writer http.ResponseWriter, req *http.Request) {
	customers := customer.GetAll()
	writeResponse(writer, customers, http.StatusOK)
}

func GetSingleCustomer(writer http.ResponseWriter, req *http.Request) {
	// Handler logic
	id := extractID(req)

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

func extractID(req *http.Request) string {
	params := mux.Vars(req)
	id := params["id"]
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
	id, database := extractID(req), customer.GetAll()
	if _, ok := database[id]; ok {
		delete(database, id)
		writeResponse(writer, database, http.StatusNoContent)
		return
	}
	writeResponse(writer, database, http.StatusNotFound)
}

//
//// TODO: - update a customer by id
//func UpdateCustomer(writer http.ResponseWriter, req *http.Request) {
//	writer.Header().Set("Content-Type", "application/json")
//	params := mux.Vars(req)
//	id := params["id"]
//
//	var newEntry models.Customer
//	if _, ok := database[id]; ok {
//		reqBody, _ := io.ReadAll(req.Body)
//		json.Unmarshal(reqBody, &newEntry)

//		database[newEntry.id] = newEntry
//		database[newEntry.id].Name = newEntry.Name
//		writer.WriteHeader(http.StatusAccepted)
//		json.NewEncoder(writer).Encode(database)
//	} else {
//		writer.WriteHeader(http.StatusConflict)
//		json.NewEncoder(writer).Encode(database)
//	}
//}
