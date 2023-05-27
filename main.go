package main

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Struct used to hold individual items read in from receipt jsons
type item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Struct used to hold information from incoming receipt jsons while doing point calculation
type rec struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

// Struct used to build jsons for the processed receipts
type procRec struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
}

// list of all alphanumeric characters, used for point tallying and id creation
const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Array to hold all processed receipts
var procRecs = []procRec{{ID: "1", Points: 50}}

func main() {
	router := gin.Default()

	router.GET("/receipts/processed", getProcRecs)
	router.GET("/receipts/:id/points", getPoints)
	router.POST("/receipts/process", processReceipt)

	router.Run("localhost:8080")
}

/*
Get the points from the processed receipt with the given id
Response: json containing the total points of the given processed receipt
*/
func getPoints(c *gin.Context) {
	//grab the id to identify the correct receipt
	id := c.Param("id")

	//loop through processed receipts looking for matching id
	for _, r := range procRecs {
		//if we find a matching id, return the points of the
		if r.ID == id {
			c.IndentedJSON(http.StatusOK, r.Points)
			return
		}
	}
	//if no match was found, return with a not found message
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}

/*
Get all the processed receipts and print them to a json for testing purposes
Response: all processed receipts in a json file
*/
func getProcRecs(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, procRecs)
}

/*
Processes a given receipt json, giving it a point value and storing the processed receipt data in the proper array
Payload: Receipt json
Response: id of the created processed receipt
*/
func processReceipt(c *gin.Context) {
	var newReceipt rec

	//Attempt to bind incoming json to rec struct, return if fail
	if err := c.BindJSON(&newReceipt); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "Receipt not accepted"})
		return
	}

	//Initialize a point counter
	points := 0

	//----------------------------------------------------------Count points-----------------------------------------------------------------

	//Add one point for each alphanumeric character in the retailer name TODO::Make this function not suck
	alNumCount := 0
	for _, c := range characters {
		alNumCount += strings.Count(newReceipt.Retailer, string(c))
	}
	points += alNumCount

	//If the total is a round dollar amount, add 50 points to the counter
	tFloat, _ := strconv.ParseFloat(newReceipt.Total, 64)
	tCents := tFloat - math.Floor(tFloat)
	if tCents == 0.0 {
		points += 50
	}

	//If the total is a multiple of .25, add 25 points to the counter
	if math.Mod(tCents, 0.25) == 0 {
		points += 25
	}

	//for every two items on the receipt, add 5 points
	listLen := len(newReceipt.Items)
	points += 5 * (listLen / 2)

	//If the trimmed length of an item's name is a multiple of three, add points equal to 0.2 * the price, rounded up
	for _, i := range newReceipt.Items {
		if len(strings.TrimSpace(i.ShortDescription))%3 == 0 {
			pFloat, _ := strconv.ParseFloat(i.Price, 64)
			points += int(math.Ceil(0.2 * pFloat))
		}
	}

	//if the purchase day is odd, add 6 points
	date := strings.Split(newReceipt.PurchaseDate, "-")
	if day, err := strconv.Atoi(date[len(date)-1]); err == nil && day%2 == 1 {
		points += 6
	}

	//if the time purchased is between 2 and 4 pm (14 and 16), add 10 points
	//Split the time up and convert the pieces to integers
	time := strings.Split(newReceipt.PurchaseTime, ":")
	var intTime []int
	for _, t := range time {
		if newTime, err := strconv.Atoi(t); err == nil {
			intTime = append(intTime, newTime)
		}
	}
	//check if the time lands between 13:00 and 15:59
	if intTime[0] >= 14 && intTime[0] < 16 {
		points += 10
	} else if intTime[0] == 16 && intTime[1] == 0 { //additional check for 16:00
		points += 10
	}

	//----------------------------------------------------------point tally finished-----------------------------------------------------------------

	//make an ID for the new tally
	newId := generateId()

	//Add new processed receipt to array
	var newProcRec procRec
	newProcRec.ID = newId
	newProcRec.Points = points
	procRecs = append(procRecs, newProcRec)

	//Return a json with the id of the new processed receipt
	c.IndentedJSON(http.StatusOK, newProcRec.ID)
}

// Makes an id
func generateId() string {

	id1 := make([]byte, 8)
	for i := range id1 {
		id1[i] = characters[rand.Intn(len(characters))]
	}
	id2 := make([]byte, 4)
	for i := range id2 {
		id2[i] = characters[rand.Intn(len(characters))]
	}
	id3 := make([]byte, 4)
	for i := range id3 {
		id3[i] = characters[rand.Intn(len(characters))]
	}
	id4 := make([]byte, 4)
	for i := range id4 {
		id4[i] = characters[rand.Intn(len(characters))]
	}
	id5 := make([]byte, 12)
	for i := range id5 {
		id5[i] = characters[rand.Intn(len(characters))]
	}
	id := string(id1) + "-" + string(id2) + "-" + string(id3) + "-" + string(id4) + "-" + string(id5)
	return id
}
