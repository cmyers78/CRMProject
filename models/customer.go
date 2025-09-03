package models

type Customer struct {
	ID        string // can't be empty
	Name      string // can't be empty
	Role      string `json: "role"`
	Email     string
	Phone     string
	Contacted bool // 0'th value is false
}

var myVar string            //emptyString (0'th value)
var myMap map[string]string // nil
var myMap2 map[int]int      // nil

// TODO: - What here needs to be validate?  i.e name can't be empty.  what can and cannot be empty?
