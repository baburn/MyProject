package contracts

import (
	"encoding/json"
	"fmt"
	"time"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ResultContract defines the smart contract for managing student results
type ResultContract struct {
	contractapi.Contract
}

// Result represents the structure of a student's academic result
type Result struct {
	ResultId      string `json:"resultId"`      // Unique identifier for the result
	StudentId     string `json:"studentId"`     // Identifier for the student
	TotalMarks    string `json:"totalMarks"`    // Total possible marks
	ObtainedMarks string `json:"obtainedMarks"` // Marks obtained by the student
	Percentage    string `json:"percentage"`    // Calculated percentage
	Status        string `json:"status"`        // Pass/Fail status
}

// HistoryQueryResult contains the result along with its transaction history
type HistoryQueryResult struct {
	Record    *Result `json:"record"`    // The result record
	TxId      string  `json:"txId"`      // Transaction ID
	Timestamp string  `json:"timestamp"` // Timestamp of the transaction
	IsDelete  bool    `json:"isDelete"`  // Indicates if the record was deleted
}

// PaginatedQueryResult supports paginated queries of results
type PaginatedQueryResult struct {
	Records             []*Result `json:"records"`               // List of result records
	FetchedRecordsCount int32     `json:"fetchedRecordsCount"`   // Number of records fetched
	Bookmark            string    `json:"bookmark"`              // Bookmark for pagination
}

// EventData represents metadata for blockchain events
type EventData struct {
	Type   string // Type of event
	Status string // Status of the event
}

// ResultExists checks if a result with the given ID already exists in the blockchain
func (r *ResultContract) ResultExists(ctx contractapi.TransactionContextInterface, resultId string) (bool, error) {
	data, err := ctx.GetStub().GetState(resultId)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

// CreateResult adds a new result to the blockchain with access control
func (r *ResultContract) CreateResult(ctx contractapi.TransactionContextInterface, resultId string, studentId string, totalMarks string, obtainedMarks string, percentage string, status string) (string, error) {
	// Validate input parameters
	if strings.TrimSpace(resultId) == "" || strings.TrimSpace(studentId) == "" {
		return "", fmt.Errorf("resultId and studentId cannot be empty")
	}

	// Verify client organization identity
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve client identity: %v", err)
	}

	// Restrict result creation to specific organization
	if clientOrgId != "UniversityMSP" {
		return "", fmt.Errorf("unauthorized organization %v cannot create results", clientOrgId)
	}

	// Check if result already exists
	exists, err := r.ResultExists(ctx, resultId)
	if err != nil {
		return "", fmt.Errorf("error checking result existence: %v", err)
	}

	if exists {
		return "", fmt.Errorf("result with ID %s already exists", resultId)
	}

	// Create result object
	result := Result{
		ResultId:      resultId,
		StudentId:     studentId,
		TotalMarks:    totalMarks,
		ObtainedMarks: obtainedMarks,
		Percentage:    percentage,
		Status:        status,
	}

	// Serialize result to JSON
	bytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize result: %v", err)
	}

	// Store result in blockchain
	err = ctx.GetStub().PutState(resultId, bytes)
	if err != nil {
		return "", fmt.Errorf("could not create result for student %s: %v", studentId, err)
	} else {
		// Emit blockchain event for result creation
		addResultEventData := EventData{
			Type:   "Result creation",
			Status: status,
		}
		eventDataByte, _ := json.Marshal(addResultEventData)
		ctx.GetStub().SetEvent("CreateResult", eventDataByte)

		return fmt.Sprintf("successfully added result %v", resultId), nil
	}
}

// ReadResult retrieves a specific result by its ID
func (r *ResultContract) ReadResult(ctx contractapi.TransactionContextInterface, resultId string) (*Result, error) {
	bytes, err := ctx.GetStub().GetState(resultId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the result does not exist for result id %v", resultId)
	}

	var result Result

	// Deserialize JSON to Result struct
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type result")
	}

	return &result, nil
}

// DeleteResult removes a result from the blockchain with access control
func (r *ResultContract) DeleteResult(ctx contractapi.TransactionContextInterface, resultId string) (string, error) {
	// Verify client organization identity
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	// Restrict deletion to specific organization
	if clientOrgId == "UniversityMSP" {
		// Check result history to prevent duplicate deletions
		history, err := r.GetResultHistory(ctx, resultId)
		if err != nil {
			return "", fmt.Errorf("could not fetch history: %v", err)
		}
		if len(history) > 0 && history[0].IsDelete {
			return "", fmt.Errorf("result with id %v has been deleted", resultId)
		}

		// Verify result exists
		exists, err := r.ResultExists(ctx, resultId)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if !exists {
			return "", fmt.Errorf("the result, %s does not exist", resultId)
		}

		// Delete result from blockchain
		err = ctx.GetStub().DelState(resultId)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("result with id %v is deleted from the world state.", resultId), nil
		}
	}

	return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgId)
}

// GetResultsByRange retrieves results within a specified key range
func (r *ResultContract) GetResultsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Result, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the data by range. %s", err)
	}
	defer resultsIterator.Close()

	return resultIteratorFunction(resultsIterator)
}

// GetAllResults retrieves all passing results, sorted by percentage in descending order
func (r *ResultContract) GetAllResults(ctx contractapi.TransactionContextInterface) ([]*Result, error) {
	// Query to select only passing results and sort by percentage
	queryString := `{"selector":{"status":"Pass"}, "sort":[{ "percentage": "desc"}]}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return resultIteratorFunction(resultsIterator)
}

// resultIteratorFunction is a helper function to process query iterators and convert results
func resultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Result, error) {
	var results []*Result
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var result Result
		err = json.Unmarshal(queryResult.Value, &result)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		results = append(results, &result)
	}

	return results, nil
}

// GetResultHistory retrieves the transaction history for a specific result
func (r *ResultContract) GetResultHistory(ctx contractapi.TransactionContextInterface, resultId string) ([]*HistoryQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(resultId)
	if err != nil {
		return nil, fmt.Errorf("could not get the data. %s", err)
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not get the value of resultsIterator. %s", err)
		}

		var result Result
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &result)
			if err != nil {
				return nil, err
			}
		} else {
			result = Result{
				ResultId: resultId,
			}
		}

		// Format timestamp
		timestamp := response.Timestamp.AsTime()
		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &result,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}

// GetResultsWithPagination supports retrieving results with pagination
func (r *ResultContract) GetResultsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {
	// Query to select only passing results
	queryString := `{"selector":{"status":"Pass"}}`

	// Retrieve results with pagination
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve result records: %v", err)
	}
	defer resultsIterator.Close()

	// Process results
	results, err := resultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not process result records: %v", err)
	}

	// Return paginated query result
	return &PaginatedQueryResult{
		Records:             results,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}