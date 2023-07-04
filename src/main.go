package main

import (
	"log"
	"math"
	"net/http"
	"os"
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
	log.Printf("POST : %v\n", c.Request.URL)
	log.Println("processReceiptHandler(c *gin.Context) handling request")
	var newReceipt Receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		log.Fatalf("Error: %v\n", err)
		return
	}
	if isEmpty(newReceipt) {
		log.Println("Error: Empty payload")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Empty Payload"})
	} else {
		newReceipt.Id = uuid.New().String()
		receipts = append(receipts, newReceipt)
		receiptPoints[newReceipt.Id] = processReceiptPoints(newReceipt)
		c.IndentedJSON(http.StatusCreated, gin.H{"id": newReceipt.Id})
		log.Printf("Response sent: {id:%v}\n", newReceipt.Id)
	}
}

func getReceiptPointsByIdHandler(c *gin.Context) {
	id := c.Param("id")
	log.Printf("GET : %v\n", c.Request.URL)
	points, exists := receiptPoints[id]
	pointsStr := strconv.Itoa(points)
	if exists {
		c.IndentedJSON(http.StatusFound, gin.H{"points": pointsStr})
		log.Printf("Response Sent: {\"points\":\"%v\"}\n", pointsStr)
	} else {
		log.Printf("Error: No receipt found with id: %v", id)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No receipt found with id: " + id})
		log.Printf("Response Sent: {\"error\" : \"No receipt found with id: %v\"} ", id)
	}

}

func getAllReceiptsHandler(c *gin.Context) {
	log.Printf("GET : %v\n", c.Request.URL)
	c.IndentedJSON(http.StatusFound, gin.H{"receipts": receipts})
	log.Printf("Response Sent: {\"receipts\": {%v}\n}", receipts)
}

func processReceiptPoints(receipt Receipt) int {
	log.Println("Processing receipt points...")
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
	log.Printf("Total Points: %v\n", points)
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
	// log.Fatal("")
	dateTimeStr := date + " " + timeStr
	purchaseDateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		log.Printf("Error: %v", err)
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

func getLogFileName() string {
	logDir := "./logs"
	currentTime := time.Now()
	date := currentTime.Format("01-02-2006 15:04:05")
	dateStr := strings.ReplaceAll(date, " ", "_")
	logFileName := logDir + "/app-log_" + dateStr + ".log"
	return logFileName
}

func isEmpty(r Receipt) bool {
	if (r.Id == "") && (len(r.Items) == 0) && (r.PurchaseDate == "") && (r.PurchaseTime == "") && (r.Retailer == "") && (r.Total == "") {
		return true
	}
	return false
}

func main() {
	router := gin.Default()
	logFileName := getLogFileName()

	logfile, err := os.Create(logFileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Log File Created: %v\n ", logFileName)
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Println("Server is starting...")
	router.POST("/receipts/process", processReceiptHandler)
	router.GET("/receipts/:id/points", getReceiptPointsByIdHandler)
	router.GET("/receipts", getAllReceiptsHandler)
	router.Run("localhost:8080")

}
