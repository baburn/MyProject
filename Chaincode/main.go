package main

import (
	"log"

	"hiring/contracts"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	UniContract := new(contracts.UniContract)
	//OfferContract := new(contracts.OfferContract)

	chaincode, err := contractapi.NewChaincode(UniContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
