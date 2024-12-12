package contracts

import (
	"encoding/json"

	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type UniContract struct {
	contractapi.Contract
}

type Student struct {
	EntityType   string `json:"entityType"`
	StudentId    string `json:"studentId"`
	Name         string `json:"name"`
	Degree       string `json:"degree"`
	Branch       string `json:"branch"`
	Percentage   string `json:"percentage"`
	Status       string `json:"status"` //passed or failed or on-going
	DateOfIssued string `json:"dateOfissued"`
}

type HistoryQueryResult struct {
	Record    *Student `json:"record"`
	Txid      string   `json:"txId"`
	TimeStanp string   `json:"timestap"`
	IsDelete  bool     `json:"isDelete"`
}

type PaginatedQueryResult struct {
	Records             []*Student `json:"records"`
	FetchedRecordsCount int32      `json:"fetchedRecordsCount"`
	Bookmark            string     `json:"bookmark"`
}

func (u *UniContract) StudentExists(ctx contractapi.TransactionContextInterface, studentId string) (bool, error) {
	data, err := ctx.GetStub().GetState(studentId)

	if err != nil {
		return false, fmt.Errorf("error finding the Student with matching Id: %v", err)
	}
	return data != nil, nil

}

func (u *UniContract) CreateStudent(ctx contractapi.TransactionContextInterface, studentId string, name string, Degree string, branch string, percentage string, status string, dateOfissued string) (string, error) {
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgId == "Org1MSP" {
		exists, err := u.StudentExists(ctx, studentId)
		if err != nil {
			return "", fmt.Errorf("could not fetch the details from world state. %s", err)
		}
		if exists {
			return "", fmt.Errorf("the student, %s already exists", studentId)
		}

		student := Student{
			EntityType:   "degree",
			StudentId:    studentId,
			Name:         name,
			Degree:       Degree, // Capitalize the parameter to match the struct field
			Branch:       branch,
			Percentage:   percentage,
			Status:       status,
			DateOfIssued: dateOfissued,
		}

		// Handle potential error from json.Marshal
		bytes, err := json.Marshal(student)
		if err != nil {
			return "", fmt.Errorf("could not marshal student data. %s", err)
		}

		err = ctx.GetStub().PutState(studentId, bytes)
		if err != nil {
			return "", fmt.Errorf("could not create student. %s", err)
		}

		return fmt.Sprintf("successfully added student %v", student), nil
	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgId) // Fix the variable name here
	}
}

func (u *UniContract) DeleteStudent(ctx contractapi.TransactionContextInterface, studentId string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch the client identity: %v", err)
	}

	if clientOrgID == "Org1MSP" {
		exists, err := u.StudentExists(ctx, studentId)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if !exists {
			return "", fmt.Errorf("the student, %s does not exist", studentId)
		}

		err = ctx.GetStub().DelState(studentId)
		if err != nil {
			return "", fmt.Errorf("could not delete student. %s", err)
		}

		return fmt.Sprintf("student with id %v is deleted from the world state.", studentId), nil
	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

func (u *UniContract) ReadStudent(ctx contractapi.TransactionContextInterface, studentId string) (*Student, error) {
	bytes, err := ctx.GetStub().GetState(studentId)

	if err != nil {
		return nil, fmt.Errorf("failed to read the Student with The Id, %v", err)
	}

	if bytes == nil {
		return nil, fmt.Errorf("the student %s does not exist ", studentId)
	}

	var student Student
	err = json.Unmarshal(bytes, &student)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarhal from the world state to type student")
	}

	return &student, nil

}

func (u *UniContract) GetStudentByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Student, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the  data by range. %s", err)
	}
	defer resultsIterator.Close()

	return studentResultIteratorFunction(resultsIterator)

}

func (c *UniContract) GetAllStudent(ctx contractapi.TransactionContextInterface) ([]*Student, error) {

	queryString := `{"selector":{"EntityType":"degree"}, "sort":[{ "percentage": "desc"}]}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return studentResultIteratorFunction(resultsIterator)
}

func studentResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Student, error) {
	var students []*Student
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var student Student
		err = json.Unmarshal(queryResult.Value, &student)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		students = append(students, &student)
	}

	return students, nil
}

func (u *UniContract) GetStudentHistory(ctx contractapi.TransactionContextInterface, studentId string) ([]*HistoryQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(studentId)
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

		var student Student
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &student)
			if err != nil {
				return nil, err
			}
		} else {
			student = Student{
				StudentId: studentId,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			Txid:      response.TxId,
			TimeStanp: formattedTime,
			Record:    &student,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}
func (c *UniContract) GetStudentWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	queryString := `{"selector":{"EntityType":"degree"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the car records. %s", err)
	}
	defer resultsIterator.Close()

	cars, err := studentResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the car records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             cars,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}
