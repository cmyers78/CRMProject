# CRM Backend Project - Go - Udacity

## Overview
This project serves to demonstrate how Go can be used to leverage updating a database through routes via the gorillaMux package.  

Users can make API calls via Postman, RapidAPI, or some other such tool to get a list of customers, add/delete customers, update customer information, or retrieve information about a single customer.

## How to run
The main.go file is located in the `cmd` folder.  To run the application locally, navigate to the project and type `go run cmd/main.go`

This will prompt you to decided whether you want to start the sever. Selecting yes will proceed.

Using Postman or RapidApi you can make various calls.  Endpoints are as follows:

### BaseURL: `http://localhost:3000/` Static webpage with information about the projects
### GET `/customers ` Retrieve all customers
### GET `/customers/{id}` Retrieves a specific customer by their unique ID
### POST `/customers` Will take a JSON body and add a new customer
`Adding a customer example`
```json
{
"987CMY": {"ID":"987CMY","Name":"Hank Hill","Role":"Assistant Manager - Propane & Propane accessories","Email":"hank.hill@nosuchco.com","Phone":"765-678-2342","Contacted":true}
}
```
### DELETE `/customers/{id}` Will delete a specific customer by their unique ID
### PUT `/customers/{id}` Will update information for a specific customer by their unique ID
