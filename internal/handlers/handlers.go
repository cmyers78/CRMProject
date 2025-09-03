package handlers

// Business Logic and Handler currently
import (
	"CRMBackendProject/internal/customer"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ShowHomePage(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "./static/static.html")
}

func GetAllCustomers(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	customers := customer.GetAll()
	json.NewEncoder(writer).Encode(customers)
}

func GetSingleCustomer(writer http.ResponseWriter, req *http.Request) {
	// Handler logic
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id := params["id"]

	customer, err := customer.Get(id)

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(customer)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	json.NewEncoder(writer).Encode(customer)
}

//
//func CreateNewCustomer(writer http.ResponseWriter, req *http.Request) {
//	// 1. set content-type to JSON
//	writer.Header().Set("Content-Type", "application/json")
//
//	// 2. keep track of new entry so that it can be added to dictionary map
//	var newEntry map[string]models.Customer
//
//	// 3. Read the request
//	reqBody, _ := io.ReadAll(req.Body)
//
//	// 4. Parse JSON Body
//	json.Unmarshal(reqBody, &newEntry)
//
//	// 5. Add new entry to dictionary map if it doesn't already exist
//	for key, value := range newEntry {
//		// - Respond with conflict if entry exists
//		if _, ok := database[key]; ok {
//			writer.WriteHeader(http.StatusConflict)
//		} else {
//			// - Respond with OK if entry does not exist
//			database[key] = value
//			writer.WriteHeader(http.StatusCreated)
//		}
//	}
//
//	// 6. Return updated dictionary
//	json.NewEncoder(writer).Encode(database)
//}
//
//func DeleteCustomer(writer http.ResponseWriter, req *http.Request) {
//	// 1. Set Content Type
//	writer.Header().Set("Content-Type", "application/json")
//	// 2. Grab the member id from the url params
//	params := mux.Vars(req)
//	id := params["id"]
//	if _, ok := database[id]; ok {
//		// delete the entry, return successful response
//		delete(database, id)
//		writer.WriteHeader(http.StatusAccepted)
//		json.NewEncoder(writer).Encode(database)
//	} else { // 4. If not, return an error, but still return the dictionary
//		writer.WriteHeader(http.StatusNotFound)
//		json.NewEncoder(writer).Encode(database)
//	}
//}
//
//// TODO: - update a customer by id
//func UpdateCustomer(writer http.ResponseWriter, req *http.Request) {
//	writer.Header().Set("Content-Type", "application/json")
//	params := mux.Vars(req)
//	id := params["id"]
//
//	var newEntry map[string]models.Customer
//	if _, ok := database[id]; ok {
//		reqBody, _ := io.ReadAll(req.Body)
//		json.Unmarshal(reqBody, &newEntry)
//		value := newEntry[id]
//		database[id] = value
//		writer.WriteHeader(http.StatusAccepted)
//		json.NewEncoder(writer).Encode(database)
//	} else {
//		writer.WriteHeader(http.StatusConflict)
//		json.NewEncoder(writer).Encode(database)
//	}
//}
