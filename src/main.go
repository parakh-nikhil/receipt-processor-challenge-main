package main

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Item struct {
	Id               string `json:"id"`
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Id           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

var receipts []Receipt

var receiptPoints = make(map[string]int)

func processReceiptHandler(c *gin.Context) {
	var newReceipt Receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	newReceipt.Id = uuid.New().String()
	receipts = append(receipts, newReceipt)
	receiptPoints[newReceipt.Id] = processReceiptPoints(newReceipt)
	c.IndentedJSON(http.StatusCreated, gin.H{"id": newReceipt.Id})

}

func getReceiptPointsByIdHandler(c *gin.Context) {
	id := c.Param("id")
	points, exists := receiptPoints[id]
	pointsStr := strconv.Itoa(points)
	if exists {
		c.IndentedJSON(http.StatusFound, gin.H{"points": pointsStr})
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No receipt found with id: " + id})
	}

}

func getAllReceiptsHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusFound, gin.H{"receipts": receipts})
}

func main() {
	router := gin.Default()

	router.POST("/receipts/process", processReceiptHandler)
	router.GET("/receipts/:id/points", getReceiptPointsByIdHandler)
	router.GET("/receipts", getAllReceiptsHandler)
	router.Run("localhost:8080")
}

func processReceiptPoints(receipt Receipt) int {
	var points int
	points += getNumberOfAlphanumericCharInString(receipt.Retailer)
	_, cents, err := getDollarAndCentFromPrice(receipt.Total)
	if err == nil {
		if cents == 0 {
			points += 50
		}
		if cents%25 == 0 {
			points += 25
		}
	}
	points += (len(receipt.Items) / 2) * 5
	for _, item := range receipt.Items {
		point, err := getPointsFromItemDescription(item)
		if err == nil {
			points += point
		}
	}
	pointsFromDateTime, err := getPointsFromDateTime(receipt.PurchaseDate, receipt.PurchaseTime)
	if err == nil {
		points += pointsFromDateTime
	}
	return points
}

func getNumberOfAlphanumericCharInString(str string) int {
	totalAlphaNumeric := 0
	var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_]*$")
	for i := 0; i < len(str); i++ {
		if alphanumeric.Match([]byte(str[i : i+1])) {
			totalAlphaNumeric++
		}
	}
	return totalAlphaNumeric
}

func getDollarAndCentFromPrice(price string) (int, int, error) {
	priceSplit := strings.Split(price, ".")
	dollar, err := strconv.Atoi(priceSplit[0])
	if err != nil {
		//TODO: check for this status
		return -1, -1, err
	}
	cents, err := strconv.Atoi(priceSplit[1])
	if err != nil {
		//TODO: check for this status
		return -1, -1, err
	}
	return dollar, cents, nil
}

func getPointsFromItemDescription(item Item) (int, error) {
	itemDesc := strings.TrimSpace(item.ShortDescription)
	if len(itemDesc)%3 == 0 {
		priceFloat, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return 0, err
		}
		roundedPrice := int(math.Round(priceFloat * 0.2))
		if float64(roundedPrice) < priceFloat {
			roundedPrice++
		}
		return roundedPrice, nil
	}
	return 0, nil
}

func getPointsFromDateTime(date string, timeStr string) (int, error) {
	dateTimeStr := date + " " + timeStr
	purchaseDateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return 0, err
	}
	var points int
	if purchaseDateTime.Day()%2 != 0 {
		points += 6
	}

	if purchaseDateTime.Hour() >= 14 && purchaseDateTime.Hour() < 16 {
		points += 10
	}
	return points, nil
}
