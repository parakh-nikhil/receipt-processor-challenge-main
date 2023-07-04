# Code Documentation

Primary Language: Go

The code provides an API for processing receipts, calculating points, and retrieving receipt information, making it useful for building a receipt tracking and reward system.

The following documentation provides an overview and explanation of the functionality of the provided Go code.

## Overview

The code is a Go program that implements a simple API for processing receipts and calculating points based on various criteria. It utilizes the Gin framework (`github.com/gin-gonic/gin`) for handling HTTP requests and routing.

## Struct Types

### Item

```go
type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}
```

Defines a struct type **Item** with two fields: ShortDescription and Price. The json tags are used for JSON serialization and specify the field names in the JSON representation.

### Receipt

```go
type Receipt struct {
    Id           string `json:"id"`
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []Item `json:"items"`
    Total        string `json:"total"`
}
```

Defines a struct type Receipt representing a receipt. It contains the following fields: Id, Retailer, PurchaseDate, PurchaseTime, Items, and Total. The Items field is a slice of Item structs. The json tags are used for JSON serialization.

Receipt examples can be found in the [examples](https://github.com/parakh-nikhil/receipt-processor-challenge-main/tree/master/examples) directory

## Global Variables

* **receipts**: A global variable that stores a slice of Receipt structs to store multiple receipts.
* **receiptPoints**: A global variable that stores a map associating receipt IDs with their corresponding points.

## Request Handlers
The program includes the following request handlers:

### processReceiptHandler
Function to handle the HTTP POST request for processing a new receipt.

### getReceiptPointsByIdHandler
Function to handle the HTTP GET request for retrieving receipt points by ID.

### getAllReceiptsHandler
Function to handle the HTTP GET request for retrieving all receipts.

## Utility Functions
The code includes several utility functions used by the request handlers to perform specific tasks:

* ```processReceiptPoints```: Calculates the points earned from a receipt based on various criteria.
* ```getNumberOfAlphanumericCharInString```: Counts the number of alphanumeric characters in a string.
* ```getDollarAndCentFromPrice```: Extracts dollar and cent values from a price string.
* ```getPointsFromItemDescription```: Calculates points based on the item description.
* ```getPointsFromDateTime```: Calculates points based on the purchase date and time.
* ```getLogFileName```: Generates a log file name with the current date and time.

## Entry Point
The ```main``` function serves as the entry point of the program. It performs the following tasks:

- Sets up the Gin router.
- Creates a log file based on the current date and time.
- Sets the log output to the file.
- Starts the server on localhost:8080.

## Logs
Logs are stored under the ```~/src/logs``` directory.

The structure of the log files is ```app-log_{DATE}_{TIME}.log```. Every time the server is restarted, it creates a log file and stores the server events on the file.

Example: 

*app-log_07-04-2023_17:23:48.log*

This is the log file for the server when it was started on (07-04-2023) at (17:23:48).
