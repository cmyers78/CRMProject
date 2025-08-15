package customer

import (
	"CRMBackendProject/models"
	"fmt"
	"github.com/google/uuid"
)

var database = seedCustomerDatabase()

func GetAll() map[string]models.Customer {
	customers := database
	return customers // works for now, but need to be able to return an empty map
}

func Get(id string) (models.Customer, error) {
	customer, ok := database[id]
	if !ok {
		return models.Customer{}, fmt.Errorf("%s not found", id)
	}

	return customer, nil
}

func seedCustomerDatabase() map[string]models.Customer {
	customers := make(map[string]models.Customer)

	custID := uuid.New().String()
	customers[custID] = models.Customer{
		ID:        custID,
		Name:      "Chris Myers",
		Role:      "Engineer",
		Email:     "chris.myers@nosuchco.com",
		Phone:     "765-897-0099",
		Contacted: false,
	}

	custID2 := uuid.New().String()
	customers[custID2] = models.Customer{
		ID:        custID2,
		Name:      "Neville Myers",
		Role:      "Chief Dog Officer",
		Email:     "give.me.a.bone@nosuchco.com",
		Phone:     "000-000-0000",
		Contacted: false,
	}
	return customers
}
