package database

import (
	"CRMBackendProject/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func GetAllCustomers() (map[string]models.Customer, error) {
	query := "SELECT id, name, role, email, phone, contacted FROM customers"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %s", err.Error())
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}()

	customers := make(map[string]models.Customer)
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone, &c.Contacted)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err.Error())
		}
		customers[c.ID] = c
	}
	return customers, nil
}

func GetCustomer(id string) (models.Customer,
	error) {
	query := "SELECT id, name, role, email, phone,contacted FROM customers WHERE id = ?"
	row := DB.QueryRow(query, id)

	var c models.Customer
	err := row.Scan(&c.ID, &c.Name, &c.Role,
		&c.Email, &c.Phone, &c.Contacted)
	if err == sql.ErrNoRows {
		return models.Customer{}, fmt.Errorf("%s notfound", id)
	}
	if err != nil {
		return models.Customer{}, fmt.Errorf("failed toget customer: %w", err)
	}
	return c, nil
}

func InsertCustomer(c models.Customer) (string,
	error) {
	id := uuid.New().String()
	query := "INSERT INTO customers (id, name, role,email, phone, contacted) VALUES (?, ?, ?, ?, ?,?)"

	_, err := DB.Exec(query, id, c.Name, c.Role, c.Email, c.Phone, c.Contacted)
	if err != nil {
		return "", fmt.Errorf("failed to insert customer: %w", err)
	}
	return id, nil
}

func UpdateCustomer(c models.Customer) error {
	query := "UPDATE customers SET name = ?, role =?, email = ?, phone = ?, contacted = ? WHERE id= ?"

	_, err := DB.Exec(query, c.Name, c.Role,
		c.Email, c.Phone, c.Contacted, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func DeleteCustomer(id string) error {
	query := "DELETE FROM customers WHERE id = ?"

	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}
