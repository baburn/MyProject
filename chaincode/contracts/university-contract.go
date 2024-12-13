package contracts

import (
	"encoding/json"
	"fmt"
	"time"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ResultContract struct {
	contractapi.Contract
}

type Result struct {
	ResultId      string `json:"resultId"`
	StudentId     string `json:"studentId"`
	TotalMarks    string `json:"totalMarks"`
	ObtainedMarks string `json:"obtainedMarks"`
	Percentage    string `json:"percentage"`
	Status    string `json:"status"`
}

type HistoryQueryResult struct {
	Record    *Result `json:"record"`
	TxId      string  `json:"txId"`
	Timestamp string  `json:"timestamp"`
	IsDelete  bool    `json:"isDelete"`
}

type PaginatedQueryResult struct {
	Records             []*Result `json:"records"`
	FetchedRecordsCount int32     `json:"fetchedRecordsCount"`
	Bookmark            string    `json:"bookmark"`
}

func (r *ResultContract) ResultExists(ctx contractapi.TransactionContextInterface, resultId string) (bool, error) {
	data, err := ctx.GetStub().GetState(resultId)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

func (r *ResultContract) CreateResult(ctx contractapi.TransactionContextInterface, resultId string, studentId string, totalMarks string, obtainedMarks string, percentage string, status string) (string, error) {
	// Input validation
	if strings.TrimSpace(resultId) == "" || strings.TrimSpace(studentId) == "" {
		return "", fmt.Errorf("resultId and studentId cannot be empty")
	}

	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve client identity: %v", err)
	}

	if clientOrgId != "UniversityMSP" {
		return "", fmt.Errorf("unauthorized organization %v cannot create results", clientOrgId)
	}

	exists, err := r.ResultExists(ctx, resultId)
	if err != nil {
		return "", fmt.Errorf("error checking result existence: %v", err)
	}

	if exists {
		return "", fmt.Errorf("result with ID %s already exists", resultId)
	}

	result := Result{
		ResultId:      resultId,
		StudentId:     studentId,
		TotalMarks:    totalMarks,
		ObtainedMarks: obtainedMarks,
		Percentage:    percentage,
		Status:    status,
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize result: %v", err)
	}

	err = ctx.GetStub().PutState(resultId, bytes)
	if err != nil {
		return "", fmt.Errorf("could not create result for student %s: %v", studentId, err)
	}

	return fmt.Sprintf("Successfully created result for student %s with result ID %s", studentId, resultId), nil
}

func (r *ResultContract) ReadResult(ctx contractapi.TransactionContextInterface, resultId string) (*Result, error) {
	bytes, err := ctx.GetStub().GetState(resultId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the result does not exist for result id %v", resultId)
	}

	var result Result

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type result")
	}

	return &result, nil
}

func (r *ResultContract) DeleteResult(ctx contractapi.TransactionContextInterface, resultId string) (string, error) {
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgId == "UniversityMSP" {

		history, err := r.GetResultHistory(ctx, resultId)
		if err != nil {
			return "", fmt.Errorf("could not fetch history: %v", err)
		}
		if len(history) > 0 && history[0].IsDelete {
			return "", fmt.Errorf("result with id %v has been deleted", resultId)
		}

		exists, err := r.ResultExists(ctx, resultId)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if !exists {
			return "", fmt.Errorf("the result, %s does not exist", resultId)
		}

		err = ctx.GetStub().DelState(resultId)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("result with id %v is deleted from the world state.", resultId), nil
		}
	}

	return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgId)
}

func (r *ResultContract) GetResultsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Result, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the data by range. %s", err)
	}
	defer resultsIterator.Close()

	return resultIteratorFunction(resultsIterator)
}

func (r *ResultContract) GetAllResults(ctx contractapi.TransactionContextInterface) ([]*Result, error) {

	queryString := `{"selector":{"status":"Pass"}, "sort":[{ "percentage": "desc"}]}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return resultIteratorFunction(resultsIterator)
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

func (r *ResultContract) GetResultsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {
	// Fixed query string with proper closing
	queryString := `{"selector":{"status":"Pass"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve result records: %v", err)
	}
	defer resultsIterator.Close()

	results, err := resultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not process result records: %v", err)
	}

	return &PaginatedQueryResult{
		Records:             results,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}
