package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ResultContract struct {
	contractapi.Contract
}

type Result struct {
	ResultId      string `json:"resultId"`
	ResultType    string `json:"resultType"`
	StudentId     string `json:"studentId"`
	StudentRollNo string `json:"studentRollNo"`
	TotalMarks    string `json:"totalMarks"`
	ObtainedMarks string `json:"obtainedMarks"`
	Percentage    string `json:"percentage"`
	Grade         string `json:"grade"`
	Conclusion    string `json:"conclusion"`
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

// Check result for particular id exists or not
func (r *ResultContract) ResultExists(ctx contractapi.TransactionContextInterface, resultId string) (bool, error) {
	data, err := ctx.GetStub().GetState(resultId)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

// Create result for the student
func (r *ResultContract) CreateResult(ctx contractapi.TransactionContextInterface, resultId string, resultType string, studentId string, studentRollNo string, totalMarks string, obtainedMarks string, percentage string, grade string, conclusion string) (string, error) {
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "", fmt.Errorf("cannot fetch client identity. %s", err)
	}

	if clientOrgId == "InstitutionMSP" {
		exists, err := r.ResultExists(ctx, resultId)

		if err != nil {
			return "", fmt.Errorf("error: %v", err)
		}

		if !exists {
			result := Result{
				ResultId:      resultId,
				ResultType:    resultType,
				StudentId:     studentId,
				StudentRollNo: studentRollNo,
				TotalMarks:    totalMarks,
				ObtainedMarks: obtainedMarks,
				Percentage:    percentage,
				Grade:         grade,
				Conclusion:    conclusion,
			}

			bytes, _ := json.Marshal(result)
			fmt.Println(bytes)

			err = ctx.GetStub().PutState(resultId, bytes)
			if err != nil {
				return "", fmt.Errorf("could not create result for %s", studentId)
			}

			return fmt.Sprintf("Successfully created result for student %s and result id %s", studentId, resultId), nil
		} else {
			return "", fmt.Errorf("Result with id %v already exists", resultId)
		}

	}

	return "", fmt.Errorf("user %v cannot perform action", clientOrgId)
}

// Read result of the student
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

// Delete result from the world state
func (r *ResultContract) DeleteResult(ctx contractapi.TransactionContextInterface, resultId string) (string, error) {
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgId == "InstitutionMSP" {

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

// get result by range from start to end
func (r *ResultContract) GetResultsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Result, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the data by range. %s", err)
	}
	defer resultsIterator.Close()

	return resultIteratorFunction(resultsIterator)
}

// get all results at a time
func (r *ResultContract) GetAllResults(ctx contractapi.TransactionContextInterface) ([]*Result, error) {

	queryString := `{"selector":{"resultType":"12_MARKSHEET"}, "sort":[{ "percentage": "desc"}]}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return resultIteratorFunction(resultsIterator)
}

// helper function to iterate result
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

// to get history of particular result
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

	queryString := `{"selector":{"resultType":"12_MARKSHEET"}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the result records. %s", err)
	}
	defer resultsIterator.Close()

	results, err := resultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the result records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             results,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}
