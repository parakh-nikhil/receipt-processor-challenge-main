# Unit Test Documentation

The following provides an overview and explanation of the unit tests in the provided Go code :

All the tests sets up a test environment using the Gin framework's `TestMode` and makes respective REST calls to the server.

## **getPointsById_test.go**

### TestGetReceiptPointsByIdHandler

This test function verifies the behavior of the `getReceiptPointsByIdHandler` request handler.

1. Creates a sample receipt payload and sends a POST request to the `/receipts/process` endpoint to add a new receipt.
2. Extracts the generated receipt ID from the response.
3. Constructs a GET request to the `/receipts/:id/points` endpoint with the extracted ID.
4. Sends the request and asserts the expected HTTP status code (StatusFound).
5. Parses the response body and verifies the existence and value of the `points` field.

### TestGetReceiptPointsWithWrongId

This test function verifies the behavior of the `getReceiptPointsByIdHandler` request handler when using a wrong receipt ID.

1. Constructs a GET request to the `/receipts/:id/points` endpoint with a wrong receipt ID.
2. Sends the request and asserts the expected HTTP status code (StatusNotFound).
3. Parses the response body and verifies the error message.

### TestGetReceiptPointsWithNoId

This test function verifies the behavior of the `getReceiptPointsByIdHandler` request handler when no receipt is provided.

1. Constructs a GET request to the `/receipts/:id/points` endpoint without a receipt ID.
2. Sends the request and asserts the expected HTTP status code (StatusNotFound).
3. Parses the response body and verifies the error message.

---

## **receiptProcess_test.go**

### TestProcessReceiptHandler
This test function verifies the behavior of the `processReceiptHandler` request handler.

1. Creates a sample receipt payload and sends a POST request to the `/receipts/process` endpoint to add a new receipt.
2. Sends the request and asserts the expected HTTP status code (StatusCreated).
3. Extracts the generated receipt ID from the response and asserts that Id exists in the response

### TestProcessReceiptHandlerWithEmptyPayload
This test function verifies the behavior of the `processReceiptHandler` request handler when an empty payload is provided.

1. Creates a sample receipt payload and sends a POST request to the `/receipts/process` endpoint to add a new receipt.
2. Sends the request and asserts the expected HTTP status code (BadRequest).
3. Parses the response body and verifies the error message.

---

This concludes the documentation for the provided unit test file. The tests ensure the expected behavior of the `getReceiptPointsByIdHandler` request handler in different scenarios.

Please refer to the code comments and function definitions in the source code for more details.