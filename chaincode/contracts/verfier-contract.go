package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// OfferContract defines the smart contract for managing job offer letters
type OfferContract struct {
	contractapi.Contract
}

// Offer represents the structure of a job offer letter
type Offer struct {
	OfferId        string `json:"offerId"`        // Unique identifier for the offer
	AssetType      string `json:"assetType"`      // Type of asset (e.g., "OfferLetter")
	Ctc            string `json:"ctc"`            // Cost to Company (compensation details)
	DateOfJoining  string `json:"dateOfJoining"`  // Date when the employee will start
	DateOfRelease  string `json:"dateOfRelease"`  // Date of offer letter release
	Name           string `json:"name"`           // Name of the offer recipient
	Email          string `json:"email"`          // Email of the offer recipient
	CompanyName    string `json:"companyName"`    // Name of the company making the offer
}

// Collection name for private data storage
const collectionName string = "Offers"

// OfferExists checks if an offer with the given ID exists in the private data collection
func (o *OfferContract) OfferExists(ctx contractapi.TransactionContextInterface, offerId string) (bool, error) {
	// Retrieve the private data hash to check existence
	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, offerId)

	if err != nil {
		return false, fmt.Errorf("could not fetch the private data hash. %s", err)
	}

	return data != nil, nil
}

// CreateOffer adds a new offer letter to the private data collection
func (o *OfferContract) CreateOffer(ctx contractapi.TransactionContextInterface, offerId string) (string, error) {
	// Verify client organization identity
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	// Restrict offer creation to CompanyMSP
	if clientOrgID == "CompanyMSP" {
		// Check if offer already exists
		exists, err := o.OfferExists(ctx, offerId)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the asset %s already exists", offerId)
		}

		var offer Offer

		// Retrieve transient data (sensitive information)
		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data. %s", err)
		}

		// Validate transient data is not empty
		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of ctc, fixed, variable, date of joining, date of release, name of person, address of person, contact of person, company name")
		}

		// Extract and validate each piece of transient data
		ctc, exists := transientData["ctc"]
		if !exists {
			return "", fmt.Errorf("the ctc was not specified in transient data. Please try again")
		}
		offer.Ctc = string(ctc)

		dateOfJoining, exists := transientData["dateOfJoining"]
		if !exists {
			return "", fmt.Errorf("the dealer was not specified in transient data. Please try again")
		}
		offer.DateOfJoining = string(dateOfJoining)

		dateOfRelease, exists := transientData["dateOfRelease"]
		if !exists {
			return "", fmt.Errorf("the date of release was not specified in transient data. Please try again")
		}
		offer.DateOfRelease = string(dateOfRelease)

		companyName, exists := transientData["companyName"]
		if !exists {
			return "", fmt.Errorf("the companyName was not specified in transient data. Please try again")
		}
		offer.CompanyName = string(companyName)

		// Set additional offer details
		offer.AssetType = "OfferLetter"
		offer.OfferId = offerId

		// Serialize and store offer in private data collection
		bytes, _ := json.Marshal(offer)
		err = ctx.GetStub().PutPrivateData(collectionName, offerId, bytes)
		if err != nil {
			return "", fmt.Errorf("could not able to write the data")
		}
		return fmt.Sprintf("offer with id %v added successfully", offerId), nil
	} else {
		return fmt.Sprintf("offer cannot be created by organisation with MSPID %v ", clientOrgID), nil
	}
}


// ReadOffer retrieves an offer letter from the private data collection
func (o *OfferContract) ReadOffer(ctx contractapi.TransactionContextInterface, offerId string) (*Offer, error) {
	// Verify client organization identity
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("could not read the client identity. %s", err)
	}

	// Allow reading for CompanyMSP and StudentMSP
	if clientOrgID == "CompanyMSP" || clientOrgID == "StudentMSP" {
		// Check if offer exists
		exists, err := o.OfferExists(ctx, offerId)
		if err != nil {
			return nil, fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return nil, fmt.Errorf("the asset %s does not exist", offerId)
		}

		// Retrieve private data
		bytes, err := ctx.GetStub().GetPrivateData(collectionName, offerId)
		if err != nil {
			return nil, fmt.Errorf("could not get the private data. %s", err)
		}
		var offer Offer

		// Deserialize offer data
		err = json.Unmarshal(bytes, &offer)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal private data collection data to type Offer")
		}

		return &offer, nil
	}

	return nil, fmt.Errorf("%v not allowed to read.", clientOrgID)
}

// DeleteOffer removes an offer letter from the private data collection
func (o *OfferContract) DeleteOffer(ctx contractapi.TransactionContextInterface, offerId string) error {
	// Verify client organization identity
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("could not read the client identity. %s", err)
	}

	// Restrict deletion to CompanyMSP
	if clientOrgID == "CompanyMSP" {
		// Check if offer exists
		exists, err := o.OfferExists(ctx, offerId)
		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the offer %s does not exist", offerId)
		}

		// Delete offer from private data collection
		return ctx.GetStub().DelPrivateData(collectionName, offerId)
	} else {
		return fmt.Errorf("organisation with %v cannot delete the offer", clientOrgID)
	}
}

// GetAllOffers retrieves all offers of a specific asset type
func (o *OfferContract) GetAllOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, error) {
	// Query string to select offers by asset type
	queryString := `{"selector":{"assetType":"Offer"}}`
	
	// Execute query on private data collection
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()

	// Process and return results
	return OfferResultIteratorFunction(resultsIterator)
}

// GetOffersByRange retrieves offers within a specified key range
func (o *OfferContract) GetOffersByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*Offer, error) {
	// Retrieve offers from private data collection within the specified key range
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the private data by range. %s", err)
	}
	defer resultsIterator.Close()

	// Process and return results
	return OfferResultIteratorFunction(resultsIterator)
}

// OfferResultIteratorFunction is a helper function to process query iterators and convert results
func OfferResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Offer, error) {
	var offers []*Offer
	
	// Iterate through query results
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of result iterator. %s", err)
		}
		
		// Deserialize each offer
		var offer Offer
		err = json.Unmarshal(queryResult.Value, &offer)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		
		// Append to results
		offers = append(offers, &offer)
	}

	return offers, nil
}