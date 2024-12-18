package main

import (
	"encoding/json"
	"log"
	"sync"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Result struct {
	ResultId      string `json:"resultId"`
	StudentId     string `json:"studentId"`
	TotalMarks    string `json:"totalMarks"`
	ObtainedMarks string `json:"obtainedMarks"`
	Percentage    string `json:"percentage"`
	Status        string `json:"status"`
}

type Offer struct {
	OfferId       string `json:"offerId"`
	StudentId     string `json:"studentId"`
	AssetType     string `json:"assetType"`
	Ctc           string `json:"ctc"`
	DateOfJoining string `json:"dateOfJoining"`
	DateOfRelease string `json:"dateOfRelease"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	CompanyName   string `json:"companyName"`
}

type ResultData struct {
	AssetType     string `json:"AssetType"`
	ResultId      string `json:"ResultId"`
	StudentId     string `json:"StudentId"`
	DateOfResult  string `json:"DateOfResult"`
	University    string `json:"University"`
	ObtainedMarks string `json:"ObtainedMarks"`
	Percentage    string `json:"Percentage"`
	Status        string `json:"Status"`
}

type OfferData struct {
	OfferId     string `json:"offerId"`
	StudentId   string `json:"studentId"`
	AssetType   string `json:"assetType"`
	Status      string `json:"status"`
	CompanyName string `json:"companyName"`
	Ctc         string `json:"ctc"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

type Match struct {
	OfferId  string `json:"offerId"`
	ResultId string `json:"resultId"`
}

type ResultHistory struct {
	Record    *ResultData `json:"record"`
	TxId      string      `json:"txId"`
	Timestamp string      `json:"timestamp"`
	IsDelete  bool        `json:"isDelete"`
}

// Placeholder for ChaincodeEventListener function
func ChaincodeEventListener(org, channel, chaincodeName string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Implement actual event listener logic
	log.Println("Listening for chaincode events")
}

// Placeholder for getEvents function
func getEvents() []string {
	// Implement actual event retrieval logic
	return []string{}
}

func main() {
	router := gin.Default()

	var wg sync.WaitGroup
	wg.Add(1)
	go ChaincodeEventListener("university", "mychannel", "Credential-Verification", &wg)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to project",
		})
	})

	// Result-related routes
	router.GET("/api/results", func(ctx *gin.Context) {
		result := submitTxnFn("university", "mychannel", "Credential-Verification", "ResultContract", "query", make(map[string][]byte), "GetAllResults")
		var results []ResultData
		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &results); err != nil {
				log.Println("Error:", err)
				ctx.JSON(500, gin.H{"error": "Failed to parse results"})
				return
			}
		}

		ctx.JSON(200, results)
	})

	router.POST("/api/result", func(ctx *gin.Context) {
		var req Result
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"message": "Bad request"})
			return
		}

		log.Printf("Result received: %+v", req)
		res := submitTxnFn("university", "mychannel", "Credential-Verification", "ResultContract", "invoke", make(map[string][]byte), "CreateResult",
			req.ResultId, req.StudentId, req.TotalMarks, req.ObtainedMarks, req.Percentage, req.Status)

		ctx.JSON(200, res)
	})

	router.GET("/api/result/:id", func(ctx *gin.Context) {
		resultId := ctx.Param("id")
		result := submitTxnFn("university", "mychannel", "Credential-Verification", "ResultContract", "query", make(map[string][]byte), "ReadResult", resultId)

		var singleResult Result

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &singleResult); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.JSON(200, gin.H{"singleData": singleResult})
	})

	// Offer-related routes
	router.POST("/api/offer", func(ctx *gin.Context) {
		var req Offer
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"message": "Invalid request format"})
			return
		}

		if req.OfferId == "" || req.StudentId == "" {
			ctx.JSON(400, gin.H{"message": "OfferId and StudentId are required"})
			return
		}

		privateData := map[string][]byte{
			"ctc":           []byte(req.Ctc),
			"dateOfJoining": []byte(req.DateOfJoining),
			"dateOfRelease": []byte(req.DateOfRelease),
			"name":          []byte(req.Name),
			"email":         []byte(req.Email),
			"companyName":   []byte(req.CompanyName),
		}

		log.Printf("Creating offer with data: %+v", privateData)
		res := submitTxnFn("company", "mychannel", "Credential-Verification", "OfferContract", "private", privateData, "CreateOffer", req.OfferId, req.StudentId)
		ctx.JSON(200, gin.H{"response": res})
	})

	router.GET("/api/offer/:id", func(ctx *gin.Context) {
		offerId := ctx.Param("id")
		result := submitTxnFn("company", "mychannel", "Credential-Verification", "OfferContract", "query", nil, "ReadOffer", offerId)

		var offer Offer
		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &offer); err != nil {
				log.Printf("Error unmarshalling offer: %v", err)
				ctx.JSON(500, gin.H{"error": "Failed to parse offer"})
				return
			}
		}

		// New check: Ensure the student can only see their own offer
		studentId := ctx.DefaultQuery("studentId", "") // Assuming student ID is passed as a query parameter
		if studentId != offer.StudentId {
			ctx.JSON(403, gin.H{"error": "You are not authorized to view this offer"})
			return
		}

		ctx.JSON(200, gin.H{"offer": offer})
	})

	router.GET("/api/offers", func(ctx *gin.Context) {
		result := submitTxnFn("company", "mychannel", "Credential-Verification", "OfferContract", "query", nil, "GetAllOffers")

		var offers []OfferData
		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &offers); err != nil {
				log.Printf("Error parsing offers: %v", err)
				ctx.JSON(500, gin.H{"error": "Failed to parse offers"})
				return
			}
		}
		ctx.JSON(200, offers)
	})

	// Matching and Events
	router.POST("/api/result/match-offer", func(ctx *gin.Context) {
		var req Match
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"message": "Bad request"})
			return
		}

		log.Printf("Match request: %+v", req)
		submitTxnFn("university", "mychannel", "Credential-Verification", "ResultContract", "invoke", make(map[string][]byte), "MatchOffer", req.ResultId, req.OfferId)

		ctx.JSON(200, req)
	})

	router.GET("/api/events", func(ctx *gin.Context) {
		result := getEvents()
		ctx.JSON(200, gin.H{"events": result})
	})

	// Start the server
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for chaincode event listener to complete
	wg.Wait()
}
