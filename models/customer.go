package models

import (
	"errors"
	"strings"
)

type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"` // can't be empty - VALIDATE
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"` // `json0'th value is false
}

var myVar string            //emptyString (0'th value)
var myMap map[string]string // nil
var myMap2 map[int]int      // nil

// TODO: - What here needs to be validate?  i.e name can't be empty.  what can and cannot be empty?

func (c Customer) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("Name is required and cannot be empty")
	}
	return nil
}
