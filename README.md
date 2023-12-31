# Receipt Processor

A webservice that processes receipt points and stores the receipts.  Data does not survive an application restart. Any data generated by this endpoint is stored in-memory.

## Instructions to run the program

*Assuming you have go installed locally. If not, follow these [download instructions](https://go.dev/doc/install) on the office goLang website.*

- Go to the project home directory
- ```cd src```
- ```go install``` (this commands downloads and installs the code dependencies)
- ```go run .``` (this commands starts the server with localhost as the IP address of the web server and 8080 as the port)

## Code and Unit Tests Documentation

The code documentation can be found in [DOCUMENTATION.md](https://github.com/parakh-nikhil/receipt-processor-challenge-main/blob/master/DOCUMENTATION.md) fle.

The Unit tests documentation can be found in [UNIT_TESTS_DOCUMENTATION.md](https://github.com/parakh-nikhil/receipt-processor-challenge-main/blob/master/UNIT_TESTS_DOCUMENTATION.md) file

## Summary of API Specification

The API is documented in [api.yaml](https://github.com/parakh-nikhil/receipt-processor-challenge-main/blob/master/api.yml) file

### Endpoint: Process Receipts

* Path: `/receipts/process`
* Method: `POST`
* Payload: Receipt JSON
* Response: JSON containing an id for the receipt.

Description:

Takes in a JSON receipt and returns a JSON object with an ID generated by the code. The ID returned is the unique identifier of the receipt object.

*Example*

Endpoint call: `localhost:8080/receipts/process`

Request Body: 
```json
{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
    ]
}
```
Response:
```json
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```

## Endpoint: Get Points

* Path: `/receipts/{id}/points`
* Method: `GET`
* Response: A JSON object containing the number of points awarded.

A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.

**Example**

Endpoint call: `localhost:8080/receipts/{id}/points`

Response:
```json
{ "points": 32 }
```

## Endpoint: Get All Receipts

* Path: `/receipts/`
* Method: `GET`
* Response: A JSON object containing the list of all receipts.

A simple Getter endpoint that returns an object specifying the list of all receipts.

**Example**

Endpoint call: `localhost:8080/receipts/`

Response:
```json
{"receipts": [{"Receipt 1"},{"Receipt 2"},{"Receipt 3"}]}
```

---

# Rules for calculating points

These rules collectively define how many points should be awarded to a receipt.

* One point for every alphanumeric character in the retailer name.
* 50 points if the total is a round dollar amount with no cents.
* 25 points if the total is a multiple of `0.25`.
* 5 points for every two items on the receipt.
* If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
* 6 points if the day in the purchase date is odd.
* 10 points if the time of purchase is after 2:00pm and before 4:00pm.
