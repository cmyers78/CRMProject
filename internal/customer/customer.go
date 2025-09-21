package customer

import (
	"CRMBackendProject/models"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var customerdb = seedCustomerDatabase()

// Go
// Implementing a thread-safe map to store customer data
// Notes:
// - SafeMap struct contains a mutex and a map to store customer data
// - CRUD operations are defined as methods on the SafeMap struct
// - Mutex is used to ensure thread safety during read and write operations
// - getAll, get, set, and delete methods provide basic CRUD functionality
// - GetDB function returns a pointer to the SafeMap instance

type SafeMap struct {
	mu       sync.RWMutex
	database map[string]models.Customer
}

// CRUD operations that will be used by Public functions in customer.go to manipulate the customer data
func (s *SafeMap) getAll() map[string]models.Customer {
	s.mu.RLock()         // multithreaded - read only, so multiple threads can read at the same time
	defer s.mu.RUnlock() // unlock after function is done
	return s.database    // returns a SafeMap database
}

func (s *SafeMap) get(key string) (models.Customer, bool) {
	s.mu.RLock() // multithreaded
	defer s.mu.RUnlock()
	val, ok := s.database[key]
	return val, ok
}

func (s *SafeMap) set(key string, val models.Customer) {
	s.mu.Lock()           // single threaded - write only, so only one thread can write at a time
	defer s.mu.Unlock()   // unlock after function is done
	s.database[key] = val // adds or updates the value in the map using the key as the identifier
}

func (s *SafeMap) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.database, key) // built-in delete function to remove a key/value pair from the "database" map
}
func GetDB() *SafeMap {
	return customerdb // returns a pointer to the SafeMap
}

// Go
// This section is public-facing methods for handlers to get/insert customer data
// Notes:
// - GetAll returns a map of all customers
// - Get returns a single customer based on the provided ID
// - Insert adds a new customer to the database and returns the new customer's ID
// - Update modifies an existing customer's information
// - Delete removes a customer from the database based on the provided ID

func GetAll() map[string]models.Customer {
	customers := customerdb.getAll()
	return customers
}

func Get(id string) (models.Customer, error) {
	customer, ok := customerdb.get(id)
	if !ok {
		return models.Customer{}, fmt.Errorf("%s not found", id)
	}
	return customer, nil
}
func Insert(c models.Customer) (string, error) {
	id := uuid.New().String()
	newEntry := models.Customer{
		ID:        id,
		Name:      c.Name,
		Role:      c.Role,
		Email:     c.Email,
		Phone:     c.Phone,
		Contacted: c.Contacted,
	}
	// check to see if value already exists
	_, ok := customerdb.get(id)
	if ok {
		return id, fmt.Errorf("customer with id %s already exists", id)
	}
	customerdb.set(id, newEntry)
	return id, nil
}

func Update(c models.Customer) error {
	customerdb.set(c.ID, c)
	return nil
}

func Delete(id string) error {
	customerdb.delete(id)
	return nil
}

// This creates our pointer to a SafeMap and seeds it with some initial data
func seedCustomerDatabase() *SafeMap {
	customers := SafeMap{
		database: make(map[string]models.Customer),
	}
	// add some seed data
	custID := uuid.New().String()
	customer1 := models.Customer{
		ID:        custID,
		Name:      "Chris Myers",
		Role:      "Engineer",
		Email:     "chris.myers@nosuchco.com",
		Phone:     "765-897-0099",
		Contacted: false,
	}
	customers.set(custID, customer1)

	custID2 := uuid.New().String()
	customer2 := models.Customer{
		ID:        custID2,
		Name:      "Neville Myers",
		Role:      "Chief Dog Officer",
		Email:     "give.me.a.bone@nosuchco.com",
		Phone:     "000-000-0000",
		Contacted: false,
	}
	customers.set(custID2, customer2)
	return &customers
}
