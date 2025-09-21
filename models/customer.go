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

func (c Customer) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("Name is required and cannot be empty")
	}

	if strings.TrimSpace(c.Email) == "" && strings.TrimSpace(c.Phone) == "" {
		return errors.New("Email or Phone Number is required and cannot be empty")
	}
	return nil
}
