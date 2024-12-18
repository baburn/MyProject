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

// PaginatedQueryResult supports paginated queries of results
type PaginatedQueryResult struct {
	Records             []*Result `json:"records"`               // List of result records
	FetchedRecordsCount int32     `json:"fetchedRecordsCount"`   // Number of records fetched
	Bookmark            string    `json:"bookmark"`              // Bookmark for pagination
}

// HistoryQueryResult contains the result along with its transaction history
type HistoryQueryResult struct {
	Record    *Result `json:"record"`    // The result record
	TxId      string  `json:"txId"`      // Transaction ID
	Timestamp string  `json:"timestamp"` // Timestamp of the transaction
	IsDelete  bool    `json:"isDelete"`  // Indicates if the record was deleted
}

// Result represents the structure of a student's academic result
type Result struct {
	AssetType      string `json:"assetType"`      // Asset type ("Result")
	ResultId       string `json:"resultId"`       // Unique identifier for the result
	StudentId      string `json:"studentId"`      // Identifier for the student
	TotalMarks     string `json:"totalMarks"`     // Total possible marks
	ObtainedMarks  string `json:"obtainedMarks"`  // Marks obtained by the student
	Percentage     string `json:"percentage"`     // Calculated percentage
	Status         string `json:"status"`         // Pass/Fail status
}

// EventData represents metadata for blockchain events
type EventData struct {
	Type   string // Type of event
	Percentage string // Percentage of the event
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
		return "", fmt.Errorf("could not fetch result: %s", err)
	}
	if exists {
		return "", fmt.Errorf("result with ID %s already exists", resultId)
	}

	// Create a new Result
	result := Result{
		AssetType:     "Result",
		ResultId:      resultId,
		StudentId:     studentId,
		TotalMarks:    totalMarks,
		ObtainedMarks: obtainedMarks,
		Percentage:    percentage,
		Status:        status,
	}

	// Store the result in the world state
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}

	err = ctx.GetStub().PutState(resultId, resultBytes)
	if err != nil {
		return "", fmt.Errorf("failed to store result in world state: %v", err)
	}

	// Trigger an event after creating the result
	eventData := EventData{
		Type:   "Result creation",
		Percentage: percentage,
	}
	eventBytes, _ := json.Marshal(eventData)
	ctx.GetStub().SetEvent("CreateResult", eventBytes)

	return fmt.Sprintf("Successfully added result %v", resultId), nil
}

// ReadResult retrieves an instance of Result from the world state
func (r *ResultContract) ReadResult(ctx contractapi.TransactionContextInterface, resultId string) (*Result, error) {
	resultBytes, err := ctx.GetStub().GetState(resultId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if resultBytes == nil {
		return nil, fmt.Errorf("the result with ID %s does not exist", resultId)
	}

	var result Result
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal result: %v", err)
	}

	return &result, nil
}

// DeleteResult removes the result from the world state
func (r *ResultContract) DeleteResult(ctx contractapi.TransactionContextInterface, resultId string) (string, error) {
	// Verify client organization identity
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity: %s", err)
	}

	// Allow deletion only for UniversityMSP
	if clientOrgID != "UniversityMSP" {
		return "", fmt.Errorf("user under the following MSPID: %v can't perform this action", clientOrgID)
	}

	// Check if the result exists
	exists, err := r.ResultExists(ctx, resultId)
	if err != nil {
		return "", fmt.Errorf("could not check result existence: %s", err)
	} else if !exists {
		return "", fmt.Errorf("the result with ID %s does not exist", resultId)
	}

	// Delete the result from the ledger
	err = ctx.GetStub().DelState(resultId)
	if err != nil {
		return "", fmt.Errorf("failed to delete result: %v", err)
	}

	return fmt.Sprintf("Successfully deleted result with ID %s", resultId), nil
}

func (r *ResultContract) GetResultsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Result, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the data by range. %s", err)
	}
	defer resultsIterator.Close()

	return resultIteratorFunction(resultsIterator)
}
// GetAllResults retrieves all results
func (r *ResultContract) GetAllResults(ctx contractapi.TransactionContextInterface) ([]*Result, error) {
	queryString := `{"selector":{"assetType":"Result"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch results: %s", err)
	}
	defer resultsIterator.Close()

	var results []*Result
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch result: %s", err)
		}

		var result Result
		err = json.Unmarshal(queryResult.Value, &result)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal result: %v", err)
		}

		results = append(results, &result)
	}

	return results, nil
}

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

// GetResultHistory retrieves the history of a result
func (r *ResultContract) GetResultHistory(ctx contractapi.TransactionContextInterface, resultId string) ([]*HistoryQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(resultId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch result history: %v", err)
	}
	defer resultsIterator.Close()

	var history []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch result history: %v", err)
		}

		var result Result
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &result)
			if err != nil {
				return nil, fmt.Errorf("could not unmarshal result history: %v", err)
			}
		} else {
			result = Result{ResultId: resultId}
		}

		timestamp := response.Timestamp.AsTime()
		formattedTime := timestamp.Format(time.RFC1123)

		historyRecord := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &result,
			IsDelete:  response.IsDelete,
		}
		history = append(history, &historyRecord)
	}

	return history, nil
}

// GetResultsWithPagination retrieves results with pagination
func (r *ResultContract) GetResultsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {
	queryString := `{"selector":{"assetType":"Result"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not fetch results with pagination: %v", err)
	}
	defer resultsIterator.Close()

	var results []*Result
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch result: %s", err)
		}

		var result Result
		err = json.Unmarshal(queryResult.Value, &result)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal result: %v", err)
		}

		results = append(results, &result)
	}

	return &PaginatedQueryResult{
		Records:             results,
		FetchedRecordsCount: int32(len(results)),
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (r *ResultContract) GetMatchingResults(ctx contractapi.TransactionContextInterface, resultID string) ([]*Result, error) {
	// Read the base result
	baseResult, err := r.ReadResult(ctx, resultID)
	if err != nil {
		return nil, fmt.Errorf("error reading result %v", err)
	}

	// Construct the query string to match the base result fields
	queryString := fmt.Sprintf(
		`{"selector":{"assetType":"Result","TotalMarks":"%s", "ObtainedMarks": "%s", "Percentage":"%s", "Status":"%s"}}`,
		baseResult.TotalMarks, baseResult.ObtainedMarks, baseResult.Percentage, baseResult.Status,
	)

	// Execute the query in the collection
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not get the data: %v", err)
	}
	defer resultsIterator.Close()

	// Process the query results
	return resultIteratorFunction(resultsIterator)
}

// MatchResult matches result with a matching target result
func (r *ResultContract) MatchResult(ctx contractapi.TransactionContextInterface, resultID string, targetResultID string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity: %s", err)
	}

	// Ensure action is restricted to UniversityMSP
	if clientOrgID != "UniversityMSP" {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}

	// Fetch target result from private data
	bytes, err := ctx.GetStub().GetPrivateData(collectionName, targetResultID)
	if err != nil {
		return "", fmt.Errorf("could not fetch private data: %s", err)
	}

	var targetResult Result
	err = json.Unmarshal(bytes, &targetResult)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal target result data: %s", err)
	}

	// Read primary result
	result, err := r.ReadResult(ctx, resultID)
	if err != nil {
		return "", fmt.Errorf("could not read result data: %s", err)
	}

	// Match results based on attributes
	if result.TotalMarks == targetResult.TotalMarks && result.ObtainedMarks == targetResult.ObtainedMarks && result.Percentage == targetResult.Percentage {
		result.Status = "Assigned"
		result.StudentId = targetResult.StudentId

		// Serialize the updated result
		bytes, _ := json.Marshal(result)

		// Delete the matched target result from private data
		err = ctx.GetStub().DelPrivateData(collectionName, targetResultID)
		if err != nil {
			return "", fmt.Errorf("could not delete target result: %s", err)
		}

		// Update the result state
		err = ctx.GetStub().PutState(resultID, bytes)
		if err != nil {
			return "", fmt.Errorf("could not update result state: %s", err)
		}

		return fmt.Sprintf("Deleted target result %v and assigned result %v to student %v", targetResultID, resultID, targetResult.StudentId), nil
	} else {
		return "", fmt.Errorf("target result does not match")
	}
}

// ConfirmResult confirms the student's result for a company
func (r *ResultContract) ConfirmResult(ctx contractapi.TransactionContextInterface, resultID string, companyName string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not get the MSPID: %s", err)
	}

	// Ensure only authorized organizations can perform this action
	if clientOrgID != "UniversityMSP" {
		return "", fmt.Errorf("user under following MSPID: %v cannot perform this action", clientOrgID)
	}

	// Read the result
	result, err := r.ReadResult(ctx, resultID)
	if err != nil {
		return "", fmt.Errorf("could not read the result: %s", err)
	}

	// Update the result status and add company details
	result.Status = fmt.Sprintf("Confirmed for %v", companyName)

	bytes, _ := json.Marshal(result)

	// Save the updated result
	err = ctx.GetStub().PutState(resultID, bytes)
	if err != nil {
		return "", fmt.Errorf("could not update result state: %s", err)
	}

	return fmt.Sprintf("Result %v successfully confirmed for %v", resultID, companyName), nil
}

// AddStudentResult stores the student's result in the blockchain
func (r *ResultContract) AddStudentResult(ctx contractapi.TransactionContextInterface, studentId string, percentage string, status string) error {
    studentResult := Result{
        StudentId: studentId,
        Percentage: percentage,
        Status: status,
    }

    // Convert the studentResult to JSON
    studentResultJSON, err := json.Marshal(studentResult)
    if err != nil {
        return fmt.Errorf("failed to marshal student result: %v", err)
    }

    // Store the student result in the ledger
    return ctx.GetStub().PutState(studentId, studentResultJSON)
}




