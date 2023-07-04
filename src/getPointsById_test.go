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

func TestGetReceiptPointsByIdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/receipts/:id/points", getReceiptPointsByIdHandler)
	router.POST("/receipts/process", processReceiptHandler)

	sampleReceipt := Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
		Total: "9.00",
	}

	payload, _ := json.Marshal(sampleReceipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	var response map[string]string
	_ = json.Unmarshal(res.Body.Bytes(), &response)
	id := response["id"]

	//TEST for getting receipt points
	path := "/receipts/" + id + "/points"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusFound, res.Code)
	var pointsResponse map[string]string
	_ = json.Unmarshal(res.Body.Bytes(), &pointsResponse)
	points, exists := pointsResponse["points"]
	assert.True(t, exists)
	assert.NotEmpty(t, points)
	assert.Equal(t, "109", points)

}

func TestGetReceiptPointsWithWrongId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/receipts/:id/points", getReceiptPointsByIdHandler)
	wrongId := "81f68990-ebd7-43f9-9884-6007ba5d0138"
	path := "/receipts/" + wrongId + "/points"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
	var response map[string]string
	_ = json.Unmarshal(res.Body.Bytes(), &response)
	errMsg := response["error"]
	assert.Equal(t, errMsg, "No receipt found with id: "+wrongId)
}

func TestGetReceiptPointsWithNoId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/receipts/:id/points", getReceiptPointsByIdHandler)
	wrongId := ""
	path := "/receipts/" + wrongId + "/points"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
	var response map[string]string
	_ = json.Unmarshal(res.Body.Bytes(), &response)
	errMsg := response["error"]
	assert.Equal(t, errMsg, "No receipt found with id: "+wrongId)
}
