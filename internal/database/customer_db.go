package database

import (
	"CRMBackendProject/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func GetAllCustomers(db *sql.DB) ([]models.Customer, error) {
	const queryStatement = "SELECT id, name, role, email, phone, contacted FROM customers" // we don't want this to change, also more performant
	rows, err := db.Query(queryStatement)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %s", err.Error())
	}
	// Close "closes" the stream of data from the database connection
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}()

	// it is more efficient to reuse database pointers than to delete or re-size db pointers
	// example: 100 open connections and just reuse those same 100 rather than expand and contract

	customers := make([]models.Customer, 0)
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone, &c.Contacted)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err.Error())
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func GetCustomer(id string, db *sql.DB) (models.Customer,
	error) {
	const queryStatement = "SELECT id, name, role, email, phone,contacted FROM customers WHERE id = ?"
	row := db.QueryRow(queryStatement, id)

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

func InsertCustomer(c models.Customer, db *sql.DB) (string,
	error) {
	id := uuid.New().String()
	const queryStatement = "INSERT INTO customers (id, name, role,email, phone, contacted) VALUES (?, ?, ?, ?, ?,?)"

	_, err := db.Exec(queryStatement, id, c.Name, c.Role, c.Email, c.Phone, c.Contacted)
	if err != nil {
		return "", fmt.Errorf("failed to insert customer: %w", err)
	}
	return id, nil
}

func UpdateCustomer(c models.Customer, db *sql.DB) error {
	const queryStatement = "UPDATE customers SET name = ?, role =?, email = ?, phone = ?, contacted = ? WHERE id= ?"

	_, err := db.Exec(queryStatement, c.Name, c.Role,
		c.Email, c.Phone, c.Contacted, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func DeleteCustomer(id string, db *sql.DB) error {
	const queryStatement = "DELETE FROM customers WHERE id = ?"

	_, err := db.Exec(queryStatement, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}
