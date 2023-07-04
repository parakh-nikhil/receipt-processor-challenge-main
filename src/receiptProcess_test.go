package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test to check if endpoint returns non empty id
func TestProcessReceiptHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/receipts/process", processReceiptHandler)

	sampleReceipt := getASampleReceipt()

	payload, _ := json.Marshal(sampleReceipt)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusCreated, res.Code)
	var response map[string]string
	_ = json.Unmarshal(res.Body.Bytes(), &response)
	id, exists := response["id"]

	assert.True(t, exists)
	assert.NotEmpty(t, id)
}

func TestProcessReceiptHandlerWithEmptyPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/receipts/process", processReceiptHandler)
	sampleReceipt := Receipt{}
	payload, _ := json.Marshal(sampleReceipt)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
	var response map[string]string
	_ = json.Unmarshal(res.Body.Bytes(), &response)
	id, exists := response["id"]
	assert.False(t, exists)
	assert.Empty(t, id)

	error := response["error"]
	assert.Equal(t, error, "Empty Payload")

}

func getASampleReceipt() Receipt {
	sampleReceipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Items: []Item{
			{
				ShortDescription: "Pepsi - 12-oz",
				Price:            "1.25",
			},
		},
		Total: "1.25",
	}

	return sampleReceipt
}
