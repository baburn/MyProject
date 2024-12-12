package contracts

import (
	"fmt"

	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	
)

type OfferContract struct {
	contractapi.Contract
}

type Offer struct {
	CompanyName string `json:"companyName"`
	StudentID   string `json:"studentId"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Email       string `json:"email"`
	LPA         string `json:"lpa"`
	OfferId     string `json:"offerId"`
	Status      string `json:"status"`
}

const collectionName string = "Offers"

// Check if offer exists in the collection
func (v *OfferContract) OfferExists(ctx contractapi.TransactionContextInterface, offerId string) (bool, error) {
	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, offerId)
	if err != nil {
		return false, fmt.Errorf("could not fetch the private data hash. %s", err)
	}
	return data != nil, nil
}

// Create a new offer and store it in the private data collection
func (v *OfferContract) CreateOffer(ctx contractapi.TransactionContextInterface, offerId string) (string, error) {
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity, %v", err)
	}

	if clientOrgId != "Org2MSP" {
		return "", fmt.Errorf("only users from Org2MSP can perform this action")
	}

	// Check if the offer already exists
	exists, err := v.OfferExists(ctx, offerId)
	if err != nil {
		return "", fmt.Errorf("could not check offer existence. %s", err)
	}
	if exists {
		return "", fmt.Errorf("the offer %s already exists", offerId)
	}

	// Create new Offer object
	var newOffer Offer
	transientData, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "", fmt.Errorf("could not fetch transient data. %s", err)
	}

	// Ensure transient data contains all necessary fields
	studentID, exists := transientData["studentId"]
	if !exists {
		return "", fmt.Errorf("the studentId was not specified in transient data. Please try again")
	}
	newOffer.StudentID = string(studentID)

	name, exists := transientData["name"]
	if !exists {
		return "", fmt.Errorf("the name was not specified in transient data. Please try again")
	}
	newOffer.Name = string(name)

	email, exists := transientData["email"]
	if !exists {
		return "", fmt.Errorf("the email was not specified in transient data. Please try again")
	}
	newOffer.Email = string(email)

	lpa, exists := transientData["lpa"]
	if !exists {
		return "", fmt.Errorf("the lpa was not specified in transient data. Please try again")
	}
	newOffer.LPA = string(lpa)

	status, exists := transientData["status"]
	if !exists {
		return "", fmt.Errorf("the status was not specified in transient data. Please try again")
	}
	newOffer.Status = string(status)

	// Optional: The company name and role can be set statically or extracted from elsewhere
	// Here we assume they come from transient data too, or default values can be set
	companyName, exists := transientData["companyName"]
	if exists {
		newOffer.CompanyName = string(companyName)
	} else {
		newOffer.CompanyName = "Unknown" // Default or fallback value
	}

	role, exists := transientData["role"]
	if exists {
		newOffer.Role = string(role)
	} else {
		newOffer.Role = "Intern" // Default or fallback value
	}

	// Serialize the offer data and save it to private data collection
	offerBytes, err := json.Marshal(newOffer)
	if err != nil {
		return "", fmt.Errorf("could not marshal the offer data. %s", err)
	}

	err = ctx.GetStub().PutPrivateData(collectionName, offerId, offerBytes)
	if err != nil {
		return "", fmt.Errorf("could not create the offer. %s", err)
	}

	return fmt.Sprintf("Successfully added offer %s for student %s", offerId, newOffer.StudentID), nil
}
