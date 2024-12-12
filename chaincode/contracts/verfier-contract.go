package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


type OfferContract struct {
	contractapi.Contract
}

type Offer struct {
	OfferId        string `json:"offerId"`
	AssetType      string `json:"assetType"`
	Ctc            string `json:"ctc"`
	FixedComp      string `json:"fixedComp"`
	VariableComp   string `json:"variableComp"`
	DateOfJoining  string `json:"dateOfJoining"`
	DateOfRelease  string `json:"dateOfRelease"`
	NameOfPerson   string `json:"nameOfPerson"`
	ContactDetails string `json:"contactDetails"`
	CompanyName    string `json:"companyName"`
}

const collectionName string = "OfferCollection"

// OfferExists returns true when asset with given ID exists in private data collection
func (o *OfferContract) OfferExists(ctx contractapi.TransactionContextInterface, offerId string) (bool, error) {

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, offerId)

	if err != nil {
		return false, fmt.Errorf("could not fetch the private data hash. %s", err)
	}

	return data != nil, nil
}

// CreateOfferLetter creates a new instance of Offer
func (o *OfferContract) CreateOfferLetter(ctx contractapi.TransactionContextInterface, offerId string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "EmployerMSP" {
		exists, err := o.OfferExists(ctx, offerId)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the asset %s already exists", offerId)
		}

		var offer Offer

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data. %s", err)
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of ctc, fixed, variable, date of joining, date of release, name of person, address of person, contact of person, company name")
		}

		ctc, exists := transientData["ctc"]
		if !exists {
			return "", fmt.Errorf("the ctc was not specified in transient data. Please try again")
		}
		offer.Ctc = string(ctc)

		fixedComp, exists := transientData["fixedComp"]
		if !exists {
			return "", fmt.Errorf("the fixedComp was not specified in transient data. Please try again")
		}
		offer.FixedComp = string(fixedComp)

		variableComp, exists := transientData["variableComp"]
		if !exists {
			return "", fmt.Errorf("the variableComp was not specified in transient data. Please try again")
		}
		offer.VariableComp = string(variableComp)

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

		nameOfPerson, exists := transientData["nameOfPerson"]
		if !exists {
			return "", fmt.Errorf("the name of person was not specified in transient data. Please try again")
		}
		offer.NameOfPerson = string(nameOfPerson)

		contactDetails, exists := transientData["contactDetails"]
		if !exists {
			return "", fmt.Errorf("the contactDetails was not specified in transient data. Please try again")
		}
		offer.ContactDetails = string(contactDetails)

		companyName, exists := transientData["companyName"]
		if !exists {
			return "", fmt.Errorf("the companyName was not specified in transient data. Please try again")
		}
		offer.CompanyName = string(companyName)

		offer.AssetType = "OfferLetter"
		offer.OfferId = offerId

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

// ReadOffer retrieves an instance of Offer from the private data collection
func (o *OfferContract) ReadOffer(ctx contractapi.TransactionContextInterface, offerId string) (*Offer, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("could not read the client identity. %s", err)
	}

	if clientOrgID == "EmployerMSP" || clientOrgID == "StudentMSP" {
		exists, err := o.OfferExists(ctx, offerId)
		if err != nil {
			return nil, fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return nil, fmt.Errorf("the asset %s does not exist", offerId)
		}

		bytes, err := ctx.GetStub().GetPrivateData(collectionName, offerId)
		if err != nil {
			return nil, fmt.Errorf("could not get the private data. %s", err)
		}
		var offer Offer

		err = json.Unmarshal(bytes, &offer)

		if err != nil {
			return nil, fmt.Errorf("could not unmarshal private data collection data to type Offer")
		}

		return &offer, nil
	}

	return nil, fmt.Errorf("%v not allowed to read.", clientOrgID)
}

// DeleteOffer deletes an instance of Offer from the private data collection
func (o *OfferContract) DeleteOffer(ctx contractapi.TransactionContextInterface, offerId string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("could not read the client identity. %s", err)
	}

	if clientOrgID == "EmployerMSP" {

		exists, err := o.OfferExists(ctx, offerId)

		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the offer %s does not exist", offerId)
		}

		return ctx.GetStub().DelPrivateData(collectionName, offerId)
	} else {
		return fmt.Errorf("organisation with %v cannot delete the offer", clientOrgID)
	}
}

func (o *OfferContract) GetAllOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, error) {
	queryString := `{"selector":{"assetType":"Offer"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return OfferResultIteratorFunction(resultsIterator)
}

func (o *OfferContract) GetOffersByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*Offer, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)

	if err != nil {
		return nil, fmt.Errorf("could not fetch the private data by range. %s", err)
	}
	defer resultsIterator.Close()

	return OfferResultIteratorFunction(resultsIterator)

}

// iterator function

func OfferResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Offer, error) {
	var offers []*Offer
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of result iterator. %s", err)
		}
		var offer Offer
		err = json.Unmarshal(queryResult.Value, &offer)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		offers = append(offers, &offer)
	}

	return offers, nil
}
